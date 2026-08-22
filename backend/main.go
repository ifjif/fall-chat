package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"zzc/fall-chat/src/controller"
	"zzc/fall-chat/src/database"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/server"
	"zzc/fall-chat/src/service"
	"zzc/fall-chat/src/utils"

	_ "zzc/fall-chat/src/filestore"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func main() {
	model.Init(database.DB)
	utils.Init(database.DB)
	controller.Init(server.APP, controller.WithBeforeInjectHook(func() {
		service.Init()
	}), controller.WithBeforeSetRoute(func() fiber.Router {
		return server.APP.Group("/api/v1")
	}))
	listenQuit(database.DB)

	log.Fatal(server.APP.Listen(":8080"))
}

func listenQuit(db *gorm.DB) {
	quit := make(chan os.Signal, 1)
	// ctrl+c 和 kill 命令
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		utils.LogInfo("shutting down server...")
		utils.ShutdownSnowFlake(db)
		os.Exit(1)
	}()
}
