package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zzc/fall-chat/src/controller"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getEnvWithDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

func getDBDSN() string {
	dbhost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		panic("未设置DB_HOST环境变量")
	}
	dbport, ok := os.LookupEnv("DB_PORT")
	if !ok {
		panic("未设置DB_PORT环境变量")
	}
	dbuser, ok := os.LookupEnv("DB_USER")
	if !ok {
		panic("未设置DB_USER环境变量")
	}
	dbpass, ok := os.LookupEnv("DB_PASS")
	if !ok {
		panic("未设置DB_PASS环境变量")
	}
	dbname, ok := os.LookupEnv("DB_NAME")
	if !ok {
		panic("未设置DB_NAME环境变量")
	}
	dsn := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	return dsn
}

func initDB() *gorm.DB {
	dsn := getDBDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	sqlDb, err := db.DB()
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)

	return db
}

func listenQuit(db *gorm.DB) {
	quit := make(chan os.Signal, 1)
	// ctrl+c 和 kill 命令
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		fmt.Println("shutting down server...")
		utils.ShutdownSnowFlake(db)
		os.Exit(1)
	}()
}

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func main() {
	db := initDB()
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			// 业务错误
			if err, ok := err.(*errs.BizError); ok {
				return c.Status(err.Status).JSON(fiber.Map{
					"code":    err.Code,
					"status":  err.Status,
					"message": err.Message,
				})
			}

			// fiber 错误
			if err, ok := err.(*fiber.Error); ok {
				return c.Status(err.Code).JSON(fiber.Map{
					"code":    err.Code,
					"status":  err.Code,
					"message": err.Message,
				})
			}

			// 未知错误 500
			return c.Status(500).JSON(fiber.Map{
				"code":    9999,
				"status":  500,
				"message": "服务器内部错误",
			})
		},
	})

	model.Init(db)
	utils.Init(db)
	controller.Init(app, db)
	listenQuit(db)

	log.Fatal(app.Listen(":8080"))
}
