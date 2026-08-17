# GradLink 本地运行指南

除非代码块明确说明其他起始目录，每个独立命令块都从仓库根目录开始执行。

## 环境要求

- Go 1.26.6 或更高的兼容版本。
- Node.js 20.19+、22.12+ 或 24+（推荐当前 LTS）。
- npm。

## 启动后端

```powershell
cd backend
$env:GRADLINK_ADDR="127.0.0.1:8080"
Remove-Item Env:GRADLINK_ADMIN_PASSWORD -ErrorAction SilentlyContinue
go run ./cmd/server
```

此方式会让新建本地数据库写入一个禁用的管理员种子账号，不会启用管理员登录。

### 可选：启用管理员或轮换管理员密码

PowerShell：

```powershell
cd backend
$env:GRADLINK_ADDR="127.0.0.1:8080"
$env:GRADLINK_ADMIN_PASSWORD = Read-Host "请输入强管理员密码" -MaskInput
go run ./cmd/server
```

Bash：

```bash
cd backend
export GRADLINK_ADDR="127.0.0.1:8080"
read -rsp "请输入强管理员密码: " GRADLINK_ADMIN_PASSWORD
echo
export GRADLINK_ADMIN_PASSWORD
go run ./cmd/server
```

默认监听：

- `http://localhost:8080`
- 健康检查：`http://localhost:8080/api/health`

`GRADLINK_ADDR` 的默认值为 `:8080`，会监听所有网络接口。本地仅自己使用时，建议如上例设置为 `127.0.0.1:8080`。

默认 SQLite 文件：

- `backend/data/gradlink.db`

说明：

- 新建数据库且不设置 `GRADLINK_ADMIN_PASSWORD` 时，种子管理员账号会存在但默认禁用，不能登录。
- 在后续任一启动时设置 `GRADLINK_ADMIN_PASSWORD`，会启用 `admin@gradlink.local` 或轮换其现有密码。
- 管理员一旦启用，后续启动时不设置该变量不会将其禁用，也不会轮换密码。需要在本地禁用或重置管理员时，请按下方步骤重置数据库。
- 本地 HTTP 开发时保持 `GRADLINK_COOKIE_SECURE=false`，否则浏览器不会在非 HTTPS 下保存安全 Cookie。

## 启动前端

```bash
cd frontend
npm install
npm run dev
```

默认前端地址：

- `http://localhost:5173`

### 手机真机测试

电脑和手机连接同一个 Wi-Fi 后，在前端目录启动局域网开发服务：

```powershell
npm run dev -- --host 0.0.0.0
```

然后在手机浏览器访问 `http://<电脑的局域网IP>:5173`。仅应在可信任的局域网中使用此方式，避免将开发服务暴露到公共网络。

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
go test ./... -count=1
```

前端：

```bash
cd frontend
npm test -- --run
npm run build
```
