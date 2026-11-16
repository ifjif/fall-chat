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

type RTC struct {
	Base
	rtcService *service.RTC
}

func NewRTCController(name string) *RTC {
	bc := Base{name: name}
	return &RTC{
		Base: bc,
	}
}

func (rc *RTC) Inject() {
	rs := service.GlobalService[service.RTC_SERVICE_NAME].(*service.RTC)
	rc.rtcService = rs
}

func (rc *RTC) SetRoute(r fiber.Router) {
	group := r.Group("/rtc")
	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Post("/room/create", rc.createRoom)
	authGroup.Post("/room/invite", rc.roomInvite)
	authGroup.Post("/room/invite/accept", rc.inviteAccept)
}

func (rc *RTC) createRoom(c fiber.Ctx) error {
	userID := ctxkey.GetUserID(c)

	req, err := bind.Body[request.RTCCreateRoom](c)
	if err != nil {
		return rc.BindBoyErr(err)
	}

	cmd := converter.ToRTCCreateRoomCommand(req, userID)

	result, err := rc.rtcService.CreateRoom(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(result)
}

func (rc *RTC) roomInvite(c fiber.Ctx) error {
	userID := ctxkey.GetUserID(c)

	req, err := bind.Body[request.RTCRoomInvite](c)
	if err != nil {
		return rc.BindBoyErr(err)
	}

	cmd, err := converter.ToRTCRoomInviteCommand(req, userID)
	if err != nil {
		return err
	}

	err = rc.rtcService.RoomInvite(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(nil)
}

func (rc *RTC) inviteAccept(c fiber.Ctx) error {
	userID := ctxkey.GetUserID(c)

	req, err := bind.Body[request.RTCRoomInviteAccept](c)
	if err != nil {
		return rc.BindBoyErr(err)
	}

	cmd, err := converter.ToRTCRoomInviteAcceptCommand(req, userID)
	if err != nil {
		return err
	}

	err = rc.rtcService.InviteAccept(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.New(c).
		OK().
		JSON(nil)
}
