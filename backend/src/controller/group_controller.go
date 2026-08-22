package controller

import (
	"zzc/fall-chat/src/bind"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/ctxkey"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/pkg/resp"
	"zzc/fall-chat/src/service"

	"github.com/gofiber/fiber/v3"
)

type GroupInfoController struct {
	BaseController
	groupService *service.GroupService
}

func NewGroupInfoController(name string) *GroupInfoController {
	bs := BaseController{name: name}
	return &GroupInfoController{
		BaseController: bs,
	}
}

func (gi *GroupInfoController) Inject() {
	gs := service.GlobalService[service.GROUP_SERVICE_NAME].(*service.GroupService)
	gi.groupService = gs
}

func (gi *GroupInfoController) SetRoute(r fiber.Router) {
	group := r.Group("/groupInfo")
	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Post("/createGroup", gi.createGroup)
}

func (gic *GroupInfoController) createGroup(c fiber.Ctx) error {
	userID := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.CreateGroupRequest](c, func(err error) {
		gic.LogWarn("createGroup body绑定失败", err,
			"username", username,
		)
	})
	if err != nil {
		return err
	}

	cmd := converter.ToCreateGroupCommand(req)
	cmd.CreatorID = userID
	cmd.Username = username

	err = gic.groupService.CreateGroup(cmd)
	if err != nil {
		gic.LogWarn("createGroup", err)
		return err
	}

	return resp.
		New(c).
		OK().
		Message("创建成功").
		JSON(nil)
}
