package controller

import (
	"zzc/fall-chat/src/bind"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/resp"
	"zzc/fall-chat/src/service"

	"github.com/gofiber/fiber/v3"
)

type GroupRequest struct {
	Base
	groupRequestService *service.GroupRequest
}

func NewGroupRequest(name string) *GroupRequest {
	bc := Base{name: name}
	return &GroupRequest{
		Base: bc,
	}
}

func (grc *GroupRequest) Inject() {
	s := service.GlobalService[service.GROUPREQUEST_SERVICE_NAME].(*service.GroupRequest)
	grc.groupRequestService = s
}

func (grc *GroupRequest) SetRoute(r fiber.Router) {
	group := r.Group("/groupRequest")
	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Post("/apply", grc.apply)
	authGroup.Post("/list", grc.list)
	authGroup.Post("/handle", grc.handle)
}

func (sc *GroupRequest) apply(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.AddGroup](c)
	if err != nil {
		return sc.Wrap("bind body", err, "username", username)
	}

	if req.InvitedID != "" && req.ApplyType != 2 {
		return sc.Wrap("applyAddGroup", errs.SocialErrGroupRequestApplyType, "username", username)
	}

	cmd, err := converter.ToAddGroupRequestCommand(req)
	if err != nil {
		return err
	}
	cmd.UserID = id
	cmd.Username = username

	err = sc.groupRequestService.Apply(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("发送成功").
		JSON(nil)
}

func (sc *GroupRequest) list(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.ListGroup](c)
	if err != nil {
		return sc.Wrap("body bind", err, "username", username)
	}

	qu := converter.ToListGroupRequestQuery(req)
	qu.UserID = id
	qu.Username = username

	result, err := sc.groupRequestService.List(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}

func (sc *GroupRequest) handle(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.HandleGroup](c)
	if err != nil {
		return sc.Wrap("body bind faild.", err, "username", username)
	}

	cmd, err := converter.ToHandleGroupRequestCommand(req)
	if err != nil {
		return err
	}
	cmd.HandleByUserID = id
	cmd.Username = username

	err = sc.groupRequestService.Handle(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("操作成功").
		JSON(nil)
}
