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

# 5. 后台启动后端（在 backend 目录下运行，与单独 cd backend && go run main.go 一致；不注入 CONFIG_PATH/POSTS_PATH）
echo "[启动] 后端 (http://127.0.0.1:8080)，日志见 backend.log ..."
(cd backend && go run main.go) 2>&1 | tee "$ROOT/backend.log" &
BACKEND_PID=$!

cleanup() {
  echo ""
  echo "[停止] 后端 (PID $BACKEND_PID)..."
  kill "$BACKEND_PID" 2>/dev/null || true
  exit 0
}
trap cleanup SIGINT SIGTERM

# 6. 等待后端就绪（首次可能需 clone 远程仓库，最多等 60 秒）；绕过代理直连本机
echo "[启动] 等待后端就绪（首次可能需 clone 远程仓库，约 30–60 秒）..."
if command -v curl &>/dev/null; then
  for i in $(seq 1 60); do
    if curl -s -o /dev/null -w "%{http_code}" --noproxy '*' "http://127.0.0.1:8080/health" 2>/dev/null | grep -q 200; then
      echo "[启动] 后端已就绪"
      break
    fi
    [ "$i" -eq 60 ] && echo "[启动] 错误: 后端 60 秒内未就绪，请查看上方 backend 输出或 backend.log"
    sleep 1
  done
else
  sleep 15
fi

# 7. 前台启动前端（直连后端，用 127.0.0.1 避免 localhost 解析到 IPv6 导致连不上）
echo "[启动] 前端 (http://127.0.0.1:3000)，API 直连 http://127.0.0.1:8080 ..."
(cd frontend && VITE_API_URL=http://127.0.0.1:8080 npm run dev)
