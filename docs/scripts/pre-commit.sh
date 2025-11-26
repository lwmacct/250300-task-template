#!/usr/bin/env bash

# 用于在文档文件变更时构建 VitePress 文档的 pre-commit 钩子入口。
# 此脚本可以直接调用或通过 pre-commit 框架调用。

set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel 2>/dev/null || (cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd))"
cd "$ROOT_DIR"

# 仅在 docs/ 有已暂存改动时运行（不受 --all-files 影响）
DOCS_CHANGED_FILES="$(git diff --cached --name-only -- 'docs/**' 'docs/*' || true)"

# 静默跳过：无文档变更
[ -z "$DOCS_CHANGED_FILES" ] && exit 0

# 静默跳过：docs 目录或配置文件不存在
[ ! -d "$ROOT_DIR/docs" ] && exit 0
[ ! -f "$ROOT_DIR/docs/package.json" ] && exit 0

# 安装依赖（仅在缺失时）
if [ ! -d "$ROOT_DIR/docs/node_modules" ]; then
  echo "📦 安装文档依赖..."
  (cd "$ROOT_DIR/docs" && npm install --silent)
fi

# 禁用 npm 进度条和交互式输出，避免在 pre-commit 环境中卡住
export CI=true
export npm_config_progress=false
export npm_config_loglevel=silent
# 提升 Node 可用内存，避免在 CI/预提交阶段 OOM (337/137)
: "${NODE_OPTIONS:=--max-old-space-size=4096}"
export NODE_OPTIONS

# 运行构建（捕获输出，仅在失败时显示）
BUILD_OUTPUT=$(cd "$ROOT_DIR/docs" && npm run build --silent 2>&1)
BUILD_EXIT_CODE=$?

if [ $BUILD_EXIT_CODE -eq 0 ]; then
  echo "✅ 文档构建成功"
  exit 0
else
  echo "❌ 文档构建失败"
  echo "$BUILD_OUTPUT"
  exit $BUILD_EXIT_CODE
fi
