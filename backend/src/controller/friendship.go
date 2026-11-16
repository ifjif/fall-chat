package controller

import (
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/resp"
	"zzc/fall-chat/src/service"

	"github.com/gofiber/fiber/v3"
)

type Friendship struct {
	Base
	friendshipService *service.Friendship
}

func NewFriendship(name string) *Friendship {
	bc := Base{name: name}
	return &Friendship{
		Base: bc,
	}
}

func (fc *Friendship) Inject() {
	s := service.GlobalService[service.FRIENDSHIP_SERVICE_NAME].(*service.Friendship)
	fc.friendshipService = s
}

func (fc *Friendship) SetRoute(r fiber.Router) {
	group := r.Group("/friendship")
	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Post("/list", fc.friendshipList)
}

func (fc *Friendship) friendshipList(c fiber.Ctx) error {
	ctx := c.Context()
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	qu := query.FriendshipList{
		ID:       id,
		Username: username,
	}

	result, err := fc.friendshipService.FriendshipList(ctx, &qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}
