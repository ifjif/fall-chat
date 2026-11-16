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

type FriendshipRequest struct {
	Base
	friendshipRequestService *service.FriendshipRequest
}

func NewFriendshipRequest(name string) *FriendshipRequest {
	bc := Base{name: name}
	return &FriendshipRequest{
		Base: bc,
	}
}

func (frc *FriendshipRequest) Inject() {
	s := service.GlobalService[service.FRIENDSHIPREQUEST_SERVICE_NAME].(*service.FriendshipRequest)
	frc.friendshipRequestService = s
}

func (frc *FriendshipRequest) SetRoute(r fiber.Router) {
	group := r.Group("/friendRequest")
	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Post("/apply", frc.apply)
	authGroup.Post("/list", frc.list)
	authGroup.Post("/handle", frc.handle)
}

func (frc *FriendshipRequest) apply(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.AddFriend](c)
	if err != nil {
		return frc.Wrap("body bind", err, "username", username)
	}

	cmd, err := converter.ToAddFriendCommand(req)
	if err != nil {
		return err
	}
	cmd.UserID = id
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

func (frc *FriendshipRequest) list(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.ListFriend](c)
	if err != nil {
		return frc.Wrap("body bind", err, "username", username)
	}

	qu := converter.ToListFriendRequestQuery(req)
	qu.UserID = uint64(id)
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

func (frc *FriendshipRequest) handle(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.HandleFriend](c)
	if err != nil {
		return frc.Wrap("body bind", err, "username", username)
	}

	cmd, err := converter.ToHandleFriendRequestCommand(req)
	if err != nil {
		return nil
	}
	cmd.UserID = id
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
