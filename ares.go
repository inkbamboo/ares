package ares

import (
	"bufio"
	"fmt"
	"github.com/duke-git/lancet/v2/validator"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/inkbamboo/ares/internal/config"
	"github.com/inkbamboo/ares/internal/logger/cls"
	"github.com/inkbamboo/ares/internal/logger/sls"
	"github.com/inkbamboo/ares/internal/store"
	"github.com/patrickmn/go-cache"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
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
	mongos      map[string]*store.MongoDB
	redis       map[string]*redis.Client
	memoryCache *cache.Cache
	gin         *gin.Engine
	logger      *logrus.Logger
	logs        map[string]*logrus.Logger
}

func NewAres() *Ares {
	cfg := config.GetBaseConfig()
	a := &Ares{}
	orms := make(map[string]*store.Orm)
	mongos := make(map[string]*store.MongoDB)
	if !validator.IsZeroValue(cfg.Databases) {
		for _, item := range cfg.Databases {
			if validator.IsZeroValue(item) {
				continue
			}
			if item.Dialect != "mongodb" {
				orms[item.Alias] = store.NewOrm(item, cfg.Debug)
			} else {
				mongos[item.Alias] = store.NewMongo(item, cfg.Debug)
			}
		}
	}
	a.orms = orms
	a.mongos = mongos
	redisClients := make(map[string]*redis.Client)
	if !validator.IsZeroValue(cfg.Caches) {
		for _, item := range cfg.Caches {
			if validator.IsZeroValue(item) {
				continue
			}
			if item.Adapter == "redis" {
				redisClients[item.Alias] = store.NewRedis(item)
			}
		}
	}
	a.redis = redisClients
	// logs
	logs := make(map[string]*logrus.Logger)
	if !validator.IsZeroValue(cfg.Logs) {
		for _, item := range cfg.Logs {
			if validator.IsZeroValue(item) {
				continue
			}
			logs[item.Alias] = NewLog(item)
		}
	}
	a.logs = logs
	if !validator.IsZeroValue(cfg.MemoryCache) {
		a.memoryCache = store.NewMemoryCache(cfg.MemoryCache)
	}
	return a
}
func (a *Ares) GetGin() *gin.Engine {
	return gin.Default()
}

// Run run ripple application
func (a *Ares) Run() {
	a.RunWith(GetBaseConfig().Domain)
}

// RunWith run ripple application
func (a *Ares) RunWith(domain string) {
	// autoMigrate all orms
	if GetBaseConfig().AutoMigrate {
		for alias := range a.orms {
			a.orms[alias].AutoMigrateAll()
		}
	}
	a.logger.Info(color.GreenString("Ripple ListenAndServe: %s", domain))
	err := a.gin.Run(domain)
	if err != nil {
		a.logger.Error(color.RedString("Ripple Start error: %s", err))
	}
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
func (a *Ares) GetMongo(alias string) *store.MongoDB {
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
func GetBaseConfig() *config.BaseConfig {
	return config.GetBaseConfig()
}

func NewLog(cfg config.LogConfig) *logrus.Logger {
	if lo.IsEmpty(cfg) {
		return nil
	}
	std := logrus.New()
	formatter := &logrus.JSONFormatter{
		DisableHTMLEscape: true,
	}
	if "sls" == cfg.Type {
		h := sls.NewSLSHook(
			cfg.AccessKeyId,
			cfg.AccessKeySecret,
			cfg.Endpoint,
			cfg.AllowLogLevel,
			sls.SetProject(cfg.Project),
			sls.SetLogstore(cfg.Logstore),
			sls.SetTopic(cfg.Topic),
			sls.SetSource(cfg.Source),
		)
		if cfg.CloseStdout {
			f, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				fmt.Println("SLS.CloseStdout Open file err: ", err)
			}
			std.SetOutput(bufio.NewWriter(f))
		}
		std.SetFormatter(formatter)
		std.AddHook(h)
	} else if "cls" == cfg.Type {
		h := cls.NewCLSHook(
			cfg.AccessKeyId,
			cfg.AccessKeySecret,
			cfg.Endpoint,
			cfg.AllowLogLevel,
			cls.SetTopic(cfg.Topic),
		)
		std.AddHook(h)
		if cfg.CloseStdout {
			f, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				fmt.Println("CLS.CloseStdout Open file err: ", err)
			}
			std.SetOutput(bufio.NewWriter(f))
		}
		std.SetFormatter(formatter)
		std.AddHook(h)
	}
	return std
}
