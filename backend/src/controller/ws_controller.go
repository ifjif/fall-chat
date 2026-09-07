package controller

import (
	"encoding/json"
	"fmt"
	"net/url"
	"zzc/fall-chat/src/ctxkey"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/pkg/resp"
	"zzc/fall-chat/src/service"
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

type WSController struct {
	BaseController
	wsService *service.WsService
}

func NewWSController(name string) *WSController {
	bs := BaseController{name: name}
	return &WSController{
		BaseController: bs,
	}
}

func (wsc *WSController) Inject() {
	wsService := service.GlobalService[service.WS_SERVICE_NAME].(*service.WsService)
	wsc.wsService = wsService
}

func (wsc *WSController) SetRoute(r fiber.Router) {
	group := r.Group("/ws")

	auth := group.Group("")
	auth.Use(middleware.AuthMiddleware)
	auth.Get("/accessUrl", wsc.accessURL)

	upgrade := auth.Group("")
	upgrade.Use(middleware.WsUpgrade)

	upgrade.Get("/:id", wsc.websocketHandler())
}

func (wsc *WSController) accessURL(c fiber.Ctx) error {
	pID := ctxkey.GetUserPubID(c)
	username := ctxkey.GetUsername(c)

	token, _, err := utils.GenerateAccessToken(pID, username, "")
	if err != nil {
		wsc.LogWarn("accessUrl", err, "username", username,
			"reason", "生成accessToken失败")
		return err
	}
	protocol := "ws"
	host := "localhost"
	port := "8080"
	urlPrefix := "/api/v1"
	url := fmt.Sprintf("%s://%s:%s%s/ws/%s?token=%s",
		protocol, host, port, urlPrefix, pID, url.QueryEscape("Bearer "+token))

	c.Set("Cache-Control", "private, max-age=300")

	return resp.
		New(c).
		OK().
		JSON(url)
}

func (wsc *WSController) websocketHandler() func(fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		id := c.Params("id", "")
		if id == "" {
			wsc.LogWarn("websocketHandler", fiber.ErrUpgradeRequired, "reason", "用户id不存在")
			return
		}
		wsc.LogInfo("websocketHandler", "连接建立", id)
		wsc.wsService.Add(id, c)
		defer wsc.wsService.Remove(id)

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				wsc.LogWarn("websocketHandler", err, "reason", "Read error")
				break
			}

			var message request.WSMessage
			if err := json.Unmarshal(msg, &message); err != nil {
				wsc.LogWarn("websocketHandler", err, "user_public_id:", id,
					"reason", "消息解码为json错误")
				continue
			}

			wsc.wsService.HandleMessage(id, message)
		}
	})
}
