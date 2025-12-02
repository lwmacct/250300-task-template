// Package config 提供应用配置管理
//
// 配置加载优先级 (从低到高) ：
// 1. 默认值 - defaultConfig() 函数中定义
// 2. 配置文件 - config.yaml 或 config/config.yaml
// 3. 环境变量 - 以 APP_ 为前缀 (如 APP_SERVER_ADDR)
//
// 重要提示：
// - 如果修改了 defaultConfig() 中的默认值，请同步更新 config/config.example.yaml 示例文件
// - 配置文件路径硬编码在 Load() 函数中：[]string{"config.yaml", "config/config.yaml"}
package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Addr     string `koanf:"addr" comment:"监听地址，格式: host:port，例如 '0.0.0.0:8080' 或 ':8080'"`
	Env      string `koanf:"env" comment:"运行环境: development | production"`
	DistWeb  string `koanf:"dist_web" comment:"静态资源目录路径，用于提供前端文件服务 (如 SPA 应用)"`
	DistDocs string `koanf:"dist_docs" comment:"文档目录路径，用于提供 VitePress 构建的文档服务，通过 /docs 路由访问"`
}

// DataConfig 数据源配置
type DataConfig struct {
	PgsqlURL  string `koanf:"pgsql_url" comment:"PostgreSQL 连接 URL，格式: postgresql://user:password@host:port/dbname?sslmode=disable"`
	RedisURL  string `koanf:"redis_url" comment:"Redis 连接 URL，格式: redis://:password@host:port/db"`
	KeyPrefix string `koanf:"key_prefix" comment:"Redis key 前缀，所有 key 读写都会自动拼接此前缀，例如 'app:'"`
}

// Config 应用配置
type Config struct {
	Server ServerConfig `koanf:"server" comment:"服务器配置"`
	Data   DataConfig   `koanf:"data" comment:"数据源配置"`
}

// defaultConfig 返回默认配置
func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Addr:     "0.0.0.0:8080",
			Env:      "development",
			DistWeb:  "web/dist",
			DistDocs: "docs/.vitepress/dist",
		},
		Data: DataConfig{
			PgsqlURL:  "postgresql://postgres@localhost:5432/app?sslmode=disable",
			RedisURL:  "redis://localhost:6379/0",
			KeyPrefix: "app:",
		},
	}
}

// Load 加载配置，按优先级合并：
// 1. 默认值 (最低优先级)
// 2. 配置文件 (config.yaml)
// 3. 环境变量 (APP_*，最高优先级)
func Load() (*Config, error) {
	k := koanf.New(".")

	// 1️⃣ 加载默认配置 (最低优先级)
	if err := k.Load(structs.Provider(defaultConfig(), "koanf"), nil); err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}

	// 2️⃣ 加载 YAML 配置文件 (可选，文件不存在不报错)
	// 按优先级搜索：当前目录 -> config 目录
	configPaths := []string{"config.yaml", "config/config.yaml"}
	configLoaded := false
	for _, path := range configPaths {
		if err := k.Load(file.Provider(path), yaml.Parser()); err == nil {
			slog.Info("Loaded config from file", "path", path)
			configLoaded = true
			break
		}
	}
	if !configLoaded {
		slog.Info("No config file found, using defaults and env vars")
	}

	// 3️⃣ 加载环境变量 (APP_SERVER_ADDR → server.addr)
	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: "APP_",
		TransformFunc: func(key, value string) (string, any) {
			// APP_SERVER_ADDR → server.addr
			// APP_DATA_PGSQL_URL → data.pgsql.url
			key = strings.ReplaceAll(
				strings.ToLower(strings.TrimPrefix(key, "APP_")),
				"_", ".",
			)
			return key, value
		},
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	// 解析到结构体
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
