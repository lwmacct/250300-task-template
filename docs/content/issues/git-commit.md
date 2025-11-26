# Git 提交常见问题

## 工作目录路径问题

### 问题描述

在使用 `cd` 切换到子目录后执行 `git add`，可能会遇到路径不匹配的错误：

```bash
cd /apps/data/workspace/project/docs && git add docs/.vitepress/config.ts
# fatal: pathspec 'docs/.vitepress/config.ts' did not match any files
```

### 原因分析

当前工作目录已经是 `docs/`，但 `git add` 的路径仍然以 `docs/` 开头，导致实际查找的路径变成了 `docs/docs/.vitepress/config.ts`。

### 解决方案

**方案一：使用 `-C` 参数指定仓库路径（推荐）**

```bash
# 无需 cd，直接指定仓库路径执行命令
git -C /apps/data/workspace/project add docs/.vitepress/config.ts
git -C /apps/data/workspace/project commit -m "message"
```

> `-C <path>` 参数让 Git 在指定目录下执行命令，避免了工作目录切换带来的路径问题。

**方案二：从项目根目录执行**

```bash
# 不要 cd 到子目录，直接从根目录执行
git add docs/.vitepress/config.ts
```

**方案三：使用相对路径**

```bash
cd /apps/data/workspace/project/docs
git add .vitepress/config.ts  # 相对于当前目录
```

**方案四：使用绝对路径**

```bash
git add /apps/data/workspace/project/docs/.vitepress/config.ts
```

## pre-commit 检查失败

### 问题描述

提交时 pre-commit 钩子检查失败：

```bash
git commit -m "message"
# Build VitePress docs when docs/ changes..................................Failed
```

### 解决方案

1. **查看具体错误**：运行构建命令查看详细错误

   ```bash
   cd docs && npm run build
   ```

2. **修复问题后重新提交**

3. **跳过检查（不推荐）**：仅在确认问题与当前任务无关时使用
   ```bash
   git commit --no-verify -m "message"
   ```

## 暂存区文件冲突

### 问题描述

暂存区包含已移动/删除的文件，导致状态混乱：

```bash
# Changes to be committed:
#   new file:   docs/.vitepress/theme/nav.json
# Changes not staged for commit:
#   deleted:    docs/.vitepress/theme/nav.json
```

### 解决方案

```bash
# 重置暂存区，重新添加
git reset HEAD
git add <正确的文件路径>
```
