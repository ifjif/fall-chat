package errs

import "github.com/gofiber/fiber/v3"

const (
	CommonStartCode    = 2000
	UserStartCode      = 4000
	EmailStartCode     = 4200
	FileStartCode      = 4400
	SocialStartCode    = 4800
	EmojiStartCode     = 5000
	ImMessageStartCode = 5200
	RTCStartCode       = 5400
)

var (
	CommonErrJWTTokenGeneration = NewBizError(CommonStartCode, fiber.StatusInternalServerError, "JWT token生成失败")
	CommonErrJWTTokenExpired    = NewBizError(CommonStartCode+1, fiber.StatusUnauthorized, "凭证已过期")
	CommonErrJWTTokenType       = NewBizError(CommonStartCode+2, fiber.StatusUnauthorized, "token类型不正确")
	CommonErrSortKeyNotExist    = NewBizError(CommonStartCode+3, fiber.StatusBadRequest, "未知排序行为")
)
