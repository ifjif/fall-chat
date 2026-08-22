package controller

import (
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/pkg/resp"
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	BaseController
}

func NewAuthController(name string) *AuthController {
	return &AuthController{
		BaseController: BaseController{
			name: name,
		},
	}
}

func (ac *AuthController) SetRoute(f fiber.Router) {
	group := f.Group("/auth")
	group.Get("/refresh", ac.refresh)
}

// /auth/refresh
// 得到新 access token
func (ac *AuthController) refresh(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return resp.
			New(c).
			Err(fiber.StatusUnauthorized).
			Message("refresh token不存在，请重新登录").
			Failure()
	}

	claims, err := utils.ParseJWT(refreshToken)
	if err != nil {
		ac.LogWarn("refresh", err, "reason", "refresh token解析失败")
		return resp.
			New(c).
			Err(fiber.StatusUnauthorized).
			Message("refresh凭证错误").
			Failure()
	}

	if claims.Type != utils.TOKEN_REFRESH_TYPE {
		ac.LogWarn("refresh", errs.CommonErrJWTTokenType,
			"username", claims.Username,
		)

		return errs.CommonErrJWTTokenType
	}

	// 生成新 access token
	accessToken, expiresIn, err := utils.GenerateAccessToken(claims.UserID, claims.Username, claims.DeviceID)
	if err != nil {
		ac.LogWarn("refresh", err,
			"username", claims.Username,
			"reason", "生成access token失败",
		)

		return errs.CommonErrJWTTokenGeneration
	}

	return resp.
		New(c).
		OK().
		JSON(response.RefreshResponse{
			AccessToken: accessToken,
			ExpiresIn:   expiresIn,
		})
}
