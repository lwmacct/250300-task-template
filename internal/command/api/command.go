package api

import (
	"github.com/lwmacct/251128-workspace/internal/version"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:   "api",
	Usage:  "简单的 Http 服务器",
	Action: action,
	Commands: []*cli.Command{
		version.Command,
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "addr",
			Aliases: []string{"a"},
			Value:   ":8080",
			Usage:   "服务器监听地址",
		},
		&cli.StringFlag{
			Name:  "dist_docs",
			Value: "docs/.vitepress/dist",
			Usage: "VitePress 文档目录路径",
		},
	},
}
