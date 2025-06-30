import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import { resolve } from "path"

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": resolve(__dirname, "src"),
      "@/components": resolve(__dirname, "src/components"),
      "@/views": resolve(__dirname, "src/views"),
      "@/stores": resolve(__dirname, "src/stores"),
      "@/utils": resolve(__dirname, "src/utils"),
      "@/types": resolve(__dirname, "src/types"),
      "@/api": resolve(__dirname, "src/api"),
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: "dist",
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ["vue", "vue-router", "pinia"],
          ui: ["radix-vue", "lucide-vue-next"],
          charts: ["echarts", "vue-echarts"],
          maps: ["leaflet", "@vue-leaflet/vue-leaflet"],
        },
      },
    },
  },
  optimizeDeps: {
    include: ["vue", "vue-router", "pinia", "@vueuse/core"],
  },
})
