package api

import (
	"github.com/lwmacct/251128-workspace/internal/config"
	"github.com/lwmacct/251128-workspace/internal/version"
	"github.com/urfave/cli/v3"
)

// 默认配置 - 单一来源 (Single Source of Truth)
var defaults = config.DefaultConfig()

var Command = &cli.Command{
	Name:     "api",
	Usage:    "简单的 Http 服务器",
	Action:   action,
	Commands: []*cli.Command{version.Command},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "addr",
			Aliases: []string{"a"},
			Value:   defaults.Addr,
			Usage:   "服务器监听地址",
		},
		&cli.StringFlag{
			Name:  "dist_docs",
			Value: defaults.DistDocs,
			Usage: "VitePress 文档目录路径",
		},
	},
}
