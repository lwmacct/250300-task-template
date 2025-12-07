// Package client 提供 HTTP 客户端命令。
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/lwmacct/251125-go-mod-logger/pkg/logger"
	"github.com/lwmacct/251128-workspace/internal/command"
	"github.com/lwmacct/251128-workspace/internal/config"
	"github.com/lwmacct/251128-workspace/internal/version"
	"github.com/urfave/cli/v3"
)

// Command 客户端命令
var Command = &cli.Command{
	Name:     "client",
	Usage:    "HTTP 客户端工具",
	Action:   action,
	Commands: []*cli.Command{version.Command, healthCommand, getCommand},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "client-url",
			Aliases: []string{"s"},
			Value:   command.Defaults.Client.URL,
			Usage:   "服务器地址",
		},
		&cli.DurationFlag{
			Name:  "client-timeout",
			Value: command.Defaults.Client.Timeout,
			Usage: "请求超时时间",
		},
		&cli.IntFlag{
			Name:  "client-retries",
			Value: command.Defaults.Client.Retries,
			Usage: "重试次数",
		},
	},
}

// healthCommand 健康检查子命令
var healthCommand = &cli.Command{
	Name:   "health",
	Usage:  "检查服务器健康状态",
	Action: healthAction,
}

// getCommand GET 请求子命令
var getCommand = &cli.Command{
	Name:      "get",
	Usage:     "发送 GET 请求",
	ArgsUsage: "[path]",
	Action:    getAction,
}

func action(ctx context.Context, cmd *cli.Command) error {
	// 默认行为：显示帮助
	return cli.ShowAppHelp(cmd)
}

func healthAction(ctx context.Context, cmd *cli.Command) error {
	if err := logger.InitEnv(); err != nil {
		slog.Warn("初始化日志系统失败，使用默认配置", "error", err)
	}

	cfg, err := config.Load(cmd, "", version.AppRawName)
	if err != nil {
		return err
	}

	client := NewHTTPClient(&cfg.Client)
	resp, err := client.Health(ctx)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	fmt.Printf("Server: %s\n", cfg.Client.URL)
	fmt.Printf("Status: %s\n", resp.Status)
	return nil
}

func getAction(ctx context.Context, cmd *cli.Command) error {

	cfg, err := config.Load(cmd, "", version.AppRawName)
	if err != nil {
		return err
	}

	path := "/"
	if cmd.NArg() > 0 {
		path = cmd.Args().First()
	}

	client := NewHTTPClient(&cfg.Client)
	body, err := client.Get(ctx, path)
	if err != nil {
		return fmt.Errorf("GET request failed: %w", err)
	}

	fmt.Println(body)
	return nil
}

// HTTPClient HTTP 客户端封装
type HTTPClient struct {
	config *config.ClientConfig
	client *http.Client
}

// NewHTTPClient 创建新的 HTTP 客户端
func NewHTTPClient(cfg *config.ClientConfig) *HTTPClient {
	return &HTTPClient{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status string `json:"status"`
}

// Health 执行健康检查
func (c *HTTPClient) Health(ctx context.Context) (*HealthResponse, error) {
	url := strings.TrimSuffix(c.config.URL, "/") + "/health"

	var lastErr error
	for i := 0; i <= c.config.Retries; i++ {
		resp, err := c.doRequest(ctx, "GET", url)
		if err != nil {
			lastErr = err
			slog.Debug("Health check attempt failed", "attempt", i+1, "error", err)
			continue
		}

		var health HealthResponse
		if err := json.Unmarshal([]byte(resp), &health); err != nil {
			return nil, fmt.Errorf("failed to parse health response: %w", err)
		}
		return &health, nil
	}

	return nil, fmt.Errorf("health check failed after %d retries: %w", c.config.Retries, lastErr)
}

// Get 发送 GET 请求
func (c *HTTPClient) Get(ctx context.Context, path string) (string, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	url := strings.TrimSuffix(c.config.URL, "/") + path

	var lastErr error
	for i := 0; i <= c.config.Retries; i++ {
		resp, err := c.doRequest(ctx, "GET", url)
		if err != nil {
			lastErr = err
			slog.Debug("GET request attempt failed", "attempt", i+1, "error", err)
			continue
		}
		return resp, nil
	}

	return "", fmt.Errorf("GET request failed after %d retries: %w", c.config.Retries, lastErr)
}

// doRequest 执行 HTTP 请求
func (c *HTTPClient) doRequest(ctx context.Context, method, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log the error but don't override the return error
			_ = err
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}
