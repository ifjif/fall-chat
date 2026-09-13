package data

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BaseDao struct {
	db  *gorm.DB
	rdb *redis.Client
}

func (d *BaseDao) DB() *gorm.DB {
	return d.db
}
