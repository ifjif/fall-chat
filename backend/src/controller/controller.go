package controller

import "github.com/gofiber/fiber/v3"

type Controller interface {
	Name() string
	Inject()
	SetRoute(r fiber.Router)
}
