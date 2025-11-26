/**
 * Markdown 渲染配置
 * @see https://vitepress.dev/reference/site-config#markdown
 *
 * 基于 markdown-it 的扩展配置，支持自定义渲染规则
 */
import type { MarkdownOptions } from "vitepress";

const markdownConfig: MarkdownOptions = {
  // 扩展 markdown-it 实例
  config: (md) => {
    // Mermaid 代码块转换 - 将 ```mermaid 转换为 <pre class="mermaid">
    const fence = md.renderer.rules.fence!;
    md.renderer.rules.fence = (...args) => {
      const [tokens, idx] = args;
      const token = tokens[idx];
      if (token.info.trim() === "mermaid") {
        return `<pre class="mermaid">${token.content}</pre>`;
      }
      return fence(...args);
    };
  },
};

export default markdownConfig;
