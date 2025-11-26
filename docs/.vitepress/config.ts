import { defineConfig, type DefaultTheme } from "vitepress";
import nav from "./config/nav.json";
// config sidebar
import sidebarGuide from "./config/sidebar.guide.json";
import sidebarIssues from "./config/sidebar.issues.json";
import sidebarExamples from "./config/sidebar.examples.json";
// config other
import cfgSearch from "./config/search.json";
import viteConfig from "./config/vite";
import markdownConfig from "./config/markdown";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "My Awesome Project",
  description: "A VitePress Site",
  base: process.env.BASE || "/docs",
  srcDir: "content",

  // Vite 构建优化配置 (从 ./config/vite.ts 导入)
  vite: viteConfig,

  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav,
    sidebar: [...sidebarGuide, ...sidebarIssues, ...sidebarExamples],

    // 本地搜索 - 使用 MiniSearch 实现浏览器内索引
    search: cfgSearch as DefaultTheme.Config["search"],

    socialLinks: [{ icon: "github", link: "https://github.com/vuejs/vitepress" }],
  },

  // Markdown 渲染配置 (从 ./config/markdown.ts 导入)
  markdown: markdownConfig,
});
