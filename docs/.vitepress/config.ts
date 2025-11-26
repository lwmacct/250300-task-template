import { defineConfig } from "vitepress";
import nav from "./config/nav.json";
import sidebar from "./config/sidebar.json";
// 模板示例侧边栏 - 新项目可删除此行及 content/examples/ 目录
import sidebarExamples from "./config/sidebar.examples.json";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "My Awesome Project",
  description: "A VitePress Site",
  base: process.env.BASE || "/docs",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav,
    sidebar: [...sidebar, ...sidebarExamples],

    socialLinks: [{ icon: "github", link: "https://github.com/vuejs/vitepress" }],
  },
});
