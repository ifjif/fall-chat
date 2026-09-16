package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"
	"zzc/fall-chat/src/ctxkey"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/mess"
	messagequeue "zzc/fall-chat/src/message_queue"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type Client struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan []byte // 发送通道，防止并发写
	logger *slog.Logger
	ctx    context.Context
}

func NewClient(ctx context.Context, id string, conn *websocket.Conn) *Client {
	return &Client{
		ctx:  ctx,
		ID:   id,
		Conn: conn,
		Send: make(chan []byte, 1024),
	}
}

func (c *Client) GetContext() context.Context {
	return c.ctx
}

func (c *Client) GetLogger() *slog.Logger {
	if c.logger == nil {
		c.logger = ctxkey.GetLogger(c.ctx)
	}
	return c.logger
}

func (c *Client) LogError(err error, fields ...any) {
	fields = append(fields, "error", err.Error())
	c.GetLogger().Error("websocket", fields...)
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		// 并发安全的写
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}

func (c *Client) SendMesage(msg []byte) {
	select {
	case c.Send <- msg:
	default:
		// 通道满，客户端消费慢或断网了
		utils.LogWarn("客户端 %s 发送缓冲区已满，丢弃消息\n", c.ID)
		// 踢掉它
		c.Conn.Close()
	}
}

type Hub struct {
	clients map[string]*Client
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

func (h *Hub) AddClient(ctx context.Context, id string, conn *websocket.Conn) *Client {
	h.mu.Lock()
	defer h.mu.Unlock()
	client := NewClient(ctx, id, conn)
	h.clients[id] = client
	// 启动写通道
	go client.WritePump()

	return client
}

func (h *Hub) RemoveClient(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, id)
}

func (h *Hub) GetClient(id string) (*Client, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	client, exists := h.clients[id]
	return client, exists
}

type WsService struct {
	BaseService
	hub               *Hub
	dao               *data.WSDao
	messageSeqService *MessageSeqService
}

func NewWsService(name string) *WsService {
	bs := BaseService{name: name}
	return &WsService{
		BaseService: bs,
	}
}

func (ws *WsService) Inject() {
	ws.BaseService.Inject()
	ws.hub = NewHub()
	ws.dao = data.NewWSDao(ws.db)
	mss := GlobalService[MESSAGE_SEQ_SERVICE_NAME].(*MessageSeqService)
	ws.messageSeqService = mss
}

func (ws *WsService) Remove(id string) {
	ws.hub.RemoveClient(id)
}

func (ws *WsService) SendMessage(ctx context.Context, id string, msg response.WSMessage) {
	client, ok := ws.hub.GetClient(id)
	if !ok {
		logger := ctxkey.GetLogger(ctx)
		logger.Warn("websocket", "reason", "SendMessage 客户端不存在", "id", id)
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		err = ws.Wrap("SendMessage 消息json序列化", err)
		client.LogError(err)
		return
	}

	client.SendMesage(data)
}

func (ws *WsService) Add(ctx context.Context, id string, conn *websocket.Conn) *Client {
	return ws.hub.AddClient(ctx, id, conn)
}

/*
*

	type ChatData struct {
		From    string `json:"from"` // 发送者id
		To      string `json:"to"`   // 接收者id
		Conv    string `json:"conv"` // 会话id
		Content string `json:"content"`
		Type int			 `json:"type"` // 0:单聊 1:群
	}

*
*/
func (ws *WsService) HandleChatMessage(client *Client, chatData *request.ChatData) error {
	id := client.ID
	if id != chatData.From {
		return ws.Wrap("", errs.ErrWSClientIDSenderIDNotMatch, "from", chatData.From)
	}

	if chatData.Type != 0 && chatData.Type != 1 {
		return ws.Wrap("", errs.ErrWSChatConvType, "conversation_type", chatData.Type)
	}

	rawID, err := ws.sf.Decode(id)
	if err != nil {
		return ws.Wrap("发送方id解码", err)
	}
	rawToID, err := ws.sf.Decode(chatData.To)
	if err != nil {
		return ws.Wrap("接受方id解码", err, "toID", chatData.To)
	}
	minID, maxID := rawID, rawToID
	if rawID > rawToID {
		minID, maxID = maxID, minID
	}
	singleKey := fmt.Sprintf("%d:%d", minID, maxID)

	// 找会话
	conv, err := ws.dao.FindConversationBySingleKey(singleKey)
	if err != nil {
		return ws.Wrap("找用户会话", err, "to", chatData.To)
	}

	if conv == nil {
		conv, err = ws.createConversation(rawID, rawToID, singleKey, chatData)
		if err != nil {
			return ws.Wrap("创建用户会话", err, "to", chatData.To)
		}
	}

	newSeq, err := ws.messageSeqService.NextSeq(int64(conv.ID))
	if err != nil {
		return ws.Wrap("生成消息seq", err, "to", chatData.To)
	}

	preview := utils.TruncateBytes(chatData.Content, 256)
	lastMessage := model.IMMessage{
		SenderID: rawID,
		Content:  chatData.Content,
		SendTime: time.Now(), // TODO:
		MsgType:  0,          // TODO:
		Seq:      newSeq,
	}
	err = utils.InjectID4Model(ws.sf, &lastMessage.Model)
	if err != nil {
		return ws.Wrap("生成消息id", err, "to", chatData.To)
	}

	// 更新会话 最后一条消息 异步
	lastMessage.ConvID = int64(conv.ID)

	updateConv := &model.IMConversation{}
	updateConv.ID = conv.ID
	updateConv.SetLastMessage(&lastMessage)

	// TODO:消息发送待抽离
	ucMess := mess.UpdateConvMess{}
	ucMess.ConvID = int64(conv.ID)
	ucMess.ConvType = int(conv.ConvType)
	ucMess.PeerID = rawToID
	ucMess.MsgID = int64(lastMessage.ID)
	ucMess.MsgSenderID = rawID
	ucMess.MsgPreview = preview
	ucMess.MsgTime = lastMessage.SendTime
	ucMess.Seq = lastMessage.Seq
	messData, err := json.Marshal(ucMess)
	if err != nil {
		// TODO: 发送错误
		return err
	}
	go messagequeue.KWriter().WriteMessages(context.Background(), kafka.Message{
		Topic: "update-conv-member",
		Key:   []byte(conv.PublicID),
		Value: messData,
	})

	// 创建消息
	err = ws.dao.CreateIMMesage(&lastMessage)
	if err != nil {
		return ws.Wrap("创建消息", err, "to", chatData.To)
	}

	// 发送成功 ACK
	msg := response.WSMessage{
		Type:      "chat_ack",
		Timestamp: time.Now().UnixMilli(),
		Data: response.ChatACK{
			ConvID:      conv.PublicID,
			To:          chatData.To,
			ServerMsgID: lastMessage.PublicID,
			ClientMsgID: chatData.ClientMsgID,
			MsgSeq:      lastMessage.Seq,
		},
	}

	ws.SendMessage(client.ctx, id, msg)

	// 给对方发送
	msg2 := response.WSMessage{
		Type:      "chat_reply",
		Timestamp: time.Now().UnixMilli(),
		Data: response.ChatReply{
			ConvID:   conv.PublicID,
			To:       chatData.To,
			MsgID:    lastMessage.PublicID,
			SenderID: chatData.From,
			Content:  chatData.Content,
			MsgSeq:   lastMessage.Seq,
		},
	}
	ws.SendMessage(client.ctx, chatData.To, msg2)

	return nil
}

func (ws *WsService) createConversation(id int64, toID int64, singleKey string, message *request.ChatData) (*model.IMConversation, error) {
	conv := &model.IMConversation{
		SingleKey: singleKey,
		ConvType:  int8(message.Type),
	}
	err := utils.InjectID4Model(ws.sf, &conv.Model)
	if err != nil {
		return nil, ws.Wrap("生成会话id", err, "to", message.To)
	}
	var members []*model.IMConversationMember

	// 找发送者信息
	senderUserProfile, err := ws.dao.FindUserProfile(id)
	if err != nil {
		return nil, ws.Wrap("找发送者信息", err)
	}

	if senderUserProfile == nil {
		return nil, ws.Wrap("", errs.ErrWSChatSenderNotExists)
	}

	// 找接受者信息
	receiverProfile, err := ws.dao.FindUserProfile(toID)
	if err != nil {
		return nil, ws.Wrap("找接受者信息", err, "to", toID)
	}

	if receiverProfile == nil {
		return nil, ws.Wrap("", errs.ErrWSChatReceiverNotExists, "to", toID)
	}

	// 发送者
	ma := &model.IMConversationMember{
		ConvID:       int64(conv.ID),
		UserID:       id,
		PeerID:       toID,
		PeerAvatar:   receiverProfile.Avatar,
		PeerNickname: receiverProfile.Nickname,
	}

	err = utils.InjectID4Model(ws.sf, &ma.Model)
	if err != nil {
		return nil, ws.Wrap("生成发送者成员id", err)
	}
	members = append(members, ma)

	// 接收者
	mb := &model.IMConversationMember{
		ConvID:       int64(conv.ID),
		UserID:       toID,
		PeerID:       id,
		PeerAvatar:   senderUserProfile.Avatar,
		PeerNickname: senderUserProfile.Nickname,
	}
	err = utils.InjectID4Model(ws.sf, &mb.Model)
	if err != nil {
		return nil, ws.Wrap("生成接收者成员id", err, "to", message.To)
	}
	members = append(members, mb)

	// 创建会话，创建成员
	exists := false
	reason := ""
	err = ws.db.Transaction(func(tx *gorm.DB) error {
		reason = "创建会话失败"
		exists, err = ws.dao.CreateConversationWithTx(tx, conv)
		if err != nil {
			return err
		}

		if exists {
			reason = "插入冲突后找会话失败"
			conv, err = ws.dao.FindConversationBySingleKeyWithTx(tx, singleKey)
			if err != nil {
				return err
			}
			return nil
		}

		reason = "创建会话成员失败失败"
		err = ws.dao.CreateConversationMemberBatchWithTx(tx, members)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, ws.Wrap(reason, err)
	}

	return conv, nil
}
