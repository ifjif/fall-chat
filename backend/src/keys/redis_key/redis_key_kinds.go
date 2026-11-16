package rediskey

import "time"

// =============================
type GrpcAddr struct {
	ID string
}

func (k *GrpcAddr) String() string {
	return WithNs(tplGrpcAddr, k.ID)
}

func (k *GrpcAddr) TTL() time.Duration {
	return tplGrpcAddrTTL
}

// =============================
type ImUserConvGroupID struct {
	ID int64
}

func (k *ImUserConvGroupID) String() string {
	return WithNs(tplImUserConvGroupID, k.ID)
}

func (k *ImUserConvGroupID) TTL() time.Duration {
	return tplImUserConvGroupIDTTL
}

// =============================
type ConvSeq struct {
	ID int64
}

func (k *ConvSeq) String() string {
	return WithNs(tplConvSeq, k.ID)
}

func (k *ConvSeq) TTL() time.Duration {
	return tplConvSeqTTL
}

// =============================

type ImUserNickname struct {
	ID int64
}

func (k *ImUserNickname) String() string {
	return WithNs(tplImUserNickname, k.ID)
}

func (k *ImUserNickname) TTL() time.Duration {
	return tplImUserNicknameTTL
}

// =============================
type ImUserAvatar struct {
	ID int64
}

func (k *ImUserAvatar) String() string {
	return WithNs(tplImUserAvatar, k.ID)
}

func (k *ImUserAvatar) TTL() time.Duration {
	return tplImUserAvatarTTL
}

// =============================
type ImUserIngroupname struct {
	ConvID int64
	ID     int64
}

func (k *ImUserIngroupname) String() string {
	return WithNs(tplImUserIngroupname, k.ConvID, k.ID)
}

func (k *ImUserIngroupname) TTL() time.Duration {
	return tplImUserIngroupnameTTL
}

// =============================
type ImGroupNickname struct {
	ID int64
}

func (k *ImGroupNickname) String() string {
	return WithNs(tplImGroupNickname, k.ID)
}

func (k *ImGroupNickname) TTL() time.Duration {
	return tplImGroupNicknameTTL
}

// =============================
type ImGroupAvatar struct {
	ID int64
}

func (k *ImGroupAvatar) String() string {
	return WithNs(tplImGroupAvatar, k.ID)
}

func (k *ImGroupAvatar) TTL() time.Duration {
	return tplImGroupAvatarTTL
}

// =============================
// 一条会话基本信息
// {id,public_id,conv_type}
type Conv struct {
	SingleKey string
}

func (k *Conv) String() string {
	return WithNs(tplConv, k.SingleKey)
}

func (k *Conv) TTL() time.Duration {
	return tplConvTTL
}

// =============================
type GroupUserIDs struct {
	GroupID int64
}

func (k *GroupUserIDs) String() string {
	return WithNs(tplGroupUserIDs, k.GroupID)
}

func (k *GroupUserIDs) TTL() time.Duration {
	return tplGroupUserIDsTTL
}

// ==================================
type PublishEmojiItemCount struct {
	PackID uint64
	UserID uint64
}

func (k *PublishEmojiItemCount) String() string {
	return WithNs(tplPublishEmojiItemCount, k.PackID, k.UserID)
}

func (k *PublishEmojiItemCount) TTL() time.Duration {
	return tplPublishEmojiItemCountTTL
}

// =======================================
type ConvMembLastReadMsgSeq struct {
	UserID int64
	ConvID int64
}

func (k *ConvMembLastReadMsgSeq) String() string {
	return WithNs(tplConvMembLastReadMsgSeq, k.UserID, k.ConvID)
}

func (k *ConvMembLastReadMsgSeq) TTL() time.Duration {
	return tplConvMembLastReadMsgSeqTTL
}

// ===================================
type RtcRoom struct {
	RoomID string
}

func (k *RtcRoom) String() string {
	return WithNs(tplRtcRoom, k.RoomID)
}

func (k *RtcRoom) TTL() time.Duration {
	return tplRtcRoomTTL
}

// ==========================================
// 用户活跃的单聊房间
type RtcPeerRoom struct {
	UserID int64
}

func (k *RtcPeerRoom) String() string {
	return WithNs(tplRtcPeerRoom, k.UserID)
}

func (k *RtcPeerRoom) TTL() time.Duration {
	return tplRtcPeerRoomTTL
}
