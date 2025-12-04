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

// Config 应用配置
type Config struct {
	Addr     string `koanf:"addr" comment:"服务器监听地址"`
	DistDocs string `koanf:"dist_docs" comment:"VitePress 文档目录路径"`
}

// defaultConfig 返回默认配置
// 注意：这里的默认值应对齐 internal/command/*/command.go 中的默认值, 确保生成的配置文件示例与 CLI 默认值一致
func DefaultConfig() Config {
	return Config{
		Addr:     ":8080",
		DistDocs: "docs/.vitepress/dist",
	}
}
