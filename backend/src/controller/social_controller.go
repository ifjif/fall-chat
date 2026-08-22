package controller

import "github.com/gofiber/fiber/v3"

type SocialController struct {
	BaseController
}

func NewSocialController(name string) *SocialController {
	bs := BaseController{name: name}
	return &SocialController{
		BaseController: bs,
	}
}

func (sc *SocialController) SetRoute(r fiber.Router) {
}

func (sc *SocialController) Inject() {
}
