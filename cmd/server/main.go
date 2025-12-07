package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/lwmacct/251125-go-pkg-logger/pkg/logger"
	app "github.com/lwmacct/251128-workspace/internal/command/server"
)

func main() {
	if err := logger.InitEnv(); err != nil {
		slog.Warn("初始化日志系统失败，使用默认配置", "error", err)
	}
	if err := app.Command.Run(context.Background(), os.Args); err != nil {
		slog.Error("应用程序运行失败", "error", err)
		os.Exit(1)
	}
}
