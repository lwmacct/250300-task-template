import { defineConfig } from "vitepress";
import nav from "./config/nav.json";

import example_a from "./config/example-a.json";
import example_b from "./config/example-b.json";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "My Awesome Project",
  description: "A VitePress Site",
  base: process.env.BASE || "/docs",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: nav,
    sidebar: [...example_a, ...example_b],

    socialLinks: [{ icon: "github", link: "https://github.com/vuejs/vitepress" }],
  },
});
