import DefaultTheme from "vitepress/theme";
import { onMounted, watch, nextTick } from "vue";
import { useRoute } from "vitepress";
import mermaid from "mermaid";
import type { Theme } from "vitepress";

import "./mermaid.css";

// Mermaid 初始化配置
mermaid.initialize({
  startOnLoad: false,
  securityLevel: "loose",
  theme: "default",
});

const theme: Theme = {
  extends: DefaultTheme,
  setup() {
    const route = useRoute();

    const initMermaid = async () => {
      await nextTick();
      await mermaid.run({
        querySelector: ".mermaid",
      });
    };

    onMounted(() => {
      initMermaid();
    });

    watch(
      () => route.path,
      () => {
        initMermaid();
      }
    );
  },
};

export default theme;
