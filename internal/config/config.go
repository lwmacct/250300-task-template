// Package config 提供应用配置管理。
//
// 配置加载优先级 (从低到高)：
//   1. 默认值 - DefaultConfig() 函数中定义
//   2. 配置文件 - 通过 --config 指定，或按顺序搜索默认路径
//   3. 环境变量 - 以 <AppRawName> 为前缀，下划线分隔嵌套路径
//   4. CLI flags - 最高优先级
package config

import "time"

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
// 注意：这里的默认值应对齐 internal/command/*/command.go 中的默认值, 确保生成的配置文件示例与 CLI 默认值一致
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
