package db

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func NewRedis(v *viper.Viper, dbName string) (client *redis.Client, err error) {
	redisConfig := v.Sub(fmt.Sprintf("redis.%s", dbName))
	client = redis.NewClient(&redis.Options{
		Addr:     redisConfig.GetString("addr"),
		Password: redisConfig.GetString("password"),
		DB:       redisConfig.GetInt("db"),
	})
	return client, nil
}
