package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v3"
)

func action(ctx context.Context, cmd *cli.Command) error {

	addr := cmd.String("addr")
	docsDir := cmd.String("dist_docs")

	mux := http.NewServeMux()

	// 健康检查端点
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// VitePress 文档静态文件服务
	docsFS := http.FileServer(http.Dir(docsDir))
	mux.Handle("/docs/", http.StripPrefix("/docs/", docsFS))

	// 默认首页（{$} 精确匹配根路径）
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message":"Hello, World!"}`)
	})

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器（非阻塞）
	go func() {
		slog.Info("Server starting", "addr", addr)
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
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	slog.Info("Server stopped gracefully")
	return nil
}
