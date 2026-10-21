package middleware

import (
	"errors"
	"strings"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/resp"
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

const MIDDLEWARE_NAME = "AuthMiddleware"

func AuthMiddleware(c fiber.Ctx) error {
	// 优先从 header
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		// 再从 query
		authHeader = c.Query("token")
		if authHeader == "" {
			return resp.
				New(c).
				Err(fiber.StatusUnauthorized).
				Message("缺少认证凭证").
				Failure()
		}
	}

	// 获取Bearer 认证方案
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return resp.
			New(c).
			Err(fiber.StatusUnauthorized).
			Message("凭证格式错误,要求Bearer").
			Failure()
	}

	tokenString := parts[1]

	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		// 判断是否是过期
		if errors.Is(err, jwt.ErrTokenExpired) {
			utils.LogWarn("%s: token已过期", MIDDLEWARE_NAME)
			return errs.CommonErrJWTTokenExpired
		}

		utils.LogWarn("%s: token解析失败：%v", MIDDLEWARE_NAME, err)
		return resp.
			New(c).
			Err(fiber.StatusUnauthorized).
			Message("凭证无效").
			Failure()
	}

	pID := claims.UserID
	id, err := utils.SF.Decode(pID)
	if err != nil {
		utils.LogWarn("%s: user id解码失败：%v", MIDDLEWARE_NAME, err)
		return resp.
			New(c).
			Err(fiber.StatusUnauthorized).
			Message("凭证信息损坏").
			Failure()
	}

	ctxkey.SetUserID(c, id)
	ctxkey.SetUserPubID(c, pID)
	ctxkey.SetUsername(c, claims.Username)
	deviceID := claims.DeviceID
	utils.LogInfo("%s: %s:%s 放行", MIDDLEWARE_NAME, deviceID, claims.Username)

	return c.Next()
}
