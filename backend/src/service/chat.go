package service

import (
	"context"
	"fmt"
	"time"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/mess"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/mq/producer"
	"zzc/fall-chat/src/utils"
)

type Chat struct {
	BaseService
	ws                  *Ws
	producer            *producer.UpdateConvMember
	conversationService *Conversation
	messageSeqService   *MessageSeq
}

func NewChat(
	ws *Ws,
	producer *producer.UpdateConvMember,
	conversationService *Conversation,
	messageSeqService *MessageSeq,
) *Chat {
	bs := BaseService{name: "ChatService"}
	bs.Inject()
	return &Chat{
		BaseService:         bs,
		ws:                  ws,
		producer:            producer,
		conversationService: conversationService,
		messageSeqService:   messageSeqService,
	}
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
//
// id => client.ID
//
// BUG:客户端连续发送消息，接收端有时会断
func (cs *Chat) HandleChatMessage(ctx context.Context, clientID string, chatData *request.ChatData) error {
	if clientID != chatData.From {
		return cs.Wrap("", errs.ErrWSClientIDSenderIDNotMatch, "from", chatData.From)
	}

	temConv := model.IMConversation{ConvType: int8(chatData.ConvType)}
	if !temConv.IsInvalidConvType() {
		return cs.Wrap("", errs.ErrWSChatConvType, "conversation_type", chatData.ConvType)
	}
	if !model.IsInvalidIMMsgType(chatData.MsgType) {
		return cs.Wrap("", errs.ErrWSChatMessageType, "message_type", chatData.MsgType)
	}

	rawID, err := model.DecodeID(clientID)
	if err != nil {
		return cs.Wrap("发送方id解码", err)
	}
	rawToID, err := model.DecodeID(chatData.To)
	if err != nil {
		return cs.Wrap("接受方id解码", err, "toID", chatData.To)
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
	conv, err := cs.conversationService.FindBrief(ctx, singleKey)
	if err != nil {
		return cs.Wrap("找用户会话", err, "to", chatData.To)
	}

	if conv == nil {
		// 群类型，不存在报错
		if temConv.IsGroupConType() {
			return cs.Wrap("", errs.ErrWSChatConvGroupNotExists, "single_key", singleKey)
		}
		err = cs.conversationService.CreateSingleChat(ctx, rawID, rawToID, chatData.ConvType)
		if err != nil {
			return cs.Wrap("创建用户会话", err, "to", chatData.To)
		}

		// 重新找
		conv, err = cs.conversationService.FindBrief(ctx, singleKey)
		if err != nil {
			return cs.Wrap("创建后找用户会话", err, "to", chatData.To)
		}

		utils.LogDebug("单聊会话创建成功: %v", conv)
	}

	// 生成会话 next sequence
	newSeq, err := cs.messageSeqService.NextSeq(ctx, int64(conv.ID))
	if err != nil {
		return cs.Wrap("生成消息seq", err, "to", chatData.To)
	}

	lastMessage := model.IMMessage{
		ConvID:   int64(conv.ID),
		SenderID: rawID,
		Content:  chatData.Content,
		SendTime: time.Now(),
		MsgType:  int8(chatData.MsgType),
		Seq:      newSeq,
	}
	err = model.InjectID(&lastMessage.Model)
	if err != nil {
		return cs.Wrap("生成消息id", err, "to", chatData.To)
	}

	// 更新会话 最后一条消息 异步
	ucMess := mess.UpdateConv{
		ConvID:      int64(conv.ID),
		ConvType:    int(conv.ConvType),
		PeerID:      rawToID,
		MsgID:       int64(lastMessage.ID),
		MsgSenderID: rawID,
		MsgPreview:  lastMessage.GetPreview(),
		MsgTime:     lastMessage.SendTime,
		Seq:         lastMessage.Seq,
	}
	cs.producer.Publish(ctx, conv.PublicID, &ucMess)

	// 创建消息
	err = cs.ws.store.WSRepos().IMMessage.CreateOne(ctx, &lastMessage)
	if err != nil {
		return cs.Wrap("创建消息", err, "to", chatData.To)
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

	cs.ws.SendMessage(ctx, clientID, msg)

	// 给对方发送
	replyData := response.ChatReply{
		ConvID:   conv.PublicID,
		To:       chatData.To,
		MsgID:    lastMessage.PublicID,
		SenderID: chatData.From,
		Content:  chatData.Content,
		MsgSeq:   lastMessage.Seq,
		MsgType:  chatData.MsgType,
	}

	if temConv.IsGroupConType() {
		cs.ws.SendGroupMessage(ctx, chatData.From, rawToID, &replyData)
	} else {
		cs.ws.SendMessageByMQ(ctx, &replyData, replyData.To)
	}

	return nil
}
