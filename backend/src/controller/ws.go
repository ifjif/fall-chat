package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"
	"zzc/fall-chat/src/conns"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/resp"
	"zzc/fall-chat/src/service"
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type MessageHandler func(client *conns.Client, message *request.WSMessage)

type WSController struct {
	BaseController
	wsService       *service.Ws
	messageHandlers map[string]MessageHandler
	room            map[string]bool
}

func NewWSController(name string) *WSController {
	bs := BaseController{name: name}
	return &WSController{
		BaseController:  bs,
		messageHandlers: make(map[string]MessageHandler),
		room:            make(map[string]bool),
	}
}

func (wsc *WSController) Inject() {
	wsService := service.GlobalService[service.WS_SERVICE_NAME].(*service.Ws)
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

	wsc.messageHandlers["ping"] = wsc.pingHandler
	wsc.messageHandlers["chat"] = wsc.chatMessageHandler
	wsc.messageHandlers["offer"] = wsc.offerMessageHandler
	wsc.messageHandlers["answer"] = wsc.answerMessageHandler
	wsc.messageHandlers["candidate"] = wsc.candidateMessageHandler
	wsc.messageHandlers["join"] = wsc.joinMessageHandler
	wsc.messageHandlers["invite"] = wsc.inviteMessageHandler
	wsc.messageHandlers["leave"] = wsc.leaveMessageHandler
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
	protocol := "wss"
	// host := "localhost"
	host := "192.168.1.106"
	port := "443"
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

		client := wsc.wsService.Add(clientCtx, id, c)
		defer wsc.wsService.Remove(clientCtx, id)

		// TODO: 将目标 所在服务器ip:port存入redis, 供后续grpc调用
		wsc.wsService.RegisterGrpcAddr(clientCtx, id)

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

func (wsc *WSController) dispatchMessage(client *conns.Client, message *request.WSMessage) {
	mType := message.Type

	messageHandler, ok := wsc.messageHandlers[mType]

	if ok {
		messageHandler(client, message)
		return
	}
	wsc.unmatchHandler(client, message)
}

func (wsc *WSController) unmatchHandler(client *conns.Client, message *request.WSMessage) {
	wsc.errorHandler(client, "unmatch", &response.WSErr{Error: errs.ErrWSNotFound})
}

func (wsc *WSController) errorHandler(client *conns.Client, ty string, data *response.WSErr) {
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

func (wsc *WSController) pingHandler(client *conns.Client, message *request.WSMessage) {
	wsc.wsService.HandlePingMessage(client)
}

func (wsc *WSController) chatMessageHandler(client *conns.Client, message *request.WSMessage) {
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

func (wsc *WSController) offerMessageHandler(client *conns.Client, message *request.WSMessage) {
	var offerData request.OfferData
	json.Unmarshal(message.Data, &offerData)
	fmt.Printf("收到offer: %v\n", offerData)

	wsc.wsService.SendMessage(context.Background(), offerData.PeerID, response.WSMessage{
		Type: message.Type,
		Data: response.WebRTCData{
			PeerID: client.ID,
			Data:   offerData.Offer,
		},
	})
}

func (wsc *WSController) answerMessageHandler(client *conns.Client, message *request.WSMessage) {
	var answerData request.AnswerData
	json.Unmarshal(message.Data, &answerData)
	fmt.Printf("收到answer: %v\n", answerData)

	wsc.wsService.SendMessage(context.Background(), answerData.PeerID, response.WSMessage{
		Type: message.Type,
		Data: response.WebRTCData{
			PeerID: client.ID,
			Data:   answerData.Answer,
		},
	})
}

func (wsc *WSController) candidateMessageHandler(client *conns.Client, message *request.WSMessage) {
	var candidateDate request.CandidateData
	json.Unmarshal(message.Data, &candidateDate)
	fmt.Printf("收到candidate: %v\n", candidateDate)

	wsc.wsService.SendMessage(context.Background(), candidateDate.PeerID, response.WSMessage{
		Type: message.Type,
		Data: response.WebRTCData{
			PeerID: client.ID,
			Data:   candidateDate.Candidate,
		},
	})
}

func (wsc *WSController) joinMessageHandler(client *conns.Client, message *request.WSMessage) {
	var joinData request.JoinData
	json.Unmarshal(message.Data, &joinData)

	// 获取房间
	room := make([]string, 0)
	for k := range wsc.room {
		room = append(room, k)
	}

	wsc.wsService.SendMessage(context.Background(), joinData.PeerID, response.WSMessage{
		Type: "joined",
		Data: fiber.Map{
			"members": room,
		},
	})

	for _, id := range room {
		wsc.wsService.SendMessage(context.Background(), id, response.WSMessage{
			Type: "peer-joined",
			Data: fiber.Map{
				"peerId": joinData.PeerID,
			},
		})
	}

	wsc.room[joinData.PeerID] = true
}

func (wsc *WSController) inviteMessageHandler(client *conns.Client, message *request.WSMessage) {
	var inviteData request.InviteData
	json.Unmarshal(message.Data, &inviteData)

	wsc.wsService.SendMessage(context.Background(), inviteData.PeerID, response.WSMessage{
		Type: message.Type,
		Data: fiber.Map{
			"room_id": inviteData.RoomID,
			"id":      inviteData.PeerID,
		},
	})
}

func (wsc *WSController) leaveMessageHandler(client *conns.Client, message *request.WSMessage) {
	var leaveData request.LeaveData
	json.Unmarshal(message.Data, &leaveData)

	if wsc.room != nil {
		delete(wsc.room, leaveData.PeerID)
	}

	for id := range wsc.room {
		wsc.wsService.SendMessage(context.Background(), id, response.WSMessage{
			Type: "peer-left",
			Data: fiber.Map{
				"peerId": leaveData.PeerID,
			},
		})
	}
}
