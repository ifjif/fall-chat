package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"zzc/fall-chat/src/conns"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/jobs"
	"zzc/fall-chat/src/keys/ctxkey"
	rediskey "zzc/fall-chat/src/keys/redis_key"
	"zzc/fall-chat/src/mess"
	messagequeue "zzc/fall-chat/src/message_queue"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/pool"
	"zzc/fall-chat/src/proto/chatpt"
	taskmanager "zzc/fall-chat/src/task_manager"
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type WsService struct {
	BaseService
	hub               *conns.Hub
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
	ws.hub = conns.NewHub()
	ws.dao = data.NewWSDao(ws.db)
	mss := GlobalService[MESSAGE_SEQ_SERVICE_NAME].(*MessageSeqService)
	ws.messageSeqService = mss

	joinGroup := messagequeue.GetDispatcher().GetOrCreateGroup("join-group", "join-group-group")
	joinGroup.RegisterHandler(ws.joinGroupHandler)
	joinGroup.StartInstance()

	chatGroup := messagequeue.GetDispatcher().GetOrCreateGroup("chat", "chat-group")
	chatGroup.RegisterHandler(ws.chatHandler)
	chatGroup.StartInstance()
}

// 注册grpc 地址
func (ws *WsService) RegisterGrpcAddr(id string) {
	ws.dao.AddGrpcAddr2Redis(id)
}

func (ws *WsService) Remove(id string) {
	ws.hub.RemoveClient(id)
	ws.dao.DelGrpcAddr2Redis(id)
}

// TODO:
// 向mq发送消息
// key是会话id
func (ws *WsService) SendMessageByMQ(reply *response.ChatReply, to string) {
	replyMsg := response.WSMessage{
		Type:      "chat_reply",
		Timestamp: time.Now().UnixMilli(),
		Data:      reply,
	}

	d, err := json.Marshal(replyMsg)
	if err != nil {
		utils.LogError("json 序列化错误：%v", err)
		return
	}

	msg := mess.ChatMess{
		Message: d,
		ConvID:  reply.ConvID,
		To:      to,
	}
	value, _ := json.Marshal(msg)

	go messagequeue.KWriter().WriteMessages(context.Background(), kafka.Message{
		Topic: "chat",
		Key:   []byte(msg.ConvID),
		Value: value,
	})
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

	client.SendMessage(data)
}

func (ws *WsService) SendBytes(ctx context.Context, id string, data []byte) {
	client, ok := ws.hub.GetClient(id)
	if !ok {
		logger := ctxkey.GetLogger(ctx)
		logger.Warn("websocket", "reason", "SendMessage 客户端不存在", "id", id)
		return
	}
	client.SendMessage(data)
}

func (ws *WsService) SendGroupMessage(ctx context.Context, senderID string, groupID int64, msg *response.ChatReply) {
	logger := ctxkey.GetLogger(ctx)
	// 找群中所有用户id
	ids, err := ws.dao.FindGroupUserIDs(groupID)
	if err != nil {
		logger.Error("websocket", "op", "找群中所有成员id", "groupID", groupID, "error", err)
		return
	}

	// 除发送者外的所有
	for _, id := range ids {
		encID := ws.sf.Encode(id)
		if encID != senderID {
			ws.SendMessageByMQ(msg, encID)
		}
	}
}

func (ws *WsService) Add(ctx context.Context, id string, conn *websocket.Conn) (*conns.Client, error) {
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

func (ws *WsService) HandlePingMessage(client *conns.Client) {
	// 延长 grpcAddr的 expire
	job := jobs.NewUpdateGrpcaddrExpireJob()
	job.ID = client.ID
	key := rediskey.GrpcAddr{}
	job.Duration = key.TTL() / 2
	taskmanager.TaskManager.Submit(job)
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
func (ws *WsService) HandleChatMessage(client *conns.Client, chatData *request.ChatData) error {
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

	ws.SendMessage(client.GetContext(), id, msg)

	// 给对方发送
	replyData := response.ChatReply{
		ConvID:   conv.PublicID,
		To:       chatData.To,
		MsgID:    lastMessage.PublicID,
		SenderID: chatData.From,
		Content:  chatData.Content,
		MsgSeq:   lastMessage.Seq,
	}

	if temConv.IsGroupConType() {
		ws.SendGroupMessage(client.GetContext(), chatData.From, rawToID, &replyData)
	} else {
		ws.SendMessageByMQ(&replyData, replyData.To)
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
	var msg mess.JoinGroupMess
	err := json.Unmarshal(message.Value, &msg)
	if err != nil {
		// TODO: 记录错误
		return err
	}
	utils.LogInfo("收到加入群消息：%v\n", msg)
	ws.hub.AddGroup(msg.UserID, msg.GroupID)
	commitFunc()
	return nil
}

func (ws *WsService) chatHandler(ctx context.Context, message *kafka.Message, commitFunc func()) error {
	var chatMess mess.ChatMess
	_ = json.Unmarshal(message.Value, &chatMess)
	utils.LogInfo("收到消息：to=%s\n", chatMess.To)

	commit := true
	defer func() {
		if commit {
			commitFunc()
		}
	}()

	_, ok := ws.hub.GetClient(chatMess.To)
	if !ok {
		// client.SendMessage(chatMess.Message)
		return nil
	}

	addr, err := ws.dao.GetGrpcAddrFromRedis(chatMess.To)
	if err != nil {
		utils.LogError("获取grpc地址从redis 出错：%v\n", err)
		return nil
	}

	// 不存在
	if addr == "" {
		return nil
	}

	utils.LogDebug("addr:%s", addr)
	// 进行拨号
	conn, err := pool.GrpcConnPool.Get(addr)
	if err != nil {
		log.Fatalf("dial failed: %v", err)
		return nil
	}

	stub := chatpt.NewChatClient(conn)
	stub.Chat(context.Background(), &chatpt.ChatpMessage{
		Content: chatMess.Message,
		ToID:    chatMess.To,
	})

	// 得到消息
	// 通过本地找
	// 没有找到，从redis找到rpc地址
	// 然后 直接发送
	return nil
}
