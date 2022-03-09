package commands

import "github.com/urfave/cli/v2"

type Flags struct {
	ServiceName string
	Path        string
	ProjectName string
	ModPrefix   string
	RedisList   []DbClient
	MysqlList   []DbClient
}
type DbClient struct {
	ConfigName string
	ClientName string
}

var f *Flags

func (f *Flags) ToNewAction() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "d",
			Value:       "",
			Usage:       "指定项目所在目录",
			Destination: &f.Path,
		},
	}
}
func (f *Flags) ToServiceAction() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "d",
			Value:       "",
			Usage:       "指定项目所在目录",
			Destination: &f.Path,
		},
	}
}
