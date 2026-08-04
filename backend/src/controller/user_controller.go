package controller

import (
	"log"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/service"

	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	name        string
	userService *service.UserService
}

func NewUserController(name string) *UserController {
	userService := service.GlobalService[service.USER_SERVICE_NAME].(*service.UserService)
	return &UserController{
		userService: userService,
		name:        name,
	}
}

func UserGroup(c *fiber.App, uc *UserController) {
	group := c.Group("/user")
	group.Post("/register", uc.register)
}

type RegisterUserReq struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

// /user/register
func (uc *UserController) register(c fiber.Ctx) error {
	var registerUserReq RegisterUserReq
	if err := c.Bind().Body(&registerUserReq); err != nil {
		log.Printf("UserController[请求参数错误]：%v\n", err)
		// return &fiber.Error{Code: fiber.StatusBadRequest, Message: "请求参数错误"}
		return errs.NewBizError(4000, fiber.StatusBadRequest, "请求参数错误")
	}
	pid, err := uc.userService.RegisterUser(registerUserReq.Name, registerUserReq.Password)
	if err != nil {
		log.Printf("注册失败: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("注册失败")
	}

	return c.JSON(fiber.Map{"id": pid})
}
