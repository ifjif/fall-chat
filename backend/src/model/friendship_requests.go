package model

/*
* id bigint primary key
* public_id varchar(32) unique key
* created_at datetime
* updated_at datetime
* deleted_at datetime index
*
* from_user_id bigint
* to_user_id bigint
* message varchar(255)
* status tinyint default 0 				//0:待处理 1:已通过 2:已拒绝 3:已过期
	apply_type tinyint default 0    // 申请类型 0:搜索申请 1:扫码 2:通过群

	pair_key varchar(255) generated always as(concat(least(from_user_id,to_user_id), '_', greatest(from_user_id,to_user_id))) stored
	unique_key  	varchar(50)				//动态唯一键，插入始终为 active, 当状态为2,3为时间戳
*
* index idx_from_to_user_id_status(from_user_id,to_user_id,status)
* index idx_to_user_id_status(to_user_id, status)
* unique key idx_active_apply_uk(pair_key,unique_key)
*
*/

var friendshipRequestStatusMap = map[int8]string{
	0: "pending",
	1: "approved",
	2: "rejected",
	3: "expired",
}

var friendshipRequestApplyType = map[int8]string{
	0: "search",
	1: "qrcode",
	2: "group",
}

type FriendshipRequest struct {
	Model
	FromUserID int64  `gorm:"index:idx_from_to_user_id_status"`
	ToUserID   int64  `gorm:"index:idx_from_to_user_id_status;index:idx_to_user_id_status"`
	Message    string `gorm:"size:255"`                                                               // 申请消息
	Status     int8   `gorm:"index:idx_from_to_user_id_status;index:idx_to_user_id_status;default:0"` // 0:待处理 1:已同意 2:已拒绝 3:已过期
	ApplyType  int8   `gorm:"default:0"`                                                              // 申请类型 0:搜索申请 1:扫码 2:通过群
	PairKey    string `gorm:"uniqueIndex:idx_active_apply_uk;type:varchar(255) generated always as(concat(least(from_user_id, to_user_id), '_', greatest(from_user_id, to_user_id))) stored"`
	UniqueKey  string `gorm:"size:50;uniqueIndex:idx_active_apply_uk"` // 动态唯一键，插入始终为 active, 当状态为2,3为时间戳
}

func GenerateFriendRequest(fromID, ToID int64, applyType int, Message string) *FriendshipRequest {
	return &FriendshipRequest{
		FromUserID: fromID,
		ToUserID:   ToID,
		Message:    Message,
		ApplyType:  int8(applyType),
		UniqueKey:  "active",
	}
}

func (f *FriendshipRequest) GetStatusText() string {
	if text, ok := friendshipRequestStatusMap[f.Status]; ok {
		return text
	}

	return "unknown"
}

func (f *FriendshipRequest) GetApplyTypeText() string {
	if text, ok := friendshipRequestApplyType[f.ApplyType]; ok {
		return text
	}

	return "unknown"
}
