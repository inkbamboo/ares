package store

import (
	"context"
	"fmt"
	"github.com/inkbamboo/ares/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoDB struct {
	Client  *mongo.Client
	Context context.Context
	*mongo.Database
}

// Close closes the mongo-go-driver connection.
func (d *MongoDB) Close() {
	d.Client.Disconnect(d.Context)
}

func NewMongo(database config.DatabaseConfig, debug bool) *MongoDB {
	mongodb := &MongoDB{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dsn := "mongodb://"
	if database.Username != "" {
		dsn += database.Username + ":" + database.Password + "@"
	}
	dsn += database.Host + ":" + fmt.Sprintf("%d", database.Port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil
	}
	db := client.Database(database.DbName)
	ctxPing, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctxPing, readpref.Primary())
	if err != nil {
		return nil
	}
	mongodb.Context = ctx
	mongodb.Client = client
	mongodb.Database = db
	return mongodb
}
