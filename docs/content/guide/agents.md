# AI AGENTS.md

- 本文件为 AI Agent 在此仓库中工作时提供指导。

## 🔑 基本原则

- 在执行任何动作之前, 请确保你已经理解了任务要求并制定了清晰的计划

## 🔨 开发环境

- 系统环境

  - 当前系统环境为 `ubuntu 22.04` 沙盒模式, 你可以使用 apt 安装任意软件包来完成工作
  - 你可以使用常用工具如 `rg fd perl tree psql redis-cli` 来辅助完成任务

- Python 环境

  - Python 相关指令必须通过 `uv` 执行, 这是正确使用 Python 环境的唯一方式
  - Python 脚本首行必须使用 `#!/usr/bin/env -S uv run python` 来确保使用正确的环境
  - 如果需要安装包, 请使用 `uv add <package-name>`
  - 如果需要运行 Python 包命令或文件, 请使用 `uv run <command>`, 例如:
    - `uv run pip list`
    - `uv run apps/<app_name>/main.py`

## 📚 项目文档

- 本项目使用 VitePress 进行文档编写和展示

  - 内容目录: `docs/content/`
  - 配置文件: `docs/.vitepress/{config.ts,config/}`
  - 文档索引: `docs/.vitepress/config/sidebar.*.json`
  - 在任务完成后需要更同步更新 VitePress 文档

- 信息来源
  - 如果需要了解更多关于 项目结构/使用方法/开发规范 等信息, 请优先查阅本项目的文档
  - 对于不了解的 框架/库, 优先调用 MCP 工具 Context7 来获取最新的帮助文档, 其次才是 WebSearch

## 📝 Git 提交约定

- 本项目使用 [pre-commit](https://pre-commit.com/) 框架 (已安装在环境中)
  - 在完成每一个任务后进行 `git commit` 来提交工作报告, 如果 pre-commit 检查失败, 请继续修改直到通过
  - 如果你需要跳过检查 (与当前任务不相关的错误)，可以使用 `git commit --no-verify`
  - 环境中可能有多个 AI Agent 在工作，`git commit` 时不必在意其他被修改的文件
