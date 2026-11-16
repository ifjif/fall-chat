package service

import (
	"context"
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/utils"
)

type ConversationMember struct {
	BaseService
	store store.ConvMembStore
}

func NewConversationMember(name string) *ConversationMember {
	bs := BaseService{name: name}

	return &ConversationMember{
		BaseService: bs,
	}
}

func (cms *ConversationMember) Inject() {
	cms.BaseService.Inject()
	cms.store = data.GetBaseDao()
}

func (cms *ConversationMember) ListBrief4Conv(ctx context.Context, qu *query.ConvMembListBrief4Conv) (map[string]*model.IMConvMembBriefRow, error) {
	convID := qu.ConvID

	rows, err := cms.store.
		ConvMembRepos().
		ConvMemb.
		ListBriefByConvID(ctx, convID)
	if err != nil {
		return nil, cms.Wrap("找会话成员基本信息", err, "id", qu.UserID)
	}

	result := make(map[string]*model.IMConvMembBriefRow, len(rows))
	for i := range rows {
		row := rows[i]
		result[row.PublicID] = &row
	}

	return result, nil
}

func (cms *ConversationMember) ListBriefBatch(ctx context.Context, qu *query.ConvMembListBriefBatch) (map[string]*model.IMConvMembBriefRow, error) {
	logger := ctxkey.GetLogger(ctx)
	// 从redis中找 批次找
	// 剩下没有的从数据库找
	// 数据库中找到的回写到redis
	// 最后返回
	convID := qu.ConvID
	memberIDs := qu.MemberIDs
	memberPIDs := qu.MemberPIDs

	// 从缓存中找
	rows, miss, err := cms.store.
		ConvMembCache().
		FindBriefInfo(ctx, convID, memberIDs, memberPIDs)
	if err != nil {
		logger.Error("[redis]",
			"reason", "从缓存中找会话成员失败",
			"conv_id", convID,
			"memberIDs", memberIDs,
			"error", err)
	}

	result := make(map[string]*model.IMConvMembBriefRow)
	for i := range rows {
		row := rows[i]
		result[row.PublicID] = &row
	}

	if len(miss) == 0 {
		return result, nil
	}

	// 找忽略的
	missRows, err := cms.store.
		ConvMembRepos().
		ConvMemb.
		ListBriefBatchByConvIDAndMemberIDs(ctx, convID, miss)
	if err != nil {
		return nil, cms.Wrap("找会话成员基本信息", err, "id", qu.UserID,
			"conv_id", qu.ConvID, "member_pids", memberPIDs)
	}

	// 存入缓存
	err = cms.store.
		ConvMembCache().
		SetBriefInfo(ctx, convID, missRows)
	if err != nil {
		logger.Error("[redis]", "reason", "将会话成员信息存入缓存失败", "conv_id", convID, "error", err)
	}

	for i := range missRows {
		row := missRows[i]
		result[row.PublicID] = &row
	}

	return result, nil
}

// 从redis中找，找到看是否大于，否结束
// 然后判断会话是否存在，消息是否存在，然后进行seq更新
// 然后更新redis计数
//
// BUG: 提交可能比 会话的快，导致用户刷新和拉 会话和它的last_msg_seq时，会话提交的比用户读的慢出现-情况
// 因为会话 提交是在延迟队列中进行的
// 只有当 UpdateConversation 有一个提交，这里也有一个提交时, 会触发比会话快
// UpdateConversation中是和会话匹配的，但这是是单个的
// 这里需拦截
//
// TODO: 并发操作 (同一个和 不同的) 放入job中进行
// 最终更新放入job中
func (cms *ConversationMember) CommitRead(ctx context.Context, cmd *command.ConvMembCommitRead) error {
	convID := cmd.ConvID
	userID := cmd.UserID
	msgID := cmd.MsgID
	msgSeq := cmd.MsgSeq
	seq := cms.store.ConvMembCache().GetLastReadMsgSeq(ctx, convID, userID)
	if msgSeq <= seq {
		utils.LogDebug("拦截: seq")
		return nil
	}

	exists, err := cms.store.
		ConvMembRepos().
		ConvMemb.
		Exists(ctx, convID, userID)
	if err != nil {
		return cms.Wrap("找用户是否存在会话的成员", err, "user_id", userID, "conv_id", convID)
	}
	if !exists {
		return errs.SocialErrGroupConvMembNotExists
	}

	seqInDB, err := cms.store.ConvMembRepos().ImMess.FindSeq(ctx, msgID, convID)
	if err != nil {
		return cms.Wrap("找im消息seq", err, "error", "conv_id", convID, "message_id", msgID)
	}
	if seqInDB != msgSeq {
		return errs.ImMessErrSeqNotMatch
	}

	success, err := cms.store.
		ConvMembRepos().
		ConvMemb.
		UpdateLastReadSeq(ctx, &model.IMConversationMember{
			ConvID:        convID,
			UserID:        userID,
			LastReadSeq:   msgSeq,
			LastReadMsgID: cmd.MsgID,
		})
	if err != nil {
		return err
	}

	if !success {
		// 找 last_msg_seq
		msgSeq, err = cms.store.
			ConvMembRepos().
			ConvMemb.
			FinsLastReadSeq(ctx, convID, userID)
		if err != nil {
			return err
		}
	}

	// 设置缓存
	if msgSeq > 0 {
		cms.store.ConvMembCache().SetLastReadMsgSeq(ctx, convID, userID, msgSeq)
	}
	return nil
}
