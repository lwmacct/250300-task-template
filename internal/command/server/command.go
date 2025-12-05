// Package server 提供 HTTP 服务器命令
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/lwmacct/251125-go-mod-logger/pkg/logger"
	"github.com/lwmacct/251128-workspace/internal/command"
	"github.com/lwmacct/251128-workspace/internal/config"
	"github.com/lwmacct/251128-workspace/internal/version"
	"github.com/urfave/cli/v3"
)

// Command 服务器命令
var Command = &cli.Command{
	Name:     "server",
	Usage:    "启动 HTTP 服务器",
	Action:   action,
	Commands: []*cli.Command{version.Command},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "server-addr",
			Aliases: []string{"a"},
			Value:   command.Defaults.Server.Addr,
			Usage:   "服务器监听地址",
		},
		&cli.StringFlag{
			Name:  "server-docs",
			Value: command.Defaults.Server.Docs,
			Usage: "VitePress 文档目录路径",
		},
		&cli.DurationFlag{
			Name:  "server-timeout",
			Value: command.Defaults.Server.Timeout,
			Usage: "HTTP 读写超时",
		},
		&cli.DurationFlag{
			Name:  "server-idletime",
			Value: command.Defaults.Server.Idletime,
			Usage: "HTTP 空闲超时",
		},
	},
}

func action(ctx context.Context, cmd *cli.Command) error {
	if err := logger.InitEnv(); err != nil {
		slog.Warn("初始化日志系统失败，使用默认配置", "error", err)
	}

	// 加载配置：默认值 → 配置文件 → 环境变量 → CLI flags
	cfg, err := config.Load(cmd, "", version.AppRawName)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	// 健康检查端点
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"status":"ok"}`)
	})

	// VitePress 文档静态文件服务
	docsFS := http.FileServer(http.Dir(cfg.Server.Docs))
	mux.Handle("/docs/", http.StripPrefix("/docs/", docsFS))

	// 默认首页（{$} 精确匹配根路径）
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"message":"Hello, World!"}`)
	})

	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      mux,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.Idletime,
	}

	// 启动服务器（非阻塞）
	go func() {
		slog.Info("Server starting", "addr", cfg.Server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down")

	// 优雅关闭，最多等待 10 秒
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.Timeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	slog.Info("Server stopped gracefully")
	return nil
}
