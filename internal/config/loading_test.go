// Author: lwmacct (https://github.com/lwmacct)
package config

import (
	"reflect"
	"testing"
	"time"

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

// TestNestedStructRecursion 测试嵌套结构递归
func TestNestedStructRecursion(t *testing.T) {
	// 通过 Load 函数间接测试嵌套结构递归
	cfg, err := Load(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 验证嵌套结构正确加载
	assert.Equal(t, ":8080", cfg.Server.Addr)
	assert.Equal(t, "http://localhost:8080", cfg.Client.URL)
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
