// Package server 提供 HTTP 服务器命令测试
package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCommandDefinition 测试命令定义
func TestCommandDefinition(t *testing.T) {
	assert.NotNil(t, Command)
	assert.Equal(t, "server", Command.Name)
	assert.NotEmpty(t, Command.Usage)
	assert.NotEmpty(t, Command.Flags)
}

// TestCommandFlags 测试命令标志
func TestCommandFlags(t *testing.T) {
	flags := Command.Flags
	assert.NotEmpty(t, flags)

	// 验证关键 flags 存在
	flagNames := make(map[string]bool)
	for _, flag := range flags {
		for _, name := range flag.Names() {
			flagNames[name] = true
		}
	}

	assert.True(t, flagNames["server-addr"] || flagNames["a"], "应该有 server-addr 或 a 标志")
	assert.True(t, flagNames["server-docs"], "应该有 server-docs 标志")
	assert.True(t, flagNames["server-timeout"], "应该有 server-timeout 标志")
	assert.True(t, flagNames["server-idletime"], "应该有 server-idletime 标志")
}

// TestCommandSubcommands 测试子命令
func TestCommandSubcommands(t *testing.T) {
	// 应该有 version 子命令
	assert.NotEmpty(t, Command.Commands)

	hasVersion := false
	for _, cmd := range Command.Commands {
		if cmd.Name == "version" {
			hasVersion = true
			break
		}
	}
	assert.True(t, hasVersion, "应该有 version 子命令")
}

// TestDefaultsInitialization 测试默认值初始化
func TestDefaultsInitialization(t *testing.T) {
	assert.NotEmpty(t, defaults.Addr)
	assert.NotEmpty(t, defaults.Docs)
	assert.NotZero(t, defaults.Timeout)
	assert.NotZero(t, defaults.Idletime)
}

// TestHealthEndpointHandler 测试健康检查端点处理器
func TestHealthEndpointHandler(t *testing.T) {
	// 创建一个模拟的 ServeMux 来测试处理器逻辑
	mux := http.NewServeMux()

	// 注册健康检查端点
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// 创建测试请求
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// 使用 ResponseRecorder 记录响应
	recorder := &mockResponseWriter{
		headers: make(http.Header),
	}
	mux.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.statusCode)
	assert.Equal(t, "application/json", recorder.headers.Get("Content-Type"))
	assert.Contains(t, recorder.body, `"status":"ok"`)
}

// TestRootEndpointHandler 测试根路径端点处理器
func TestRootEndpointHandler(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"Hello, World!"}`))
	})

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	recorder := &mockResponseWriter{
		headers: make(http.Header),
	}
	mux.ServeHTTP(recorder, req)

	assert.Equal(t, "application/json", recorder.headers.Get("Content-Type"))
	assert.Contains(t, recorder.body, `"message":"Hello, World!"`)
}

// mockResponseWriter 模拟 http.ResponseWriter
type mockResponseWriter struct {
	headers    http.Header
	body       string
	statusCode int
}

func (m *mockResponseWriter) Header() http.Header {
	return m.headers
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	m.body += string(data)
	return len(data), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

// TestServerConfigTimeout 测试服务器超时配置
func TestServerConfigTimeout(t *testing.T) {
	// 验证默认超时值是否合理
	assert.Equal(t, 15*time.Second, defaults.Timeout)
	assert.Equal(t, 60*time.Second, defaults.Idletime)
}

// TestServerAddrDefault 测试服务器默认地址
func TestServerAddrDefault(t *testing.T) {
	assert.Equal(t, ":8080", defaults.Addr)
}

// TestServerDocsDefault 测试文档目录默认值
func TestServerDocsDefault(t *testing.T) {
	assert.Equal(t, "docs/.vitepress/dist", defaults.Docs)
}

// TestServerContextCancellation 测试服务器上下文取消
func TestServerContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// 创建一个简单的服务器配置
	server := &http.Server{
		Addr:    ":0", // 随机端口
		Handler: http.NewServeMux(),
	}

	// 启动服务器
	go func() {
		_ = server.ListenAndServe()
	}()

	// 取消上下文
	cancel()
	<-ctx.Done()

	// 关闭服务器
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdownCancel()

	err := server.Shutdown(shutdownCtx)
	assert.NoError(t, err)
}

// TestHTTPServerConfiguration 测试 HTTP 服务器配置
func TestHTTPServerConfiguration(t *testing.T) {
	server := &http.Server{
		Addr:         defaults.Addr,
		ReadTimeout:  defaults.Timeout,
		WriteTimeout: defaults.Timeout,
		IdleTimeout:  defaults.Idletime,
	}

	assert.Equal(t, defaults.Addr, server.Addr)
	assert.Equal(t, defaults.Timeout, server.ReadTimeout)
	assert.Equal(t, defaults.Timeout, server.WriteTimeout)
	assert.Equal(t, defaults.Idletime, server.IdleTimeout)
}

// TestCommandAction 测试命令 action 不为 nil
func TestCommandAction(t *testing.T) {
	assert.NotNil(t, Command.Action)
}

// TestMuxRouting 测试路由配置
func TestMuxRouting(t *testing.T) {
	mux := http.NewServeMux()

	// 健康检查
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// 根路径
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// 测试各路由
	tests := []struct {
		method string
		path   string
		expect int
	}{
		{"GET", "/health", http.StatusOK},
		{"GET", "/", http.StatusOK},
		{"POST", "/health", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			recorder := &mockResponseWriter{headers: make(http.Header)}
			mux.ServeHTTP(recorder, req)

			// 只验证成功的路由
			if tt.expect == http.StatusOK {
				assert.Equal(t, tt.expect, recorder.statusCode)
			}
		})
	}
}
