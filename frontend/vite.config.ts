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
      // Order Management Service
      '/api/v1/orders': {
        target: 'http://localhost:8082',
        changeOrigin: true,
      },
      // Tracking Service
      '/api/v1/tracking': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      // Auth Service
      '/api/v1/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // Pricing Service
      '/api/v1/pricing': {
        target: 'http://localhost:8084',
        changeOrigin: true,
      },
      // Dispatch Service
      '/api/v1/dispatch': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
      // Settlement Service
      '/api/v1/settlement': {
        target: 'http://localhost:8085',
        changeOrigin: true,
      },
      // ePOD Service
      '/api/v1/epod': {
        target: 'http://localhost:8086',
        changeOrigin: true,
      },
      // Warehouse Service
      '/api/v1/warehouse': {
        target: 'http://localhost:8087',
        changeOrigin: true,
      },
    }
  }
})

