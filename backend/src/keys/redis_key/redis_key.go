package rediskey

import (
	"fmt"
	"time"
)

const Namespace = "fallchat:"

// 集中管理模板
const (
	tplGrpcAddr          = "grpcaddr:%s"
	tplImUserConvGroupID = "im:user:conv:groupid:%d"
	tplConvSeq           = "conv:seq:%d"
	tplImUserNickname    = "im:user:nickname:%d"
	tplImUserAvatar      = "im:user:avatar:%d"
	tplImUserIngroupname = "im:user:ingroupname:%d:%d"
	tplImGroupNickname   = "im:group:nickname:%d"
	tplImGroupAvatar     = "im:group:avatar:%d"
	tplConv              = "conv:%s"
	tplGroupUserIDs      = "group:user:ids:%d"
)

// 集中管理ttl
const (
	tplGrpcAddrTTL          = 5 * time.Minute
	tplImUserConvGroupIDTTL = 8 * time.Hour
	tplConvSeqTTL           = 7 * 24 * time.Hour
	tplImUserNicknameTTL    = 1 * time.Hour
	tplImUserAvatarTTL      = 1 * time.Hour
	tplImUserIngroupnameTTL = 1 * time.Hour
	tplImGroupNicknameTTL   = 1 * time.Hour
	tplImGroupAvatarTTL     = 1 * time.Hour
	tplConvTTL              = 7 * 24 * time.Hour
	tplGroupUserIDsTTL      = 7 * 24 * time.Hour
)

// 带命名空间
func WithNs(tpl string, args ...any) string {
	return Namespace + fmt.Sprintf(tpl, args...)
}
