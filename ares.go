package ares

import (
	"fmt"
	"github.com/duke-git/lancet/v2/validator"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/inkbamboo/ares/internal/config"
	"github.com/inkbamboo/ares/internal/logger/cls"
	"github.com/inkbamboo/ares/internal/logger/sls"
	"github.com/inkbamboo/ares/internal/store"
	log "github.com/inkbamboo/ares/logger"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sync"
)

var ares *Ares
var once sync.Once

// Default 返回 Ares 的全局单例实例。
// 通过 sync.Once 保证只初始化一次，是框架的主要入口。
func Default() *Ares {
	once.Do(func() {
		ares = NewAres()
	})
	return ares
}

// Ares 是框架的核心结构体，管理所有基础设施资源，
// 包括 ORM 连接、MongoDB 连接、Redis 客户端、内存缓存、Gin 引擎和日志实例。
type Ares struct {
	orms        map[string]*store.Orm
	mongos      map[string]*store.MongoDB
	redis       map[string]*redis.Client
	memoryCache *cache.Cache
	gin         *gin.Engine
	logger      *logrus.Logger
	logs        map[string]*logrus.Logger
}

// NewAres 创建一个新的 Ares 实例。
// 根据配置文件初始化 Gin 引擎、数据库连接、Redis 客户端、日志实例和内存缓存。
// 在 debug 模式下 Gin 使用 DebugMode，否则使用 ReleaseMode。
func NewAres() *Ares {
	cfg := config.GetBaseConfig()
	a := &Ares{}
	// 根据配置设置 Gin 模式，release 模式下不会打印路由日志
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	a.gin = gin.New()
	a.gin.Use(gin.Recovery())
	a.logger = log.StandardLogger()
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

// GetGin 返回 Gin 引擎实例，用于注册路由和中间件。
func (a *Ares) GetGin() *gin.Engine {
	return a.gin
}

// Run 启动 Web 服务，使用配置文件中的 Domain 地址监听。
// 如果配置了 AutoMigrate，会自动执行所有已注册模型的数据库迁移。
func (a *Ares) Run() {
	a.RunWith(GetBaseConfig().Domain)
}

// RunWith 启动 Web 服务，使用指定的 domain 地址监听。
// 如果配置了 AutoMigrate，会自动执行所有已注册模型的数据库迁移。
func (a *Ares) RunWith(domain string) {
	if GetBaseConfig().AutoMigrate {
		for alias := range a.orms {
			a.orms[alias].AutoMigrateAll()
		}
	}
	a.logger.Info(color.GreenString("Ares ListenAndServe: %s", domain))
	err := a.gin.Run(domain)
	if err != nil {
		a.logger.Error(color.RedString("Ares Start error: %s", err))
	}
}

// GetOrm 根据别名获取 ORM 实例。如果别名不存在则 panic。
func (a *Ares) GetOrm(alias string) *store.Orm {
	if _, ok := a.orms[alias]; !ok {
		panic(fmt.Errorf("GetOrm: cannot get orm alias '%s'", alias))
	}
	return a.orms[alias]
}

// GetRedis 根据别名获取 Redis 客户端实例。如果别名不存在则 panic。
func (a *Ares) GetRedis(alias string) *redis.Client {
	if _, ok := a.redis[alias]; !ok {
		panic(fmt.Errorf("GetRedis: cannot get redis alias '%s'", alias))
	}
	return a.redis[alias]
}

// GetMongo 根据别名获取 MongoDB 连接实例。如果别名不存在则 panic。
func (a *Ares) GetMongo(alias string) *store.MongoDB {
	if _, ok := a.mongos[alias]; !ok {
		panic(fmt.Errorf("GetMongo: cannot get mongo alias '%s'", alias))
	}
	return a.mongos[alias]
}

// GetMemoryCache 返回内存缓存实例（基于 go-cache）。
func (a *Ares) GetMemoryCache() *cache.Cache {
	return a.memoryCache
}

// InitConfigWithPath 使用指定的环境名称和配置路径初始化配置。
// env 为环境标识（如 local、test、prod），configPath 为配置文件所在目录。
func InitConfigWithPath(env string, configPath string) {
	config.InitConfigWithPath(env, configPath)
}

// GetEnv 返回当前运行环境名称（如 local、test、prod）。
func GetEnv() string {
	return config.GetEnv()
}

// GetConfig 返回 Viper 配置实例，可用于读取自定义配置项。
func GetConfig() *viper.Viper {
	return config.GetConfig()
}

// GetBaseConfig 返回框架基础配置结构体指针。
func GetBaseConfig() *config.BaseConfig {
	return config.GetBaseConfig()
}

// NewLog 根据日志配置创建一个 logrus.Logger 实例。
// 支持 "sls"（阿里云日志服务）和 "cls"（腾讯云日志服务）两种类型。
// 可通过 CloseStdout 配置关闭标准输出，仅将日志发送到云端。
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
			} else {
				std.SetOutput(f)
			}
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
		if cfg.CloseStdout {
			f, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				fmt.Println("CLS.CloseStdout Open file err: ", err)
			} else {
				std.SetOutput(f)
			}
		}
		std.SetFormatter(formatter)
		std.AddHook(h)
	}
	return std
}
