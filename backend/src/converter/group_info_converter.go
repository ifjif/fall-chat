package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
)

func ToCreateGroupCommand(req *request.CreateGroupRequest) *command.CreateGroupCommand {
	return &command.CreateGroupCommand{
		Name:        req.Name,
		Description: req.Description,
		MaxMembers:  req.MaxMembers,
		Avatar:      req.Avatar,
		SearchFlag:  req.SearchFlag,
		JoinMode:    req.JoinMode,
	}
}
