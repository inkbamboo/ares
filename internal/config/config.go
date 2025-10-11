package config

import (
	"bytes"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
)

type BaseConfig struct {
	AutoMigrate bool              `mapstructure:"autoMigrate,omitempty"` // debug, default false
	Debug       bool              `mapstructure:"debug,omitempty"`       // debug, default false
	Domain      string            `mapstructure:"domain,omitempty"`      // domain=127.0.0.1:8090
	Databases   []DatabaseConfig  `mapstructure:"databases,omitempty"`   // databases
	Caches      []CacheConfig     `mapstructure:"caches,omitempty"`      // caches
	Logs        []LogConfig       `mapstructure:"logs,omitempty"`        // logs
	MemoryCache MemoryCacheConfig `mapstructure:"memoryCache,omitempty"` // caches
}

type DatabaseConfig struct {
	Alias        string `mapstructure:"alias"`        // alias=forum
	Dialect      string `mapstructure:"dialect"`      // dialect=mysql
	Host         string `mapstructure:"host"`         // host=127.0.0.1
	Port         int    `mapstructure:"port"`         // port=3306
	DbName       string `mapstructure:"dbName"`       // name=forum
	Username     string `mapstructure:"username"`     // username=root
	Password     string `mapstructure:"password"`     // password=123456
	MaxIdleConns int    `mapstructure:"maxIdleConns"` // maxIdleConns
	MaxOpenConns int    `mapstructure:"maxOpenConns"` // maxOpenConns
}

type CacheConfig struct {
	Alias    string `mapstructure:"alias"`    // alias=forum
	Section  string `mapstructure:"section"`  // section=forum
	Adapter  string `mapstructure:"adapter"`  // adapter=redis
	Host     string `mapstructure:"host"`     // host=127.0.0.1
	Port     int    `mapstructure:"port"`     // port=6379
	Password string `mapstructure:"password"` // password=123456
	DB       int    `mapstructure:"db"`       // db, select db
}
type MemoryCacheConfig struct {
	DefaultExpiration int `mapstructure:"defaultExpiration"`
	CleanupInterval   int `mapstructure:"defaultExpiration"`
}
type LogConfig struct {
	Alias           string `mapstructure:"alias"`
	Type            string `mapstructure:"type"` // sls, cls
	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	Endpoint        string `mapstructure:"endpoint"`
	AllowLogLevel   string `mapstructure:"allowLogLevel"`
	CloseStdout     bool   `mapstructure:"closeStdout"`
	Project         string `mapstructure:"project"`
	Logstore        string `mapstructure:"logstore"`
	Topic           string `mapstructure:"topic"`
	Source          string `mapstructure:"source"`
}

var (
	v  *viper.Viper
	bc *BaseConfig
)

func InitConfig(env string) {
	InitConfigWithPath(env, "./config/")
}

func InitConfigWithPath(env string, configPath string) {
	fmt.Println(fmt.Sprintf("配置文件路径: %s", configPath))
	fmt.Println(fmt.Sprintf("执行环境: %s", env))
	v = viper.New()
	configName := fmt.Sprintf("config.%s.yaml", env)
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
	configs := packr.New("configs", configPath)
	data, err := configs.Find(configName)
	if err != nil {
		panic(err)
	}
	viper.ReadConfig(bytes.NewBuffer(data))
	if err = v.ReadInConfig(); err != nil {
		fmt.Println(fmt.Sprintf("Viper ReadInConfig err:%s\n", err))
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
func GetConfig() *viper.Viper {
	if v == nil {
		panic("Please init Config")
	}
	return v
}
func GetBaseConfig() *BaseConfig {
	if bc == nil {
		panic("Please init Config")
	}
	return bc
}
func GetEnv() string {
	return v.GetString("env")
}
