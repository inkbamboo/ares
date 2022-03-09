package db

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func NewMysqlDB(v *viper.Viper, dbName string) (db *gorm.DB, err error) {
	config := v.Sub(fmt.Sprintf("mysql.%s", dbName))
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetString("user"),
		config.GetString("password"),
		config.GetString("host"),
		config.GetInt("port"),
		config.GetString("dbname"))
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxLifetime(config.GetDuration("connMaxLifetime") * time.Minute)
	sqlDB.SetMaxIdleConns(config.GetInt("maxIdleConns"))
	sqlDB.SetMaxOpenConns(config.GetInt("maxOpenConns"))
	return
}
