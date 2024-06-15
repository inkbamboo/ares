package store

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/inkbamboo/ares/internal/config"
	"github.com/labstack/gommon/color"
)

func NewRedis(config config.CacheConfig) *redis.Client {
	nc := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
	fmt.Println(fmt.Sprintf("%s: %s, db: %d", color.Green("Connect.redis"), config.Host, config.DB))
	return nc
}
