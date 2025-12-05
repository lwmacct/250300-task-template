// Package version 提供应用程序的版本信息管理。
//
// 版本信息优先通过 go build -ldflags 在构建时注入，
// 若未注入则自动从 runtime/debug.BuildInfo 读取（Go 1.18+）。
// Author: lwmacct (https://github.com/lwmacct)
package version

import (
	"fmt"
	"path"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// 构建时注入的版本信息变量。
// 通过 go build -ldflags "-X package.Variable=value" 注入。
var (
	AppRawName string = "Unknown" // 应用原始名称
	AppProject string = "Unknown" // 项目名称（通常为 Git 仓库名）
	AppVersion string = "Unknown" // 应用版本号（语义化版本）
	GitCommit  string = "Unknown" // Git 提交哈希
	BuildTime  string = "Unknown" // 构建时间
	Developer  string = "Unknown" // 开发者/维护者
	Workspace  string = "Unknown" // 构建时工作目录，用于去除堆栈中的绝对路径
)

func init() {
	initFromBuildInfo()
}

// initFromBuildInfo 从 runtime/debug.BuildInfo 读取版本信息作为后备。
// 当 ldflags 未注入时，自动从 Go 模块信息和 VCS 设置中提取。
func initFromBuildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	// 从模块路径提取项目名称
	if AppProject == "Unknown" && info.Main.Path != "" {
		AppProject = path.Base(info.Main.Path)
	}

	// 从模块版本提取应用版本
	if AppVersion == "Unknown" && info.Main.Version != "" && info.Main.Version != "(devel)" {
		AppVersion = info.Main.Version
	}

	// 从 VCS 设置中提取 Git 信息
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			if GitCommit == "Unknown" && setting.Value != "" {
				// 使用短 hash（7位），与 git log -n 1 --format=%h 一致
				if len(setting.Value) > 7 {
					GitCommit = setting.Value[:7]
				} else {
					GitCommit = setting.Value
				}
			}
		case "vcs.time":
			if BuildTime == "Unknown" && setting.Value != "" {
				BuildTime = formatBuildTime(setting.Value)
			}
		case "vcs.modified":
			// 如果有未提交的修改，标记为 dirty
			if setting.Value == "true" && GitCommit != "Unknown" && !strings.HasSuffix(GitCommit, "-dirty") {
				GitCommit = GitCommit + "-dirty"
			}
		}
	}
}

// formatBuildTime 将 VCS 时间（RFC3339 格式）转换为 UTC+8 时区格式。
func formatBuildTime(vcsTime string) string {
	t, err := time.Parse(time.RFC3339, vcsTime)
	if err != nil {
		return vcsTime // 解析失败则返回原始值
	}
	// 使用固定偏移量 UTC+8，避免依赖系统时区数据库
	cst := time.FixedZone("CST", 8*60*60)
	return t.In(cst).Format("2006-01-02 15:04:05 MST")
}

// PrintVersion 打印版本信息
func PrintVersion() {
	fmt.Printf("AppRawName:   %s\n", AppRawName)
	fmt.Printf("AppVersion:   %s\n", AppVersion)
	fmt.Printf("Go Version:   %s\n", runtime.Version())
	fmt.Printf("Git Commit:   %s\n", GitCommit)
	fmt.Printf("Build Time:   %s\n", BuildTime)
	fmt.Printf("AppProject:   %s\n", AppProject)
	fmt.Printf("Developer :   %s\n", Developer)
}

// GetVersion 返回应用版本号
func GetVersion() string {
	if AppVersion == "Unknown" && GitCommit != "Unknown" && len(GitCommit) > 7 {
		return fmt.Sprintf("dev-%s", GitCommit[:7])
	}
	return AppVersion
}

// GetShortVersion 返回简短版本号 (兼容性函数)
func GetShortVersion() string {
	return AppVersion
}

// GetBuildInfo 返回构建相关信息 (用于健康检查等)
func GetBuildInfo() string {
	return fmt.Sprintf("版本: %s, 提交: %s, 构建时间: %s", AppVersion, GitCommit, BuildTime)
}
