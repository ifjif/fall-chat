package controller

import (
	"zzc/fall-chat/src/bind"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/resp"
	"zzc/fall-chat/src/service"

	"github.com/gofiber/fiber/v3"
)

type GroupInfo struct {
	Base
	groupService *service.Group
}

func NewGroupInfo(name string) *GroupInfo {
	bs := Base{name: name}
	return &GroupInfo{
		Base: bs,
	}
}

func (gi *GroupInfo) Inject() {
	gs := service.GlobalService[service.GROUP_SERVICE_NAME].(*service.Group)
	gi.groupService = gs
}

func (gi *GroupInfo) SetRoute(r fiber.Router) {
	group := r.Group("/groupInfo")
	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Post("/createGroup", gi.createGroup)
	authGroup.Post("/groups", gi.groups)
}

func (gic *GroupInfo) createGroup(c fiber.Ctx) error {
	userID := ctxkey.GetUserID(c)
	userPublicID := ctxkey.GetUserPubID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.CreateGroup](c)
	if err != nil {
		return gic.Wrap("body绑定", err,
			"username", username,
		)
	}

	cmd := converter.ToCreateGroupCommand(req)
	cmd.CreatorID = userID
	cmd.Username = username
	cmd.CreatorPublicID = userPublicID

	err = gic.groupService.CreateGroup(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("创建成功").
		JSON(nil)
}

func (gic *GroupInfo) groups(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.GroupInfoGroups](c)
	if err != nil {
		return gic.Wrap("body bind", err, "username", username)
	}

	qu := converter.ToGroupInfoGroupsQuery(req)
	qu.ID = id
	qu.Username = username

	result, err := gic.groupService.Groups(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}
