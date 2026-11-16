package service

import (
	"context"
	"encoding/json"
	"time"
	"zzc/fall-chat/src/conns"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/delivery"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/job"
	"zzc/fall-chat/src/keys/ctxkey"
	rediskey "zzc/fall-chat/src/keys/redis_key"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/mq/consumer"
	"zzc/fall-chat/src/mq/producer"
	"zzc/fall-chat/src/pool"
	taskmanager "zzc/fall-chat/src/task_manager"
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/contrib/v3/websocket"
)

type Ws struct {
	BaseService
	hub                   *conns.Hub
	store                 store.WSStore
	chatService           *Chat
	messageSeqService     *MessageSeq
	conversationService   *Conversation
	groupMemberService    *GroupMember
	chatConsumer          *consumer.Chat
	joinGroupConsumer     *consumer.JoinGroup
	roomInviteConsumer    *consumer.RoomInvite
	chatProducer          *producer.Chat
	updateConvMembProduer *producer.UpdateConvMember
}

func NewWs(name string) *Ws {
	bs := BaseService{name: name}
	return &Ws{
		BaseService: bs,
	}
}

func (ws *Ws) Inject() {
	ws.BaseService.Inject()
	ws.store = data.GetBaseDao()
	ws.hub = conns.NewHub()
	mss := GlobalService[MESSAGE_SEQ_SERVICE_NAME].(*MessageSeq)
	cs := GlobalService[CONVERSATION_SERVICE_NAME].(*Conversation)
	gms := GlobalService[GROUPMEMBER_SERVICE_NAME].(*GroupMember)
	ws.messageSeqService = mss
	ws.conversationService = cs
	ws.groupMemberService = gms

	deliverer := delivery.NewDeliverer(delivery.DelivererDeps{
		Hub:      ws.hub,
		Store:    ws.store,
		GrpcPool: pool.GrpcConnPool,
	})

	ws.chatProducer = producer.GetChat()
	ws.updateConvMembProduer = producer.GetUpdateConvMember()
	ws.chatService = NewChat(ws, ws.updateConvMembProduer, cs, mss)

	ws.chatConsumer = consumer.GetChat(deliverer)
	ws.joinGroupConsumer = consumer.GetJoinGroup()
	ws.roomInviteConsumer = consumer.GetRoomInvite(deliverer)
}

func (ws *Ws) Remove(ctx context.Context, id string) {
	ws.hub.RemoveClient(id)
	ws.store.WSCache().DelGrpcAddr(ctx, id)
}

// 注册grpc 地址
func (ws *Ws) RegisterGrpcAddr(ctx context.Context, id string) {
	ws.store.WSCache().AddGrpcAddr(ctx, id)
}

func (ws *Ws) SendMessage(ctx context.Context, id string, msg response.WSMessage) {
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

func (ws *Ws) SendBytes(ctx context.Context, id string, data []byte) {
	client, ok := ws.hub.GetClient(id)
	if !ok {
		logger := ctxkey.GetLogger(ctx)
		logger.Warn("websocket", "reason", "SendMessage 客户端不存在", "id", id)
		return
	}
	client.SendMessage(data)
}

// 向mq发送消息
// key是会话id
func (ws *Ws) SendMessageByMQ(ctx context.Context, reply *response.ChatReply, to string) {
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

	ws.chatProducer.Publish(ctx, reply.ConvID, to, d)
}

func (ws *Ws) SendGroupMessage(ctx context.Context, senderID string, groupID int64, msg *response.ChatReply) {
	logger := ctxkey.GetLogger(ctx)
	// 找群中所有用户id
	ids, err := ws.groupMemberService.FindUserIDs(ctx, groupID)
	if err != nil {
		logger.Error("websocket", "op", "找群中所有成员id", "groupID", groupID, "error", err)
		return
	}

	// 除发送者外的所有
	for _, id := range ids {
		encID := model.EncodeID(id)
		if encID != senderID {
			ws.SendMessageByMQ(ctx, msg, encID)
		}
	}
}

func (ws *Ws) Add(ctx context.Context, id string, conn *websocket.Conn) *conns.Client {
	client := ws.hub.AddClient(ctx, id, conn)
	return client
}

func (ws *Ws) HandlePingMessage(client *conns.Client) {
	// 延长 grpcAddr的 expire
	jo := job.NewUpdateGrpcaddrExpire(ws.store.WSCache())
	jo.ID = client.ID
	key := rediskey.GrpcAddr{}
	jo.Duration = key.TTL() / 2
	taskmanager.TaskManager.Submit(jo)
}

func (ws *Ws) HandleChatMessage(client *conns.Client, chatData *request.ChatData) error {
	return ws.chatService.HandleChatMessage(client.GetContext(), client.ID, chatData)
}
