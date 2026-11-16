package controller

import (
	"zzc/fall-chat/src/bind"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/resp"
	"zzc/fall-chat/src/service"

	"github.com/gofiber/fiber/v3"
)

type Email struct {
	Base
	emailService *service.Email
}

func NewEmail(name string) *Email {
	bs := Base{name: name}
	return &Email{Base: bs}
}

func (ec *Email) Inject() {
	es := service.GlobalService[service.EMAIL_SERVICE_NAME].(*service.Email)
	ec.emailService = es
}

func (ec *Email) SetRoute(r fiber.Router) {
	group := r.Group("/email")
	group.Post("/passRestCode", ec.passResetCode)
}

// /email/passResetCode
func (ec *Email) passResetCode(c fiber.Ctx) error {
	req, err := bind.Body[request.PassResetCode](c)
	if err != nil {
		return ec.Wrap("Body绑定", err)
	}

	id, err := ec.emailService.PassResetCode(c.Context(), req.Email)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		Message("发送成功").
		JSON(response.PassResetCode{
			ID: id,
		})
}
