# GradLink

GradLink 是一个面向国内应届生的移动端优先求职交流网站。MVP 包含邮箱密码登录、应届生资料门槛、帖子、评论、回复、点赞、收藏、举报、关键词审核和轻量管理后台 API。

## 项目结构

- `backend/`: Go + Gin + GORM + SQLite REST API。
- `frontend/`: Vue 3 + Vite + TypeScript 移动端优先 Web 前端。
- `docs/`: 产品设计、实现计划和本地运行文档。

## 快速启动

### 依赖与 PATH

需要 Go、Node.js 和 npm。若系统 PATH 已包含它们，可直接执行：

```powershell
go version
node -v
npm -v
```

如果命令不存在，但本机已使用项目缓存工具链，可在当前 PowerShell 临时加入 PATH：

```powershell
$env:PATH='C:\Users\16021\.cache\gradlink-toolchains\go1.26.4\go\bin;C:\Users\16021\.cache\gradlink-toolchains\v24.17.0;' + $env:PATH
```

该设置只影响当前终端窗口。换新窗口后需要重新执行，或自行把 Go/Node 安装目录加入系统 PATH。

后端：

```powershell
cd backend
$env:GRADLINK_ADMIN_PASSWORD="admin-secret"
go run ./cmd/server
```

前端：

```powershell
cd frontend
npm install
npm run dev
```

默认地址：

- 前端：`http://localhost:5173`
- 后端：`http://localhost:8080`
- 健康检查：`http://localhost:8080/api/health`

### 手机真机测试

电脑和手机连接同一个 Wi-Fi 后，把前端开发服务以局域网可访问方式启动：

```powershell
cd frontend
npm run dev -- --host 0.0.0.0
```

然后在手机浏览器打开：

```text
http://172.20.10.10:5173
```

如果手机无法打开，优先检查 Windows 防火墙是否允许 Node/Vite 访问专用网络。

管理员账号
    邮箱：admin@gradlink.local
    密码：admin-secret

更完整的本地运行、管理员账号和重置数据说明见 [docs/run-local.md](docs/run-local.md)。

## 验证命令

后端：

```powershell
cd backend
go test ./... -count=1
```

前端：

```powershell
cd frontend
npm test -- --run
npm run build
```
