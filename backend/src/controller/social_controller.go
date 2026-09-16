package controller

import (
	"zzc/fall-chat/src/bind"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/ctxkey"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/pkg/resp"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/service"

	"github.com/gofiber/fiber/v3"
)

type SocialController struct {
	BaseController
	socialService *service.SocialService
}

func NewSocialController(name string) *SocialController {
	bs := BaseController{name: name}
	return &SocialController{
		BaseController: bs,
	}
}

func (sc *SocialController) SetRoute(r fiber.Router) {
	group := r.Group("/social")
	auth := group.Group("")
	auth.Use(middleware.AuthMiddleware)

	auth.Post("/applyAddFriend", sc.applyAddFriend)
	auth.Post("/applyAddGroup", sc.applyAddGroup)
	auth.Post("/searchFriendRequest", sc.searchFriendRequest)
	auth.Post("/searchGroupRequest", sc.searchGroupRequest)
	auth.Post("/handleFriendRequest", sc.handleFriendRequest)
	auth.Post("/handleGroupRequest", sc.handleGroupRequest)
	auth.Post("/friendshipList", sc.friendshipList)
}

func (sc *SocialController) Inject() {
	ss := service.GlobalService[service.SOCIAL_SERVICE_NAME].(*service.SocialService)
	sc.socialService = ss
}

func (sc *SocialController) applyAddFriend(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.SocialAddFriendRequest](c)
	if err != nil {
		return sc.Wrap("body bind", err, "username", username)
	}

	cmd := converter.ToSocialAddFriendCommand(req)
	cmd.ID = id
	cmd.Username = username

	err = sc.socialService.ApplyAddFriend(cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("发送成功").
		JSON(nil)
}

func (sc *SocialController) applyAddGroup(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.SocialAddGroupRequestRequest](c)
	if err != nil {
		return sc.Wrap("bind body", err, "username", username)
	}

	if req.InvitedID != "" && req.ApplyType != 2 {
		return sc.Wrap("applyAddGroup", errs.SocialErrGroupRequestApplyType, "username", username)
	}

	cmd := converter.ToSocialAddGroupRequestCommand(req)
	cmd.ID = id
	cmd.Username = username

	err = sc.socialService.ApplyAddGroup(cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("发送成功").
		JSON(nil)
}

func (sc *SocialController) searchFriendRequest(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.SocialSearchFriendRequestRequest](c)
	if err != nil {
		return sc.Wrap("body bind", err, "username", username)
	}

	qu := converter.ToSocialSearchFriendRequestQuery(req)
	qu.ID = id
	qu.Username = username

	result, err := sc.socialService.SearchFriendRequest(qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}

func (sc *SocialController) searchGroupRequest(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.SocialSearchGroupRequestRequest](c)
	if err != nil {
		return sc.Wrap("body bind", err, "username", username)
	}

	qu := converter.ToSocialSearchGroupRequestQuery(req)
	qu.ID = id
	qu.Username = username

	result, err := sc.socialService.SearchGroupRequest(qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}

func (sc *SocialController) handleFriendRequest(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.SocialHandleFriendRequestRequest](c)
	if err != nil {
		return sc.Wrap("body bind", err, "username", username)
	}

	cmd := converter.ToSocialHandleFriendRequestCommand(req)
	cmd.ID = id
	cmd.Username = username

	err = sc.socialService.HandleFriendRequest(cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("操作成功").
		JSON(nil)
}

func (sc *SocialController) handleGroupRequest(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.SocialHandleGroupRequestRequest](c)
	if err != nil {
		return sc.Wrap("body bind faild.", err, "username", username)
	}

	cmd := converter.ToSocialHandleGroupRequestCommand(req)
	cmd.HandleByUserID = id
	cmd.Username = username

	err = sc.socialService.HandleGroupRequest(cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("操作成功").
		JSON(nil)
}

func (sc *SocialController) friendshipList(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	qu := query.SocialFriendshipListQuery{
		ID:       id,
		Username: username,
	}

	result, err := sc.socialService.FriendshipList(&qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}
