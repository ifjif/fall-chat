package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/mess"
	messagequeue "zzc/fall-chat/src/message_queue"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte // 发送通道，防止并发写
}

func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
		Send: make(chan []byte, 1024),
	}
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

func (h *Hub) AddClient(id string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	client := NewClient(id, conn)
	h.clients[id] = client
	// 启动写通道
	go client.WritePump()
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

func (ws *WsService) SendMessage(id string, msg response.WSMessage) {
	client, ok := ws.hub.GetClient(id)
	if !ok {
		ws.LogWarn("SendMessage", nil, "id", id, "reason", "客户端不存在")
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		ws.LogError("SendMessage", err, "id", id, "reason", "消息序列化失败")
		return
	}

	client.SendMesage(data)
}

func (ws *WsService) Add(id string, conn *websocket.Conn) {
	ws.hub.AddClient(id, conn)
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
func (ws *WsService) HandleChatMessage(id string, chatData *request.ChatData) error {
	if id != chatData.From {
		ws.LogWarn("handleChatMessage", nil, "id", id, "from", chatData.From,
			"reason", "消息携带的发送id与连接的id不匹配")
		// TODO: 发送错误
		return nil
	}

	if chatData.Type != 0 && chatData.Type != 1 {
		ws.LogWarn("handleChatMessage", nil, "id", id, "reason", "会话类型错误")
		// TODO: 发送错误
		return nil
	}

	rawID, err := ws.sf.Decode(id)
	if err != nil {
		ws.LogError("handleChatMessage", err, "id", id, "reason", "发送方id解码失败")
		// TODO: 发送错误
		return err
	}
	rawToID, err := ws.sf.Decode(chatData.To)
	if err != nil {
		ws.LogError("handleChatMessage", err, "id", id, "toID", chatData.To,
			"reason", "接受方id解码失败")
		// TODO: 发送错误
		return err
	}
	minID, maxID := rawID, rawToID
	if rawID > rawToID {
		minID, maxID = maxID, minID
	}
	singleKey := fmt.Sprintf("%d:%d", minID, maxID)

	// 找会话
	conv, err := ws.dao.FindConversationBySingleKey(singleKey)
	if err != nil {
		ws.LogWarn("createConversation", err, "id", id, "to", chatData.To, "reason", "找用户会话失败")
		// TODO: 发送错误
		return err
	}

	if conv == nil {
		conv, err = ws.createConversation(rawID, rawToID, singleKey, chatData)
		if err != nil {
			// TODO: 发送错误
			return err
		}
	}

	newSeq, err := ws.messageSeqService.NextSeq(int64(conv.ID))
	if err != nil {
		ws.LogError("createConversation", err, "id", id, "to", chatData.To,
			"reason", "生成消息seq失败")
		// TODO: 发送错误
		return err
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
		ws.LogError("createConversation", err, "id", id, "to", chatData.To, "reason", "生成消息id失败")
		// TODO: 发送错误
		return err
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
		ws.LogError("handleChatMessage", err, "id", id, "to", chatData.To, "reason", "创建消息失败")
		// TODO: 发送错误
		return err
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

	ws.SendMessage(id, msg)

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
	ws.SendMessage(chatData.To, msg2)

	return nil
}

func (ws *WsService) createConversation(id int64, toID int64, singleKey string, message *request.ChatData) (*model.IMConversation, error) {
	conv := &model.IMConversation{
		SingleKey: singleKey,
		ConvType:  int8(message.Type),
	}
	err := utils.InjectID4Model(ws.sf, &conv.Model)
	if err != nil {
		ws.LogError("createConversation", err, "id", id, "to", message.To, "reason", "生成会话id失败")
		// TODO: 发送错误
		return nil, err
	}
	var members []*model.IMConversationMember

	// 找发送者信息
	senderUserProfile, err := ws.dao.FindUserProfile(id)
	if err != nil {
		ws.LogWarn("createConversation", err, "id", id, "reason", "找发送者信息失败")
		// TODO: 发送错误
		return nil, err
	}

	if senderUserProfile == nil {
		ws.LogError("createConversation", nil, "id", id, "reason", "发送者不存在")
		// TODO: 发送错误
		return nil, nil
	}

	// 找接受者信息
	receiverProfile, err := ws.dao.FindUserProfile(toID)
	if err != nil {
		ws.LogError("createConversation", err, "id", id, "to", toID, "reason", "找接受者信息失败")
		// TODO: 发送错误
		return nil, err
	}

	if receiverProfile == nil {
		ws.LogError("createConversation", nil, "id", id, "to", toID, "reason", "接受者不存在")
		// TODO: 发送错误
		return nil, nil
	}

	ma := &model.IMConversationMember{
		ConvID:       int64(conv.ID),
		UserID:       id,
		PeerID:       toID,
		PeerAvatar:   receiverProfile.Avatar,
		PeerNickname: receiverProfile.Nickname,
	}

	err = utils.InjectID4Model(ws.sf, &ma.Model)
	if err != nil {
		ws.LogError("createConversation", err, "id", id, "to", message.To, "reason", "生成成员id失败")
		return nil, err
	}
	members = append(members, ma)

	mb := &model.IMConversationMember{
		ConvID:       int64(conv.ID),
		UserID:       toID,
		PeerID:       id,
		PeerAvatar:   senderUserProfile.Avatar,
		PeerNickname: senderUserProfile.Nickname,
	}
	err = utils.InjectID4Model(ws.sf, &mb.Model)
	if err != nil {
		ws.LogError("createConversation", err, "id", id, "to", message.To, "reason", "生成成员id失败")
		return nil, err
	}
	members = append(members, mb)

	exists := false
	reason := ""
	// 创建会话，创建成员
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
		ws.LogWarn("createConversation", err, "reason", reason)
		return nil, err
	}

	return conv, nil
}
