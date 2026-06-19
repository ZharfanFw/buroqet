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
      // Jembatan untuk Pricing Service (Biarkan saja, jangan dihapus)
      '/api/pricing': {
        target: 'http://localhost:8080', 
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/pricing/, '/pricing')
      },
      
      // 👇 JEMBATAN UNTUK ORDER SERVICE (SUDAH DIPERBAIKI):
      '/api/orders': {
        target: 'http://localhost:8082', 
        changeOrigin: true,
        // Ini akan mengubah fetch('/api/orders') dari React 
        // menjadi ('/api/v1/orders') saat masuk ke Golang
        rewrite: (path) => path.replace(/^\/api\/orders/, '/api/v1/orders') 
      }, 

      // Tambahkan di dalam object proxy server:
      '/api/epod': {
        target: 'http://localhost:8080', // Sesuai dengan r.Run(":8080") di Golang-mu
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/epod/, '') 
      }
    }
  }
})