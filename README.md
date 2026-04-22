# Ares Framework

Ares 是一个轻量级的 Go Web 应用框架，提供了开箱即用的数据库连接、缓存、日志等基础设施封装，帮助开发者快速构建高性能的 Web 应用。

## 📋 目录

- [项目简介](#项目简介)
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
- [示例项目](#示例项目)
- [常见问题](#常见问题)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

## 项目简介

Ares 是一个基于 Gin 框架的 Go Web 应用基础框架，主要目标是：

- **简化配置**：通过 YAML 配置文件统一管理数据库、缓存、日志等配置
- **多数据源支持**：支持 MySQL、PostgreSQL、MongoDB 等多种数据库
- **灵活缓存**：支持 Redis 分布式缓存和内存缓存
- **强大日志**：集成 Logrus，支持阿里云 SLS 和腾讯云 CLS 日志服务
- **单例模式**：全局单例访问，避免重复初始化

## 核心特性

✅ **多数据库支持** - MySQL、PostgreSQL、MongoDB  
✅ **多缓存策略** - Redis、内存缓存  
✅ **云日志集成** - 阿里云 SLS、腾讯云 CLS  
✅ **配置管理** - 基于 Viper 的多环境配置  
✅ **信号处理** - 优雅退出和资源清理  
✅ **自动迁移** - GORM 自动数据库迁移  
✅ **JWT 中间件** - 内置 JWT 认证支持  
✅ **工具函数** - 常用辅助函数集合  

## 技术栈

### 核心依赖
- **Go**: 1.25+
- **Gin**: 1.11.0 - Web 框架
- **GORM**: 1.31.0 - ORM 框架
- **Viper**: 1.21.0 - 配置管理
- **Logrus**: 1.9.3 - 日志框架

### 数据库驱动
- **MySQL**: gorm.io/driver/mysql v1.6.0
- **PostgreSQL**: gorm.io/driver/postgres v1.6.0
- **MongoDB**: go.mongodb.org/mongo-driver v1.17.4

### 缓存
- **Redis**: go-redis/redis/v8 v8.11.5
- **Memory Cache**: patrickmn/go-cache v2.1.0

### 云服务
- **阿里云 SLS**: aliyun/aliyun-log-go-sdk v0.1.109
- **腾讯云 CLS**: tencentcloud/tencentcloud-cls-sdk-go v1.0.12

### 工具库
- **Lancet**: duke-git/lancet/v2 - Go 工具函数库
- **Samber/lo**: samber/lo - Lodash 风格工具库
- **Fatih/color**: fatih/color - 终端彩色输出
- **Packr**: gobuffalo/packr/v2 - 静态资源打包

## 项目结构

```
ares/
├── helper/                     # 辅助函数
│   └── helper.go               # 通用工具函数
├── internal/                   # 内部包（不对外暴露）
│   ├── config/                 # 配置管理
│   │   └── config.go           # 配置初始化和读取
│   ├── logger/                 # 日志钩子
│   │   ├── cls/                # 腾讯云 CLS 钩子
│   │   │   └── hook.go
│   │   └── sls/                # 阿里云 SLS 钩子
│   │       └── hook.go
│   ├── mdw/                    # 中间件
│   │   └── jwt.go              # JWT 中间件
│   └── store/                  # 数据存储
│       ├── memory.go           # 内存缓存
│       ├── mongo.go            # MongoDB 连接
│       ├── orm.go              # ORM 封装（MySQL/PostgreSQL）
│       └── redis.go            # Redis 连接
├── logger/                     # 日志模块（对外暴露）
│   ├── cls/                    # CLS 日志钩子
│   │   └── hook.go
│   ├── sls/                    # SLS 日志钩子
│   │   └── hook.go
│   └── logger.go               # 标准日志接口
├── middlewares/logger/         # 日志中间件
│   ├── cls/
│   ├── sls/
│   └── logger.go
├── utils/                      # 工具函数
│   ├── datetime/               # 日期时间工具
│   │   └── datetime.go
│   └── util.go                 # 通用工具
├── test/                       # 测试文件
│   └── test.go
├── ares.go                     # 框架核心入口
├── go.mod                      # Go 模块依赖
├── go.sum                      # 依赖校验文件
├── LICENSE                     # 开源许可证
├── .gitignore                  # Git 忽略配置
└── swagger.bat                 # Swagger 生成脚本（Windows）
```

## 快速开始

### 安装

```bash
go get github.com/inkbamboo/ares
```

或在 `go.mod` 中添加：

```go
require github.com/inkbamboo/ares latest
```

### 最小示例

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/inkbamboo/ares"
    "github.com/inkbamboo/ares/helper"
)

func main() {
    // 初始化配置
    ares.InitConfigWithPath("local", "./config")
    
    // 获取默认实例
    app := ares.Default()
    
    // 获取 Gin 引擎
    engine := app.GetGin()
    
    // 定义路由
    engine.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })
    
    // 优雅退出
    helper.OnStop(func() {
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

### BaseConfig 结构

```go
type BaseConfig struct {
    AutoMigrate bool              // 是否自动迁移数据库
    Debug       bool              // 调试模式
    Domain      string            // 服务监听地址
    Databases   []DatabaseConfig  // 数据库配置列表
    Caches      []CacheConfig     // 缓存配置列表
    Logs        []LogConfig       // 日志配置列表
    MemoryCache MemoryCacheConfig // 内存缓存配置
}
```

### DatabaseConfig

```go
type DatabaseConfig struct {
    Alias        string // 别名，用于获取连接
    Dialect      string // 数据库类型：mysql, postgresql, mongodb
    Host         string // 主机地址
    Port         int    // 端口
    DbName       string // 数据库名
    Username     string // 用户名
    Password     string // 密码
    MaxIdleConns int    // 最大空闲连接数
    MaxOpenConns int    // 最大打开连接数
}
```

### CacheConfig

```go
type CacheConfig struct {
    Alias    string // 别名
    Section  string // 分区
    Adapter  string // 适配器类型：redis
    Host     string // 主机地址
    Port     int    // 端口
    Password string // 密码
    DB       int    // Redis 数据库编号
}
```

### LogConfig

```go
type LogConfig struct {
    Alias           string // 日志别名
    Type            string // 日志类型：sls, cls
    AccessKeyId     string // 访问密钥 ID
    AccessKeySecret string // 访问密钥 Secret
    Endpoint        string // 服务端点
    AllowLogLevel   string // 允许的日志级别
    CloseStdout     bool   // 是否关闭标准输出
    Project         string // 项目名称（SLS）
    Logstore        string // 日志存储（SLS）
    Topic           string // 主题
    Source          string // 来源
}
```

### MemoryCacheConfig

```go
type MemoryCacheConfig struct {
    DefaultExpiration int // 默认过期时间（秒）
    CleanupInterval   int // 清理间隔（秒）
}
```

## 使用指南

### 数据库操作

#### 获取 ORM 实例

```go
// 获取名为 "master" 的数据库连接
orm := ares.Default().GetOrm("master")

// 添加模型
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string `gorm:"column:name;type:varchar(100)"`
    Email    string `gorm:"column:email;type:varchar(100)"`
    Age      int    `gorm:"column:age"`
}

err := orm.AddModels(&User{})
if err != nil {
    panic(err)
}

// 自动迁移（如果配置了 autoMigrate: true）
orm.AutoMigrateAll()
```

#### CRUD 操作

```go
orm := ares.Default().GetOrm("master")

// 创建
user := User{Name: "John", Email: "john@example.com", Age: 25}
orm.Create(&user)

// 查询
var users []User
orm.Where("age > ?", 18).Find(&users)

// 更新
orm.Model(&User{}).Where("id = ?", 1).Update("name", "Jane")

// 删除
orm.Delete(&User{}, 1)
```

#### 多数据库支持

```yaml
databases:
  - alias: "master"
    dialect: "mysql"
    host: "127.0.0.1"
    port: 3306
    dbName: "master_db"
    username: "root"
    password: "password"
    
  - alias: "slave"
    dialect: "postgresql"
    host: "127.0.0.1"
    port: 5432
    dbName: "slave_db"
    username: "postgres"
    password: "password"
```

```go
// 获取不同的数据库连接
master := ares.Default().GetOrm("master")
slave := ares.Default().GetOrm("slave")
```

#### MongoDB 支持

```yaml
databases:
  - alias: "mongo"
    dialect: "mongodb"
    host: "127.0.0.1"
    port: 27017
    dbName: "mydb"
    username: "admin"
    password: "password"
```

```go
// 获取 MongoDB 连接
mongo := ares.Default().GetMongo("mongo")
```

### Redis 缓存

#### 获取 Redis 客户端

```go
// 获取名为 "default" 的 Redis 连接
redisClient := ares.Default().GetRedis("default")
```

#### 基本操作

```go
import "context"

ctx := context.Background()
redisClient := ares.Default().GetRedis("default")

// 设置键值
redisClient.Set(ctx, "key", "value", 0).Err()

// 获取值
val, err := redisClient.Get(ctx, "key").Result()

// 设置带过期时间的键值
redisClient.Set(ctx, "temp_key", "temp_value", 5*time.Minute).Err()

// 删除键
redisClient.Del(ctx, "key").Err()

// Hash 操作
redisClient.HSet(ctx, "user:1", "name", "John").Err()
redisClient.HGet(ctx, "user:1", "name").Result()
```

#### 多 Redis 实例

```yaml
caches:
  - alias: "session"
    adapter: "redis"
    host: "127.0.0.1"
    port: 6379
    db: 0
    
  - alias: "cache"
    adapter: "redis"
    host: "127.0.0.1"
    port: 6379
    db: 1
```

```go
sessionRedis := ares.Default().GetRedis("session")
cacheRedis := ares.Default().GetRedis("cache")
```

### 内存缓存

#### 配置

```yaml
memoryCache:
  defaultExpiration: 300  # 默认过期时间 5 分钟
  cleanupInterval: 600    # 清理间隔 10 分钟
```

#### 使用

```go
import "time"

cache := ares.Default().GetMemoryCache()

// 设置键值（使用默认过期时间）
cache.Set("key", "value", cache.DefaultExpiration)

// 设置键值（自定义过期时间）
cache.Set("temp_key", "temp_value", 5*time.Minute)

// 获取值
value, found := cache.Get("key")
if found {
    fmt.Println(value)
}

// 删除键
cache.Delete("key")

// 清空所有缓存
cache.Flush()
```

### 日志系统

#### 标准日志

```go
import log "github.com/inkbamboo/ares/logger"

// 基本日志
log.Info("Application started")
log.Warn("This is a warning")
log.Error("This is an error")

// 格式化日志
log.Infof("User %s logged in", username)
log.Errorf("Failed to connect to database: %v", err)

// 带字段的日志
log.WithFields(log.Fields{
    "user_id": 123,
    "action":  "login",
}).Info("User action")

// 带错误信息的日志
log.WithError(err).Error("Operation failed")
```

#### 自定义 Logger

```go
// 获取特定别名的 logger
logger := ares.Default().GetLogger("app")

// 使用该 logger
logger.Info("Custom logger message")
```

#### 阿里云 SLS 日志

```yaml
logs:
  - alias: "sls_logger"
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

```go
// 获取 SLS logger
slsLogger := ares.NewLog(config.LogConfig{
    Type:            "sls",
    AccessKeyId:     "your-access-key-id",
    AccessKeySecret: "your-access-key-secret",
    Endpoint:        "cn-hangzhou.log.aliyuncs.com",
    Project:         "my-project",
    Logstore:        "my-logstore",
    Topic:           "app-log",
    Source:          "my-app",
    AllowLogLevel:   "info",
})

slsLogger.Info("This log will be sent to Aliyun SLS")
```

#### 腾讯云 CLS 日志

```yaml
logs:
  - alias: "cls_logger"
    type: "cls"
    accessKeyId: "your-secret-id"
    accessKeySecret: "your-secret-key"
    endpoint: "ap-guangzhou.cls.tencentcs.com"
    topic: "app-log"
    allowLogLevel: "info"
    closeStdout: false
```

```go
// 获取 CLS logger
clsLogger := ares.NewLog(config.LogConfig{
    Type:            "cls",
    AccessKeyId:     "your-secret-id",
    AccessKeySecret: "your-secret-key",
    Endpoint:        "ap-guangzhou.cls.tencentcs.com",
    Topic:           "app-log",
    AllowLogLevel:   "info",
})

clsLogger.Info("This log will be sent to Tencent Cloud CLS")
```

### 工具函数

#### 判断空值

```go
import "github.com/inkbamboo/ares/helper"

// 判断是否为空
if helper.IsEmpty("") {
    fmt.Println("String is empty")
}

if helper.IsNotEmpty("hello") {
    fmt.Println("String is not empty")
}

// 支持各种类型
helper.IsEmpty(0)           // true
helper.IsEmpty([]int{})     // true
helper.IsEmpty(map[string]int{}) // true
```

#### 字符串截取

```go
// 支持中文等多字节字符
result := helper.Substring("Hello 世界", 0, 5)
fmt.Println(result) // "Hello"

result = helper.Substring("Hello 世界", 6, 2)
fmt.Println(result) // "世界"
```

#### 获取当前方法名

```go
func MyFunction() {
    methodName := helper.CurrentMethodName()
    fmt.Println(methodName) // "MyFunction"
}
```

#### 优雅退出

```go
import "github.com/inkbamboo/ares/helper"

// 注册退出时的清理函数
helper.OnStop(func() {
    fmt.Println("Cleaning up resources...")
    // 关闭数据库连接
    // 保存状态
    // 其他清理工作
})
```

#### 日期时间常量

```go
import "github.com/inkbamboo/ares/helper"

// 日期格式
format1 := helper.DateShortLayout  // "2006-01-02"
format2 := helper.DateFullLayout   // "2006-01-02 15:04:05"
location := helper.TimeLocationName // "Asia/Shanghai"
```

## API 参考

### Ares 核心 API

#### `Default() *Ares`

获取 Ares 单例实例。

```go
app := ares.Default()
```

#### `InitConfigWithPath(env string, configPath string)`

初始化配置。

```go
ares.InitConfigWithPath("local", "./config")
```

参数：
- `env`: 环境名称（local, test, prod 等）
- `configPath`: 配置文件路径

#### `GetConfig() *viper.Viper`

获取 Viper 配置实例。

```go
config := ares.GetConfig()
value := config.GetString("some.key")
```

#### `GetBaseConfig() *config.BaseConfig`

获取基础配置结构。

```go
baseConfig := ares.GetBaseConfig()
fmt.Println(baseConfig.Domain)
```

#### `GetEnv() string`

获取当前环境名称。

```go
env := ares.GetEnv() // "local", "test", "prod"
```

### Ares 实例方法

#### `GetGin() *gin.Engine`

获取 Gin 引擎实例。

```go
engine := app.GetGin()
engine.GET("/ping", handler)
```

#### `GetOrm(alias string) *store.Orm`

获取指定别名的 ORM 实例。

```go
orm := app.GetOrm("master")
```

#### `GetRedis(alias string) *redis.Client`

获取指定别名的 Redis 客户端。

```go
redis := app.GetRedis("default")
```

#### `GetMongo(alias string) *store.MongoDB`

获取指定别名的 MongoDB 连接。

```go
mongo := app.GetMongo("mongo")
```

#### `GetMemoryCache() *cache.Cache`

获取内存缓存实例。

```go
cache := app.GetMemoryCache()
```

#### `Run()`

启动 Web 服务，使用配置文件中的 domain。

```go
app.Run()
```

#### `RunWith(domain string)`

启动 Web 服务，使用指定的地址。

```go
app.RunWith("0.0.0.0:9090")
```

### Orm 方法

#### `AddModels(values ...interface{}) error`

注册模型。

```go
orm.AddModels(&User{}, &Product{})
```

#### `AutoMigrateAll()`

自动迁移所有注册的模型。

```go
orm.AutoMigrateAll()
```

#### `Close() error`

关闭数据库连接。

```go
orm.Close()
```

### Logger API

Ares 的 logger 完全兼容 Logrus API，支持以下方法：

- `Trace()`, `Debug()`, `Info()`, `Warn()`, `Error()`, `Fatal()`, `Panic()`
- `Tracef()`, `Debugf()`, `Infof()`, `Warnf()`, `Errorf()`, `Fatalf()`, `Panicf()`
- `Traceln()`, `Debugln()`, `Infoln()`, `Warnln()`, `Errorln()`, `Fatalln()`, `Panicln()`
- `WithField()`, `WithFields()`, `WithError()`, `WithContext()`, `WithTime()`

## 示例项目

查看 [risk-monitoring-system-server](../risk-monitoring-system-server) 项目，这是一个使用 Ares 框架构建的完整应用示例。

主要特点：
- 多环境配置管理
- JWT 认证
- RBAC 权限控制
- 多种数据库操作
- Redis 缓存
- 云日志集成

## 常见问题

### 1. 如何切换不同环境的配置？

```go
// 开发环境
ares.InitConfigWithPath("local", "./config")

// 测试环境
ares.InitConfigWithPath("test", "./config")

// 生产环境
ares.InitConfigWithPath("prod", "./config")
```

### 2. 如何处理多个数据库连接？

在配置文件中定义多个数据库：

```yaml
databases:
  - alias: "master"
    dialect: "mysql"
    # ... master 配置
    
  - alias: "slave"
    dialect: "mysql"
    # ... slave 配置
```

使用时通过别名获取：

```go
master := ares.Default().GetOrm("master")
slave := ares.Default().GetOrm("slave")
```

### 3. 如何实现读写分离？

结合多数据源和中间件实现：

```go
// 写操作使用 master
master := ares.Default().GetOrm("master")
master.Create(&user)

// 读操作使用 slave
slave := ares.Default().GetOrm("slave")
slave.Find(&users)
```

### 4. 日志没有输出到云日志服务？

检查以下几点：
1. 确认 AccessKeyId 和 AccessKeySecret 正确
2. 确认 Endpoint、Project、Logstore 配置正确
3. 确认网络可以访问云服务
4. 检查 `allowLogLevel` 设置，确保日志级别匹配

### 5. 如何自定义日志格式？

```go
import log "github.com/inkbamboo/ares/logger"
import "github.com/sirupsen/logrus"

// 设置为文本格式
log.SetFormatter(&logrus.TextFormatter{
    FullTimestamp: true,
})

// 或保持 JSON 格式（默认）
log.SetFormatter(&logrus.JSONFormatter{})
```

### 6. 内存缓存的过期时间如何设置？

```go
cache := ares.Default().GetMemoryCache()

// 使用默认过期时间
cache.Set("key", "value", cache.DefaultExpiration)

// 自定义过期时间
cache.Set("key", "value", 10*time.Minute)

// 永不过期
cache.Set("key", "value", -1)
```

### 7. 如何优雅地关闭服务？

使用 `helper.OnStop` 注册清理函数：

```go
helper.OnStop(func() {
    // 关闭数据库连接
    orm := ares.Default().GetOrm("master")
    orm.Close()
    
    // 其他清理工作
    fmt.Println("Resources cleaned up")
})
```

当收到 SIGINT、SIGTERM 等信号时，会自动执行清理函数。

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 开发规范

- 遵循 [Effective Go](https://go.dev/doc/effective_go) 规范
- 所有导出的函数和类型需要注释
- 添加单元测试
- 保持代码简洁清晰

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 联系方式

- 作者: liuxinlei
- 邮箱: [your-email@example.com]

## 致谢

- [Gin Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [Viper](https://github.com/spf13/viper)
- [Logrus](https://github.com/sirupsen/logrus)
- [Redis](https://redis.io/)
- [MongoDB](https://www.mongodb.com/)
