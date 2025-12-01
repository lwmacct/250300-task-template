# 快速开始

## 环境要求

- [Docker](https://www.docker.com/) 或 [Podman](https://podman.io/)
- [VS Code](https://code.visualstudio.com/) + [Dev Containers 扩展](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

## 开始使用

### 1. 克隆项目

```shell
git clone <repository-url>
cd <project-name>
```

### 2. 启动开发容器

使用 VS Code 打开项目，按 `F1` 输入 `Dev Containers: Reopen in Container`，等待容器构建完成。

### 3. 初始化开发环境

```shell
pre-commit install
```

### 4. 查看可用命令

```shell
task -a
```

## 下一步

- 阅读 [项目介绍](/readme) 了解项目结构
- 查看 [AI Agent](/guide/agents) 了解 AI 辅助开发
