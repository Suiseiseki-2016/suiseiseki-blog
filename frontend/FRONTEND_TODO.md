# 前端实现与待办

对照 PRD 与当前代码整理：已实现项与可选项。

## 已实现（符合 PRD）

| 需求 | 实现 |
|------|------|
| 文章列表 | 按日期排序，展示标题、摘要、发布时间、分类；加载/错误/空状态 |
| 文章详情 | Markdown 转 HTML 展示，`.prose` 样式；返回列表入口 |
| 简历页 | 硬编码结构；`@media print` 隐藏 nav/footer；「打印/导出 PDF」按钮 |
| 布局 | Layout：导航（首页、简历）、主内容区、页脚 |
| 开发代理 | Vite 将 `/api` 代理到后端 `localhost:8080` |
| 路由 | `/` 列表、`/posts/:slug` 详情、`/resume` 简历 |

## 本次补充

| 项目 | 说明 |
|------|------|
| 404 页 | 新增 `NotFound.jsx`，`Route path="*"` 捕获未匹配路由 |
| 页面标题 (SEO) | 各页 `useEffect` 设置 `document.title`（列表/详情/简历/404），离开时恢复 |
| 分页 | 列表页「加载更多」：`limit=10&offset=...`，后端已支持 |

## 可选后续

- **SEO 增强**：详情页用 `post.summary` 写 `<meta name="description">`（需 react-helmet 或类似，或服务端/SSR）  
- **favicon**：替换 `index.html` 中的 `/vite.svg` 为站点图标  
- **无障碍**：为导航、按钮等加 `aria-label`、焦点样式  
- **简历内容**：将 Resume 中占位文案改为真实信息  

## 运行与构建

```bash
# 开发（需后端 8080 已启动）
npm run dev

# 构建
npm run build
```

生产环境需保证与后端同源或 Caddy 将 `/api` 反向代理到后端，前端请求 `/api/posts` 等即可。
