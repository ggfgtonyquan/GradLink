# GradLink 本地运行指南

## 环境要求

- Go 1.24 或更高版本。
- Node.js 20 或更高版本。
- npm。

## 启动后端

```bash
cd backend
go run ./cmd/server
```

默认监听：

- `http://localhost:8080`
- 健康检查：`http://localhost:8080/api/health`

默认 SQLite 文件：

- `backend/data/gradlink.db`

可选环境变量：

```bash
GRADLINK_ADDR=:8080
GRADLINK_DB_PATH=data/gradlink.db
GRADLINK_ADMIN_PASSWORD=admin-secret
GRADLINK_COOKIE_SECURE=false
```

说明：

- 不设置 `GRADLINK_ADMIN_PASSWORD` 时，种子管理员账号会存在但默认禁用，不能登录。
- 设置 `GRADLINK_ADMIN_PASSWORD` 后，可用 `admin@gradlink.local` 和该密码登录管理员账号。
- 本地 HTTP 开发时保持 `GRADLINK_COOKIE_SECURE=false`，否则浏览器不会在非 HTTPS 下保存安全 Cookie。

## 启动前端

```bash
cd frontend
npm install
npm run dev
```

默认前端地址：

- `http://localhost:5173`

Vite 已配置 `/api` 代理到 `http://127.0.0.1:8080`，前端开发时可以直接请求 `/api/...`。

## 重置本地数据

停止后端服务后删除 SQLite 文件即可：

```bash
rm backend/data/gradlink.db
```

Windows PowerShell：

```powershell
Remove-Item backend\data\gradlink.db -ErrorAction SilentlyContinue
Remove-Item backend\data\gradlink.db-shm -ErrorAction SilentlyContinue
Remove-Item backend\data\gradlink.db-wal -ErrorAction SilentlyContinue
```

再次启动后端会自动迁移表结构并写入标签、审核关键词和管理员种子数据。

## 本地验证清单

1. 启动后端和前端。
2. 注册邮箱和密码。
3. 登录。
4. 未完善资料时尝试发帖，应看到资料门槛。
5. 完善应届生资料。
6. 发布一篇正常帖子，应直接发布。
7. 发布一篇命中关键词的帖子，例如包含“加微信”，应进入待审核。
8. 评论、回复、点赞、收藏和举报。
9. 使用管理员账号处理审核和举报。

## 常用测试命令

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
