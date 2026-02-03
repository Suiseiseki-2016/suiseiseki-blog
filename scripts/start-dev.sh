#!/bin/bash
# Dev startup: backend in background, frontend in foreground. Ctrl+C stops backend.
# Ports read from config.yaml (server.port, frontend.port).

set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

# Read ports from config.yaml (defaults: 8080, 3000)
BACKEND_PORT=8080
FRONTEND_PORT=3000
if [ -f config.yaml ]; then
  BACKEND_PORT=$(sed -n '/^server:/,/^[a-z]/p' config.yaml | grep -E '^\s*port:' | head -1 | grep -oE '[0-9]+' | head -1 || echo 8080)
  FRONTEND_PORT=$(sed -n '/^frontend:/,/^[a-z]/p' config.yaml | grep -E '^\s*port:' | head -1 | grep -oE '[0-9]+' | head -1 || echo 3000)
fi
[ -z "$BACKEND_PORT" ] && BACKEND_PORT=8080
[ -z "$FRONTEND_PORT" ] && FRONTEND_PORT=3000

echo "Starting backend (http://127.0.0.1:$BACKEND_PORT), log: backend.log ..."
# Log to file only so backend never blocks on terminal (tee can cause ECONNREFUSED when frontend proxy hits backend)
(cd backend && go run main.go >> "$ROOT/backend.log" 2>&1) &
BACKEND_PID=$!

cleanup() {
  echo ""
  echo "Stopping backend (PID $BACKEND_PID)..."
  kill "$BACKEND_PID" 2>/dev/null || true
  pkill -f "go run main.go" 2>/dev/null || true
  exit 0
}
trap cleanup SIGINT SIGTERM

echo "Waiting for backend (first run may clone remote repo, up to 60s)..."
BACKEND_READY=0
if command -v curl &>/dev/null; then
  for i in $(seq 1 60); do
    if curl -s -o /dev/null -w "%{http_code}" --noproxy '*' "http://127.0.0.1:$BACKEND_PORT/health" 2>/dev/null | grep -q 200; then
      echo "Backend ready"
      BACKEND_READY=1
      break
    fi
    sleep 1
  done
  if [ "$BACKEND_READY" -eq 0 ]; then
    echo "Error: backend did not become ready in 60s. Check backend.log and port $BACKEND_PORT."
    kill "$BACKEND_PID" 2>/dev/null || true
    exit 1
  fi
else
  sleep 15
fi

echo "Starting frontend (http://127.0.0.1:$FRONTEND_PORT), /api -> backend :$BACKEND_PORT (proxy, no 502/504) ..."
(cd frontend && npm run dev)
