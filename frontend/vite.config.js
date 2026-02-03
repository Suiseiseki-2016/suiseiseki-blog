import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { readFileSync, existsSync } from 'fs'
import { resolve } from 'path'

// 从项目根 config.yaml 读取 frontend.api_base_url，供前端直连后端（代理不生效时）
const configPath = resolve(__dirname, '../config.yaml')
if (existsSync(configPath)) {
  const content = readFileSync(configPath, 'utf8')
  const m = content.match(/api_base_url:\s*["']?([^"'\s#]+)["']?/)
  if (m) process.env.VITE_API_URL = m[1].trim()
}

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    strictPort: true,  // 端口被占用时直接报错，不等待
    // host: '0.0.0.0',  // 需要本机 IP 访问时再取消注释；部分环境下会卡住启动
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
  },
})
