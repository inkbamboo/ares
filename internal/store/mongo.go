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

// MongoDB 封装 MongoDB 客户端连接，包含客户端实例、上下文和数据库对象。
type MongoDB struct {
	Client  *mongo.Client
	Context context.Context
	*mongo.Database
}

// Close 关闭 MongoDB 连接。
func (d *MongoDB) Close() {
	d.Client.Disconnect(d.Context)
}

// NewMongo 根据数据库配置创建 MongoDB 连接。
// 连接超时时间为 10 秒，Ping 超时时间为 5 秒。
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
	ctxPing, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelPing()
	err = client.Ping(ctxPing, readpref.Primary())
	if err != nil {
		client.Disconnect(context.Background())
		return nil
	}
	// 使用不会被取消的 context，确保后续 Close() 可正常调用 Disconnect
	mongodb.Context = context.Background()
	mongodb.Client = client
	mongodb.Database = db
	return mongodb
}
