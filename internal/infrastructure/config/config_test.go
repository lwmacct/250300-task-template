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
)

// TestGenerateExample 生成 config/config.example.yaml 配置示例文件
//
// 此测试会从 defaultConfig() 提取默认配置值，并生成符合规范的 YAML 配置示例文件。
// 设计为可被 pre-commit hook 调用，在 config.go 变更时自动执行。
//
// 运行方式:
//
//	go test -v -run TestGenerateExample ./internal/infrastructure/config/...
func TestGenerateExample(t *testing.T) {
	// 获取项目根目录
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("无法找到项目根目录: %v", err)
	}

	// 获取默认配置
	cfg := defaultConfig()

	// 生成 YAML 内容
	var buf bytes.Buffer
	writeConfigYAML(&buf, cfg)

	// 写入文件
	outputPath := filepath.Join(projectRoot, "config", "config.example.yaml")
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
	cfg := defaultConfig()
	cfgType := reflect.TypeOf(cfg)

	// 验证所有字段都有 koanf 标签
	checkKoanfTags(t, cfgType, "Config")
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
	buf.WriteString(`# 示例配置文件
# 复制此文件为 config.yaml 并根据需要修改
#
# 此文件与 internal/infrastructure/config/config.go 中的 defaultConfig() 保持同步
# 所有配置项都可以通过环境变量覆盖 (环境变量前缀：APP_)
# 例如：APP_SERVER_ADDR=:8080 会覆盖 server.addr 的值

`)

	// 通过反射遍历 Config 结构体的字段
	cfgVal := reflect.ValueOf(cfg)
	cfgType := cfgVal.Type()

	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		fieldVal := cfgVal.Field(i)

		koanfKey := field.Tag.Get("koanf")
		sectionComment := field.Tag.Get("comment")

		// 写入配置段注释和名称
		fmt.Fprintf(buf, "# %s\n", sectionComment)
		fmt.Fprintf(buf, "%s:\n", koanfKey)

		// 写入该配置段下的所有字段
		writeStructFields(buf, fieldVal)

		// 配置段之间添加空行（最后一个除外）
		if i < cfgType.NumField()-1 {
			buf.WriteString("\n")
		}
	}
}

// writeStructFields 通过反射写入结构体的所有字段
func writeStructFields(buf *bytes.Buffer, structVal reflect.Value) {
	structType := structVal.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldVal := structVal.Field(i)

		koanfKey := field.Tag.Get("koanf")
		comment := field.Tag.Get("comment")

		// 根据字段类型输出不同格式
		switch fieldVal.Kind() {
		case reflect.String:
			fmt.Fprintf(buf, "  %s: %q # %s\n", koanfKey, fieldVal.String(), comment)
		case reflect.Bool:
			fmt.Fprintf(buf, "  %s: %t # %s\n", koanfKey, fieldVal.Bool(), comment)
		case reflect.Int64:
			// 处理 time.Duration 类型
			if field.Type == reflect.TypeOf(time.Duration(0)) {
				duration := time.Duration(fieldVal.Int())
				fmt.Fprintf(buf, "  %s: %q # %s\n", koanfKey, formatDuration(duration), comment)
			} else {
				fmt.Fprintf(buf, "  %s: %d # %s\n", koanfKey, fieldVal.Int(), comment)
			}
		default:
			// 其他类型使用默认格式
			fmt.Fprintf(buf, "  %s: %v # %s\n", koanfKey, fieldVal.Interface(), comment)
		}
	}
}

// formatDuration 将 time.Duration 转换为人类可读的字符串格式
func formatDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	// 尝试简化为最大单位
	if d%(24*time.Hour) == 0 {
		hours := d / time.Hour
		return fmt.Sprintf("%dh", hours)
	}
	if d%time.Hour == 0 {
		return fmt.Sprintf("%dh", d/time.Hour)
	}
	if d%time.Minute == 0 {
		return fmt.Sprintf("%dm", d/time.Minute)
	}
	if d%time.Second == 0 {
		return fmt.Sprintf("%ds", d/time.Second)
	}

	return d.String()
}

// checkKoanfTags 递归检查结构体字段的 koanf 标签
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

		// 递归检查嵌套结构体
		if field.Type.Kind() == reflect.Struct && !strings.HasPrefix(field.Type.String(), "time.") {
			checkKoanfTags(t, field.Type, fieldPath)
		}
	}
}
