package controller

// TODO: 表情包的 图片上传需单走某些验证

import (
	"zzc/fall-chat/src/bind"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/resp"
	"zzc/fall-chat/src/service"

	"github.com/gofiber/fiber/v3"
)

type Emoji struct {
	Base
	emojiService *service.Emoji
}

func NewEmoji(name string) *Emoji {
	bc := Base{name: name}
	return &Emoji{
		Base: bc,
	}
}

func (ec *Emoji) Inject() {
	es := service.GlobalService[service.EMOJI_SERVICE_NAME].(*service.Emoji)
	ec.emojiService = es
}

func (ec *Emoji) SetRoute(r fiber.Router) {
	group := r.Group("/emoji")
	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Get("/system", ec.system)
	authGroup.Get("/pack/:id", ec.packDetail)
	authGroup.Post("/listPack", ec.listPack)
	authGroup.Post("/createPack", ec.createPack)
	authGroup.Post("/publish", ec.publish)
	authGroup.Post("/publish/verifyCount", ec.publishVerifyCount)
	authGroup.Post("/listItem", ec.listItem)
	authGroup.Post("/listPack/own", ec.ownListPack)
}

func (ec *Emoji) system(c fiber.Ctx) error {
	return nil
}

func (ec *Emoji) createPack(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.CreateEmojiPack](c)
	if err != nil {
		return ec.Wrap("body bind", err)
	}

	cmd := converter.ToCreateEmojiPackCommand(req)

	cmd.UserID = uint64(id)
	cmd.Username = username

	err = ec.emojiService.CreatePack(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.New(c).
		OK().
		Message("创建成功").
		JSON(nil)
}

func (ec *Emoji) packDetail(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.URI[request.EmojiPackDetail](c)
	if err != nil {
		return ec.Wrap("uri bind", err)
	}

	packageID, err := model.DecodeID(req.ID)
	if err != nil {
		return err
	}
	qu := &query.EmojiPackDetail{
		ID:        uint64(id),
		Username:  username,
		PackageID: packageID,
	}

	result, err := ec.emojiService.PackDetail(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.New(c).
		OK().
		JSON(result)
}

func (ec *Emoji) listPack(c fiber.Ctx) error {
	req, err := bind.Body[request.ListEmojiPack](c)
	if err != nil {
		return ec.Wrap("body bind", err)
	}

	qu := converter.ToListEmojiPackQuery(req)

	result, err := ec.emojiService.ListPack(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.New(c).
		OK().
		JSON(result)
}

func (ec *Emoji) publish(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.PublishEmoji](c)
	if err != nil {
		return ec.Wrap("body bind", err)
	}

	cmd, err := converter.ToPublishEmojiCommand(req)
	if err != nil {
		return err
	}
	cmd.UserID = uint64(id)
	cmd.Username = username

	err = ec.emojiService.Publish(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.New(c).
		OK().
		Message("发布成功").
		JSON(nil)
}

func (ec *Emoji) publishVerifyCount(c fiber.Ctx) error {
	userID := ctxkey.GetUserID(c)

	req, err := bind.Body[request.PublishEmojiVerifyCount](c)
	if err != nil {
		return ec.Wrap("body bind", err)
	}

	qu, err := converter.ToPublishEmojiVerifyCountQuery(req, uint64(userID))
	if err != nil {
		return err
	}

	err = ec.emojiService.PublishVerifyCount(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(nil)
}

func (ec *Emoji) listItem(c fiber.Ctx) error {
	req, err := bind.Body[request.ListEmojiItem](c)
	if err != nil {
		return ec.Wrap("body bind", err)
	}

	qu, err := converter.ToListEmojiItemQuery(req)
	if err != nil {
		return err
	}

	result, err := ec.emojiService.ListItem(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.New(c).
		OK().
		JSON(result)
}

func (ec *Emoji) ownListPack(c fiber.Ctx) error {
	userID := ctxkey.GetUserID(c)
	req, err := bind.Body[request.OwnListEmojiPack](c)
	if err != nil {
		return ec.Wrap("body bind", err)
	}

	qu := converter.ToOwnListEmojiPackQuery(req)
	qu.UserID = uint64(userID)

	result, err := ec.emojiService.OwnListPack(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}
