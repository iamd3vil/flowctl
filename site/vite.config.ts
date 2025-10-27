import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [tailwindcss(), sveltekit()],
  cacheDir: "node_modules/.vite",
  optimizeDeps: {
    force: false,
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:7000",
        changeOrigin: true,
        secure: false,
      },
      "/login": {
        target: "http://localhost:7000",
        changeOrigin: true,
        secure: false,
      },
      "/logout": {
        target: "http://localhost:7000",
        changeOrigin: true,
        secure: false,
      },
    },
  },
});
