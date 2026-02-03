#!/bin/bash
# Setup and start: install deps, optional .env, then run backend and frontend in background (nohup). Safe to close terminal after.

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

echo "=== Setup and start (background) ==="

# 1. .env (copy from .env.example if missing)
if [ ! -f .env ]; then
  if [ -f .env.example ]; then
    echo "[config] No .env, copied from .env.example"
    cp .env.example .env
  else
    echo "[config] No .env or .env.example, using config.yaml defaults"
  fi
else
  echo "[config] .env exists"
fi

# 2. Backend deps
echo "[config] Installing backend deps..."
(cd backend && go mod download)

# 3. Frontend deps
echo "[config] Installing frontend deps..."
(cd frontend && npm install)

# 4. Load .env if present
if [ -f .env ]; then
  set -a
  source .env
  set +a
fi

# 5. Start backend in background (nohup)
echo "[start] Backend (http://127.0.0.1:$BACKEND_PORT), log: backend.log ..."
(cd "$ROOT/backend" && nohup go run main.go >> "$ROOT/backend.log" 2>&1) &
BACKEND_PID=$!

# 6. Wait for backend (first run may clone remote repo, up to 60s)
echo "[start] Waiting for backend (first run may clone repo, ~30-60s)..."
if command -v curl &>/dev/null; then
  for i in $(seq 1 60); do
    if curl -s -o /dev/null -w "%{http_code}" --noproxy '*' "http://127.0.0.1:$BACKEND_PORT/health" 2>/dev/null | grep -q 200; then
      echo "[start] Backend ready"
      break
    fi
    [ "$i" -eq 60 ] && echo "[start] Error: backend not ready in 60s, check backend.log"
    sleep 1
  done
else
  sleep 15
fi

# 7. Start frontend in background (nohup)
echo "[start] Frontend (http://127.0.0.1:$FRONTEND_PORT), log: frontend.log ..."
(cd "$ROOT/frontend" && nohup npm run dev >> "$ROOT/frontend.log" 2>&1) &
FRONTEND_PID=$!

sleep 3

# Show LAN IP so other devices can access (frontend is on 0.0.0.0)
LAN_IP=""
if command -v ipconfig &>/dev/null; then
  LAN_IP=$(ipconfig getifaddr en0 2>/dev/null || ipconfig getifaddr en1 2>/dev/null || true)
fi
[ -z "$LAN_IP" ] && command -v hostname &>/dev/null && LAN_IP=$(hostname -I 2>/dev/null | awk '{print $1}')

echo ""
echo "=== Running in background (serves LAN) ==="
echo "  Local:    http://127.0.0.1:$FRONTEND_PORT"
if [ -n "$LAN_IP" ]; then
  echo "  LAN:      http://${LAN_IP}:$FRONTEND_PORT  (other devices on same network)"
fi
echo "  Backend:  http://127.0.0.1:$BACKEND_PORT (internal)"
echo "  Logs:     backend.log  frontend.log"
echo ""
echo "To stop:"
echo "  Backend:  kill $BACKEND_PID   or  pkill -f 'go run main.go'"
echo "  Frontend: kill $FRONTEND_PID   or  pkill -f 'vite'"
echo ""
