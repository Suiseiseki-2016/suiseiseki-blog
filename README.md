# Blog Suiseiseki

个人全栈博客：**Go + SQLite** 后端，**React** 前端，文章存于独立 **GitHub 仓库**，支持 Webhook 与定期拉取自动同步。

## 仓库用途

- **本仓库**：博客的代码与配置（后端、前端、脚本、`config.yaml` 等）。
- **文章**：来自单独的「文章仓库」（仅 Markdown），通过 Webhook 或定期同步拉取到服务器。

## 快速开始

```bash
# 依赖
cd backend && go mod download && cd ..
cd frontend && npm install && cd ..

# 启动（先后端、再前端，或使用脚本）
./scripts/setup-and-start.sh   # 一键配置 + 启动
# 或
./scripts/start-dev.sh         # 仅启动
```

浏览器访问 **http://localhost:3000**。配置与 GitHub 同步说明见 [CONFIG.md](CONFIG.md)。

**`npm run dev` 卡住不动**：先看 3000 端口是否被占（`lsof -i :3000`）；若被占则关掉对应进程或改 `vite.config.js` 里 `port`。若仍卡在 `> vite` 无输出，可试 Node 18/20 LTS（部分 Node 23 环境会卡住），或删掉 `frontend/node_modules` 和 `package-lock.json` 后重新 `npm install`。

**页面一直「加载中」且后端收不到请求**：多半是 Vite 代理没把 `/api` 转到后端。在 **config.yaml** 里设置 `frontend.api_base_url: "http://localhost:8080"`，Vite 启动时会自动读取，前端将直连后端；后端开发模式已开 CORS。

