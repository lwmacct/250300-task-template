package main

import (
	"context"
	"log/slog"
	"os"

	app "github.com/lwmacct/251128-workspace/internal/command/client"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm"
)

func main() {
	logm.MustInit(logm.PresetAuto()...)

	if err := app.Command.Run(context.Background(), os.Args); err != nil {
		slog.Error("应用程序运行失败", "error", err)
		os.Exit(1)
	}
}
