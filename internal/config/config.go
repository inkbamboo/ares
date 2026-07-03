// Package config 提供配置管理功能，支持基于 YAML 的多环境配置加载。
package config

import (
	"bytes"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
)

// BaseConfig 是框架的基础配置结构，包含所有核心组件的配置信息。
type BaseConfig struct {
	AutoMigrate bool              `mapstructure:"autoMigrate,omitempty"` // 是否自动执行数据库迁移，默认 false
	Debug       bool              `mapstructure:"debug,omitempty"`       // 是否开启调试模式，默认 false
	Domain      string            `mapstructure:"domain,omitempty"`      // 服务监听地址，如 127.0.0.1:8090
	Databases   []DatabaseConfig  `mapstructure:"databases,omitempty"`   // 数据库配置列表
	Caches      []CacheConfig     `mapstructure:"caches,omitempty"`      // 缓存配置列表
	Logs        []LogConfig       `mapstructure:"logs,omitempty"`        // 日志配置列表
	MemoryCache MemoryCacheConfig `mapstructure:"memoryCache,omitempty"` // 内存缓存配置
}

// DatabaseConfig 定义数据库连接配置，支持 MySQL、PostgreSQL 和 MongoDB。
type DatabaseConfig struct {
	Alias        string `mapstructure:"alias"`        // 数据库别名，用于获取连接
	Dialect      string `mapstructure:"dialect"`      // 数据库类型：mysql、postgresql、mongodb
	Host         string `mapstructure:"host"`         // 主机地址
	Port         int    `mapstructure:"port"`         // 端口号
	DbName       string `mapstructure:"dbName"`       // 数据库名称
	Username     string `mapstructure:"username"`     // 用户名
	Password     string `mapstructure:"password"`     // 密码
	MaxIdleConns int    `mapstructure:"maxIdleConns"` // 最大空闲连接数
	MaxOpenConns int    `mapstructure:"maxOpenConns"` // 最大打开连接数
}

// CacheConfig 定义缓存连接配置，目前支持 Redis 适配器。
type CacheConfig struct {
	Alias    string `mapstructure:"alias"`    // 缓存别名，用于获取连接
	Section  string `mapstructure:"section"`  // 分区标识
	Adapter  string `mapstructure:"adapter"`  // 适配器类型，如 redis
	Host     string `mapstructure:"host"`     // 主机地址
	Port     int    `mapstructure:"port"`     // 端口号
	Password string `mapstructure:"password"` // 密码
	DB       int    `mapstructure:"db"`       // Redis 数据库编号
}

// MemoryCacheConfig 定义内存缓存配置，基于 go-cache 实现。
type MemoryCacheConfig struct {
	DefaultExpiration int `mapstructure:"defaultExpiration"` // 默认过期时间（秒）
	CleanupInterval   int `mapstructure:"cleanupInterval"`   // 清理间隔（秒）
}

// LogConfig 定义云日志服务配置，支持阿里云 SLS 和腾讯云 CLS。
type LogConfig struct {
	Alias           string `mapstructure:"alias"`           // 日志别名，用于获取 logger
	Type            string `mapstructure:"type"`            // 日志服务类型：sls（阿里云）、cls（腾讯云）
	AccessKeyId     string `mapstructure:"accessKeyId"`     // 访问密钥 ID
	AccessKeySecret string `mapstructure:"accessKeySecret"` // 访问密钥 Secret
	Endpoint        string `mapstructure:"endpoint"`        // 服务端点地址
	AllowLogLevel   string `mapstructure:"allowLogLevel"`   // 允许的最低日志级别
	CloseStdout     bool   `mapstructure:"closeStdout"`     // 是否关闭标准输出
	Project         string `mapstructure:"project"`         // SLS 项目名称
	Logstore        string `mapstructure:"logstore"`        // SLS 日志库名称
	Topic           string `mapstructure:"topic"`           // 日志主题
	Source          string `mapstructure:"source"`          // 日志来源标识
}

var (
	v  *viper.Viper
	bc *BaseConfig
)

// InitConfig 使用默认配置路径 ./config/ 初始化指定环境的配置。
func InitConfig(env string) {
	InitConfigWithPath(env, "./config/")
}

// InitConfigWithPath 使用指定的环境名称和配置路径初始化配置。
// 配置文件命名格式为 config.{env}.yaml，如 config.local.yaml。
// 支持通过 packr 打包配置文件到二进制中。
func InitConfigWithPath(env string, configPath string) {
	fmt.Println(fmt.Sprintf("配置文件路径: %s", configPath))
	fmt.Println(fmt.Sprintf("执行环境: %s", env))
	v = viper.New()
	configName := fmt.Sprintf("config.%s.yaml", env)
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
	configs := packr.New("configs", configPath)
	var data []byte
	var err error
	if data, err = configs.Find(configName); err != nil {
		panic(err)
	}
	if err = v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		panic(err)
	}
	v.Set("env", env)
	baseConfig := BaseConfig{}
	err = v.Unmarshal(&baseConfig)
	if err != nil {
		fmt.Println("yaml parse err: ", err)
		panic(err)
	}
	bc = &baseConfig
}
// GetConfig 返回 Viper 配置实例。如果未初始化则 panic。
func GetConfig() *viper.Viper {
	if v == nil {
		panic("Please init Config")
	}
	return v
}

// GetBaseConfig 返回基础配置结构体指针。如果未初始化则 panic。
func GetBaseConfig() *BaseConfig {
	if bc == nil {
		panic("Please init Config")
	}
	return bc
}

// GetEnv 返回当前环境名称（如 local、test、prod）。
func GetEnv() string {
	if v == nil {
		panic("Please init Config")
	}
	return v.GetString("env")
}
