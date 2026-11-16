package service

import (
	"context"
	"fmt"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/repos"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/mq/consumer"
	"zzc/fall-chat/src/query"
)

type Conversation struct {
	BaseService
	store                    store.ConversationStore
	updateConvMemberConsumer *consumer.UpdateConvMember
}

func NewConversation(name string) *Conversation {
	bs := BaseService{name: name}
	return &Conversation{
		BaseService: bs,
	}
}

func (cs *Conversation) Inject() {
	cs.BaseService.Inject()
	cs.store = data.GetBaseDao()
	cs.updateConvMemberConsumer = consumer.GetUpdateConvMember()
}

func (cs *Conversation) List(ctx context.Context, id int64, name string) ([]*response.IMConversation, error) {
	// 携带
	// peer_avatar
	// peer_nickname
	var result []*response.IMConversation
	rows, err := cs.store.
		ConvRepos().
		Conversation.
		ListByUserID(ctx, id)
	if err != nil {
		return nil, cs.Wrap("找用户会话列表", err, "id", id, "username", name)
	}

	// 查缓存并附加信息
	miss, err := cs.store.
		ConvCache().
		AppendInfo(ctx, rows)
	if err != nil {
		// TODO:
	}

	var missIDs []int64
	for k := range miss {
		missIDs = append(missIDs, k)
	}

	// 找忽略的
	info, err := cs.store.
		ConvRepos().
		Conversation.
		GetUserOrGroupBaseInfoBatch(ctx, missIDs)

	for i := range info {
		peerID := info[i].ID
		idx := miss[peerID]
		avatar := info[i].Avatar
		nickname := info[i].PeerNickname
		rows[idx].PeerAvatar = avatar
		rows[idx].PeerNickname = nickname
	}

	// 设置忽略缓存
	err = cs.store.
		ConvCache().
		SetInfo(ctx, info)
	if err != nil {
		// TODO:
	}

	for i := range rows {
		row := rows[i]
		v := converter.ToIMConversationResponse(row)
		v.PeerID = model.EncodeID(row.PeerID)
		result = append(result, v)
	}
	return result, err
}

func (cs *Conversation) Single(ctx context.Context, qu *query.ConvSingle) (*response.IMConversation, error) {
	logger := ctxkey.GetLogger(ctx)
	// 携带
	// peer_avatar
	// peer_nickname
	convID := qu.ConvID

	row, err := cs.store.
		ConvRepos().
		Conversation.
		SingleByConvID(ctx, convID, qu.UserID)
	if err != nil {
		return nil, cs.Wrap("找用户的单条会话", err, "id", qu.UserID, "conv_id", convID)
	}

	// 没有会话
	if row == nil {
		return nil, nil
	}

	// 从缓存中找到额外信息，并附加
	miss, err := cs.store.
		ConvCache().
		AppendInfo(ctx, []*model.IMConversationRow{row})
	if err != nil {
		logger.Error("[redis]", "reason", "找会话额外信息失败", "error", err)
	}

	var ids []int64

	for k := range miss {
		ids = append(ids, k)
	}

	// 从数据库中找 忽略的信息
	res, err := cs.store.
		ConvRepos().
		Conversation.
		GetUserOrGroupBaseInfoBatch(ctx, ids)
	if err != nil {
		logger.Error("[database]", "reason", "从数据库中找额外信息", "error", err)
	}

	if len(res) > 0 {
		row.PeerAvatar = res[0].Avatar
		row.PeerNickname = res[0].PeerNickname

		// 设置缓存
		cs.store.
			ConvCache().
			SetInfo(ctx, res)
	}

	result := converter.ToIMConversationResponse(row)
	result.PeerID = model.EncodeID(row.PeerID)
	return result, nil
}

// 找会话简单信息
func (cs *Conversation) FindBrief(ctx context.Context, singleKey string) (*model.IMConversation, error) {
	// 从缓存中找
	conv, err := cs.store.ConvCache().GetBrief(ctx, singleKey)
	if err != nil {
		return nil, err
	}
	if conv != nil {
		return conv, err
	}

	// 从数据库
	conv, err = cs.store.
		ConvRepos().
		Conversation.
		FindBriefBySingleKey(ctx, singleKey)
	if err != nil {
		return nil, err
	}
	if conv != nil {
		cs.store.ConvCache().SetBrief(ctx, *conv)
	}

	return conv, nil
}

// 单聊的首次会话创建
func (cs *Conversation) CreateSingleChat(ctx context.Context, id int64, toID int64, ConvType int) error {
	minID, maxID := id, toID
	if id > toID {
		minID, maxID = maxID, minID
	}
	singleKey := fmt.Sprintf("%d:%d", minID, maxID)
	conv := &model.IMConversation{
		SingleKey: singleKey,
		ConvType:  int8(ConvType),
	}
	err := model.InjectID(&conv.Model)
	if err != nil {
		return cs.Wrap("生成会话id", err, "to", toID)
	}
	var members []*model.IMConversationMember

	// 找发送者信息
	senderUserProfile, err := cs.store.
		ConvRepos().
		User.
		FindNicknameAndAvatar(ctx, id)
	if err != nil {
		return cs.Wrap("找发送者信息", err)
	}

	if senderUserProfile == nil {
		return cs.Wrap("", errs.ErrWSChatSenderNotExists)
	}

	// 找接受者信息
	receiverProfile, err := cs.store.
		ConvRepos().
		User.
		FindNicknameAndAvatar(ctx, toID)
	if err != nil {
		return cs.Wrap("找接受者信息", err, "to", toID)
	}

	if receiverProfile == nil {
		return cs.Wrap("", errs.ErrWSChatReceiverNotExists, "to", toID)
	}

	// 发送者
	ma := model.GenerateConvMember(int64(conv.ID), id, toID)
	err = model.InjectID(&ma.Model)
	if err != nil {
		return cs.Wrap("生成发送者成员id", err)
	}
	members = append(members, ma)

	// 接收者
	mb := model.GenerateConvMember(int64(conv.ID), toID, id)
	err = model.InjectID(&mb.Model)
	if err != nil {
		return cs.Wrap("生成接收者成员id", err, "to", toID)
	}
	members = append(members, mb)

	// 创建会话，创建成员，返回会话id
	exists := false
	reason := ""
	cs.store.WithinConvTX(ctx, func(repo repos.ConvRepositories) error {
		reason = "创建会话失败"
		exists, err = repo.Conversation.CreateOneDup(ctx, conv)
		if err != nil {
			return err
		}

		reason = "会话已存在"
		if exists {
			return nil
		}

		reason = "创建会话成员失败失败"
		err = repo.ConvMemb.Create(ctx, members)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return cs.Wrap(reason, err)
	}

	return nil
}
