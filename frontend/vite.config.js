import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { readFileSync, existsSync } from 'fs'
import { resolve } from 'path'

// Read server.port and frontend.port from config.yaml
// Use proxy so browser only talks to same origin (localhost:3000) — avoids 502/504 when system proxy intercepts 127.0.0.1
let serverPort = '8080'
let devPort = 3000
const configPath = resolve(__dirname, '../config.yaml')
if (existsSync(configPath)) {
  const content = readFileSync(configPath, 'utf8')
  serverPort = content.match(/^server:\s*\n\s*port:\s*["']?(\d+)/m)?.[1] || '8080'
  const frontPort = content.match(/^frontend:\s*\n[\s\S]*?port:\s*["']?(\d+)/m)?.[1] || '3000'
  devPort = parseInt(frontPort, 10)
}

const backendTarget = `http://127.0.0.1:${serverPort}`

export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    port: devPort,
    strictPort: true,
    allowedHosts: ['www.aeoluswu.info', 'aeoluswu.info'],
    proxy: {
      '/api': {
        target: backendTarget,
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
  },
})
