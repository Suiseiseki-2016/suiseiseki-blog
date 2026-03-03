# 配置说明

包含：**config.yaml / .env / 环境变量** 与 **GitHub 同步**。

优先级：**环境变量 > .env > config.yaml > 内置默认**。

---

## 〇、推荐先设置远程文章仓库

博客内容来自一个**独立的 GitHub 仓库**（只放 Markdown）。建议先建好这个仓库，本地和生产都指向它，后续定期同步才有「远程」可拉。

### 步骤概要

1. **在 GitHub 新建一个空仓库**  
   例如：`my-blog-posts`（或任意名称），不要勾选「Add a README」，保持空仓库。

2. **在本地准备好文章目录并推送到该仓库**  
   - 方式 A：用项目自带的 `posts/`  
     ```bash
     cd /path/to/blog-suiseiseki/posts
     git init
     git add .
     git commit -m "init posts"
     git remote add origin https://github.com/你的用户名/my-blog-posts.git
     git branch -M main
     git push -u origin main
     ```
   - 方式 B：新建一个目录，放一篇示例 `.md` 再 `git init`、`add`、`commit`、`remote add origin`、`push`（同上）。

3. **之后本地开发**  
   - 若已在项目 `posts/` 里推送到远程：直接在本项目里 `cd posts && git pull` 即可更新。  
   - 若文章在别的目录：把该仓库 **clone 到项目的 `posts/`**（先备份或清空现有 `posts/` 再 `git clone <文章仓库地址> posts`），再启动后端。

4. **之后生产部署**
   在服务器上把该仓库 **clone 到博客的数据目录**（如 `/var/lib/blog/posts`），在 config / 环境变量里把 `POSTS_PATH` 或 `GIT_REPO_PATH` 指过去，并设置 `sync.interval_minutes`（见下文）。

这样你就有一个固定的「远程文章仓库」；本地和生产都从它拉取，便于统一管理和自动同步。

---

## 一、config.yaml

放在**项目根目录**，后端在 `backend/` 下启动时读取 `../config.yaml`；可用环境变量 `CONFIG_PATH` 指定路径。

```yaml
server:
  port: "8080"
  mode: dev   # dev | prod

database:
  path: "./blog.db"

posts:
  path: "../posts"
  remote_url: "https://github.com/你的用户名/blog-posts.git"   # 当 posts 无文章时自动 clone 此仓库；留空则不自动 clone

webhook:
  secret: ""
  git_repo_path: ""

sync:
  interval_minutes: 5

frontend:
  port: "3000"   # 前端开发服务器端口
```

| 字段 | 说明 | 默认 |
|------|------|------|
| `server.port` | 后端端口 | `8080` |
| `server.mode` | `dev` / `prod`（prod 下同步时会 `git pull`） | `dev` |
| `database.path` | SQLite 路径（相对 `backend/`） | `./blog.db` |
| `posts.path` | 文章目录；生产多为文章仓库 clone 路径 | dev: `../posts`，prod: `/var/lib/blog/posts` |
| `posts.remote_url` | 当 posts 无文章时自动 clone 的远程仓库 URL（如 `https://github.com/xxx/blog-posts.git`）；留空则不自动 clone | 空 |
| `webhook.secret` | (已废弃) |
| `webhook.git_repo_path` | (已废弃) |
| `sync.interval_minutes` | 定期同步间隔（分钟）；0 表示不启用，仅靠 Webhook 触发 | `0` |
| `frontend.port` | 前端开发服务器端口；前端直连后端 127.0.0.1:server.port（同机） | `3000` |

---

## 二、.env 与环境变量

**主配置在 config.yaml**，一般只需改 `config.yaml` 即可。  
**.env / 环境变量** 仅用于：① 覆盖 config.yaml 中的某项；② 存放敏感信息。`WEBHOOK_SECRET`、`GIT_REPO_PATH` 等相关设置已废弃。

复制 `.env.example` 为 `.env` 后，按需取消注释并填写。后端不自动加载 `.env`，需自行注入（如 `source .env` 后再运行，或 systemd 的 `EnvironmentFile`）。

| 变量 | 说明 |
|------|------|
| `WEBHOOK_SECRET` | (已废弃) |
| `PORT` | 覆盖 server.port |
| `MODE` | 覆盖 server.mode |
| `DB_PATH` | 覆盖 database.path |
| `POSTS_PATH` | 覆盖 posts.path |
| `POSTS_REMOTE_URL` | 覆盖 posts.remote_url（当 posts 无文章时自动 clone 的仓库） |
| `GIT_REPO_PATH` | 覆盖 生产文章目录 |
| `CONFIG_PATH` | 指定 config.yaml 路径 |
| `SYNC_INTERVAL_MINUTES` | 覆盖 sync.interval_minutes |

---

## 三、按环境示例

- **开发**：`config.yaml` 里 `mode: dev`、`posts.path: "../posts"`，`webhook.secret` 可留空。
- **生产**：`mode: prod`，`posts.path` / `GIT_REPO_PATH` 指向文章仓库 clone 目录。可选地设置 `sync.interval_minutes` 来启用定期拉取。

前端端口在 config.yaml 的 `frontend.port`（默认 3000）；开发时前端直连后端 127.0.0.1:server.port（同机，无代理）。

---

## 四、与 GitHub 同步

### 4.1 流程

文章来自**独立「文章仓库」**（只放 Markdown），不是博客代码仓库。

```
定期从远程文章仓库执行 `git pull` → 扫描 .md 更新数据库 → 博客显示最新
```

### 4.2 开发环境

- 不会自动和 GitHub 同步；用本地 `posts/`。
- 若要用 GitHub 上的文章：把文章仓库 clone 到 `posts/`，之后在 `posts/` 里 `git pull`，再重启后端。

> 启动时后端会执行一次同步以扫描 `posts/` 并写入数据库；开发环境只扫描本地目录，不会执行 `git pull`。

### 4.3 生产环境

1. **服务器**：把文章仓库 clone 到博客数据目录，例如  
   `git clone https://github.com/你的用户名/文章仓库名.git /var/lib/blog/posts`  
   配置 `POSTS_PATH=/var/lib/blog/posts`（或 `GIT_REPO_PATH`）、`MODE=prod`。  
   若不希望手动拉取，可以额外在配置/环境变量中设置 `sync.interval_minutes`。

<!-- webhook instructions removed -->

3. 之后每次 push 到文章仓库，生产只需依赖定时拉取（`sync.interval_minutes`），Webhook 配置可在不使用时删除。

> 部署后服务器也会在启动时立即同步并填充数据库，因此无需额外手动操作来获取第一次的内容。

**定期同步（可选）**：在 config.yaml 中设置 `sync.interval_minutes`（如 `5`），后端会每隔 N 分钟在文章目录执行 `git pull` 并更新数据库；生产模式下 Sync 会先执行 `git pull`。定期同步可与手动触发共存，用于确保仓库始终最新。


---

## 五、远程同步仓库（文章仓库）格式

远程「文章仓库」就是**一个普通 Git 仓库**，里面只放 **Markdown 文件**（`.md`）。后端会递归扫描该目录下所有 `.md` 并同步到数据库，便于搜索与索引。

### 5.1 仓库结构

- **任意目录结构**：`.md` 可以放在根目录，也可以放在子目录（如 `2024/01/my-post.md`），都会被扫描。
- **只认 `.md`**：其他文件（图片、附件等）可共存，但不会参与同步；如需在文章里引用图片，可用相对路径或图床链接。

### 5.2 单篇 Markdown 格式

每篇 `.md` 建议带 **Front-matter**（文件开头用 `---` 包起来的 YAML），后面是正文。  
没有 Front-matter 也可以：标题会为空，`slug` 会从文件名生成（如 `my-post.md` → `my-post`），发布日期用文件修改时间。

**推荐格式示例：**

```markdown
---
title: 文章标题
slug: my-post
summary: 文章摘要，用于列表页
category: 技术
published_at: 2024-01-01
---

# 正文标题

这里是 Markdown 正文...
```

### 5.3 Front-matter 字段说明

| 字段 | 必填 | 说明 |
|------|------|------|
| `title` | 否 | 文章标题；不填则列表/详情显示为空 |
| `slug` | 否 | URL 别名，需**唯一**；不填则从文件名生成（如 `hello-world.md` → `hello-world`） |
| `summary` | 否 | 摘要，列表页展示 |
| `category` | 否 | 分类标签 |
| `published_at` | 否 | 发布日期；支持 `2006-01-02`、`2006-01-02 15:04:05`、RFC3339；不填则用文件修改时间 |

- **slug 唯一性**：数据库里 `slug` 唯一，两篇若填相同 `slug` 会互相覆盖（后同步的为准）。建议每篇显式写不同 `slug`。

### 5.4 示例仓库结构

```
你的文章仓库/
├── 2024-01-01-hello.md
├── 2024-01-15-second-post.md
└── draft/
    └── wip.md
```

每篇按上面格式写 Front-matter 和正文即可；clone 到服务器并配置 Webhook 后，push 即自动同步。
