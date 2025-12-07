package version

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

// Command 版本信息命令
var Command = &cli.Command{
	Name:  "version",
	Usage: "显示版本信息",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "short",
			Aliases: []string{"s"},
			Usage:   "显示简短版本信息",
		},
		&cli.BoolFlag{
			Name:    "json",
			Aliases: []string{"j"},
			Usage:   "以JSON格式显示版本信息",
		},
	},
	Action: action, // 根据标志显示不同格式
}

// action 根据标志显示不同格式的版本信息
func action(ctx context.Context, c *cli.Command) error {
	switch {
	case c.Bool("json"):
		PrintVersionJSON()
	case c.Bool("short"):
		fmt.Println(GetVersion())
	default:
		PrintBuildInfo()
	}
	return nil
}
