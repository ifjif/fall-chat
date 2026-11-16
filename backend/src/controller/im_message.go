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

type IMMessage struct {
	Base
	imMessageService *service.IMMessage
}

func NewIMMessage(name string) *IMMessage {
	bc := Base{name: name}
	return &IMMessage{
		Base: bc,
	}
}

func (immc *IMMessage) Inject() {
	s := service.GlobalService[service.IMMESSAGE_SERVICE_NAME].(*service.IMMessage)
	immc.imMessageService = s
}

func (immc *IMMessage) SetRoute(r fiber.Router) {
	group := r.Group("/imMessage")
	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Post("/list", immc.list)
	authGroup.Post("/patchList", immc.patchList)
}

func (immc *IMMessage) list(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)
	req, err := bind.Body[request.ListIMMessage](c)
	if err != nil {
		return immc.Wrap("body bind", err)
	}

	qu, err := converter.ToListIMMessageQuery(req)
	if err != nil {
		return err
	}
	qu.UserID = id
	qu.Username = username

	result, err := immc.imMessageService.List(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}

func (immc *IMMessage) patchList(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.PatchListIMMessage](c)
	if err != nil {
		return immc.BindBoyErr(err)
	}
	if *req.MaxSeq < *req.MinSeq {
		return immc.Wrap("maxseq < minseq", fiber.ErrBadRequest,
			"max_seq", req.MaxSeq,
			"min_seq", req.MinSeq)
	}

	if (*req.MaxSeq - *req.MinSeq) > 50 {
		return immc.Wrap("[max_seq,min_seq] out of bound: 50", fiber.ErrBadRequest,
			"max_seq", req.MaxSeq,
			"min_seq", req.MinSeq)
	}

	qu, err := converter.ToPatchListIMMessageQuery(req)
	if err != nil {
		return err
	}
	qu.UserID = uint64(id)
	qu.Username = username

	result, err := immc.imMessageService.PatchList(c.Context(), qu)
	if err != nil {
		return err
	}
	return resp.New(c).
		OK().
		JSON(result)
}
