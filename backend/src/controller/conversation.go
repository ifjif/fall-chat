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

type ConversationController struct {
	BaseController
	conversationService *service.Conversation
}

func NewConversationController(name string) *ConversationController {
	bs := BaseController{name: name}

	return &ConversationController{
		BaseController: bs,
	}
}

func (cc *ConversationController) Inject() {
	cs := service.GlobalService[service.CONVERSATION_SERVICE_NAME].(*service.Conversation)
	cc.conversationService = cs
}

func (cc *ConversationController) SetRoute(r fiber.Router) {
	group := r.Group("/conversation")
	auth := group.Group("")
	auth.Use(middleware.AuthMiddleware)

	auth.Get("/list", cc.list)
	auth.Get("/single", cc.single)
}

func (cc *ConversationController) list(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	result, err := cc.conversationService.List(c.Context(), id, username)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}

func (cc *ConversationController) single(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Query[request.ConvSingle](c)
	if err != nil {
		return cc.Wrap("query bind faild.", err, "id", id)
	}

	qu := converter.ToConvSingleQuery(req)
	qu.ID = id
	qu.Username = username

	result, err := cc.conversationService.Single(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}
