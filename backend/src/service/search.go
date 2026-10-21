package service

import (
	"context"
	"strings"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

type Search struct {
	BaseService
	store store.SearchStore
}

func NewSearch(name string) *Search {
	bs := BaseService{name: name}

	return &Search{
		BaseService: bs,
	}
}

func (ss *Search) Inject() {
	ss.BaseService.Inject()
	ss.store = data.GetBaseDao()
}

// 模糊匹配
//
// user -> 10
// group - 10
// user_cursor = last_user_id    infinity end
// group_cursor = last_group_id  infinity end
//
//	already_friend	显示"发消息"	friendships 中 status=1
//	request_pending	显示"等待验证"	friend_requests 中有待处理申请
//	can_add	显示"添加好友"	无任何关系记录
//	already_joined	显示"进入群聊"	group_members 中有记录
//	need_approval	显示"申请加入"	join_mode=1 且未加入
//	cannot_join	显示"不可加入"	join_mode=2
//	can_join 显示加入群聊 join_mode=1
func (ss *Search) UserOrGroup(ctx context.Context, qu *query.UserOrGroup) (*response.UserOrGroupResult, error) {
	userLimit := 10
	groupLimit := 10
	// 先找人
	id := qu.ID
	searchName := qu.Name
	keyword := "%" + searchName + "%"
	// userID_groupID
	nextCursor := qu.NextCursor
	userCursor, groupCursor := "inf", "inf"
	lastUserID, lastGroupID := int64(0), int64(0)

	cursors := strings.Split(nextCursor, "_")
	if len(cursors) >= 2 {
		userCursor = cursors[0]
		groupCursor = cursors[1]
	}

	if userCursor != "inf" && userCursor != "end" {
		lastUserID, _ = ss.sf.Decode(userCursor)
	}
	if groupCursor != "inf" && groupCursor != "end" {
		lastGroupID, _ = ss.sf.Decode(groupCursor)
	}

	finalResult := &response.UserOrGroupResult{}
	var results []*response.UserOrGroupDataResult

	// 找人
	var userRows []*model.UserOrGroupRelationRow
	if userCursor != "end" {
		result, err := ss.store.
			SearchRepos().
			User.
			SearchUsers(ctx, id, keyword, lastUserID, int64(userLimit))
		if err != nil {
			return nil, ss.Wrap("模糊找用户", err,
				"username", qu.Username,
				"searchName", searchName)
		}

		userRows = result
	}

	for i := range userRows {
		results = append(results, converter.ToUserOrGroupDataResult(userRows[i]))
	}
	nextUserID := "end"
	if len(userRows) == userLimit {
		newID := userRows[len(userRows)-1].ID
		if newID != userCursor {
			nextUserID = newID
		}
	}

	// 找群
	var groupRows []*model.UserOrGroupRelationRow

	if groupCursor != "end" {
		result, err := ss.store.
			SearchRepos().
			Group.
			SearchGroups(ctx, id, keyword, lastGroupID, int64(groupLimit))
		if err != nil {
			return nil, ss.Wrap("模糊找群", err,
				"username", qu.Username, "searchName", searchName)
		}

		groupRows = result
	}

	for i := range groupRows {
		results = append(results, converter.ToUserOrGroupDataResult(groupRows[i]))
	}

	nextGroupID := "end"
	if len(groupRows) == groupLimit {
		newID := groupRows[len(groupRows)-1].ID
		if newID != groupCursor {
			nextGroupID = newID
		}
	}

	cursor := nextUserID + "_" + nextGroupID

	finalResult.Result = results
	finalResult.NextCursor = cursor

	return finalResult, nil
}
