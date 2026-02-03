# 🚀 个人博客系统项目需求文档 (PRD)

## 1. 项目概述

构建一个极简、高性能的个人全栈博客系统。采用 **Go + SQLite** 作为后端，**React** 作为前端，文章存储于 **GitHub**，通过 **Webhook** 实现全自动同步更新。

---

## 2. 核心架构设计

* **内容存储：** GitHub 仓库 (Markdown 文件)。
* **后端 (Go):** 监听 Webhook、拉取 Git 内容、解析 Markdown、提供 REST API。
* **数据库 (SQLite 3):** 存储文章元数据（标题、标签、路径、发布日期）。
* **前端 (React + Vite):** 负责页面展示，包括响应式博客列表、详情页和硬编码简历页。
* **部署层:** Caddy 作为反向代理，负责 HTTPS 和静态文件托管。

---

## 3. 功能需求

### 3.1 自动化同步 (Core)

* **Webhook 监听：** 接收来自 GitHub 的 `push` 事件。
* **安全验证：** 使用 GitHub Webhook Secret 进行签名校验。
* **同步逻辑：** 1.  触发后执行 `git pull` 更新本地文章目录。
2.  扫描目录下的 `.md` 文件。
3.  解析文件头部的 YAML（Front-matter）。
4.  对比数据库，增量更新/删除文章索引。

### 3.2 博客功能

* **文章列表：** 支持按日期排序，展示标题、摘要和发布时间。
* **文章详情：** 后端将 Markdown 转换为 HTML，前端负责样式渲染。
* **无评论系统：** 纯粹的内容展示。

### 3.3 简历页面

* **硬编码展示：** 在 React 中直接编写简历组件（以便精细控制排版）。
* **PDF 导出：** 适配 `@media print` 样式，支持用户通过浏览器“另存为 PDF”。

---

## 4. 技术栈细节

| 模块 | 技术选型 | 理由 |
| --- | --- | --- |
| **后端** | **Go (Gin)** | 内存占用低，单二进制文件，并发能力强。 |
| **数据库** | **SQLite 3** | 无需独立进程，读性能极高。 |
| **MD 解析** | **Goldmark** | Go 生态最标准、最快的 Markdown 库。 |
| **前端** | **React + Tailwind CSS** | 开发效率高，简历页排版极其简单。 |
| **反向代理** | **Caddy** | 自动 HTTPS，配置极简。 |

---

## 5. 接口设计 (API)

* `POST /api/webhook`: 供 GitHub 调用，触发同步。
* `GET /api/posts`: 获取文章列表（分页可选）。
* `GET /api/posts/:slug`: 获取单篇文章的 HTML 内容和元数据。

---

## 6. 数据库表结构 (SQLite)

```sql
CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT UNIQUE,          -- URL 别名
    title TEXT,                -- 标题
    summary TEXT,              -- 摘要
    category TEXT,             -- 分类
    published_at DATETIME,     -- 发布日期
    content_path TEXT          -- 对应的本地文件路径
);

```

---

## 7. 非功能性需求

* **性能：** 文章加载时间 (TTFB) 需小于 100ms。
* **资源：** 整体内存占用控制在 100MB 以内（不含操作系统）。
* **SEO：** 确保每篇文章有独立的 URL 和清晰的 HTML 结构。

---

## 8. 开发路线图 (Roadmap)

1. **Phase 1 (MVP):** 编写 Go 后端，手动触发 `git pull` 并解析 MD 存入 SQLite。
2. **Phase 2 (Webhook):** 实现 GitHub 自动通知同步。
3. **Phase 3 (Frontend):** React 编写列表和详情页，接入 API。
4. **Phase 4 (Resume):** 完成简历页及 PDF 打印样式适配。
5. **Phase 5 (Ops):** 配置 Caddyfile，部署到 1GB 服务器并开启 HTTPS。

---

**你想让我先为你展示哪一部分的具体实现代码？比如：Go 处理 Webhook 的逻辑，或者是 React 简历页的打印适配方案？**