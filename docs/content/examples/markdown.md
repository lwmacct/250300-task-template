# Markdown 示例

本页展示 VitePress 支持的 Markdown 扩展语法。

## 代码高亮

VitePress 使用 [Shiki](https://github.com/shikijs/shiki) 实现代码高亮，支持行号标记。

**输入**

````md
```js{4}
export default {
  data () {
    return {
      msg: 'Highlighted!'
    }
  }
}
```
````

**输出**

```js{4}
export default {
  data () {
    return {
      msg: 'Highlighted!'
    }
  }
}
```

## 自定义容器

**输入**

```md
::: info
这是一条信息
:::

::: tip
这是一条提示
:::

::: warning
这是一条警告
:::

::: danger
这是一条危险警告
:::

::: details
这是一个详情块
:::
```

**输出**

::: info
这是一条信息
:::

::: tip
这是一条提示
:::

::: warning
这是一条警告
:::

::: danger
这是一条危险警告
:::

::: details
这是一个详情块
:::

## 代码组

**输入**

````md
::: code-group

```js [config.js]
/**
 * @type {import('vitepress').UserConfig}
 */
const config = {
  // ...
}

export default config
```

```ts [config.ts]
import type { UserConfig } from 'vitepress'

const config: UserConfig = {
  // ...
}

export default config
```

:::
````

**输出**

::: code-group

```js [config.js]
/**
 * @type {import('vitepress').UserConfig}
 */
const config = {
  // ...
}

export default config
```

```ts [config.ts]
import type { UserConfig } from 'vitepress'

const config: UserConfig = {
  // ...
}

export default config
```

:::

## 更多功能

| 功能 | 支持 |
| ---- | ---- |
| 表格 | ✅ |
| 链接 | ✅ |
| Emoji | ✅ |
| TOC | ✅ |

更多语法请参考 [VitePress Markdown 文档](https://vitepress.dev/guide/markdown)
