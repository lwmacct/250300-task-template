// Package version 提供版本信息显示命令。
//
// 支持多种输出格式：
//   - 默认：详细版本信息
//   - short：仅版本号
//   - json：JSON 格式的完整信息
//
// Author: lwmacct (https://github.com/lwmacct)
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
	Commands: []*cli.Command{
		{
			Name:   "short",
			Usage:  "显示简短版本信息",
			Action: shortAction,
		},
		{
			Name:   "json",
			Usage:  "以JSON格式显示版本信息",
			Action: jsonAction,
		},
	},
	Action: showAction, // 默认显示详细信息
}

// showAction 显示详细版本信息
func showAction(ctx context.Context, c *cli.Command) error {
	PrintVersion()
	return nil
}

// shortAction 显示简短版本信息
func shortAction(ctx context.Context, c *cli.Command) error {
	fmt.Println(GetVersion())
	return nil
}

// jsonAction 以JSON格式显示版本信息
func jsonAction(ctx context.Context, c *cli.Command) error {
	fmt.Printf(`{
  "appRawName": "%s",
  "appProject": "%s",
  "appVersion": "%s",
  "gitCommit": "%s",
  "buildTime": "%s",
  "developer": "%s",
  "workspace": "%s"
}
`, AppRawName, AppProject, AppVersion, GitCommit, BuildTime, Developer, Workspace)
	return nil
}
