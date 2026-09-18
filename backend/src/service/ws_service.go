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
	groups  map[string]map[string]bool
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		groups:  make(map[string]map[string]bool),
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

func (h *Hub) AddGroups(id string, groupIds []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, groupId := range groupIds {
		if h.groups[groupId] == nil {
			h.groups[groupId] = make(map[string]bool)
		}
		h.groups[groupId][id] = true
	}
}

func (h *Hub) RemoveClient(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, id)
	h.removeGroup(id)
}

func (h *Hub) removeGroup(id string) {
	for gID, group := range h.groups {
		delete(group, id)

		if len(group) == 0 {
			delete(h.groups, gID)
		}
	}
}

func (h *Hub) GetClient(id string) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, exists := h.clients[id]
	return client, exists
}

func (h *Hub) GetGroup(id string) (map[string]bool, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	group, exists := h.groups[id]
	return group, exists
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

	messageGroup := messagequeue.GetDispatcher().GetOrCreateGroup("join-group", "join-group-group")
	messageGroup.RegisterHandler(ws.joinGroupHandler)
	messageGroup.StartInstance()
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

func (ws *WsService) SendGroupMessage(ctx context.Context, senderID, groupID string, msg response.WSMessage) {
	logger := ctxkey.GetLogger(ctx)
	group, ok := ws.hub.GetGroup(groupID)
	if !ok {
		logger.Info("websocket", "op", "发送群消息", "reason", "群不存在，消息丢弃")
		return
	}
	// 除发送者外的所有
	for k := range group {
		if k != senderID {
			ws.SendMessage(ctx, k, msg)
		}
	}
}

func (ws *WsService) Add(ctx context.Context, id string, conn *websocket.Conn) (*Client, error) {
	client := ws.hub.AddClient(ctx, id, conn)

	// 找用户加入的群
	groupIds, err := ws.fetchGroupConvIDs(id)
	if err != nil {
		return nil, err
	}

	if len(groupIds) != 0 {
		ws.hub.AddGroups(id, groupIds)
	}

	return client, err
}

//	type ChatData struct {
//		From        string `json:"from"` // 发送者id
//		To          string `json:"to"`   // 接收者id
//		Conv        string `json:"conv"` // 会话id
//		Content     string `json:"content"`
//		ConvType    int    `json:"conv_type"`     // 会话类型 0:单 1:群
//		MsgType     int    `json:"msg_type"`      // 消息类型 1:文本 2:图片 3:语言 4:视频
//		ClientMsgID string `json:"client_msg_id"` // 本地的临时id
//	}
func (ws *WsService) HandleChatMessage(client *Client, chatData *request.ChatData) error {
	id := client.ID
	if id != chatData.From {
		return ws.Wrap("", errs.ErrWSClientIDSenderIDNotMatch, "from", chatData.From)
	}

	temConv := model.IMConversation{ConvType: int8(chatData.ConvType)}
	if !temConv.IsInvalidConvType() {
		return ws.Wrap("", errs.ErrWSChatConvType, "conversation_type", chatData.ConvType)
	}
	if !model.IsInvalidIMMsgType(chatData.MsgType) {
		return ws.Wrap("", errs.ErrWSChatMessageType, "message_type", chatData.MsgType)
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
	if temConv.IsGroupConType() {
		singleKey = chatData.To
	}

	// 找会话
	conv, err := ws.dao.FindConversationBySingleKey(singleKey)
	if err != nil {
		return ws.Wrap("找用户会话", err, "to", chatData.To)
	}

	if conv == nil {
		// INFO:群类型，不存在报错
		if temConv.IsGroupConType() {
			return ws.Wrap("", errs.ErrWSChatConvGroupNotExists, "single_key", singleKey)
		}
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
		MsgType:  int8(chatData.MsgType),
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
		return ws.Wrap("json.Marshal 更新会话消息数据", err)
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
	replyMsg := response.WSMessage{
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

	if temConv.IsGroupConType() {
		ws.SendGroupMessage(client.ctx, chatData.From, chatData.To, replyMsg)
	} else {
		ws.SendMessage(client.ctx, chatData.To, replyMsg)
	}

	return nil
}

// 单聊的首次会话创建
func (ws *WsService) createConversation(id int64, toID int64, singleKey string, message *request.ChatData) (*model.IMConversation, error) {
	conv := &model.IMConversation{
		SingleKey: singleKey,
		ConvType:  int8(message.ConvType),
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
	ma := model.GenerateConvMember(int64(conv.ID), id, toID)
	err = utils.InjectID4Model(ws.sf, &ma.Model)
	if err != nil {
		return nil, ws.Wrap("生成发送者成员id", err)
	}
	members = append(members, ma)

	// 接收者
	mb := model.GenerateConvMember(int64(conv.ID), toID, id)
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

func (ws *WsService) fetchGroupConvIDs(id string) ([]string, error) {
	rawID, err := ws.sf.Decode(id)
	if err != nil {
		return nil, ws.Wrap("连接id解码失败", err)
	}

	groupIDs, err := ws.dao.FindGroupConvIDsByUserID(rawID)
	if err != nil {
		return nil, ws.Wrap("找会话组id", err)
	}

	return groupIDs, nil
}

func (ws *WsService) joinGroupHandler(ctx context.Context, message *kafka.Message, commitFunc func()) error {
	var msg *mess.JoinGroupMess
	err := json.Unmarshal(message.Value, msg)
	if err != nil {
		// TODO: 记录错误
		return err
	}
	utils.LogInfo("收到加入群消息：%v\n", msg)
	return nil
}
