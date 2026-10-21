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

type FriendshipRequestController struct {
	BaseController
	friendshipRequestService *service.FriendshipRequest
}

func NewFriendshipRequestController(name string) *FriendshipRequestController {
	bc := BaseController{name: name}
	return &FriendshipRequestController{
		BaseController: bc,
	}
}

func (frc *FriendshipRequestController) Inject() {
	s := service.GlobalService[service.FRIENDSHIPREQUEST_SERVICE_NAME].(*service.FriendshipRequest)
	frc.friendshipRequestService = s
}

func (frc *FriendshipRequestController) SetRoute(r fiber.Router) {
	group := r.Group("/friendRequest")
	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Post("/apply", frc.apply)
	authGroup.Post("/list", frc.list)
	authGroup.Post("/handle", frc.handle)
}

func (frc *FriendshipRequestController) apply(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.AddFriend](c)
	if err != nil {
		return frc.Wrap("body bind", err, "username", username)
	}

	cmd := converter.ToAddFriendCommand(req)
	cmd.ID = id
	cmd.Username = username

	err = frc.friendshipRequestService.Apply(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("发送成功").
		JSON(nil)
}

func (frc *FriendshipRequestController) list(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.ListFriend](c)
	if err != nil {
		return frc.Wrap("body bind", err, "username", username)
	}

	qu := converter.ToListFriendRequestQuery(req)
	qu.ID = id
	qu.Username = username

	result, err := frc.friendshipRequestService.List(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}

func (frc *FriendshipRequestController) handle(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.HandleFriend](c)
	if err != nil {
		return frc.Wrap("body bind", err, "username", username)
	}

	cmd := converter.ToHandleFriendRequestCommand(req)
	cmd.ID = id
	cmd.Username = username

	err = frc.friendshipRequestService.Handle(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("操作成功").
		JSON(nil)
}
