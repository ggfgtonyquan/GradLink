# GradLink Open Source Release Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Produce a clean `main` branch that can be pushed directly to `github.com/ggfgtonyquan/GradLink` as a public MIT-licensed repository without exposing the original commit email or known reachable dependency vulnerabilities.

**Architecture:** Keep the application architecture unchanged. Add repository-level licensing, community health, CI, and dependency-maintenance files; make the one user-facing credential cleanup under a regression test; update vulnerable dependency locks; then rewrite all reachable Git commit metadata and prune the original objects.

**Tech Stack:** Go 1.26.6, Gin, GORM, SQLite, Vue 3, Vite, TypeScript, Vitest, GitHub Actions, Dependabot, Git.

---

## File map

- `LICENSE`: MIT grant and Tony Quan copyright notice.
- `README.md`: public project overview, portable setup, security warning, and contributor links.
- `docs/run-local.md`: portable supported toolchain and safe configuration examples.
- `frontend/src/views/LoginView.vue`: empty login form defaults.
- `frontend/src/views/__tests__/LoginView.test.ts`: regression test for empty credential fields.
- `backend/go.mod`, `backend/go.sum`: patched Go toolchain and module graph.
- `frontend/package-lock.json`: patched npm dependency graph.
- `CONTRIBUTING.md`: contributor workflow and validation requirements.
- `SECURITY.md`: supported-version policy and private reporting route.
- `CODE_OF_CONDUCT.md`: Contributor Covenant 2.1 and private enforcement route.
- `.github/workflows/ci.yml`: backend and frontend verification.
- `.github/dependabot.yml`: weekly Go, npm, and Actions dependency updates.
- `.github/ISSUE_TEMPLATE/*.yml`: structured bug, feature, and security routing.
- `.github/PULL_REQUEST_TEMPLATE.md`: pull request checklist.

### Task 1: Add the MIT license

**Files:**
- Create: `LICENSE`

- [ ] **Step 1: Create the exact MIT license**

```text
MIT License

Copyright (c) 2026 Tony Quan

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

- [ ] **Step 2: Verify the copyright identity and license detection text**

Run:

```powershell
rg -n "MIT License|Copyright \(c\) 2026 Tony Quan" LICENSE
```

Expected: both lines are present, and `git diff --check` produces no output.

- [ ] **Step 3: Commit the license**

```powershell
git add LICENSE
git commit -m "docs: add MIT license"
```

Expected: the new commit uses `136245549+ggfgtonyquan@users.noreply.github.com`.

### Task 2: Remove prefilled login credentials with TDD

**Files:**
- Create: `frontend/src/views/__tests__/LoginView.test.ts`
- Modify: `frontend/src/views/LoginView.vue:27-28`

- [ ] **Step 1: Write the failing regression test**

```ts
import { render, screen } from '@testing-library/vue'
import { createPinia } from 'pinia'
import { createMemoryHistory, createRouter } from 'vue-router'
import LoginView from '../LoginView.vue'

test('starts with empty credential fields', async () => {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: LoginView },
      { path: '/register', component: { template: '<div />' } },
    ],
  })
  await router.push('/')
  await router.isReady()

  render(LoginView, {
    global: {
      plugins: [createPinia(), router],
    },
  })

  expect((screen.getByLabelText('邮箱') as HTMLInputElement).value).toBe('')
  expect((screen.getByLabelText('密码') as HTMLInputElement).value).toBe('')
})
```

- [ ] **Step 2: Run the focused test and observe the existing behavior fail**

Run:

```powershell
cd frontend
npm test -- --run src/views/__tests__/LoginView.test.ts
```

Expected: FAIL because the current values are `student@example.com` and `password`.

- [ ] **Step 3: Make both form defaults empty**

Replace the two initializers with:

```ts
const email = ref('')
const password = ref('')
```

- [ ] **Step 4: Run the focused and full frontend test suites**

Run:

```powershell
npm test -- --run src/views/__tests__/LoginView.test.ts
npm test -- --run
```

Expected: the focused test passes and the full suite has 9 passing files and 15 passing tests.

- [ ] **Step 5: Commit the credential cleanup**

```powershell
git add frontend/src/views/LoginView.vue frontend/src/views/__tests__/LoginView.test.ts
git commit -m "fix: clear default login credentials"
```

### Task 3: Make public documentation portable and safe

**Files:**
- Modify: `README.md`
- Modify: `docs/run-local.md`

- [ ] **Step 1: Replace `README.md` with the public-facing content**

````markdown
# GradLink

[![CI](https://github.com/ggfgtonyquan/GradLink/actions/workflows/ci.yml/badge.svg)](https://github.com/ggfgtonyquan/GradLink/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

GradLink 是一个面向国内应届生的移动端优先求职交流网站。MVP 包含邮箱密码登录、应届生资料门槛、帖子、评论、回复、点赞、收藏、举报、关键词审核和轻量管理后台。

> [!IMPORTANT]
> 当前项目处于 MVP 阶段，适合本地开发、学习和功能验证。公开部署前请配置 HTTPS、强管理员密码、安全 Cookie、备份、限流及完整的内容与隐私治理措施。

## 功能

- 邮箱密码注册、登录和数据库会话。
- 应届生资料完善门槛。
- 帖子、评论、回复、点赞和收藏。
- 标签筛选、用户主页和活动记录。
- 内容举报、关键词审核和管理员处理流程。
- 移动端优先的 Vue 单页界面。

## 技术栈

- 后端：Go、Gin、GORM、SQLite。
- 前端：Vue 3、Vite、TypeScript、Pinia、Vue Router。
- 测试：Go testing、Vitest、Testing Library。

## 项目结构

```text
backend/   Go REST API 和领域模块
frontend/  Vue 3 Web 前端
docs/      产品设计、实施计划和运行说明
```

## 环境要求

- Go 1.26.6 或更高的兼容版本。
- Node.js 20.19+、22.12+ 或 24+（推荐当前 LTS）。
- npm。

## 快速启动

后端：

```powershell
cd backend
$env:GRADLINK_ADDR="127.0.0.1:8080"
Remove-Item Env:GRADLINK_ADMIN_PASSWORD -ErrorAction SilentlyContinue
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

管理员账号为 `admin@gradlink.local`。新建数据库且未设置 `GRADLINK_ADMIN_PASSWORD` 时，该账号保持禁用；在后续启动时设置变量会启用该账号或轮换其现有密码。管理员一旦启用，后续启动时取消设置该变量不会将其禁用或轮换密码；本地禁用或重置管理员需要重置本地数据库。

完整的配置、局域网访问和数据重置说明见 [本地运行指南](docs/run-local.md)。

## 配置

| 环境变量 | 默认值 | 说明 |
| --- | --- | --- |
| `GRADLINK_ADDR` | `:8080` | 默认监听所有网络接口；本地开发建议设为 `127.0.0.1:8080` |
| `GRADLINK_DB_PATH` | `data/gradlink.db` | SQLite 数据库路径 |
| `GRADLINK_ADMIN_PASSWORD` | 未设置 | 新建数据库时写入禁用的管理员；在后续启动时设置该变量会启用管理员或轮换密码 |
| `GRADLINK_COOKIE_SECURE` | `false` | HTTPS 生产环境必须设为 `true` |

## 验证

```powershell
cd backend
go test ./... -count=1
go run golang.org/x/vuln/cmd/govulncheck@latest ./...

cd ..\frontend
npm test -- --run
npm run build
npm audit --audit-level=high
```

## 贡献与安全

- 贡献代码前请阅读 [CONTRIBUTING.md](CONTRIBUTING.md)。
- 安全漏洞请按照 [SECURITY.md](SECURITY.md) 私密报告，不要创建公开 Issue。
- 参与社区表示同意遵守 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)。

## 许可证

本项目采用 [MIT License](LICENSE)。
````

- [ ] **Step 2: Update the supported Go version and safe examples in `docs/run-local.md`**

Apply these exact replacements:

```diff
-- Go 1.24 或更高版本。
+- Go 1.26.6 或更高的兼容版本。
```

```diff
+$env:GRADLINK_ADMIN_PASSWORD = Read-Host "请输入强管理员密码" -MaskInput
```

After the existing front-end address section, add:

```markdown
### 手机真机测试

电脑和手机连接同一个局域网后运行：

```bash
npm run dev -- --host 0.0.0.0
```

然后在手机浏览器访问 `http://<电脑的局域网IP>:5173`。请只在可信任的本地网络中使用该开发模式。
```

- [ ] **Step 3: Scan public documentation for machine-specific or weak examples**

Run:

```powershell
rg -n "C:\\Users\\|172\.20\.10\.10|fixed administrator password" README.md docs/run-local.md
```

Expected: no output.

- [ ] **Step 4: Check Markdown and commit**

Run:

```powershell
git diff --check
git add README.md docs/run-local.md
git commit -m "docs: prepare public project guide"
```

Expected: no whitespace errors and one documentation commit.

### Task 4: Patch Go and npm dependency vulnerabilities

**Files:**
- Modify: `backend/go.mod`
- Modify: `backend/go.sum`
- Modify: `frontend/package-lock.json`
- Modify only if required by npm resolution: `frontend/package.json`

- [ ] **Step 1: Capture the failing security baselines**

Run:

```powershell
cd backend
go run golang.org/x/vuln/cmd/govulncheck@latest ./...

cd ..\frontend
npm audit --audit-level=high
```

Expected before upgrades: `govulncheck` reports reachable standard-library, `x/text`, and `quic-go` vulnerabilities; npm reports high-severity `postcss` and `nanoid` advisories.

- [ ] **Step 2: Upgrade the Go toolchain requirement and vulnerable modules**

Run from `backend/` with a Go toolchain that supports automatic toolchain download:

```powershell
go mod edit -go=1.26.6
go get golang.org/x/text@v0.39.0
go get github.com/quic-go/quic-go@v0.59.1
go mod tidy
```

Expected: `go.mod` declares `go 1.26.6`; selected `x/text` and `quic-go` versions meet or exceed the fixed versions without an unrelated application architecture change.

- [ ] **Step 3: Verify the Go update**

Run:

```powershell
go test ./... -count=1
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

Expected: all packages pass and `govulncheck` prints `No vulnerabilities found.`

- [ ] **Step 4: Apply compatible npm audit fixes**

Run from `frontend/`:

```powershell
npm audit fix
npm ls nanoid postcss --all
```

Expected: the resolved `nanoid` and `postcss` versions no longer match the reported advisory ranges. Do not use `--force`.

- [ ] **Step 5: Verify the npm update**

Run:

```powershell
npm test -- --run
npm run build
npm audit --audit-level=high
```

Expected: 9 test files and 15 tests pass, the production build succeeds, and npm reports zero vulnerabilities at high severity or above.

- [ ] **Step 6: Review and commit only dependency files**

Run:

```powershell
git diff -- backend/go.mod backend/go.sum frontend/package.json frontend/package-lock.json
git add backend/go.mod backend/go.sum frontend/package.json frontend/package-lock.json
git commit -m "fix: update vulnerable dependencies"
```

Expected: dependency metadata changes only; `git add` safely ignores `frontend/package.json` if it did not change.

### Task 5: Add contribution and community policies

**Files:**
- Create: `CONTRIBUTING.md`
- Create: `SECURITY.md`
- Create: `CODE_OF_CONDUCT.md`

- [ ] **Step 1: Create `CONTRIBUTING.md`**

```markdown
# 贡献指南

感谢你为 GradLink 做出贡献。

## 开始之前

1. 搜索现有 Issue，确认问题或建议尚未被记录。
2. 较大的功能请先创建 Issue，说明目标、范围和方案。
3. 安全漏洞不要创建公开 Issue，请按照 [SECURITY.md](SECURITY.md) 私密报告。

## 本地开发

环境和启动方式见 [README.md](README.md) 与 [docs/run-local.md](docs/run-local.md)。

创建分支时使用简短、明确的名称，例如 `fix/login-error` 或 `feat/post-search`。不要直接向 `main` 推送功能改动。

## 代码规范

- Go 代码使用 `gofmt`，继续沿用领域模块及 handler/service/repository 分层。
- 前端使用 Vue 单文件组件和 TypeScript；组件使用 PascalCase，变量和函数使用 camelCase。
- 面向用户的 MVP 文案保持中文。
- 修复 bug 或新增行为时添加回归测试。

## 提交与验证

提交信息使用 Conventional Commits，例如：

```text
feat: add post search
fix: handle expired sessions
docs: clarify local setup
```

提交 PR 前运行：

```powershell
cd backend
go test ./... -count=1
go run golang.org/x/vuln/cmd/govulncheck@latest ./...

cd ..\frontend
npm test -- --run
npm run build
npm audit --audit-level=high
```

## Pull Request

PR 应保持单一范围，并包含：

- 变更目的和实现摘要。
- 已运行的验证命令。
- 关联 Issue 或背景。
- UI 变更的截图或录屏。

提交贡献即表示你同意按本仓库的 MIT License 授权该贡献，并遵守 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)。
```

- [ ] **Step 2: Create `SECURITY.md`**

```markdown
# 安全策略

## 支持版本

GradLink 当前处于 MVP 阶段。安全修复只应用于 `main` 分支的最新版本；旧提交和个人 Fork 不提供安全支持。

## 报告漏洞

请不要通过公开 Issue、Discussion 或 Pull Request 报告漏洞。

请使用仓库 **Security** 页面中的 **Report a vulnerability** 私密提交报告：

https://github.com/ggfgtonyquan/GradLink/security/advisories/new

报告应包含：

- 受影响的提交或版本。
- 漏洞影响和可行的攻击场景。
- 最小复现步骤或验证代码。
- 已知缓解措施。

维护者会尽快确认报告，在修复可用前请避免公开细节。若 GitHub 私密报告功能暂不可用，请先通过维护者的 GitHub 主页联系，不要公开漏洞内容。

## 部署说明

本仓库是 MVP，不代表默认配置适合公开生产环境。运营者需要自行配置 HTTPS、安全 Cookie、强管理员密码、速率限制、日志脱敏、数据备份和适用的隐私及内容治理措施。
```

- [ ] **Step 3: Create the Contributor Covenant 2.1 `CODE_OF_CONDUCT.md`**

```markdown
# Contributor Covenant Code of Conduct

## Our Pledge

We as members, contributors, and leaders pledge to make participation in our
community a harassment-free experience for everyone, regardless of age, body
size, visible or invisible disability, ethnicity, sex characteristics, gender
identity and expression, level of experience, education, socio-economic status,
nationality, personal appearance, race, caste, color, religion, or sexual
identity and orientation.

We pledge to act and interact in ways that contribute to an open, welcoming,
diverse, inclusive, and healthy community.

## Our Standards

Examples of behavior that contributes to a positive environment include:

* Demonstrating empathy and kindness toward other people
* Being respectful of differing opinions, viewpoints, and experiences
* Giving and gracefully accepting constructive feedback
* Accepting responsibility and apologizing to those affected by our mistakes,
  and learning from the experience
* Focusing on what is best not just for us as individuals, but for the overall
  community

Examples of unacceptable behavior include:

* The use of sexualized language or imagery, and sexual attention or advances
* Trolling, insulting or derogatory comments, and personal or political attacks
* Public or private harassment
* Publishing others' private information without their explicit permission
* Other conduct which could reasonably be considered inappropriate in a
  professional setting

## Enforcement Responsibilities

Community leaders are responsible for clarifying and enforcing our standards of
acceptable behavior and will take appropriate and fair corrective action in
response to behavior that they deem inappropriate, threatening, offensive, or
harmful.

Community leaders have the right and responsibility to remove, edit, or reject
comments, commits, code, wiki edits, issues, and other contributions that are not
aligned to this Code of Conduct, and will communicate reasons for moderation
decisions when appropriate.

## Scope

This Code of Conduct applies within all community spaces, and also applies when
an individual is officially representing the community in public spaces.

## Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be
reported privately through the repository's **Security** page using **Report a
vulnerability** and prefixing the report title with `[Conduct]`:

https://github.com/ggfgtonyquan/GradLink/security/advisories/new

All complaints will be reviewed and investigated promptly and fairly. Community
leaders are obligated to respect the privacy and security of the reporter.

## Enforcement Guidelines

Community leaders will follow these Community Impact Guidelines in determining
the consequences for any action they deem in violation of this Code of Conduct:

### 1. Correction

**Community Impact**: Use of inappropriate language or other behavior deemed
unprofessional or unwelcome in the community.

**Consequence**: A private, written warning providing clarity around the nature
of the violation and an explanation of why the behavior was inappropriate.

### 2. Warning

**Community Impact**: A violation through a single incident or series of actions.

**Consequence**: A warning with consequences for continued behavior. No
interaction with the people involved for a specified period of time.

### 3. Temporary Ban

**Community Impact**: A serious violation of community standards, including
sustained inappropriate behavior.

**Consequence**: A temporary ban from any sort of interaction or public
communication with the community for a specified period of time.

### 4. Permanent Ban

**Community Impact**: Demonstrating a pattern of violation of community
standards, including sustained inappropriate behavior, harassment of an
individual, or aggression toward or disparagement of classes of individuals.

**Consequence**: A permanent ban from any sort of public interaction within the
community.

## Attribution

This Code of Conduct is adapted from the [Contributor Covenant][homepage],
version 2.1, available at
[https://www.contributor-covenant.org/version/2/1/code_of_conduct.html][version].

[homepage]: https://www.contributor-covenant.org
[version]: https://www.contributor-covenant.org/version/2/1/code_of_conduct.html
```

- [ ] **Step 4: Validate policy links and commit**

Run:

```powershell
rg -n "CONTRIBUTING|SECURITY|CODE_OF_CONDUCT" README.md CONTRIBUTING.md SECURITY.md CODE_OF_CONDUCT.md
git diff --check
git add CONTRIBUTING.md SECURITY.md CODE_OF_CONDUCT.md
git commit -m "docs: add community health policies"
```

Expected: all documents cross-link correctly and the commit succeeds.

### Task 6: Add GitHub CI, Dependabot, and contribution templates

**Files:**
- Create: `.github/workflows/ci.yml`
- Create: `.github/dependabot.yml`
- Create: `.github/ISSUE_TEMPLATE/bug_report.yml`
- Create: `.github/ISSUE_TEMPLATE/feature_request.yml`
- Create: `.github/ISSUE_TEMPLATE/config.yml`
- Create: `.github/PULL_REQUEST_TEMPLATE.md`

- [ ] **Step 1: Create `.github/workflows/ci.yml`**

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

permissions:
  contents: read

jobs:
  backend:
    name: Backend
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: backend
    steps:
      - name: Check out repository
        uses: actions/checkout@v7
      - name: Set up Go
        uses: actions/setup-go@v7
        with:
          go-version-file: backend/go.mod
          cache-dependency-path: backend/go.sum
      - name: Test
        run: go test ./... -count=1
      - name: Install govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@latest
      - name: Scan reachable Go vulnerabilities
        run: govulncheck ./...

  frontend:
    name: Frontend
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: frontend
    steps:
      - name: Check out repository
        uses: actions/checkout@v7
      - name: Set up Node.js
        uses: actions/setup-node@v7
        with:
          node-version: 20
          cache: npm
          cache-dependency-path: frontend/package-lock.json
      - name: Install dependencies
        run: npm ci
      - name: Test
        run: npm test -- --run
      - name: Build
        run: npm run build
      - name: Audit high-severity dependencies
        run: npm audit --audit-level=high
```

- [ ] **Step 2: Create `.github/dependabot.yml`**

```yaml
version: 2
updates:
  - package-ecosystem: gomod
    directory: /backend
    schedule:
      interval: weekly
      day: monday
    open-pull-requests-limit: 5
  - package-ecosystem: npm
    directory: /frontend
    schedule:
      interval: weekly
      day: monday
    open-pull-requests-limit: 5
  - package-ecosystem: github-actions
    directory: /
    schedule:
      interval: weekly
      day: monday
    open-pull-requests-limit: 5
```

- [ ] **Step 3: Create `.github/ISSUE_TEMPLATE/bug_report.yml`**

```yaml
name: Bug 报告
description: 报告可以稳定复现的问题
title: "[Bug]: "
labels: []
body:
  - type: markdown
    attributes:
      value: 感谢反馈。安全漏洞请使用 Security 页面的私密报告入口。
  - type: textarea
    id: description
    attributes:
      label: 问题描述
      description: 实际发生了什么？预期行为是什么？
    validations:
      required: true
  - type: textarea
    id: reproduction
    attributes:
      label: 复现步骤
      description: 请提供最小、完整、可重复的步骤。
    validations:
      required: true
  - type: input
    id: commit
    attributes:
      label: GradLink 版本或提交
      placeholder: main 或提交哈希
    validations:
      required: true
  - type: textarea
    id: environment
    attributes:
      label: 环境
      description: 操作系统、Go、Node.js、npm 和浏览器版本。
    validations:
      required: true
  - type: textarea
    id: logs
    attributes:
      label: 日志或截图
      description: 请先删除密码、令牌、Cookie 和个人数据。
  - type: checkboxes
    id: checks
    attributes:
      label: 提交前确认
      options:
        - label: 我已搜索现有 Issue，并确认这不是安全漏洞。
          required: true
```

- [ ] **Step 4: Create `.github/ISSUE_TEMPLATE/feature_request.yml`**

```yaml
name: 功能建议
description: 建议新的能力或体验改进
title: "[Feature]: "
labels: []
body:
  - type: textarea
    id: problem
    attributes:
      label: 问题背景
      description: 当前什么场景不方便，谁会受到影响？
    validations:
      required: true
  - type: textarea
    id: solution
    attributes:
      label: 建议方案
      description: 描述期望行为，不必限定具体实现。
    validations:
      required: true
  - type: textarea
    id: alternatives
    attributes:
      label: 备选方案
      description: 已考虑过哪些替代方法？
  - type: checkboxes
    id: checks
    attributes:
      label: 提交前确认
      options:
        - label: 我已搜索现有 Issue，且该建议符合 GradLink 的求职交流定位。
          required: true
```

- [ ] **Step 5: Create `.github/ISSUE_TEMPLATE/config.yml`**

```yaml
blank_issues_enabled: false
contact_links:
  - name: 私密安全漏洞报告
    url: https://github.com/ggfgtonyquan/GradLink/security/advisories/new
    about: 请勿公开安全漏洞，使用 GitHub Private Vulnerability Reporting。
```

- [ ] **Step 6: Create `.github/PULL_REQUEST_TEMPLATE.md`**

```markdown
## 变更说明

<!-- 说明为什么需要这项变更，以及采用了什么方案。 -->

## 验证

- [ ] `cd backend; go test ./... -count=1`
- [ ] `cd backend; govulncheck ./...`
- [ ] `cd frontend; npm test -- --run`
- [ ] `cd frontend; npm run build`
- [ ] `cd frontend; npm audit --audit-level=high`

## 关联内容

<!-- 填写 Fixes #123；没有关联 Issue 时说明背景。 -->

## UI 证据

<!-- UI 变更附截图或录屏；不涉及 UI 时填写“不适用”。 -->
```

- [ ] **Step 7: Validate YAML and repository whitespace**

Run from the repository root with the bundled Python environment containing PyYAML:

```powershell
python -c "from pathlib import Path; import yaml; files=list(Path('.github').rglob('*.yml')); [yaml.safe_load(p.read_text(encoding='utf-8')) for p in files]; print(f'validated {len(files)} YAML files')"
git diff --check
```

Expected: `validated 5 YAML files` and no whitespace errors.

- [ ] **Step 8: Commit GitHub configuration**

```powershell
git add .github
git commit -m "ci: add repository automation"
```

### Task 7: Run the complete release verification

**Files:**
- Verify only; modify files only to correct a failure attributable to Tasks 1-6.

- [ ] **Step 1: Run backend tests and vulnerability scanning**

```powershell
cd backend
go test ./... -count=1
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

Expected: all backend packages pass and no reachable vulnerabilities are reported.

- [ ] **Step 2: Run frontend tests, build, and audit**

```powershell
cd frontend
npm test -- --run
npm run build
npm audit --audit-level=high
```

Expected: 9 test files and 15 tests pass, Vite builds successfully, and no high-severity audit issue remains.

- [ ] **Step 3: Run public-content scans**

```powershell
cd ..
rg -n "C:\\Users\\|172\.20\.10\.10|fixed administrator password" README.md docs/run-local.md
rg -n "student@example\.com|ref\('password'\)" frontend/src/views/LoginView.vue
git diff --check
```

Expected: both ripgrep commands and `git diff --check` produce no output.

- [ ] **Step 4: Verify required release files and clean state**

```powershell
$required = @(
  'LICENSE', 'CONTRIBUTING.md', 'SECURITY.md', 'CODE_OF_CONDUCT.md',
  '.github/workflows/ci.yml', '.github/dependabot.yml',
  '.github/ISSUE_TEMPLATE/bug_report.yml',
  '.github/ISSUE_TEMPLATE/feature_request.yml',
  '.github/ISSUE_TEMPLATE/config.yml',
  '.github/PULL_REQUEST_TEMPLATE.md'
)
$missing = $required | Where-Object { -not (Test-Path -LiteralPath $_) }
if ($missing) { throw "Missing release files: $($missing -join ', ')" }
git status --short --branch
git remote -v
```

Expected: no missing files, a clean `main` branch, and no remote output.

### Task 8: Rewrite commit emails and rename the release branch

**Files:**
- Rewrite: Git objects and references under `.git`
- Rename the current release branch to `main`

- [ ] **Step 1: Reconfirm the exact destructive scope**

Run:

```powershell
$resolvedRepo = (Resolve-Path -LiteralPath $env:USERPROFILE).Path
if (-not $resolvedRepo) { throw 'Unable to resolve the current user profile path.' }
if (git status --porcelain) { throw 'Working tree must be clean before history rewrite.' }
if (git remote) { throw 'A remote exists; stop before rewriting published history.' }
git branch --show-current
git log --all --format='%ae%n%ce' | Sort-Object -Unique
```

Expected: the current user profile path resolves, the release branch is current, and the only historical emails are the old source address and the new noreply address.

- [ ] **Step 2: Rewrite matching author and committer emails across all refs**

Run in PowerShell; `git filter-branch` executes the single-quoted filter body through Git's shell:

```powershell
$env:FILTER_BRANCH_SQUELCH_WARNING = '1'
$env:NEW_EMAIL = '136245549+ggfgtonyquan@users.noreply.github.com'
$existingEmails = @(git log --all --format='%ae%n%ce' | Sort-Object -Unique)
$sourceEmails = @($existingEmails | Where-Object { $_ -ne $env:NEW_EMAIL })
if ($sourceEmails.Count -ne 1) { throw "Expected exactly one source email, found $($sourceEmails.Count)." }
$env:OLD_EMAIL = $sourceEmails[0]
git filter-branch --force --env-filter '
if [ "$GIT_AUTHOR_EMAIL" = "$OLD_EMAIL" ]; then
  export GIT_AUTHOR_EMAIL="$NEW_EMAIL"
fi
if [ "$GIT_COMMITTER_EMAIL" = "$OLD_EMAIL" ]; then
  export GIT_COMMITTER_EMAIL="$NEW_EMAIL"
fi
' --tag-name-filter cat -- --all
if ($LASTEXITCODE -ne 0) { throw 'Git history rewrite failed.' }
```

Expected: every reachable commit is rewritten while merge topology, messages, names, and dates remain intact.

- [ ] **Step 3: Remove filter-branch backup refs and prune original objects**

Run from the already validated repository path:

```powershell
$originalRefs = @(git for-each-ref --format='%(refname)' refs/original/)
foreach ($ref in $originalRefs) {
  git update-ref -d $ref
  if ($LASTEXITCODE -ne 0) { throw "Failed to delete backup ref: $ref" }
}
git reflog expire --expire=now --all
git gc --prune=now
```

Expected: no `refs/original/*` remain and unreachable original commits are pruned.

- [ ] **Step 4: Rename the branch and verify privacy invariants**

```powershell
git branch -m main
$emails = @(git log --all --format='%ae%n%ce' | Sort-Object -Unique)
if ($emails.Count -ne 1 -or $emails[0] -ne '136245549+ggfgtonyquan@users.noreply.github.com') {
  throw "Unexpected history emails: $($emails -join ', ')"
}
if (git for-each-ref refs/original/) { throw 'filter-branch backup refs remain.' }
if (rg -n "@qq\.com" --hidden -g '!.git/**' .) { throw 'A QQ email remains in tracked content.' }
git fsck --full --unreachable --no-reflogs
git status --short --branch
git remote -v
```

Expected: `git fsck` reports no unreachable old commits, status is clean on `main`, and no remote is configured.

- [ ] **Step 5: Record final push instructions without pushing**

The final handoff must state that the user can create an empty public repository named `GradLink` and run:

```powershell
git remote add origin git@github.com:ggfgtonyquan/GradLink.git
git push -u origin main
```

It must also remind the user to enable Private Vulnerability Reporting, Dependabot alerts, secret scanning, push protection, code scanning, and `main` branch protection after the first push.
