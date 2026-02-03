# 🧪 后端测试指南

本文档说明如何测试后端服务。

## 快速开始

### 运行所有测试

```bash
cd backend
go test -v ./...
```

### 运行测试脚本

```bash
chmod +x test.sh
./test.sh
```

## 测试类型

### 1. 单元测试

#### 数据库测试 (`database/db_test.go`)
- 测试数据库初始化
- 测试表结构创建
- 测试目录自动创建

```bash
go test -v ./database
```

#### 工具函数测试 (`utils/markdown_test.go`)
- 测试Markdown解析
- 测试Front-matter提取
- 测试Markdown转HTML
- 测试Slug生成

```bash
go test -v ./utils
```

#### 处理器测试 (`handlers/posts_test.go`)
- 测试获取文章列表
- 测试获取单篇文章
- 测试404错误处理

```bash
go test -v ./handlers
```

#### 同步服务测试 (`services/sync_test.go`)
- 测试文章同步
- 测试删除已移除的文章

```bash
go test -v ./services
```

### 2. 集成测试

#### 健康检查测试 (`main_test.go`)
- 测试健康检查端点

```bash
go test -v ./main_test.go
```

## 手动测试

### 1. 启动服务器

```bash
cd backend
go run main.go
```

### 2. 测试API端点

#### 健康检查
```bash
curl http://localhost:8080/health
```

预期响应：
```json
{"status":"ok"}
```

#### 获取文章列表
```bash
curl http://localhost:8080/api/posts
```

#### 获取单篇文章
```bash
curl http://localhost:8080/api/posts/your-slug
```

### 3. 测试Webhook（需要配置secret）

```bash
curl -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -H "X-Hub-Signature-256: sha256=your-signature" \
  -d '{}'
```

## 测试覆盖率

生成覆盖率报告：

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

查看HTML报告：

```bash
go tool cover -html=coverage.out
```

## 测试最佳实践

1. **使用临时目录**：所有测试使用 `t.TempDir()` 创建临时目录
2. **清理资源**：使用 `defer` 确保测试后清理资源
3. **独立测试**：每个测试都是独立的，不依赖其他测试
4. **测试模式**：使用 `gin.TestMode` 进行HTTP测试

## 常见问题

### 测试失败：数据库锁定
- 确保所有数据库连接都已正确关闭
- 使用 `defer db.Close()` 或 `defer cleanup()`

### 测试失败：文件不存在
- 确保使用 `t.TempDir()` 创建临时目录
- 检查文件路径是否正确

### 端口占用
- 测试不会启动实际服务器，只测试逻辑
- 如果手动测试时端口被占用，修改 `PORT` 环境变量
