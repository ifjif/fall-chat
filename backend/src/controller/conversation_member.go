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

type ConversationMember struct {
	Base
	convMembService *service.ConversationMember
}

func NewConversationMember(name string) *ConversationMember {
	bs := Base{name: name}
	return &ConversationMember{
		Base: bs,
	}
}

func (cmc *ConversationMember) Inject() {
	cms := service.GlobalService[service.CONVERSATION_MEMBER_SERVICE_NAME].(*service.ConversationMember)
	cmc.convMembService = cms
}

func (cmc *ConversationMember) SetRoute(r fiber.Router) {
	group := r.Group("/conversationMember")

	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Get("/listBrief4Conv", cmc.listBrief4Conv)
	authGroup.Post("/listBriefBatch", cmc.listBriefBatch)
	authGroup.Post("/commit/read", cmc.commitRead)
}

func (cmc *ConversationMember) listBrief4Conv(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Query[request.ConvMembListBrief4Conv](c)
	if err != nil {
		return cmc.Wrap("query bind faild.", err, "id", id)
	}

	qu, err := converter.ToConvMembListBrief4ConvQuery(req)
	if err != nil {
		return err
	}
	qu.UserID = id
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

func (cmc *ConversationMember) listBriefBatch(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.ConvMembListBriefBatch](c)
	if err != nil {
		return cmc.Wrap("body bind faild.", err, "id", id)
	}

	qu, err := converter.ToConvMembListBriefBatchQuery(req)
	if err != nil {
		return err
	}
	qu.UserID = id
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

func (cmc *ConversationMember) commitRead(c fiber.Ctx) error {
	userID := ctxkey.GetUserID(c)

	req, err := bind.Body[request.ConvMembCommitRead](c)
	if err != nil {
		return cmc.BindBoyErr(err)
	}

	cmd, err := converter.ToConvMembCommitReadCommand(req, userID)
	if err != nil {
		return err
	}

	err = cmc.convMembService.CommitRead(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.New(c).
		OK().
		JSON(nil)
}
