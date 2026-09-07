package model

import (
	"strconv"
	"time"
)

/*
*  id bigint primary key
*  public_id varchar(32) unique key
*  created_at datetime
*  updated_at datetime
*  deleted_at datetime index
*
*  to_group_id bigint
*  from_user_id bigint
*  inviter_user_id bigint
*  handle_by_user_id bigint
*  apply_type tinyint  // 0:search 1:qrcode 2:invited
*  message varchar(255)
*  status tinyint    // 0:pending 1:approved 2:rejected 3:expired
*  uniqueKey varchar(50) // 当是 status(0/1) -> active ，否则 timestame str
*
*  unique key idx_active_apply_uk(to_group_id, from_user_id, uniqueKey)
*  key idx_group_id_status(to_group_id, status)
*  key idx_from_user_id_status(from_user_id, status)
*
 */

type GroupRequest struct {
	Model
	ToGroupID      int64 `gorm:"uniqueKey:idx_active_apply_uk;index:idx_group_id_status"`
	FromUserID     int64 `gorm:"uniqueKey:idx_active_apply_uk;index:idx_from_user_id_status"`
	InviterUserID  int64
	HandleByUserID int64
	ApplyType      int8   `gorm:"default:0"`
	Message        string `gorm:"size:255"`
	Status         int8   `gorm:"index:idx_group_id_status;index:idx_from_user_id_status;default:0"`
	UniqueKey      string `gorm:"uniqueKey:idx_active_apply_uk;size:50"`
}

var requestStatus = map[int8]string{
	0: "pending",
	1: "approved",
	2: "rejected",
	3: "expired",
}

var requestApplyType = map[int8]string{
	0: "search",
	1: "qrcode",
	2: "invited",
}

func (gr *GroupRequest) ActiveUniqueKey() {
	gr.UniqueKey = "active"
}

func (gr *GroupRequest) DeactiveUniqueKey() {
	timeStr := strconv.FormatInt(time.Now().UnixMilli(), 10)
	gr.UniqueKey = timeStr
}

func (gr *GroupRequest) GetStatusText() string {
	if text, ok := requestStatus[gr.Status]; ok {
		return text
	}

	return "unknown"
}

func (gr *GroupRequest) GetApplyTypeText() string {
	if text, ok := requestApplyType[gr.ApplyType]; ok {
		return text
	}

	return "unknown"
}

func (gr *GroupRequest) IsPendingStatus() bool {
	return gr.Status == 0
}

func (gr *GroupRequest) IsApprovedStatus() bool {
	return gr.Status == 1
}
