// Package config 提供应用配置管理。
//
// 配置加载优先级 (从低到高)：
//  1. 默认值 - DefaultConfig() 函数中定义
//  2. 配置文件 - 按 configPaths 顺序搜索
//  3. CLI flags - 最高优先级
package config

import (
	"time"

	"github.com/lwmacct/251207-go-pkg-config/pkg/config"
	"github.com/urfave/cli/v3"
)

// Config 应用配置
type Config struct {
	Server ServerConfig `koanf:"server" comment:"服务端配置"`
	Client ClientConfig `koanf:"client" comment:"客户端配置"`
}

// ServerConfig 服务端配置
type ServerConfig struct {
	Addr     string        `koanf:"addr" comment:"服务器监听地址"`
	Docs     string        `koanf:"docs" comment:"VitePress 文档目录路径"`
	Timeout  time.Duration `koanf:"timeout" comment:"HTTP 读写超时"`
	Idletime time.Duration `koanf:"idletime" comment:"HTTP 空闲超时"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
	URL     string        `koanf:"url" comment:"服务器地址"`
	Timeout time.Duration `koanf:"timeout" comment:"请求超时时间"`
	Retries int           `koanf:"retries" comment:"重试次数"`
}

// DefaultConfig 返回默认配置
// 注意：internal/command/command.go 中的 Defaults 变量引用此函数以实现单一配置来源。
func DefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Addr:     ":8080",
			Docs:     "docs/.vitepress/dist",
			Timeout:  15 * time.Second,
			Idletime: 60 * time.Second,
		},
		Client: ClientConfig{
			URL:     "http://localhost:8080",
			Timeout: 30 * time.Second,
			Retries: 3,
		},
	}
}

// Load 加载配置，委托给 pkg/config.Load 泛型函数
// configPaths 为可选的配置文件搜索路径
func Load(cmd *cli.Command, configPaths []string) (*Config, error) {
	return config.Load(cmd, configPaths, DefaultConfig())
}
