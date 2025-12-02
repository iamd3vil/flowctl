import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

const backendHost = process.env.VITE_BACKEND_HOST || "localhost:7000";

export default defineConfig({
  plugins: [tailwindcss(), sveltekit()],
  cacheDir: "node_modules/.vite",
  build: {
    cssCodeSplit: true,
    minify: 'esbuild',
  },
  optimizeDeps: {
    force: false,
    include: ['@tabler/icons-svelte']
  },
  server: {
    proxy: {
      "/api": {
        target: `http://${backendHost}`,
        changeOrigin: true,
        secure: false,
      },
      "/sso-providers": {
        target: `http://${backendHost}`,
        changeOrigin: true,
        secure: false,
      },
      "/login": {
        target: `http://${backendHost}`,
        changeOrigin: true,
        secure: false,
      },
      "/logout": {
        target: `http://${backendHost}`,
        changeOrigin: true,
        secure: false,
      },
      "/auth": {
        target: `http://${backendHost}`,
        changeOrigin: true,
        secure: false,
      }
    },
  },
});
