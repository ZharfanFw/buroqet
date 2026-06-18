import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
  ],
  server: {
    proxy: {
      '/api/pricing': {
        target: 'http://localhost:8080', 
        changeOrigin: true,
        // 👇 BARIS INI YANG PALING PENTING:
        rewrite: (path) => path.replace(/^\/api\/pricing/, '/pricing')
      },
    }
  }
})