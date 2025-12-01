package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lwmacct/251125-go-mod-logger/pkg/logger"
	"github.com/lwmacct/251128-workspace/internal/commands/api"
)

var version = "v0.0.0"

func main() {
	if err := logger.InitEnv(); err != nil {
		slog.Warn("初始化日志系统失败，使用默认配置", "error", err)
	}

	cmd := api.Command(version)
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
