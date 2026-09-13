package controller

import (
	"zzc/fall-chat/src/bind"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/pkg/resp"
	"zzc/fall-chat/src/service"

	"github.com/gofiber/fiber/v3"
)

type EmailController struct {
	BaseController
	emailService *service.EmailService
}

func NewEmailController(name string) *EmailController {
	bs := BaseController{name: name}
	return &EmailController{BaseController: bs}
}

func (ec *EmailController) Inject() {
	es := service.GlobalService[service.EMAIL_SERVICE_NAME].(*service.EmailService)
	ec.emailService = es
}

func (ec *EmailController) SetRoute(r fiber.Router) {
	group := r.Group("/email")
	group.Post("/passRestCode", ec.passResetCode)
}

// /email/passResetCode
func (ec *EmailController) passResetCode(c fiber.Ctx) error {
	req, err := bind.Body[request.PassResetCodeRequest](c)
	if err != nil {
		return ec.LogWarn("passResetCode Body绑定失败", err)
	}

	id, err := ec.emailService.PassResetCode(req.Email)
	if err != nil {
		return ec.LogWarn("passResetCode", err)
	}

	return resp.
		New(c).
		OK().
		Message("发送成功").
		JSON(response.PassResetCodeResponse{
			ID: id,
		})
}
