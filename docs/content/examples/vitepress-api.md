# 运行时 API 示例

本页面演示了 VitePress 运行时 API 的使用方法。

## useData

`useData()` 是 VitePress 的核心组合式 API，用于访问站点和页面数据：

```vue
<script setup>
import { useData } from 'vitepress'

const { site, page, theme, frontmatter } = useData()
</script>

<template>
  <h1>{{ page.title }}</h1>
  <p>站点标题: {{ site.title }}</p>
</template>
```

### 返回值类型

```ts
interface VitePressData {
  site: Ref<SiteData>           // 站点级别数据
  page: Ref<PageData>           // 页面级别数据
  theme: Ref<ThemeConfig>       // 主题配置
  frontmatter: Ref<PageFrontmatter>  // 页面 frontmatter
  title: Ref<string>            // 页面标题
  description: Ref<string>      // 页面描述
  lang: Ref<string>             // 当前语言
  isDark: Ref<boolean>          // 是否为暗色模式
}
```

## useRoute

`useRoute()` 返回当前路由对象：

```vue
<script setup>
import { useRoute } from 'vitepress'

const route = useRoute()
</script>

<template>
  <p>当前路径: {{ route.path }}</p>
</template>
```

## useRouter

`useRouter()` 返回 VitePress 路由器实例，用于编程式导航：

```vue
<script setup>
import { useRouter } from 'vitepress'

const router = useRouter()

function navigate() {
  router.go('/examples/markdown')
}
</script>

<template>
  <button @click="navigate">跳转到 Markdown 示例</button>
</template>
```

### 路由器方法

| 方法 | 说明 |
| ---- | ---- |
| `go(href)` | 导航到指定 URL |
| `onBeforeRouteChange` | 路由变化前的钩子 |
| `onAfterRouteChanged` | 路由变化后的钩子 |

## $frontmatter

在 Markdown 中可以直接访问 frontmatter 数据：

```md
---
title: 我的页面
description: 页面描述
custom:
  key: value
---

# {{ $frontmatter.title }}

{{ $frontmatter.description }}
```

## 更多

完整 API 文档请参考 [VitePress Runtime API](https://vitepress.dev/reference/runtime-api)。
