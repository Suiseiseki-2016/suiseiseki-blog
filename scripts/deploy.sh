#!/bin/bash

# 部署脚本
# 用于将编译好的二进制文件和前端构建产物同步到服务器

set -e

# 配置
SERVER_USER="${DEPLOY_USER:-root}"
SERVER_HOST="${DEPLOY_HOST:-your-server.com}"
SERVER_PATH="${DEPLOY_PATH:-/var/lib/blog}"
SERVICE_NAME="${SERVICE_NAME:-blog-suiseiseki}"

echo "🚀 开始部署..."

# 1. 构建前端
echo "📦 构建前端..."
cd frontend
npm install
npm run build
cd ..

# 2. 构建后端（Linux amd64）
echo "🔨 构建后端..."
cd backend
GOOS=linux GOARCH=amd64 go build -o blog-suiseiseki main.go
cd ..

# 3. 同步文件到服务器
echo "📤 同步文件到服务器..."
rsync -avz --delete \
  backend/blog-suiseiseki \
  frontend/dist/ \
  "${SERVER_USER}@${SERVER_HOST}:${SERVER_PATH}/"

# 4. 同步配置文件（如果需要）
if [ -f "Caddyfile" ]; then
  echo "📋 同步 Caddyfile..."
  rsync -avz Caddyfile "${SERVER_USER}@${SERVER_HOST}:${SERVER_PATH}/"
fi

# 5. 重启服务
echo "🔄 重启服务..."
ssh "${SERVER_USER}@${SERVER_HOST}" "
  cd ${SERVER_PATH}
  chmod +x blog-suiseiseki
  sudo systemctl restart ${SERVICE_NAME} || echo '服务未配置，请手动启动'
  sudo systemctl reload caddy || echo 'Caddy 未配置'
"

echo "✅ 部署完成！"
echo ""
echo "📝 后续步骤："
echo "1. 确保服务器上已配置 systemd 服务单元"
echo "2. 确保 Caddy 已配置并运行"
echo "3. 检查服务状态: ssh ${SERVER_USER}@${SERVER_HOST} 'systemctl status ${SERVICE_NAME}'"
