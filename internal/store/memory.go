package store

import (
	"fmt"
	"github.com/inkbamboo/ares/internal/config"
	"github.com/labstack/gommon/color"
	"github.com/patrickmn/go-cache"
	"time"
)

func NewMemoryCache(config config.MemoryCacheConfig) *cache.Cache {
	c := cache.New(time.Duration(config.DefaultExpiration)*time.Second, time.Duration(config.CleanupInterval)*time.Second)
	fmt.Println(fmt.Sprintf("%s: %s, db: %d", color.Green("Connect.go-cache"), config.DefaultExpiration, config.CleanupInterval))
	return c
}
