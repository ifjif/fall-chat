package database

import (
	"log"
	"time"
	"zzc/fall-chat/src/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDB() {
	dsn := getDBDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	sqlDb, err := db.DB()
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)

	DB = db
}

func getDBDSN() string {
	dbhost, ok := utils.GetEnv("DB_HOST")
	if !ok {
		panic("未设置DB_HOST环境变量")
	}
	dbport, ok := utils.GetEnv("DB_PORT")
	if !ok {
		panic("未设置DB_PORT环境变量")
	}
	dbuser, ok := utils.GetEnv("DB_USER")
	if !ok {
		panic("未设置DB_USER环境变量")
	}
	dbpass, ok := utils.GetEnv("DB_PASS")
	if !ok {
		panic("未设置DB_PASS环境变量")
	}
	dbname, ok := utils.GetEnv("DB_NAME")
	if !ok {
		panic("未设置DB_NAME环境变量")
	}
	dsn := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	return dsn
}
