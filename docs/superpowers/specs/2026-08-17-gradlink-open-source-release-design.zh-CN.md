# GradLink 开源发布设计

## 目标

将当前 GradLink 仓库整理成可以直接推送到 GitHub 公共仓库的版本：采用 MIT 许可证，不公开真实提交邮箱，修复已知依赖漏洞，并补齐最小但完整的社区协作、自动验证和安全维护配置。

本次工作只准备本地仓库，不创建 GitHub 远程仓库，也不执行公开推送。

## 发布身份与许可证

- 根目录新增 `LICENSE`，使用 MIT License。
- 版权声明使用 `Copyright (c) 2026 Tony Quan`。
- 当前及后续本地提交使用作者名 `tony` 和邮箱 `136245549+ggfgtonyquan@users.noreply.github.com`。
- 改写当前仓库全部历史提交的作者邮箱和提交者邮箱，将现有真实邮箱替换为上述 noreply 邮箱；保留原作者名、提交时间、提交消息和分支拓扑。
- 仓库当前没有远程地址，因此历史改写不会影响已发布的协作者克隆。

## 仓库公开内容

### README 与运行说明

更新 `README.md` 和必要的运行文档，使陌生贡献者无需本机专用配置即可启动项目：

- 删除 `C:\Users\16021\...` 本机工具链路径。
- 将固定局域网 IP 改为通用占位示例。
- 将示例管理员密码改为明确的强密码占位符。
- 增加项目状态、功能、技术栈、配置项、测试、贡献、安全和许可证入口。
- 明确当前版本是 MVP，不应在没有 HTTPS、强管理员密码和生产安全加固的情况下直接部署。

登录页不再预填 `student@example.com` 和 `password`，避免公开版本呈现不存在的演示凭据。

### 社区健康文件

新增以下文件：

- `CONTRIBUTING.md`：说明分支、Conventional Commit、后端/前端验证命令和 PR 要求。
- `SECURITY.md`：说明支持版本与私密漏洞报告方式，不公开个人邮箱。
- `CODE_OF_CONDUCT.md`：采用 Contributor Covenant，私密执法报告通过仓库的 GitHub Private Vulnerability Reporting 渠道提交。
- `.github/ISSUE_TEMPLATE/bug_report.yml`：结构化收集复现步骤、环境和日志。
- `.github/ISSUE_TEMPLATE/feature_request.yml`：结构化收集问题背景和建议方案。
- `.github/ISSUE_TEMPLATE/config.yml`：关闭空白 Issue，并提供安全报告入口。
- `.github/PULL_REQUEST_TEMPLATE.md`：要求变更说明、测试结果、关联问题和 UI 证据。

### 自动化与依赖维护

新增 `.github/workflows/ci.yml`，在 push 和 pull request 时并行执行：

- Go：安装与项目兼容的已修复 Go 工具链，运行 `go test ./... -count=1` 和 `govulncheck ./...`。
- 前端：安装 Node.js 20，使用 `npm ci`，运行 `npm test -- --run`、`npm run build` 和高危依赖审计。

工作流只授予读取仓库内容所需的最小权限，并使用 npm 缓存。

新增 `.github/dependabot.yml`，每周检查 Go Modules、npm 和 GitHub Actions 更新，限制同时打开的更新数量。

## 依赖安全处理

### Go

- 使用受支持且已修复当前标准库漏洞的 Go 补丁版本进行验证。
- 更新 `golang.org/x/text` 至不低于 `v0.39.0` 的兼容版本。
- 更新 `github.com/quic-go/quic-go` 至不低于 `v0.59.1` 的兼容版本。
- 允许 `go get`/`go mod tidy` 为满足兼容性而同步调整相关间接依赖，但不做无关的大版本迁移。

### npm

- 更新锁文件，使 `postcss` 和 `nanoid` 不再命中当前高危公告。
- 保持现有 Vue/Vite 架构和公开 API，不做无关的前端功能升级。

所有依赖变化必须通过原有测试、生产构建、`npm audit` 和 `govulncheck`。

## 历史改写与分支

历史改写在所有内容提交完成后执行，确保新提交也使用 noreply 邮箱。改写后将当前分支从 `master` 重命名为 `main`，并验证：

- 所有可达提交的作者邮箱和提交者邮箱均为 noreply 地址。
- 不再能从任何可达引用搜索到原 QQ 邮箱。
- 工作区干净，`main` 指向完整发布版本。

由于这是破坏性历史操作，执行前再次确认仓库没有远程地址和未提交改动；不会创建包含旧邮箱的仓库内备份引用。

## 验证与完成标准

发布版本必须同时满足：

1. 后端 `go test ./... -count=1` 通过。
2. `govulncheck ./...` 不再报告代码可达漏洞。
3. 前端 `npm test -- --run` 通过。
4. 前端 `npm run build` 通过。
5. 前端高危依赖审计通过。
6. README 中不存在本机绝对路径、固定局域网 IP 或可直接复制的弱管理员密码。
7. Git 历史中不存在原 QQ 邮箱，所有提交使用 GitHub noreply 邮箱。
8. `LICENSE`、社区健康文件、CI 和 Dependabot 配置齐全。
9. 最终分支为 `main`，工作区干净且没有配置远程仓库。

满足以上条件后，用户可以创建一个空 GitHub 公共仓库，添加 `origin` 并直接推送 `main`。
