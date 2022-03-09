package config

import (
	"bytes"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func NewConfig() (*viper.Viper, error) {
	pflag.String("env", "local", "env")
	pflag.String("configPath", "configs/", "configPath")
	pflag.String("data", "file", "配置来源")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	viper.AddConfigPath(viper.GetString("configPath")) //设置读取的文件路径
	viper.SetConfigType("yaml")
	box := packr.New("My Box", "./configs")
	data, err := box.Find("config.yaml")
	if err != nil {
		return nil, err
	}
	viper.ReadConfig(bytes.NewBuffer(data))
	if env := viper.GetString("env"); env != "" {
		if data, err = box.Find("config." + env + ".yaml"); err != nil {
			return nil, err
		}
		if err = viper.MergeConfig(bytes.NewBuffer(data)); err != nil {
			return nil, err
		}
	}
	return viper.GetViper(), err
}
