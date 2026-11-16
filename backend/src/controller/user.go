package controller

import (
	"fmt"
	"time"
	"zzc/fall-chat/src/bind"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/resp"
	"zzc/fall-chat/src/service"
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/fiber/v3"
)

type User struct {
	Base
	userService *service.User
}

func NewUser(name string) *User {
	bs := Base{name: name}
	return &User{
		Base: bs,
	}
}

func (uc *User) Inject() {
	userService := service.GlobalService[service.USER_SERVICE_NAME].(*service.User)
	uc.userService = userService
}

func (uc *User) SetRoute(r fiber.Router) {
	group := r.Group("/user")
	group.Post("/register", uc.register)
	group.Post("/login", uc.login)
	group.Post("/forgetPassword", uc.forgetPassword)

	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)
	authGroup.Post("/resetPassword", uc.resetPassword)
	authGroup.Post("/updateProfile", uc.updateProfile)
	authGroup.Get("/profile", uc.profile)
}

// /user/register
func (uc *User) register(c fiber.Ctx) error {
	req, err := bind.Body[request.RegisterUser](c)
	if err != nil {
		return uc.Wrap("Body绑定", err)
	}

	pid, err := uc.userService.RegisterUser(c.Context(), req.Username, req.Password)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("注册成功").
		JSON(fiber.Map{"id": pid})
}

// /user/login
func (uc *User) login(c fiber.Ctx) error {
	req, err := bind.Body[request.LoginUser](c)
	if err != nil {
		return uc.Wrap("Body绑定", err)
	}

	result, err := uc.userService.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return err
	}

	loginUserResp := converter.ToLoginUserResponse(result)
	// 设置 refresh cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		MaxAge:   int(result.RefreshExpiresIn),
		HTTPOnly: true,
		// Secure:   true, 仅https
	})

	return resp.New(c).
		OK().
		Message("登录成功").
		JSON(loginUserResp)
}

// /user/resetPassword	需auth
func (uc *User) resetPassword(c fiber.Ctx) error {
	// 已前置处理, 不再判断
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.ResetPassword](c)
	if err != nil {
		return uc.Wrap("Body绑定", err)
	}

	req.ID = id
	req.Username = username

	cmd := converter.ToResetPasswordCommand(req)

	err = uc.userService.ResetPassword(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("更新成功").
		JSON(nil)
}

// /user/forgetPassword
func (uc *User) forgetPassword(c fiber.Ctx) error {
	req, err := bind.Body[request.ForgetPassword](c)
	if err != nil {
		return uc.Wrap("Body绑定", err)
	}

	cmd := converter.ToForgetPasswordCommand(req)
	err = uc.userService.ForgetPassword(c.Context(), cmd)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("找回成功").
		JSON(nil)
}

// /user/updateProfile 需auth
func (uc *User) updateProfile(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)

	req, err := bind.Body[request.UpdateProfile](c)
	if err != nil {
		return uc.Wrap("Body绑定", err)
	}

	req.ID = id
	req.Username = username

	cmd := converter.ToUpdateProfileCommand(req)

	if err := uc.userService.UpdateProfile(c.Context(), cmd); err != nil {
		return err
	}
	return resp.
		New(c).
		OK().
		Message("更新成功").
		JSON(nil)
}

// /user/profile 需auth
func (uc *User) profile(c fiber.Ctx) error {
	id := ctxkey.GetUserID(c)

	result, err := uc.userService.Profile(c.Context(), id)
	if err != nil {
		return err
	}

	response := converter.ToUserProfileResponse(result)
	accessTokenExpires := utils.TOKEN_ACCESS_EXPIRE_DURATION / time.Second
	// access_token 访问期限 30 mins
	c.Set("Cache-Control", fmt.Sprintf("private, max-age=%d", accessTokenExpires))

	return resp.
		New(c).
		OK().JSON(response)
}
