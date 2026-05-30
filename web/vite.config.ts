import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';

const backendPort = parseInt(process.env.RAG_EVAL_BACKEND_PORT || '8772', 10);

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: `http://127.0.0.1:${backendPort}`,
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: '../internal/web/dist',
    emptyOutDir: true,
    sourcemap: true,
  },
});
