// Package client 提供 HTTP 客户端命令测试
package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lwmacct/251128-workspace/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewHTTPClient 测试创建 HTTP 客户端
func TestNewHTTPClient(t *testing.T) {
	cfg := &config.ClientConfig{
		URL:     "http://localhost:8080",
		Timeout: 30 * time.Second,
		Retries: 3,
	}

	client := NewHTTPClient(cfg)

	assert.NotNil(t, client)
	assert.Equal(t, cfg, client.config)
	assert.NotNil(t, client.client)
	assert.Equal(t, cfg.Timeout, client.client.Timeout)
}

// TestHTTPClientHealth 测试健康检查
func TestHTTPClientHealth(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := &config.ClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
		Retries: 0,
	}

	client := NewHTTPClient(cfg)
	resp, err := client.Health(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}

// TestHTTPClientHealthWithRetry 测试健康检查重试
func TestHTTPClientHealthWithRetry(t *testing.T) {
	attempts := 0

	// 创建测试服务器，第一次失败，第二次成功
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	}))
	defer server.Close()

	cfg := &config.ClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
		Retries: 3,
	}

	client := NewHTTPClient(cfg)
	resp, err := client.Health(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
	assert.Equal(t, 2, attempts)
}

// TestHTTPClientHealthFailure 测试健康检查失败
func TestHTTPClientHealthFailure(t *testing.T) {
	// 创建始终失败的测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := &config.ClientConfig{
		URL:     server.URL,
		Timeout: 1 * time.Second,
		Retries: 1,
	}

	client := NewHTTPClient(cfg)
	_, err := client.Health(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "health check failed after")
}

// TestHTTPClientGet 测试 GET 请求
func TestHTTPClientGet(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Hello, World!"}`))
	}))
	defer server.Close()

	cfg := &config.ClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
		Retries: 0,
	}

	client := NewHTTPClient(cfg)
	body, err := client.Get(context.Background(), "/")

	require.NoError(t, err)
	assert.Contains(t, body, "Hello, World!")
}

// TestHTTPClientGetWithPath 测试带路径的 GET 请求
func TestHTTPClientGetWithPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/users" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"users":[]}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := &config.ClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
		Retries: 0,
	}

	client := NewHTTPClient(cfg)

	// 测试带斜杠的路径
	body, err := client.Get(context.Background(), "/api/users")
	require.NoError(t, err)
	assert.Contains(t, body, "users")

	// 测试不带斜杠的路径
	body, err = client.Get(context.Background(), "api/users")
	require.NoError(t, err)
	assert.Contains(t, body, "users")
}

// TestHTTPClientGetFailure 测试 GET 请求失败
func TestHTTPClientGetFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	cfg := &config.ClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
		Retries: 0,
	}

	client := NewHTTPClient(cfg)
	_, err := client.Get(context.Background(), "/nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

// TestHTTPClientGetWithRetry 测试 GET 请求重试
func TestHTTPClientGetWithRetry(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))
	defer server.Close()

	cfg := &config.ClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
		Retries: 5,
	}

	client := NewHTTPClient(cfg)
	body, err := client.Get(context.Background(), "/")

	require.NoError(t, err)
	assert.Equal(t, "success", body)
	assert.Equal(t, 3, attempts)
}

// TestHTTPClientDoRequestError 测试请求错误
func TestHTTPClientDoRequestError(t *testing.T) {
	cfg := &config.ClientConfig{
		URL:     "http://invalid-host-that-does-not-exist:99999",
		Timeout: 1 * time.Second,
		Retries: 0,
	}

	client := NewHTTPClient(cfg)
	_, err := client.Health(context.Background())

	assert.Error(t, err)
}

// TestHTTPClientContextCancellation 测试上下文取消
func TestHTTPClientContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟慢请求
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.ClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
		Retries: 0,
	}

	client := NewHTTPClient(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.Get(ctx, "/")
	assert.Error(t, err)
}

// TestHealthResponseJSON 测试 HealthResponse JSON 序列化
func TestHealthResponseJSON(t *testing.T) {
	resp := HealthResponse{Status: "ok"}

	data, err := json.Marshal(resp)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"status":"ok"`)

	var parsed HealthResponse
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)
	assert.Equal(t, "ok", parsed.Status)
}

// TestCommandFlags 测试命令标志定义
func TestCommandFlags(t *testing.T) {
	assert.NotNil(t, Command)
	assert.Equal(t, "client", Command.Name)
	assert.NotEmpty(t, Command.Usage)
	assert.NotEmpty(t, Command.Flags)

	// 验证子命令
	assert.Len(t, Command.Commands, 3) // version, health, get
}

// TestHealthCommand 测试健康检查命令定义
func TestHealthCommand(t *testing.T) {
	assert.NotNil(t, healthCommand)
	assert.Equal(t, "health", healthCommand.Name)
	assert.NotEmpty(t, healthCommand.Usage)
}

// TestGetCommand 测试 GET 命令定义
func TestGetCommand(t *testing.T) {
	assert.NotNil(t, getCommand)
	assert.Equal(t, "get", getCommand.Name)
	assert.NotEmpty(t, getCommand.Usage)
}

// TestDefaultsInitialization 测试默认值初始化
func TestDefaultsInitialization(t *testing.T) {
	assert.NotEmpty(t, defaults.URL)
	assert.NotZero(t, defaults.Timeout)
	assert.NotZero(t, defaults.Retries)
}
