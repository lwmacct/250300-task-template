package api

import "github.com/urfave/cli/v3"

func Command(version string) *cli.Command {
	return &cli.Command{
		Name:    "api",
		Usage:   "简单的 Http 服务器",
		Version: version,
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
		Action: action,
	}
}
