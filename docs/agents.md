- 本文件为 AI Agent 在此仓库中工作时提供指导。

## 🔑 基本原则

- 在执行任何动作之前, 请确保你已经理解了任务要求并制定了清晰的计划
- 如果对于 框架/库 不了解, 优先调用 MCP 工具 Context7 来获取使用方法和帮助文档, 其次才是 WebSearch

## 🔨 开发环境

- 所有 Python 相关指令必须通过 `uv` 执行, 这是正确使用 Python 环境的唯一方式
  - 如果需要安装包, 请使用 `uv add <package-name>`
  - 如果需要运行 Python 包命令或脚本, 请使用 `uv run <command-name>`, 例如: `uv run pip list`, `uv run apps/<app_name>/main.py`

## 📚 项目文档

- 本项目使用 VitePress 进行文档编写和展示, 文档源码位于 `docs/` 目录下
- 文件 `docs/.vitepress/config.ts` 包含所有的配置
- 在任务完成后需要更同步更新 VitePress 文档

## 📝 Git 提交约定

- 本项目使用 [pre-commit](https://pre-commit.com/) 框架 (已安装在环境中)
- 在完成每一个任务后进行 git commit 来提交工作报告, 如果 pre-commit 检查失败, 请继续修改直到通过
- 环境中可能有多个 AI Agent 在工作，git commit 时不必在意其他被修改的文件
- 如果你需要跳过检查 (与当前任务不相关的错误)，可以使用 `git commit --no-verify`
