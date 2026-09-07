package model

/*
*  id bigint primary key
*  public_id varchar(32) index
*  created_at datetime
*  updated_at datetime
*  deleted_at datetime
*
*  group_name varchar(100) index
*  avatar varchar(255)
*  creator_id bigint
*  search_flag tinyint 					// 1: 所有人可搜 0:仅通过id
*  join_mode tinyint 						// 0:自由加入 1:需审批 2:禁止加入
*  status tinyint 							// 1:正常 0:禁用
*  max_members int
*  muted_all tinyint 						// 全员禁言 0:否 1:是
*  description varchar(255)
*  notice text 									// 群公告
*
*
 */
type GroupInfo struct {
	Model
	GroupName   string `gorm:"size:100;index"`
	Description string `gorm:"size:255"`
	Avatar      string `gorm:"size:255"`
	CreatorID   int64
	MaxMembers  int    `gorm:"default:500"`
	SearchFlag  *int8  `gorm:"default:1"` // 1: 所有人可搜 0:仅通过id
	JoinMode    *int8  `gorm:"default:1"` // 0:自由加入 1:需审批 2:禁止加入
	Status      *int8  `gorm:"default:1"` // 1:正常 0:禁用
	MutedAll    *int8  `gorm:"default:0"` // 全员禁言 0:否 1:是
	Notice      string // 群公告
}

func (g *GroupInfo) IsDisabled() bool {
	return *g.Status == 0
}

func (g *GroupInfo) IsNotJoin() bool {
	return *g.JoinMode == 2
}
