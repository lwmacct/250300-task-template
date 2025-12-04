// Author: lwmacct (https://github.com/lwmacct)
package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDefaultConfig 测试默认配置值
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// 验证 Server 配置
	assert.Equal(t, ":8080", cfg.Server.Addr)
	assert.Equal(t, "docs/.vitepress/dist", cfg.Server.Docs)
	assert.Equal(t, 15*time.Second, cfg.Server.Timeout)
	assert.Equal(t, 60*time.Second, cfg.Server.Idletime)

	// 验证 Client 配置
	assert.Equal(t, "http://localhost:8080", cfg.Client.URL)
	assert.Equal(t, 30*time.Second, cfg.Client.Timeout)
	assert.Equal(t, 3, cfg.Client.Retries)
}

// TestDefaultServerConfig 测试默认服务端配置
func TestDefaultServerConfig(t *testing.T) {
	cfg := DefaultServerConfig()

	assert.Equal(t, ":8080", cfg.Addr)
	assert.Equal(t, "docs/.vitepress/dist", cfg.Docs)
	assert.Equal(t, 15*time.Second, cfg.Timeout)
	assert.Equal(t, 60*time.Second, cfg.Idletime)
}

// TestDefaultClientConfig 测试默认客户端配置
func TestDefaultClientConfig(t *testing.T) {
	cfg := DefaultClientConfig()

	assert.Equal(t, "http://localhost:8080", cfg.URL)
	assert.Equal(t, 30*time.Second, cfg.Timeout)
	assert.Equal(t, 3, cfg.Retries)
}

// TestLoadWithDefaults 测试使用默认值加载配置
func TestLoadWithDefaults(t *testing.T) {
	cfg, err := Load(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 验证使用默认值
	expected := DefaultConfig()
	assert.Equal(t, expected.Server.Addr, cfg.Server.Addr)
	assert.Equal(t, expected.Client.URL, cfg.Client.URL)
}

// TestLoadServerConfig 测试加载服务端配置
func TestLoadServerConfig(t *testing.T) {
	cfg, err := LoadServerConfig(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	expected := DefaultServerConfig()
	assert.Equal(t, expected.Addr, cfg.Addr)
	assert.Equal(t, expected.Docs, cfg.Docs)
}

// TestLoadClientConfig 测试加载客户端配置
func TestLoadClientConfig(t *testing.T) {
	cfg, err := LoadClientConfig(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	expected := DefaultClientConfig()
	assert.Equal(t, expected.URL, cfg.URL)
	assert.Equal(t, expected.Timeout, cfg.Timeout)
}

// TestLoadWithEnvVars 测试使用环境变量加载配置
func TestLoadWithEnvVars(t *testing.T) {
	// 设置环境变量
	_ = os.Setenv("APP_SERVER_ADDR", ":9090")
	_ = os.Setenv("APP_CLIENT_URL", "http://test:8080")
	defer func() {
		_ = os.Unsetenv("APP_SERVER_ADDR")
		_ = os.Unsetenv("APP_CLIENT_URL")
	}()

	cfg, err := Load(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 环境变量应该覆盖默认值
	assert.Equal(t, ":9090", cfg.Server.Addr)
	assert.Equal(t, "http://test:8080", cfg.Client.URL)
}

// TestLoadWithCustomEnvPrefix 测试使用自定义环境变量前缀
func TestLoadWithCustomEnvPrefix(t *testing.T) {
	// 设置自定义前缀的环境变量
	_ = os.Setenv("MYAPP_SERVER_ADDR", ":7070")
	defer func() { _ = os.Unsetenv("MYAPP_SERVER_ADDR") }()

	cfg, err := Load(nil, "MYAPP_")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, ":7070", cfg.Server.Addr)
}

// TestLoadWithYAMLFile 测试从 YAML 文件加载配置
func TestLoadWithYAMLFile(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  addr: ":3000"
  docs: "/custom/path"
client:
  url: "http://custom:3000"
  retries: 5
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// 切换到临时目录
	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(oldWd) }()

	cfg, err := Load(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, ":3000", cfg.Server.Addr)
	assert.Equal(t, "/custom/path", cfg.Server.Docs)
	assert.Equal(t, "http://custom:3000", cfg.Client.URL)
	assert.Equal(t, 5, cfg.Client.Retries)
}

// TestLoadWithConfigSubdir 测试从 config/ 子目录加载配置
func TestLoadWithConfigSubdir(t *testing.T) {
	// 创建临时配置文件在 config/ 子目录
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "config")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, "config.yaml")
	configContent := `
server:
  addr: ":4000"
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// 切换到临时目录
	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(oldWd) }()

	cfg, err := Load(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, ":4000", cfg.Server.Addr)
}

// TestConfigPriority 测试配置优先级
func TestConfigPriority(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  addr: ":3000"
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// 设置环境变量（优先级高于配置文件）
	_ = os.Setenv("APP_SERVER_ADDR", ":4000")
	defer func() { _ = os.Unsetenv("APP_SERVER_ADDR") }()

	// 切换到临时目录
	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(oldWd) }()

	cfg, err := Load(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 环境变量应该覆盖配置文件
	assert.Equal(t, ":4000", cfg.Server.Addr)
}

// TestStructTags 验证 Config 结构体字段与 koanf 标签一致性
func TestStructTags(t *testing.T) {
	cfg := DefaultConfig()
	cfgType := reflect.TypeOf(cfg)

	// 验证所有字段都有 koanf 标签
	checkKoanfTags(t, cfgType, "Config")
}

// TestServerConfigStructTags 验证 ServerConfig 结构体字段标签
func TestServerConfigStructTags(t *testing.T) {
	cfg := DefaultServerConfig()
	cfgType := reflect.TypeOf(cfg)

	checkKoanfTags(t, cfgType, "ServerConfig")
}

// TestClientConfigStructTags 验证 ClientConfig 结构体字段标签
func TestClientConfigStructTags(t *testing.T) {
	cfg := DefaultClientConfig()
	cfgType := reflect.TypeOf(cfg)

	checkKoanfTags(t, cfgType, "ClientConfig")
}

// checkKoanfTags 检查结构体字段的 koanf 标签
func checkKoanfTags(t *testing.T, typ reflect.Type, path string) {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldPath := path + "." + field.Name

		// 检查 koanf 标签
		koanfTag := field.Tag.Get("koanf")
		if koanfTag == "" {
			t.Errorf("字段 %s 缺少 koanf 标签", fieldPath)
		}

		// 检查 comment 标签
		commentTag := field.Tag.Get("comment")
		if commentTag == "" {
			t.Errorf("字段 %s 缺少 comment 标签", fieldPath)
		}

		// 递归检查嵌套结构体
		if field.Type.Kind() == reflect.Struct &&
			field.Type.String() != "time.Duration" &&
			field.Type.String() != "time.Time" {
			checkKoanfTags(t, field.Type, fieldPath)
		}
	}
}
