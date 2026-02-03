#!/bin/bash
# 一键配置并启动前后端：安装依赖、可选复制 .env、后台起后端、前台起前端；Ctrl+C 会同时停止后端。

set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "=== 一键配置与启动 ==="

# 1. 配置：.env（若不存在则从 .env.example 复制）
if [ ! -f .env ]; then
  if [ -f .env.example ]; then
    echo "[配置] 未找到 .env，从 .env.example 复制（可按需修改）"
    cp .env.example .env
  else
    echo "[配置] 未找到 .env 与 .env.example，将使用 config.yaml 或内置默认"
  fi
else
  echo "[配置] 已存在 .env"
fi

# 2. 后端依赖
echo "[配置] 安装后端依赖..."
(cd backend && go mod download)

# 3. 前端依赖
echo "[配置] 安装前端依赖..."
(cd frontend && npm install)

# 4. 加载 .env 到当前 shell（若存在），供后续子进程继承
if [ -f .env ]; then
  set -a
  source .env
  set +a
fi

# 5. 后台启动后端（必须在 backend 目录下运行，以便 ../posts 指向项目根目录的 posts）
echo "[启动] 后端 (http://localhost:8080)..."
(cd backend && go run main.go) &
BACKEND_PID=$!

cleanup() {
  echo ""
  echo "[停止] 后端 (PID $BACKEND_PID)..."
  kill "$BACKEND_PID" 2>/dev/null || true
  exit 0
}
trap cleanup SIGINT SIGTERM

sleep 2

# 6. 前台启动前端（Vite 会从 config.yaml 读取 frontend.api_base_url）
echo "[启动] 前端 (http://localhost:3000)..."
(cd frontend && npm run dev)
