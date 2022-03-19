package rpc

import (
	"fmt"
	"github.com/inkbamboo/ares/libraries/utils"
	"github.com/rcrowley/go-metrics"
	"time"

	"github.com/smallnest/rpcx/serverplugin"

	"github.com/smallnest/rpcx/server"
	"github.com/spf13/viper"
)

func NewRpcPlugin(v *viper.Viper) (*serverplugin.ZooKeeperRegisterPlugin, error) {
	ip, err := utils.ExternalIP()
	if err != nil {
		return nil, err
	}
	rpcxConfig := v.Sub("rpc")
	rpcPlugin := &serverplugin.ZooKeeperRegisterPlugin{
		ServiceAddress:   fmt.Sprintf("tcp@%s:%s", ip, rpcxConfig.GetString("port")),
		ZooKeeperServers: rpcxConfig.GetStringSlice("servers"),
		BasePath:         rpcxConfig.GetString("basePath"),
		Metrics:          metrics.NewRegistry(),
		UpdateInterval:   time.Minute,
	}
	return rpcPlugin, nil
}

func NewRPCServer(plugin *serverplugin.ZooKeeperRegisterPlugin) (*server.Server, error) {
	s := server.NewServer()
	if err := plugin.Start(); err != nil {
		return nil, err
	}
	s.Plugins.Add(plugin)
	return s, nil
}
