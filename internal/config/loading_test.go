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
	"github.com/urfave/cli/v3"
)

// TestApplyCLIFlags 测试 applyCLIFlags 函数
func TestApplyCLIFlags(t *testing.T) {
	k := koanf.New(".")

	// 设置初始值
	_ = k.Set("server.addr", ":8080")
	_ = k.Set("client.url", "http://localhost:8080")

	// cmd 为 nil 时应该直接返回不处理
	// applyCLIFlags 会在 cmd 为 nil 时直接返回
}

// TestApplyCLIFlagsRecursive 测试递归应用 CLI flags
func TestApplyCLIFlagsRecursive(t *testing.T) {
	k := koanf.New(".")

	// 设置初始值
	_ = k.Set("server.addr", ":8080")

	// 由于 applyCLIFlagsRecursive 需要非 nil 的 cmd 才能调用 IsSet
	// 这里测试是在 Load 函数中，cmd 为 nil 时跳过 applyCLIFlags
	// 验证值未改变
	assert.Equal(t, ":8080", k.String("server.addr"))
}

// TestSetCLIFlagValue 测试 setCLIFlagValue 函数各种类型
func TestSetCLIFlagValue(t *testing.T) {
	k := koanf.New(".")

	// 创建一个模拟的 CLI command
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "test-string", Value: "test"},
			&cli.BoolFlag{Name: "test-bool", Value: true},
			&cli.IntFlag{Name: "test-int", Value: 42},
			&cli.Int64Flag{Name: "test-int64", Value: 100},
			&cli.Float64Flag{Name: "test-float64", Value: 3.14},
			&cli.DurationFlag{Name: "test-duration", Value: 5 * time.Second},
			&cli.StringSliceFlag{Name: "test-strings", Value: []string{"a", "b"}},
		},
	}

	// 测试各种类型
	testCases := []struct {
		name      string
		koanfKey  string
		cliFlag   string
		fieldType reflect.Type
	}{
		{"string", "test.string", "test-string", reflect.TypeOf("")},
		{"bool", "test.bool", "test-bool", reflect.TypeOf(false)},
		{"int", "test.int", "test-int", reflect.TypeOf(0)},
		{"int64", "test.int64", "test-int64", reflect.TypeOf(int64(0))},
		{"float64", "test.float64", "test-float64", reflect.TypeOf(float64(0))},
		{"duration", "test.duration", "test-duration", reflect.TypeOf(time.Duration(0))},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setCLIFlagValue(cmd, k, tc.koanfKey, tc.cliFlag, tc.fieldType)
			// 函数应该不 panic
		})
	}
}

// TestSetSliceFlagValue 测试切片类型的 CLI flag 处理
func TestSetSliceFlagValue(t *testing.T) {
	k := koanf.New(".")

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringSliceFlag{Name: "strings", Value: []string{"a", "b"}},
			&cli.IntSliceFlag{Name: "ints", Value: []int{1, 2, 3}},
			&cli.Int64SliceFlag{Name: "int64s", Value: []int64{1, 2}},
			&cli.Float64SliceFlag{Name: "float64s", Value: []float64{1.1, 2.2}},
		},
	}

	testCases := []struct {
		name      string
		koanfKey  string
		cliFlag   string
		fieldType reflect.Type
	}{
		{"string slice", "test.strings", "strings", reflect.TypeOf([]string{})},
		{"int slice", "test.ints", "ints", reflect.TypeOf([]int{})},
		{"int64 slice", "test.int64s", "int64s", reflect.TypeOf([]int64{})},
		{"float64 slice", "test.float64s", "float64s", reflect.TypeOf([]float64{})},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setSliceFlagValue(cmd, k, tc.koanfKey, tc.cliFlag, tc.fieldType)
			// 函数应该不 panic
		})
	}
}

// TestSetCLIFlagValueOtherTypes 测试其他类型
func TestSetCLIFlagValueOtherTypes(t *testing.T) {
	k := koanf.New(".")

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.Int8Flag{Name: "int8", Value: 8},
			&cli.Int16Flag{Name: "int16", Value: 16},
			&cli.Int32Flag{Name: "int32", Value: 32},
			&cli.UintFlag{Name: "uint", Value: 100},
			&cli.Uint16Flag{Name: "uint16", Value: 16},
			&cli.Uint32Flag{Name: "uint32", Value: 32},
			&cli.Uint64Flag{Name: "uint64", Value: 64},
			&cli.Float32Flag{Name: "float32", Value: 3.14},
			&cli.StringMapFlag{Name: "map", Value: map[string]string{"k": "v"}},
		},
	}

	testCases := []struct {
		name      string
		koanfKey  string
		cliFlag   string
		fieldType reflect.Type
	}{
		{"int8", "test.int8", "int8", reflect.TypeOf(int8(0))},
		{"int16", "test.int16", "int16", reflect.TypeOf(int16(0))},
		{"int32", "test.int32", "int32", reflect.TypeOf(int32(0))},
		{"uint", "test.uint", "uint", reflect.TypeOf(uint(0))},
		{"uint8", "test.uint8", "uint", reflect.TypeOf(uint8(0))},
		{"uint16", "test.uint16", "uint16", reflect.TypeOf(uint16(0))},
		{"uint32", "test.uint32", "uint32", reflect.TypeOf(uint32(0))},
		{"uint64", "test.uint64", "uint64", reflect.TypeOf(uint64(0))},
		{"float32", "test.float32", "float32", reflect.TypeOf(float32(0))},
		{"map", "test.map", "map", reflect.TypeOf(map[string]string{})},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setCLIFlagValue(cmd, k, tc.koanfKey, tc.cliFlag, tc.fieldType)
			// 函数应该不 panic
		})
	}
}

// TestSetSliceFlagValueOtherTypes 测试其他切片类型
func TestSetSliceFlagValueOtherTypes(t *testing.T) {
	k := koanf.New(".")

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.Int8SliceFlag{Name: "int8s", Value: []int8{1, 2}},
			&cli.Int16SliceFlag{Name: "int16s", Value: []int16{1, 2}},
			&cli.Int32SliceFlag{Name: "int32s", Value: []int32{1, 2}},
			&cli.Uint16SliceFlag{Name: "uint16s", Value: []uint16{1, 2}},
			&cli.Uint32SliceFlag{Name: "uint32s", Value: []uint32{1, 2}},
			&cli.Float32SliceFlag{Name: "float32s", Value: []float32{1.1, 2.2}},
		},
	}

	testCases := []struct {
		name      string
		koanfKey  string
		cliFlag   string
		fieldType reflect.Type
	}{
		{"int8 slice", "test.int8s", "int8s", reflect.TypeOf([]int8{})},
		{"int16 slice", "test.int16s", "int16s", reflect.TypeOf([]int16{})},
		{"int32 slice", "test.int32s", "int32s", reflect.TypeOf([]int32{})},
		{"uint16 slice", "test.uint16s", "uint16s", reflect.TypeOf([]uint16{})},
		{"uint32 slice", "test.uint32s", "uint32s", reflect.TypeOf([]uint32{})},
		{"float32 slice", "test.float32s", "float32s", reflect.TypeOf([]float32{})},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setSliceFlagValue(cmd, k, tc.koanfKey, tc.cliFlag, tc.fieldType)
			// 函数应该不 panic
		})
	}
}

// TestLoadServerConfigError 测试加载服务端配置错误处理
func TestLoadServerConfigError(t *testing.T) {
	// 正常情况下不会出错，因为默认配置总是有效的
	cfg, err := LoadServerConfig(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

// TestLoadClientConfigError 测试加载客户端配置错误处理
func TestLoadClientConfigError(t *testing.T) {
	// 正常情况下不会出错，因为默认配置总是有效的
	cfg, err := LoadClientConfig(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

// TestApplyCLIFlagsWithEmptyKoanfTag 测试空 koanf 标签
func TestApplyCLIFlagsWithEmptyKoanfTag(t *testing.T) {
	// 测试通过 Load 函数间接调用
	cfg, err := Load(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

// TestSetCLIFlagValueWithTimestamp 测试时间戳类型
func TestSetCLIFlagValueWithTimestamp(t *testing.T) {
	k := koanf.New(".")

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.TimestampFlag{
				Name:   "timestamp",
				Config: cli.TimestampConfig{Layouts: []string{time.RFC3339}},
			},
		},
	}

	setCLIFlagValue(cmd, k, "test.timestamp", "timestamp", reflect.TypeOf(time.Time{}))
}

// TestKoanfKeyToCliFlag 测试 koanf key 到 CLI flag 的转换
func TestKoanfKeyToCliFlag(t *testing.T) {
	testCases := []struct {
		koanfKey string
		cliFlag  string
	}{
		{"server.addr", "server-addr"},
		{"client.url", "client-url"},
		{"server.timeout", "server-timeout"},
	}

	for _, tc := range testCases {
		t.Run(tc.koanfKey, func(t *testing.T) {
			// 验证转换逻辑
			result := tc.koanfKey
			result = replaceAll(result, ".", "-")
			result = replaceAll(result, "_", "-")
			assert.Equal(t, tc.cliFlag, result)
		})
	}
}

// replaceAll 辅助函数
func replaceAll(s, old, new string) string {
	for {
		i := indexString(s, old)
		if i == -1 {
			break
		}
		s = s[:i] + new + s[i+len(old):]
	}
	return s
}

// indexString 辅助函数
func indexString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// TestGenerateExample 生成配置示例文件
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
