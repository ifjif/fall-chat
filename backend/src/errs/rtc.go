package errs

import "github.com/gofiber/fiber/v3"

var (
	RTCErrPeerRoomExists        = NewBizError(RTCStartCode, fiber.StatusBadRequest, "已存在活跃房间")
	RTCErrPeerRoomNotExists     = NewBizError(RTCStartCode+1, fiber.StatusBadRequest, "房间不存在")
	RTCErrPeerRoomInviteSelf    = NewBizError(RTCStartCode+2, fiber.StatusBadRequest, "无法邀请自己")
	RTCErrPeerRoomCreator       = NewBizError(RTCStartCode+3, fiber.StatusBadRequest, "不是房间创建者")
	RTCErrPeerRoomInviteeExists = NewBizError(RTCStartCode+4, fiber.StatusBadRequest, "请勿重复邀请")
)
