package controller

import (
	"zzc/fall-chat/src/ctxkey"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/pkg/resp"
	"zzc/fall-chat/src/service"

	"github.com/gofiber/fiber/v3"
)

type ConversationController struct {
	BaseController
	conversationService *service.ConversationService
}

func NewConversationController(name string) *ConversationController {
	bs := BaseController{name: name}

	return &ConversationController{
		BaseController: bs,
	}
}

func (cc *ConversationController) Inject() {
	cs := service.GlobalService[service.CONVERSATION_SERVICE_NAME].(*service.ConversationService)
	cc.conversationService = cs
}

func (cc *ConversationController) SetRoute(r fiber.Router) {
	group := r.Group("/conversation")
	auth := group.Group("")
	auth.Use(middleware.AuthMiddleware)

	auth.Get("/list", cc.list)
}

func (cc *ConversationController) list(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	result, err := cc.conversationService.List(id, username)
	if err != nil {
		cc.LogWarn("list", err)
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}
