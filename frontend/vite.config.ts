import { defineConfig, loadEnv } from "vite";
import vue from "@vitejs/plugin-vue";
import path from "node:path";
import postcssPresetEnv from "postcss-preset-env";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");

  return {
    base: env.VITE_PUBLIC_PATH || "/",
    plugins: [vue()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
        "@assets": path.resolve(__dirname, "./src/assets")
      }
    },
    server: {
      open: false,
      proxy: {
        "/api": {
          target: env.VITE_APP_BASE_URL,
          changeOrigin: true
        },
        "/uploads": {
          target: env.VITE_APP_BASE_URL,
          changeOrigin: true
        }
      }
    },
    css: {
      postcss: {
        plugins: [postcssPresetEnv()]
      },
      preprocessorOptions: {
        scss: {
          additionalData: `@use "@/style/var/index.scss" as *;`
        }
      }
    },
    build: {
      outDir: "dist",
      minify: "terser",
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (!id.includes("node_modules")) {
              return;
            }
            if (id.includes("@arco-design/web-vue")) {
              return "arco";
            }
            if (id.includes("vue-router")) {
              return "vue-router";
            }
            if (id.includes("pinia")) {
              return "pinia";
            }
            return "vendor";
          }
        }
      },
      terserOptions: {
        compress: {
          drop_console: true,
          drop_debugger: true
        }
      }
    }
  };
});
