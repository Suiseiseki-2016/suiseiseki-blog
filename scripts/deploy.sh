#!/bin/bash

# Deploy script: build backend and frontend, sync to server, restart services.

set -e

SERVER_USER="${DEPLOY_USER:-root}"
SERVER_HOST="${DEPLOY_HOST:-your-server.com}"
SERVER_PATH="${DEPLOY_PATH:-/var/lib/blog}"
SERVICE_NAME="${SERVICE_NAME:-blog-suiseiseki}"

echo "Deploying..."

# 1. Build frontend
echo "Building frontend..."
cd frontend
npm install
npm run build
cd ..

# 2. Build backend (Linux amd64)
echo "Building backend..."
cd backend
GOOS=linux GOARCH=amd64 go build -o blog-suiseiseki main.go
cd ..

# 3. Sync to server
echo "Syncing to server..."
rsync -avz --delete \
  backend/blog-suiseiseki \
  frontend/dist/ \
  "${SERVER_USER}@${SERVER_HOST}:${SERVER_PATH}/"

# 4. Sync Caddyfile if present
if [ -f "Caddyfile" ]; then
  echo "Syncing Caddyfile..."
  rsync -avz Caddyfile "${SERVER_USER}@${SERVER_HOST}:${SERVER_PATH}/"
fi

# 5. Restart services
echo "Restarting services..."
ssh "${SERVER_USER}@${SERVER_HOST}" "
  cd ${SERVER_PATH}
  chmod +x blog-suiseiseki
  sudo systemctl restart ${SERVICE_NAME} || echo 'Service not configured; start manually'
  sudo systemctl reload caddy || echo 'Caddy not configured'
"

echo "Deploy complete."
echo ""
echo "Next steps:"
echo "1. Ensure systemd unit is configured on the server"
echo "2. Ensure Caddy is configured and running"
echo "3. Check status: ssh ${SERVER_USER}@${SERVER_HOST} 'systemctl status ${SERVICE_NAME}'"
