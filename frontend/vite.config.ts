import path from 'path'
import fs from 'fs'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { VitePWA } from 'vite-plugin-pwa'

export default defineConfig({
  plugins: [
    react(),
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['favicon.ico', 'apple-touch-icon.png'],
      manifest: {
        name: 'Skopidom — Inventory',
        short_name: 'Skopidom',
        description: 'Automated inventory accounting system',
        theme_color: '#0f172a',
        background_color: '#ffffff',
        display: 'standalone',
        orientation: 'portrait',
        icons: [
          {
            src: 'pwa-192x192.png',
            sizes: '192x192',
            type: 'image/png',
          },
          {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
          },
        ],
      },
      workbox: {
        // Cache API responses for offline support.
        runtimeCaching: [
          {
            urlPattern: /^https?:\/\/.*\/api\/v1\/(categories|buildings|rooms)/,
            handler: 'StaleWhileRevalidate',
            options: {
              cacheName: 'lookup-data',
              expiration: { maxAgeSeconds: 60 * 60 * 24 },
            },
          },
        ],
      },
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    https: {
      key: fs.readFileSync('./ssl/192.168.1.200-key.pem'),
      cert: fs.readFileSync('./ssl/192.168.1.200.pem'),
      // key: fs.readFileSync('./ssl/10.112.252.76-key.pem'),
      // cert: fs.readFileSync('./ssl/10.112.252.76.pem'),
    },
    host: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/static': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
