/**
 * Vite 构建优化配置
 * @see https://vitejs.dev/config/build-options.html
 *
 * VitePress 不使用独立的 vite.config.ts，而是通过 .vitepress/config.ts 的 vite 选项配置
 * 这里将 Vite 配置提取为独立模块以保持代码整洁
 */
import type { UserConfig } from "vite";

const viteConfig: UserConfig = {
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
};

export default viteConfig;
