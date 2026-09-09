package store

import (
	"fmt"

	"github.com/inkbamboo/ares/internal/config"
	"github.com/labstack/gommon/color"
	"github.com/redis/go-redis/v9"
)

// NewRedis 根据缓存配置创建 Redis 客户端连接。
func NewRedis(config config.CacheConfig) *redis.Client {
	nc := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
	fmt.Println(fmt.Sprintf("%s: %s, db: %d", color.Green("Connect.redis"), config.Host, config.DB))
	return nc
}
