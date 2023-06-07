package repo

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Implementation struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewImplementation(db *gorm.DB, redis *redis.Client) *Implementation {
	return &Implementation{
		db:    db,
		redis: redis,
	}
}
