package controller

import (
	"zzc/fall-chat/src/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Init(app *fiber.App, db *gorm.DB) {
	service.Init(db)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello world")
	})

	userController := NewUserController("UserController")

	UserGroup(app, userController)
}
