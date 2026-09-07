package controller

import (
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/fiber/v3"
)

type BaseController struct {
	name string
}

func (bs *BaseController) Name() string {
	return bs.name
}

func (bs *BaseController) SetRoute(f fiber.Router) {}

func (bs *BaseController) Inject() {}

func (bc *BaseController) LogWarn(op string, err error, fields ...any) {
	utils.LogWarn("[%s] %s failed: %v | Context: %v", bc.name, op, err, fields)
}

func (bc *BaseController) LogInfo(op string, fields ...any) {
	utils.LogInfo("[%s] %s | Context: %v", bc.name, op, fields)
}
