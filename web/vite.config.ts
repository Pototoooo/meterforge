import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'node:path'

export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('recharts') || id.includes('d3-')) return 'charts'
          if (id.includes('react-dom') || id.includes('react-router') || id.includes('/react/')) return 'react-vendor'
          if (id.includes('@tanstack')) return 'query'
        },
      },
    },
  },
  server: {
    port: 3000,
    fs: { allow: [path.resolve(__dirname, '..')] },
    proxy: {
      '/api': {
        target: process.env.METERFORGE_API_URL ?? 'http://localhost:48888',
        changeOrigin: true,
      },
    },
  },
  preview: { port: 3000 },
})
