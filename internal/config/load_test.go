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

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

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

	// 创建有效键的映射表以实现 O(1) 查找
	// 这比每次遍历 exampleKeys (O(m*n)) 更高效，整体复杂度为 O(m+n)
	validKeyMap := make(map[string]bool, len(exampleKeys))
	for _, key := range exampleKeys {
		validKeyMap[key] = true
	}

	// 收集无效的配置键
	var invalidKeys []string
	for _, key := range configKeys {
		if !validKeyMap[key] {
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

// findProjectRoot 通过查找 go.mod 文件定位项目根目录
func findProjectRoot() (string, error) {
	// 获取当前测试文件所在目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("无法获取当前文件路径")
	}

	dir := filepath.Dir(filename)

	// 预先构建 go.mod 路径，避免重复拼接
	sep := string(filepath.Separator)
	for {
		goModPath := dir + sep + "go.mod"
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir { // 到达根目录
			return "", fmt.Errorf("未找到 go.mod 文件（从 %s 向上查找）", dir)
		}
		dir = parent
	}
}

// writeConfigYAML 将配置结构体转换为带注释的 YAML 格式
// 通过反射读取 koanf 和 comment tag 自动生成 YAML
func writeConfigYAML(buf *bytes.Buffer, cfg Config) {
	// 写入文件头注释
	buf.WriteString(`# 配置示例文件, 复制此文件为 config.yaml 并根据需要修改
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
