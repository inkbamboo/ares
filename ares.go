package ares

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/inkbamboo/ares/internal/config"
	"github.com/inkbamboo/ares/internal/store"
	"github.com/patrickmn/go-cache"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

var ares *Ares
var once sync.Once

func Default() *Ares {
	once.Do(func() {
		ares = NewAres()
	})
	return ares
}

type Ares struct {
	orms        map[string]*store.Orm
	mongos      map[string]*mongo.Database
	redis       map[string]*redis.Client
	memoryCache *cache.Cache
}

func NewAres() *Ares {
	cfg := config.GetBaseConfig()
	a := &Ares{}
	orms := make(map[string]*store.Orm)
	mongos := make(map[string]*mongo.Database)
	if len(cfg.Databases) > 0 {
		for _, item := range cfg.Databases {
			if item.Dialect != "mongo" {
				orms[item.Alias] = store.NewOrm(item, cfg.Debug)
			} else {
				mongos[item.Alias] = store.NewMongo(item, cfg.Debug)
			}
		}
	}
	a.orms = orms
	a.mongos = mongos
	redisClients := make(map[string]*redis.Client)
	if len(cfg.Caches) > 0 {
		for _, item := range cfg.Caches {
			redisClients[item.Alias] = store.NewRedis(item)
		}
	}
	a.redis = redisClients

	if lo.IsNotEmpty(cfg.MemoryCache) {
		a.memoryCache = store.NewMemoryCache(cfg.MemoryCache)
	}

	return a
}
func (a *Ares) GetOrm(alias string) *store.Orm {
	if _, ok := a.orms[alias]; !ok {
		panic(fmt.Errorf("GetOrm: cannot get orm alias '%s'", alias))
	}
	return a.orms[alias]
}

func (a *Ares) GetRedis(alias string) *redis.Client {
	if _, ok := a.redis[alias]; !ok {
		panic(fmt.Errorf("GetRedis: cannot get redis alias '%s'", alias))
	}
	return a.redis[alias]
}
func (a *Ares) GetMongo(alias string) *mongo.Database {
	if _, ok := a.mongos[alias]; !ok {
		panic(fmt.Errorf("GetMongo: cannot get mongo alias '%s'", alias))
	}
	return a.mongos[alias]
}
func (a *Ares) GetMemoryCache() *cache.Cache {
	return a.memoryCache
}
func InitConfigWithPath(env string, configPath string) {
	config.InitConfigWithPath(env, configPath)
}

func GetEnv() string {
	return config.GetEnv()
}

func GetConfig() *viper.Viper {
	return config.GetConfig()
}
