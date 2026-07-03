package store

import (
	"fmt"
	"github.com/inkbamboo/ares/internal/config"
	"github.com/labstack/gommon/color"
	"github.com/patrickmn/go-cache"
	"time"
)

// NewMemoryCache 根据内存缓存配置创建 go-cache 实例。
// DefaultExpiration 为默认过期时间，CleanupInterval 为自动清理间隔，单位均为秒。
func NewMemoryCache(config config.MemoryCacheConfig) *cache.Cache {
	c := cache.New(time.Duration(config.DefaultExpiration)*time.Second, time.Duration(config.CleanupInterval)*time.Second)
	fmt.Println(fmt.Sprintf("%s: %s, db: %d", color.Green("Connect.go-cache"), config.DefaultExpiration, config.CleanupInterval))
	return c
}
