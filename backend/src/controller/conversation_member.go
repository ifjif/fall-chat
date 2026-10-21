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

type ConversationMemberController struct {
	BaseController
	convMembService *service.ConversationMember
}

func NewConversationMemberController(name string) *ConversationMemberController {
	bs := BaseController{name: name}
	return &ConversationMemberController{
		BaseController: bs,
	}
}

func (cmc *ConversationMemberController) Inject() {
	cms := service.GlobalService[service.CONVERSATION_MEMBER_SERVICE_NAME].(*service.ConversationMember)
	cmc.convMembService = cms
}

func (cmc *ConversationMemberController) SetRoute(r fiber.Router) {
	group := r.Group("/conversationMember")

	auth := group.Group("")
	auth.Use(middleware.AuthMiddleware)

	auth.Get("/listBrief4Conv", cmc.listBrief4Conv)
	auth.Post("/listBriefBatch", cmc.listBriefBatch)
}

func (cmc *ConversationMemberController) listBrief4Conv(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Query[request.ConvMembListBrief4Conv](c)
	if err != nil {
		return cmc.Wrap("query bind faild.", err, "id", id)
	}

	qu := converter.ToConvMembListBrief4ConvQuery(req)
	qu.ID = id
	qu.Username = username

	result, err := cmc.convMembService.ListBrief4Conv(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}

func (cmc *ConversationMemberController) listBriefBatch(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.ConvMembListBriefBatch](c)
	if err != nil {
		return cmc.Wrap("body bind faild.", err, "id", id)
	}

	qu := converter.ToConvMembListBriefBatchQuery(req)
	qu.ID = id
	qu.Username = username

	result, err := cmc.convMembService.ListBriefBatch(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}
