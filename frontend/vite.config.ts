import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
  ],
  server: {
    proxy: {
      // Warehouse Service (port 8087) — harus lebih spesifik, di-check duluan
      '/api/warehouse': {
        target: 'http://localhost:8087',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/warehouse/, '/api'),
      },
      // Settlement Service (port 8085) — harus lebih spesifik, di-check duluan
      '/api/settlement': {
        target: 'http://localhost:8085',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/settlement/, '/api'),
      },
      // Auth Service (port 8080) — untuk login, register, dll
      '/api/v1/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/api/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    }
  }
})
