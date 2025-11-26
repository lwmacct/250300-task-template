import { defineConfig, type DefaultTheme } from "vitepress";
import nav from "./config/nav.json";
import sidebarGuide from "./config/sidebar.guide.json";
import sidebarIssues from "./config/sidebar.issues.json";
// 模板示例侧边栏 - 新项目可删除此行及 content/examples/ 目录
import sidebarExamples from "./config/sidebar.examples.json";
import cfgSearch from "./config/search.json";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "My Awesome Project",
  description: "A VitePress Site",
  base: process.env.BASE || "/docs",
  srcDir: "content",

  // Vite 构建优化 - 代码分割配置
  vite: {
    build: {
      chunkSizeWarningLimit: 600, // 提高警告阈值 (KB)
      rollupOptions: {
        output: {
          // 手动分块策略，优化首屏加载和缓存利用率
          // 使用函数形式避免 external 模块冲突
          manualChunks(id) {
            if (!id.includes("node_modules")) return;
            // MiniSearch 搜索引擎
            if (id.includes("minisearch")) {
              return "search";
            }
            // shiki 代码高亮 (体积较大)
            if (id.includes("shiki") || id.includes("oniguruma")) {
              return "syntax-highlight";
            }
          },
        },
      },
    },
  },

  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav,
    sidebar: [...sidebarGuide, ...sidebarIssues, ...sidebarExamples],

    // 本地搜索 - 使用 MiniSearch 实现浏览器内索引
    search: cfgSearch as DefaultTheme.Config["search"],

    socialLinks: [{ icon: "github", link: "https://github.com/vuejs/vitepress" }],
  },

  // Mermaid 代码块转换 - 将 ```mermaid 转换为 <pre class="mermaid">
  markdown: {
    config: (md) => {
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
  },
});
