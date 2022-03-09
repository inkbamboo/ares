package main

import (
	"github.com/inkbamboo/ares/tool/ares/commands"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Cupid"
	app.Usage = "工具集"
	app.Version = commands.GetVersion()
	app.Authors = []*cli.Author{{
		Name:  "LiuXinlei",
		Email: "939896372@qq.com",
	}}
	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help",
		Usage: "查看帮助",
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "GoCMD Version",
	}
	app.Commands = commands.InitCommands()
	err := app.Run(os.Args)
	if err != nil {
		log.Print(err)
	}
}
