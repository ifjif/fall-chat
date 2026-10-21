package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToCreateGroupCommand(req *request.CreateGroup) *command.CreateGroup {
	return &command.CreateGroup{
		Name:        req.Name,
		Description: req.Description,
		MaxMembers:  req.MaxMembers,
		Avatar:      req.Avatar,
		SearchFlag:  req.SearchFlag,
		JoinMode:    req.JoinMode,
	}
}

func ToGroupInfoGroupsQuery(req *request.GroupInfoGroups) *query.GroupInfoGroups {
	return &query.GroupInfoGroups{
		Filter: req.Filter,
	}
}

func ToGroupInfoJoinedGroupsResult(m *model.GroupInfo) *response.GroupInfoJoinedGroupsResult {
	return &response.GroupInfoJoinedGroupsResult{
		ID:          m.PublicID,
		Name:        m.GroupName,
		Avatar:      m.Avatar,
		Description: m.Description,
	}
}
