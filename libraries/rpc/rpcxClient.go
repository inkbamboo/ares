package rpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/rpcxio/libkv/store"
	"github.com/smallnest/rpcx/client"
)

type (
	XClientPool struct {
		pPrefix string
		clients cmap.ConcurrentMap
		servers []string
	}
)

func NewXClientPool(pPrefix string, servers []string) (*XClientPool, error) {
	return &XClientPool{
		pPrefix: pPrefix,
		clients: cmap.New(),
		servers: servers,
	}, nil
}

func (s *XClientPool) G(pSuffix string, service string) (client.XClient, error) {
	if c, ok := s.clients.Get(s.key(pSuffix, service)); ok {
		return c.(client.XClient), nil
	}
	d, err := client.NewZookeeperDiscovery(s.pPrefix+"_"+pSuffix, service, s.servers, &store.Config{
		PersistConnection: true,
	})
	if err != nil {
		return nil, err
	}
	xclient := client.NewXClient(service, client.Failover, client.RandomSelect, d, client.DefaultOption)
	s.clients.Set(s.key(pSuffix, service), xclient)
	return xclient, nil
}

// Call serviceMethod: {包名后缀}_{服务名}_{成员函数名} 比如: activity_TaskService_Draw
func (s *XClientPool) Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error {
	results := strings.Split(serviceMethod, "_")
	if len(results) != 3 {
		return errors.New("serviceMethod split must be 3")
	}
	client, err := s.G(results[0], results[1])
	if err != nil {
		return err
	}
	return client.Call(ctx, results[2], args, reply)
}

func (s *XClientPool) Must(c client.XClient, err error) client.XClient {
	if err != nil {
		panic(err)
	}
	return c
}

func (s *XClientPool) key(pSuffix string, service string) string {
	return fmt.Sprintf("%s_%s_%s", s.pPrefix, pSuffix, service)
}
