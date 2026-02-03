#!/bin/bash
# 开发环境一键启动：后台启动后端，前台启动前端；Ctrl+C 会同时停止后端。
# 与「分开启动」一致：不注入 CONFIG_PATH/POSTS_PATH，后端用 ../config.yaml 和 ../posts（在 backend 目录下）。

set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

# 后台启动后端（输出同时写入 backend.log 并在终端显示，便于排查 CONNECTION_REFUSED）
echo "启动后端 (http://127.0.0.1:8080)，日志见 backend.log ..."
(cd backend && go run main.go) 2>&1 | tee "$ROOT/backend.log" &
BACKEND_PID=$!

cleanup() {
  echo ""
  echo "正在停止后端 (PID $BACKEND_PID)..."
  kill "$BACKEND_PID" 2>/dev/null || true
  exit 0
}
trap cleanup SIGINT SIGTERM

# 等待后端就绪（首次可能需 clone 远程仓库，最多等 60 秒）；绕过代理直连本机
echo "等待后端就绪（首次可能需 clone 远程仓库，约 30–60 秒）..."
if command -v curl &>/dev/null; then
  for i in $(seq 1 60); do
    if curl -s -o /dev/null -w "%{http_code}" --noproxy '*' "http://127.0.0.1:8080/health" 2>/dev/null | grep -q 200; then
      echo "后端已就绪"
      break
    fi
    [ "$i" -eq 60 ] && echo "错误: 后端 60 秒内未就绪，请查看上方 backend 输出或 backend.log"
    sleep 1
  done
else
  sleep 15
fi

# 前端直连后端（用 127.0.0.1 与健康检查一致，避免 localhost 解析到 IPv6 导致连不上）
echo "启动前端 (http://127.0.0.1:3000)，API 直连 http://127.0.0.1:8080 ..."
(cd frontend && VITE_API_URL=http://127.0.0.1:8080 npm run dev)
