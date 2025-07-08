import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  build: {
    minify: false,
    sourcemap: true
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:9528',
        changeOrigin: true,
        rewrite: path => path.replace(/^\/api/, ''),
      },
    },
  },
}); 