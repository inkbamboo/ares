package db

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresqlClient(v *viper.Viper, dbName string) (*gorm.DB, error) {
	config := v.Sub(fmt.Sprintf("postgresql.%s", dbName))
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s TimeZone=%s",
		config.GetString("user"),
		config.GetString("password"),
		config.GetString("host"),
		config.GetInt("port"),
		config.GetString("dbname"),
		config.GetString("sslmode"),
		config.GetString("timezone"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
	return db, err
}
