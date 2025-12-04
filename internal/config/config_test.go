// Author: lwmacct (https://github.com/lwmacct)
package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
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

// TestGenerateExample 生成 config/config.example.yaml 配置示例文件
//
// 此测试会从 defaultConfig() 提取默认配置值，并生成符合规范的 YAML 配置示例文件。
// 设计为可被 pre-commit hook 调用，在 config.go 变更时自动执行。
//
// 运行方式:
//
//	go test -v -run TestGenerateExample ./internal/config/...
func TestGenerateExample(t *testing.T) {
	// 获取项目根目录
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("无法找到项目根目录: %v", err)
	}

	// 获取默认配置
	cfg := DefaultConfig()

	// 生成 YAML 内容
	var buf bytes.Buffer
	writeConfigYAML(&buf, cfg)

	// 确保 config 目录存在
	configDir := filepath.Join(projectRoot, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("创建 config 目录失败: %v", err)
	}

	// 写入文件
	outputPath := filepath.Join(configDir, "config.example.yaml")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		t.Fatalf("写入配置文件失败: %v", err)
	}

	t.Logf("✅ 已生成配置示例文件: %s", outputPath)
}

// TestConfigKeysValid 验证 config.yaml 不包含 config.example.yaml 中不存在的配置项
//
// 此测试确保用户的配置文件不会有未知的配置项，防止因拼写错误或过时配置导致的问题。
// 如果 config.yaml 不存在，测试会跳过。
func TestConfigKeysValid(t *testing.T) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("无法找到项目根目录: %v", err)
	}

	configPath := filepath.Join(projectRoot, "config", "config.yaml")
	examplePath := filepath.Join(projectRoot, "config", "config.example.yaml")

	// 如果 config.yaml 不存在，跳过测试
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("config.yaml 不存在，跳过验证")
	}

	// 加载 config.example.yaml 获取有效的 keys
	exampleKeys, err := loadYAMLKeys(examplePath)
	if err != nil {
		t.Fatalf("无法加载 config.example.yaml: %v", err)
	}

	// 加载 config.yaml 获取用户配置的 keys
	configKeys, err := loadYAMLKeys(configPath)
	if err != nil {
		t.Fatalf("无法加载 config.yaml: %v", err)
	}

	// 检查 config.yaml 中是否有 config.example.yaml 不存在的 keys
	var invalidKeys []string
	for _, key := range configKeys {
		if !containsKey(exampleKeys, key) {
			invalidKeys = append(invalidKeys, key)
		}
	}

	if len(invalidKeys) > 0 {
		t.Errorf("config.yaml 包含以下无效配置项 (在 config.example.yaml 中不存在):\n")
		for _, key := range invalidKeys {
			t.Errorf("  - %s", key)
		}
		t.Errorf("\n请检查拼写或从 config.example.yaml 中确认有效的配置项")
	}
}

// loadYAMLKeys 加载 YAML 文件并返回所有配置键的扁平化列表
func loadYAMLKeys(path string) ([]string, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("加载文件失败: %w", err)
	}

	return k.Keys(), nil
}

// containsKey 检查 keys 列表中是否包含指定的 key
func containsKey(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
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

// findProjectRoot 通过查找 go.mod 文件定位项目根目录
func findProjectRoot() (string, error) {
	// 获取当前测试文件所在目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("无法获取当前文件路径")
	}

	dir := filepath.Dir(filename)

	// 向上查找 go.mod 文件
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("未找到 go.mod 文件")
		}
		dir = parent
	}
}

// writeConfigYAML 将配置结构体转换为带注释的 YAML 格式
// 通过反射读取 koanf 和 comment tag 自动生成 YAML
func writeConfigYAML(buf *bytes.Buffer, cfg Config) {
	// 写入文件头注释
	buf.WriteString(`# 配置示例文件
# 复制此文件为 config.yaml 并根据需要修改
#
# 此文件与 internal/config/config.go 中的 DefaultConfig() 保持同步
# 所有配置项都可以通过环境变量覆盖 (默认的环境变量前缀：APP_)
# 例如：APP_SERVER_ADDR=0.0.0.0:8080 会覆盖 server.addr 的值
`)

	// 通过反射遍历 Config 结构体的字段
	writeStructYAML(buf, reflect.ValueOf(cfg), reflect.TypeOf(cfg), 0)
}

// writeStructYAML 递归写入结构体的 YAML 格式
func writeStructYAML(buf *bytes.Buffer, val reflect.Value, typ reflect.Type, indent int) {
	prefix := strings.Repeat("  ", indent)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		koanfKey := field.Tag.Get("koanf")
		comment := field.Tag.Get("comment")
		if koanfKey == "" {
			continue
		}

		// 处理嵌套结构体
		if field.Type.Kind() == reflect.Struct && field.Type.String() != "time.Duration" && field.Type.String() != "time.Time" {
			fmt.Fprintf(buf, "\n%s# %s\n", prefix, comment)
			fmt.Fprintf(buf, "%s%s:\n", prefix, koanfKey)
			writeStructYAML(buf, fieldVal, field.Type, indent+1)
			continue
		}

		// 根据字段类型输出不同格式
		switch fieldVal.Kind() {
		case reflect.String:
			fmt.Fprintf(buf, "%s%s: %q # %s\n", prefix, koanfKey, fieldVal.String(), comment)
		case reflect.Bool:
			fmt.Fprintf(buf, "%s%s: %t # %s\n", prefix, koanfKey, fieldVal.Bool(), comment)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// 特殊处理 time.Duration
			if field.Type.String() == "time.Duration" {
				fmt.Fprintf(buf, "%s%s: %s # %s\n", prefix, koanfKey, fieldVal.Interface(), comment)
			} else {
				fmt.Fprintf(buf, "%s%s: %d # %s\n", prefix, koanfKey, fieldVal.Int(), comment)
			}
		case reflect.Slice:
			if fieldVal.Len() == 0 {
				fmt.Fprintf(buf, "%s%s: [] # %s\n", prefix, koanfKey, comment)
			} else {
				fmt.Fprintf(buf, "%s%s: # %s\n", prefix, koanfKey, comment)
				for j := 0; j < fieldVal.Len(); j++ {
					fmt.Fprintf(buf, "%s  - %v\n", prefix, fieldVal.Index(j).Interface())
				}
			}
		case reflect.Map:
			if fieldVal.Len() == 0 {
				fmt.Fprintf(buf, "%s%s: {} # %s\n", prefix, koanfKey, comment)
			} else {
				fmt.Fprintf(buf, "%s%s: # %s\n", prefix, koanfKey, comment)
				iter := fieldVal.MapRange()
				for iter.Next() {
					fmt.Fprintf(buf, "%s  %v: %v\n", prefix, iter.Key().Interface(), iter.Value().Interface())
				}
			}
		default:
			fmt.Fprintf(buf, "%s%s: %v # %s\n", prefix, koanfKey, fieldVal.Interface(), comment)
		}
	}
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

// TestWriteConfigYAML 测试 YAML 生成功能
func TestWriteConfigYAML(t *testing.T) {
	cfg := DefaultConfig()

	var buf bytes.Buffer
	writeConfigYAML(&buf, cfg)

	content := buf.String()

	// 验证生成的 YAML 包含关键内容
	assert.Contains(t, content, "server:")
	assert.Contains(t, content, "client:")
	assert.Contains(t, content, "addr:")
	assert.Contains(t, content, "url:")
	assert.Contains(t, content, "# 服务端配置")
	assert.Contains(t, content, "# 客户端配置")
}

// TestContainsKey 测试 containsKey 函数
func TestContainsKey(t *testing.T) {
	keys := []string{"server.addr", "client.url", "debug"}

	assert.True(t, containsKey(keys, "server.addr"))
	assert.True(t, containsKey(keys, "client.url"))
	assert.True(t, containsKey(keys, "debug"))
	assert.False(t, containsKey(keys, "nonexistent"))
	assert.False(t, containsKey(keys, ""))
}

// TestFindProjectRoot 测试 findProjectRoot 函数
func TestFindProjectRoot(t *testing.T) {
	root, err := findProjectRoot()
	require.NoError(t, err)

	// 验证找到的目录包含 go.mod
	goModPath := filepath.Join(root, "go.mod")
	_, err = os.Stat(goModPath)
	assert.NoError(t, err)
}
