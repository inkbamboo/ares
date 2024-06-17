package store

import (
	"context"
	"fmt"
	"github.com/inkbamboo/ares/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo(database config.DatabaseConfig, debug bool) *mongo.Database {
	dsn := "mongodb://"
	if database.Username != "" {
		dsn += database.Username + ":" + database.Password + "@"
	}
	dsn += database.Host + ":" + fmt.Sprintf("%d", database.Port)
	opts := options.Client().ApplyURI(dsn)
	client, _ := mongo.Connect(context.TODO(), opts)
	db := client.Database(database.DbName)
	return db
}
