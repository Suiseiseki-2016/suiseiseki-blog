#!/bin/bash
# 开发环境一键启动：后台启动后端，前台启动前端；Ctrl+C 会同时停止后端。

set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

# 后台启动后端（输出重定向到 backend.log，避免与前端输出混在一起）
echo "启动后端 (http://localhost:8080)，日志见 backend.log ..."
cd backend
go run main.go &> "$ROOT/backend.log" &
BACKEND_PID=$!
cd ..

# 退出时杀掉后端
cleanup() {
  echo ""
  echo "正在停止后端 (PID $BACKEND_PID)..."
  kill "$BACKEND_PID" 2>/dev/null || true
  exit 0
}
trap cleanup SIGINT SIGTERM

# 等待后端就绪
sleep 2

# 前台启动前端（Vite 会从 config.yaml 读取 frontend.api_base_url，无需脚本同步）
echo "启动前端 (http://localhost:3000)..."
cd frontend
npm run dev
