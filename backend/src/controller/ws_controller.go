package controller

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
	"zzc/fall-chat/src/ctxkey"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/pkg/resp"
	"zzc/fall-chat/src/service"
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

type MessageHandler func(id string, message *request.WSMessage)

type WSController struct {
	BaseController
	wsService       *service.WsService
	messageHandlers map[string]MessageHandler
}

func NewWSController(name string) *WSController {
	bs := BaseController{name: name}
	return &WSController{
		BaseController:  bs,
		messageHandlers: make(map[string]MessageHandler),
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

	wsc.messageHandlers["chat"] = wsc.chatMessageHandler
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

			wsc.dispatchMessage(id, &message)
		}
	})
}

func (wsc *WSController) dispatchMessage(id string, message *request.WSMessage) {
	mType := message.Type

	if mType != "ping" {
		wsc.LogInfo("dispatchMessage", "id:", id, "type:", mType)
	}

	messageHandler, ok := wsc.messageHandlers[mType]

	if ok {
		messageHandler(id, message)
		return
	}
	wsc.unmatchHandler(id, message)
}

func (wsc *WSController) unmatchHandler(id string, message *request.WSMessage) {
	wsc.errorHandler(id, "unmatch", errs.ErrWSNotFound)
}

func (wsc *WSController) errorHandler(id string, ty string, err error) {
	// 错误的类型
	errType := ty + "_error"

	respErr := response.WSMessage{
		Type:      errType,
		Timestamp: time.Now().UnixMilli(),
	}

	if _, ok := err.(*errs.WSError); ok {
		respErr.Data = err
	} else {
		respErr.Data = errs.ErrWSInternal
	}

	wsc.wsService.SendMessage(id, respErr)
}

func (wsc *WSController) chatMessageHandler(id string, message *request.WSMessage) {
	var chatData request.ChatData
	err := json.Unmarshal(message.Data, &chatData)
	if err != nil {
		wsc.LogWarn("handleChatMessage", err, "id", id, "reason", "json消息反序列化失败")
		wsc.errorHandler(id, message.Type, err)
		return
	}

	err = wsc.wsService.HandleChatMessage(id, &chatData)
	if err != nil {
		wsc.LogWarn("handleChatMessage", err, "id", id)
		wsc.errorHandler(id, message.Type, err)
	}
}
