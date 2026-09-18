package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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
	"github.com/google/uuid"
)

type MessageHandler func(client *service.Client, message *request.WSMessage)

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

func (wsc *WSController) getLogger(c *websocket.Conn) *slog.Logger {
	traceID := uuid.New().String()
	return slog.Default().With("trace_id", traceID)
}

func (wsc *WSController) accessURL(c fiber.Ctx) error {
	pID := ctxkey.GetUserPubID(c)
	username := ctxkey.GetUsername(c)

	token, _, err := utils.GenerateAccessToken(pID, username, "")
	if err != nil {
		wsc.Wrap("生成accessToken", err, "username", username)
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
		logger := wsc.getLogger(c)
		id := c.Params("id", "")
		if id == "" {
			logger.Error("websocket upgrade", fiber.ErrUpgradeRequired, "reason", "用户id不存在")
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		clientCtx := ctxkey.SetLogger(ctx, logger.With("id", id))

		client, err := wsc.wsService.Add(clientCtx, id, c)
		if err != nil {
			logger.Error("websocket", "reason", "添加websocket client失败", "id", id, "error", err)
			return
		}
		defer wsc.wsService.Remove(id)

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				err = wsc.Wrap("client read error", err, "user_public_id", id)
				client.LogError(err)
				break
			}

			var message request.WSMessage
			if err := json.Unmarshal(msg, &message); err != nil {
				err = wsc.Wrap("消息解码为json错误", err, "user_public_id:", id)
				client.LogError(err)
				continue
			}

			wsc.dispatchMessage(client, &message)
		}
	})
}

func (wsc *WSController) dispatchMessage(client *service.Client, message *request.WSMessage) {
	mType := message.Type

	if mType != "ping" {
		// wsc.LogInfo("dispatchMessage", "id:", id, "type:", mType)
	}

	messageHandler, ok := wsc.messageHandlers[mType]

	if ok {
		messageHandler(client, message)
		return
	}
	wsc.unmatchHandler(client, message)
}

func (wsc *WSController) unmatchHandler(client *service.Client, message *request.WSMessage) {
	wsc.errorHandler(client, "unmatch", &response.WSErr{Error: errs.ErrWSNotFound})
}

func (wsc *WSController) errorHandler(client *service.Client, ty string, data *response.WSErr) {
	// 错误的类型
	errType := ty + "_error"
	err := data.Error

	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		appErr.WithContext("type", errType)
		client.LogError(err, appErr.ExtractContext()...)
	}

	respErr := response.WSMessage{
		Type:      errType,
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
	}

	data.Error = errs.ErrWSInternal
	var wsErr *errs.WSError
	if errors.As(err, &wsErr) {
		data.Error = wsErr
	}

	wsc.wsService.SendMessage(client.GetContext(), client.ID, respErr)
}

func (wsc *WSController) chatMessageHandler(client *service.Client, message *request.WSMessage) {
	var chatData request.ChatData
	err := json.Unmarshal(message.Data, &chatData)
	if err != nil {
		err = wsc.Wrap("json消息反序列化", err, "reason", "")
		wsc.errorHandler(client, message.Type, &response.WSErr{Error: err})
		return
	}

	err = wsc.wsService.HandleChatMessage(client, &chatData)
	if err != nil {
		errResp := response.WSErr{
			Error: err,
			Data: response.ChatErrData{
				ConvID: chatData.Conv,
				MsgID:  chatData.ClientMsgID,
			},
		}
		wsc.errorHandler(client, message.Type, &errResp)
	}
}
