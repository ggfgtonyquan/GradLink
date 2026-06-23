# GradLink

GradLink 是一个面向国内应届生的移动端优先求职交流网站。MVP 包含邮箱密码登录、应届生资料门槛、帖子、评论、回复、点赞、收藏、举报、关键词审核和轻量管理后台 API。

## 项目结构

- `backend/`: Go + Gin + GORM + SQLite REST API。
- `frontend/`: Vue 3 + Vite + TypeScript 移动端优先 Web 前端。
- `docs/`: 产品设计、实现计划和本地运行文档。

## 快速启动

后端：

```bash
cd backend
go run ./cmd/server
```

前端：

```bash
cd frontend
npm install
npm run dev
```

默认地址：

- 前端：`http://localhost:5173`
- 后端：`http://localhost:8080`
- 健康检查：`http://localhost:8080/api/health`

更完整的本地运行、管理员账号和重置数据说明见 [docs/run-local.md](docs/run-local.md)。

## 验证命令

后端：

```bash
cd backend
go test ./...
```

前端：

```bash
cd frontend
npm test -- --run
npm run build
```
