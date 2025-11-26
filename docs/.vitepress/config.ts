import { defineConfig } from "vitepress";
import nav from "./config/nav.json";
import sidebarGuide from "./config/sidebar.guide.json";
import sidebarIssues from "./config/sidebar.issues.json";
// 模板示例侧边栏 - 新项目可删除此行及 content/examples/ 目录
import sidebarExamples from "./config/sidebar.examples.json";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "My Awesome Project",
  description: "A VitePress Site",
  base: process.env.BASE || "/docs",
  srcDir: "content",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav,
    sidebar: [...sidebarGuide, ...sidebarIssues, ...sidebarExamples],

    // 本地搜索 - 使用 MiniSearch 实现浏览器内索引
    search: {
      provider: "local",
      options: {
        locales: {
          root: {
            translations: {
              button: {
                buttonText: "搜索",
                buttonAriaLabel: "搜索文档",
              },
              modal: {
                displayDetails: "显示详细列表",
                resetButtonTitle: "重置搜索",
                backButtonTitle: "关闭搜索",
                noResultsText: "没有找到结果",
                footer: {
                  selectText: "选择",
                  selectKeyAriaLabel: "输入",
                  navigateText: "导航",
                  navigateUpKeyAriaLabel: "上箭头",
                  navigateDownKeyAriaLabel: "下箭头",
                  closeText: "关闭",
                  closeKeyAriaLabel: "esc",
                },
              },
            },
          },
        },
      },
    },

    socialLinks: [{ icon: "github", link: "https://github.com/vuejs/vitepress" }],
  },
});
