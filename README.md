# Ares Framework

Ares 是一个基于 Gin 的轻量级 Go Web 应用框架，提供开箱即用的数据库连接、缓存、日志等基础设施封装，帮助开发者快速构建高性能的 Web 应用。

## 目录

- [核心特性](#核心特性)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [使用指南](#使用指南)
  - [数据库操作](#数据库操作)
  - [Redis 缓存](#redis-缓存)
  - [内存缓存](#内存缓存)
  - [日志系统](#日志系统)
  - [工具函数](#工具函数)
- [API 参考](#api-参考)
- [常见问题](#常见问题)
- [许可证](#许可证)

## 核心特性

- **多数据库支持** - MySQL、PostgreSQL、MongoDB，通过别名管理多数据源
- **多缓存策略** - Redis 分布式缓存、go-cache 内存缓存
- **云日志集成** - 阿里云 SLS、腾讯云 CLS，支持 logrus Hook
- **配置管理** - 基于 Viper 的多环境 YAML 配置，支持 packr 打包
- **单例模式** - 全局单例访问，避免重复初始化
- **自动迁移** - GORM 自动数据库迁移
- **优雅退出** - 信号处理（SIGINT、SIGTERM 等），支持自定义清理函数
- **时间工具** - 封装 carbon/now，支持 JSON/SQL 序列化和时间差计算

## 技术栈

| 类别 | 依赖 | 版本 |
|------|------|------|
| Web 框架 | Gin | 1.11.0 |
| ORM | GORM | 1.31.0 |
| 配置管理 | Viper | 1.21.0 |
| 日志 | Logrus | 1.9.3 |
| MySQL 驱动 | gorm.io/driver/mysql | v1.6.0 |
| PostgreSQL 驱动 | gorm.io/driver/postgres | v1.6.0 |
| MongoDB 驱动 | mongo-driver | v1.17.4 |
| Redis | go-redis/v8 | v8.11.5 |
| 内存缓存 | go-cache | v2.1.0 |
| 阿里云 SLS | aliyun-log-go-sdk | v0.1.109 |
| 腾讯云 CLS | tencentcloud-cls-sdk-go | v1.0.12 |
| 工具库 | Lancet / Samber-lo / Carbon / now | - |

## 项目结构

```
ares/
├── ares.go                          # 框架核心入口，单例模式
├── internal/                        # 内部包（不对外暴露）
│   ├── config/config.go             # 配置管理，多环境 YAML 加载
│   ├── store/
│   │   ├── orm.go                   # ORM 封装（MySQL/PostgreSQL）
│   │   ├── redis.go                 # Redis 连接
│   │   ├── mongo.go                 # MongoDB 连接
│   │   └── memory.go                # 内存缓存（go-cache）
│   ├── logger/
│   │   ├── cls/hook.go              # 腾讯云 CLS Hook（内部版本）
│   │   └── sls/hook.go              # 阿里云 SLS Hook（内部版本，使用 CredentialsProvider API）
│   └── mdw/jwt.go                   # JWT 中间件（占位实现）
├── logger/                          # 日志模块（对外暴露）
│   ├── logger.go                    # 基于 logrus 的标准日志接口
│   ├── cls/hook.go                  # 腾讯云 CLS Hook
│   └── sls/hook.go                  # 阿里云 SLS Hook
├── middlewares/logger/              # 日志中间件
│   ├── logger.go                    # 自定义彩色日志器（独立实现）
│   ├── cls/hook.go                  # 腾讯云 CLS Hook
│   └── sls/hook.go                  # 阿里云 SLS Hook
├── utils/
│   ├── util.go                      # 工具函数（ExternalIP、OnStop、CurrentMethodName）
│   └── datetime/datetime.go         # 时间处理（LocalTime，支持 JSON/SQL 序列化）
├── test/test.go                     # 测试文件
├── go.mod                           # Go 模块依赖
├── go.sum                           # 依赖校验文件
├── LICENSE                          # 开源许可证
└── swagger.bat                      # Swagger 生成脚本
```

## 快速开始

### 安装

```bash
go get github.com/inkbamboo/ares
```

### 最小示例

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/inkbamboo/ares"
    "github.com/inkbamboo/ares/utils"
)

func main() {
    // 初始化配置（环境名称 + 配置文件目录）
    ares.InitConfigWithPath("local", "./config")

    // 获取全局单例
    app := ares.Default()

    // 注册路由
    app.GetGin().GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    // 优雅退出
    utils.OnStop(func() {
        println("Server shutting down...")
    })

    // 启动服务
    app.Run()
}
```

### 配置文件

创建 `config/config.local.yaml`：

```yaml
domain: "0.0.0.0:8080"
autoMigrate: false
debug: true

databases:
  - alias: "master"
    dialect: "mysql"
    host: "127.0.0.1"
    port: 3306
    dbName: "myapp"
    username: "root"
    password: "password"
    maxIdleConns: 25
    maxOpenConns: 25

caches:
  - alias: "default"
    adapter: "redis"
    host: "127.0.0.1"
    port: 6379
    password: ""
    db: 0

memoryCache:
  defaultExpiration: 300
  cleanupInterval: 600

logs:
  - alias: "app"
    type: "sls"
    accessKeyId: "your-access-key-id"
    accessKeySecret: "your-access-key-secret"
    endpoint: "cn-hangzhou.log.aliyuncs.com"
    project: "my-project"
    logstore: "my-logstore"
    topic: "app-log"
    source: "my-app"
    allowLogLevel: "info"
    closeStdout: false
```

## 配置说明

### BaseConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| AutoMigrate | bool | 是否自动执行数据库迁移 |
| Debug | bool | 调试模式（控制 Gin 模式和 SQL 日志） |
| Domain | string | 服务监听地址，如 `0.0.0.0:8080` |
| Databases | []DatabaseConfig | 数据库配置列表 |
| Caches | []CacheConfig | 缓存配置列表 |
| Logs | []LogConfig | 日志配置列表 |
| MemoryCache | MemoryCacheConfig | 内存缓存配置 |

### DatabaseConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| Alias | string | 别名，用于 `GetOrm()` / `GetMongo()` |
| Dialect | string | 数据库类型：`mysql`、`postgresql`、`mongodb` |
| Host | string | 主机地址 |
| Port | int | 端口号 |
| DbName | string | 数据库名称 |
| Username | string | 用户名 |
| Password | string | 密码 |
| MaxIdleConns | int | 最大空闲连接数 |
| MaxOpenConns | int | 最大打开连接数 |

### CacheConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| Alias | string | 别名，用于 `GetRedis()` |
| Adapter | string | 适配器类型，目前支持 `redis` |
| Host | string | 主机地址 |
| Port | int | 端口号 |
| Password | string | 密码 |
| DB | int | Redis 数据库编号 |

### LogConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| Alias | string | 日志别名 |
| Type | string | 日志服务类型：`sls`（阿里云）、`cls`（腾讯云） |
| AccessKeyId | string | 访问密钥 ID |
| AccessKeySecret | string | 访问密钥 Secret |
| Endpoint | string | 服务端点地址 |
| AllowLogLevel | string | 允许的最低日志级别 |
| CloseStdout | bool | 是否关闭标准输出 |
| Project | string | SLS 项目名称 |
| Logstore | string | SLS 日志库名称 |
| Topic | string | 日志主题 |
| Source | string | 日志来源标识 |

### MemoryCacheConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| DefaultExpiration | int | 默认过期时间（秒） |
| CleanupInterval | int | 自动清理间隔（秒） |

## 使用指南

### 数据库操作

#### 获取 ORM 实例

```go
orm := ares.Default().GetOrm("master")

// 注册模型
type User struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"column:name;type:varchar(100)"`
}
orm.AddModels(&User{})

// 自动迁移
orm.AutoMigrateAll()
```

#### CRUD 操作

```go
orm := ares.Default().GetOrm("master")

// 创建
orm.Create(&User{Name: "John"})

// 查询
var users []User
orm.Where("name = ?", "John").Find(&users)

// 更新
orm.Model(&User{}).Where("id = ?", 1).Update("name", "Jane")

// 删除
orm.Delete(&User{}, 1)
```

#### 多数据源

```yaml
databases:
  - alias: "master"
    dialect: "mysql"
    host: "127.0.0.1"
    port: 3306
    dbName: "master_db"
    username: "root"
    password: "password"
  - alias: "mongo"
    dialect: "mongodb"
    host: "127.0.0.1"
    port: 27017
    dbName: "mydb"
```

```go
master := ares.Default().GetOrm("master")
mongo  := ares.Default().GetMongo("mongo")
```

### Redis 缓存

```go
ctx := context.Background()
rdb := ares.Default().GetRedis("default")

rdb.Set(ctx, "key", "value", 5*time.Minute)
val, err := rdb.Get(ctx, "key").Result()
rdb.Del(ctx, "key")
```

### 内存缓存

```go
cache := ares.Default().GetMemoryCache()

cache.Set("key", "value", 5*time.Minute)
val, found := cache.Get("key")
cache.Delete("key")
```

### 日志系统

#### 标准日志（基于 logrus）

```go
import log "github.com/inkbamboo/ares/logger"

log.Info("Application started")
log.Errorf("Error: %v", err)
log.WithFields(logrus.Fields{"user_id": 123}).Info("login")
log.WithError(err).Error("Operation failed")
```

#### 云日志（阿里云 SLS / 腾讯云 CLS）

通过配置文件中的 `logs` 字段定义云日志 logger，框架在初始化时自动创建对应的 Hook 实例。日志会通过 logrus Hook 异步发送到云端。

### 工具函数

```go
import "github.com/inkbamboo/ares/utils"

// 获取本机外网 IP
ip, _ := utils.ExternalIP()

// 获取当前方法名
name := utils.CurrentMethodName()

// 优雅退出
utils.OnStop(func() {
    // 清理资源
})
```

#### 时间工具

```go
import "github.com/inkbamboo/ares/utils/datetime"

t := datetime.Now()
t2 := datetime.NewLocalTime(time.Now())

// 时间差计算
days := t.DiffInDays(t2)
hours := t.DiffAbsInHours()

// 支持 JSON/SQL 序列化（可直接用于 GORM 模型字段）
```

## API 参考

### 全局函数

| 函数 | 说明 |
|------|------|
| `Default() *Ares` | 获取全局单例实例 |
| `InitConfigWithPath(env, configPath string)` | 初始化配置 |
| `GetEnv() string` | 获取当前环境名称 |
| `GetConfig() *viper.Viper` | 获取 Viper 配置实例 |
| `GetBaseConfig() *config.BaseConfig` | 获取基础配置结构体 |
| `NewLog(cfg config.LogConfig) *logrus.Logger` | 创建云日志 Logger |

### Ares 实例方法

| 方法 | 说明 |
|------|------|
| `GetGin() *gin.Engine` | 获取 Gin 引擎 |
| `GetOrm(alias string) *store.Orm` | 获取 ORM 实例 |
| `GetRedis(alias string) *redis.Client` | 获取 Redis 客户端 |
| `GetMongo(alias string) *store.MongoDB` | 获取 MongoDB 连接 |
| `GetMemoryCache() *cache.Cache` | 获取内存缓存 |
| `Run()` | 启动服务（使用配置中的 Domain） |
| `RunWith(domain string)` | 启动服务（指定地址） |

### Orm 方法

| 方法 | 说明 |
|------|------|
| `AddModels(values ...interface{}) error` | 注册模型 |
| `AutoMigrateAll()` | 自动迁移所有已注册模型 |
| `Close() error` | 关闭数据库连接 |

### Logger API

兼容 Logrus 全部 API：
- `Trace/Debug/Info/Warn/Error/Fatal/Panic`
- `Tracef/Debugf/Infof/Warnf/Errorf/Fatalf/Panicf`
- `WithField/WithFields/WithError/WithContext/WithTime`

## 常见问题

### 如何切换环境？

```go
ares.InitConfigWithPath("local", "./config")  // 开发
ares.InitConfigWithPath("test", "./config")   // 测试
ares.InitConfigWithPath("prod", "./config")   // 生产
```

配置文件命名格式：`config.{env}.yaml`

### 如何处理多数据库？

在配置中定义多个数据库，通过别名获取：

```go
master := ares.Default().GetOrm("master")
slave  := ares.Default().GetOrm("slave")
```

### 云日志没有输出？

1. 确认 AccessKeyId / AccessKeySecret 正确
2. 确认 Endpoint、Project、Logstore 配置正确
3. 确认网络可访问云服务
4. 检查 `allowLogLevel` 设置

### 如何自定义日志格式？

```go
import log "github.com/inkbamboo/ares/logger"
import "github.com/sirupsen/logrus"

log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
```

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件
