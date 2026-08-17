# 安全政策

## 支持范围

仅支持最新的 `main` 分支。GradLink 目前是 MVP，不为旧提交、旧版本或分叉项目提供安全支持。

## 报告安全问题

请勿通过公开 Issue、Discussion 或 Pull Request 报告安全问题。请在 GitHub 中进入 **Security > Report a vulnerability**，使用以下私密报告地址：

<https://github.com/ggfgtonyquan/GradLink/security/advisories/new>

报告必须包含：

- 受影响的提交或版本；
- 影响及可能的攻击场景；
- 最小化的复现步骤或验证方式；
- 已知的缓解措施或修复建议。

维护者会确认收到报告，并在处理过程中对报告内容保密；我们会尽快、公平地评估问题，并在适当范围内与报告者沟通修复进展。请避免在问题修复前公开披露细节。

如果无法使用私密报告功能，请通过维护者的 [GitHub 个人主页](https://github.com/ggfgtonyquan) 联系；不要公开发布敏感细节，也不要寻找或公开个人邮箱地址。

## 部署提醒

公开部署前请务必配置 HTTPS、将 Secure Cookie 开启、使用强管理员密码、启用限流、避免在日志中记录或暴露敏感信息，并制定可靠的备份、隐私和内容治理方案。项目默认配置主要服务于本地开发与功能验证。
