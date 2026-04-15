# Skills 使用场景与最佳实践

## 核心使用场景

### 场景 1：为 Agent 定义特定领域的工作流

**问题**：一个 AI Agent 需要按照团队特定的代码审查流程来工作，而不是通用的代码审查。

**解决方案**：创建一个 "Code Review" Skill，定义审查检查清单、关注重点、输出格式。

```
SKILL.md 内容示例：

# Code Review Expert

## When to Use
When you are assigned to review code changes on an issue.

## Review Checklist
1. Check for security vulnerabilities (OWASP Top 10)
2. Verify error handling completeness
3. Check test coverage for new code
4. Review naming conventions against team style guide

## Output Format
Post a comment with:
- Summary of changes
- Issues found (Critical/Warning/Info)
- Suggested improvements
```

### 场景 2：为不同 Agent 分配不同专长

**问题**：团队有多个 Agent，每个负责不同的工作领域。

**解决方案**：创建多个 Skills，按需分配给不同 Agent。

```
Agent "Frontend Dev" → Skills: [React Patterns, CSS Architecture, Accessibility]
Agent "Backend Dev"  → Skills: [Go Patterns, Database Design, API Security]
Agent "Reviewer"     → Skills: [Code Review, Performance Audit, Security Audit]
```

Skills 的多对多关系允许：
- 同一个 Skill 分配给多个 Agent（如 "Git Workflow"）
- 同一个 Agent 拥有多个 Skills（组合能力）

### 场景 3：标准化团队的编码规范

**问题**：AI Agent 生成的代码风格不统一，需要遵循团队编码规范。

**解决方案**：创建编码规范 Skill，附带模板文件。

```
Skill: "Go Coding Standards"
├── SKILL.md              # 编码规范主体
├── templates/
│   ├── handler.go        # 标准 handler 模板
│   └── service.go        # 标准 service 模板
└── examples/
    └── error-handling.go # 错误处理示例
```

Agent 在执行任务时，会在工作目录中看到这些模板和示例，并据此生成代码。

### 场景 4：复用社区 Skills

**问题**：不想从零编写 Skill，想使用已有的社区最佳实践。

**解决方案**：从 ClawHub 或 GitHub 导入 Skills。

```
# 从 ClawHub 导入
multica skill import --url https://clawhub.ai/skills/code-review

# 从 GitHub 仓库导入
multica skill import --url https://skills.sh/anthropics/skills/frontend-design
```

导入后可根据团队需要修改内容。

### 场景 5：为特定项目定制 Agent 行为

**问题**：同一个 Agent 在不同项目中需要遵循不同的约定。

**解决方案**：创建项目级 Skills，按任务分配。

例如，一个 Agent 在处理前端任务时使用 "React + Tailwind" Skill，处理后端任务时使用 "Go + Chi Router" Skill。

## 操作指南

### 如何创建一个 Skill

**通过 UI**：

1. 打开 Skills 页面（侧边栏 → Skills）
2. 点击右上角 "Create" 按钮
3. 在 Create 标签页中填写名称和描述
4. 创建后，在右侧详情面板编辑 SKILL.md 内容
5. 可选：添加辅助文件（模板、示例等）
6. 点击 Save 保存

**通过 CLI**：

```bash
# 创建 Skill
multica skill create --name "Code Review" --description "Standard code review workflow"

# 添加 SKILL.md 内容（通过文件）
multica skill files upsert <skill-id> --path "SKILL.md" --content "$(cat skill.md)"

# 添加辅助文件
multica skill files upsert <skill-id> --path "templates/checklist.md" --content "$(cat checklist.md)"
```

### 如何将 Skill 分配给 Agent

**通过 UI**：

1. 打开 Agent 详情页
2. 切换到 "Skills" 标签页
3. 点击 "Add Skill" 按钮
4. 从弹出的选择器中选择未分配的 Skill
5. 确认添加

**通过 CLI**：

```bash
# 查看 Agent 当前 Skills
multica agent skills list <agent-id>

# 设置 Agent 的 Skills（全量替换）
# 注意：这是全量操作，会替换所有已有的 Skill 分配
multica agent skills set <agent-id> --skill-ids <id1>,<id2>,<id3>
```

### 如何从外部导入 Skill

**通过 UI**：

1. 在 Skills 页面点击 "Create"
2. 切换到 "Import" 标签页
3. 输入 ClawHub 或 GitHub URL
4. 点击 "Import"

**通过 CLI**：

```bash
# 从 ClawHub 导入
multica skill import --url https://clawhub.ai/skills/your-skill

# 从 GitHub 导入
multica skill import --url https://skills.sh/owner/repo/skill-name
```

## SKILL.md 编写最佳实践

### 推荐结构

```markdown
# Skill 名称

## When to Use
描述什么场景下应该激活这个 Skill。
这帮助 Agent 判断当前任务是否适用此 Skill。

## How It Works
描述 Skill 的工作流程和步骤。
Agent 会按步骤执行。

## Rules / Guidelines
列出具体的规则和约束。
Agent 必须遵守这些规则。

## Examples
提供输入/输出示例。
Agent 会参照示例的格式和风格。

## References
指向辅助文件或其他资源。
```

### 编写原则

1. **明确触发条件**："When to Use" 帮助 Agent 在多个 Skills 中选择正确的
2. **步骤化工作流**：用编号列表描述步骤，Agent 按顺序执行
3. **具体规则**：避免模糊描述，用明确的 do/don't 列表
4. **包含示例**：Agent 会模仿示例的格式和风格
5. **善用辅助文件**：长模板、配置示例放在辅助文件中，保持 SKILL.md 简洁

### 示例：一个完整的 Code Review Skill

```markdown
# Code Review Expert

## When to Use
Use this skill when you are asked to review code changes, perform a PR review,
or audit code quality on an issue.

## How It Works
1. Run `multica issue get <id> --output json` to understand the task
2. Review all changed files in the working directory
3. Check each file against the review checklist below
4. Post a review comment summarizing findings
5. Update issue status based on review outcome

## Review Checklist
- [ ] Security: No SQL injection, XSS, or credential exposure
- [ ] Error handling: All error paths return meaningful messages
- [ ] Tests: New code has corresponding test coverage
- [ ] Naming: Follows team naming conventions
- [ ] Performance: No obvious N+1 queries or memory leaks
- [ ] Documentation: Public APIs are documented

## Output Format
Post a comment with this structure:

### Summary
Brief description of changes reviewed.

### Issues Found
- 🔴 **Critical**: [description] in `file:line`
- 🟡 **Warning**: [description] in `file:line`
- 🔵 **Info**: [suggestion] in `file:line`

### Verdict
APPROVE / REQUEST_CHANGES

## Examples
See `examples/review-output.md` for a sample review comment.
```

对应的辅助文件：

```
Code Review Expert/
├── SKILL.md
├── examples/
│   └── review-output.md      # 审查输出示例
└── templates/
    └── checklist.json         # 结构化检查清单
```

## 常见问题

### Q: Skill 和 Agent Instructions 有什么区别？

**Agent Instructions** 是 Agent 的"身份卡"，定义 Agent 是谁、基本行为准则。每个 Agent 只有一份。

**Skill** 是 Agent 的"能力模块"，定义在特定场景下如何行动。一个 Agent 可以拥有多个 Skills，一个 Skill 可以分配给多个 Agents。

### Q: 本地运行时的 Skills 怎么处理？

本地运行时（Local Daemon）的 Skills 会自动从 Agent 配置中加载。Agent 详情页的 Skills Tab 有提示："Local runtime skills are automatically available and don't need to be assigned."

### Q: 一个 Skill 应该多大？

推荐一个 Skill 聚焦一个独立的工作流。如果一个 Skill 超过 200 行 SKILL.md，考虑拆分成多个 Skills。

### Q: 导入的 Skill 可以修改吗？

可以。导入后的 Skill 与手动创建的完全等价，可以自由编辑内容和文件。

### Q: Agent 执行任务时，Skills 什么时候被加载？

Skills 在 Daemon 认领任务（Task Claim）时加载。每次任务执行都会重新加载最新的 Skills 内容，确保使用的是最新版本。

### Q: 多个 Skills 之间有冲突怎么办？

Agent 的运行时配置会列出所有 Skills。如果多个 Skills 的规则冲突，Agent 会根据任务上下文自行判断。建议在 SKILL.md 的 "When to Use" 部分明确区分适用场景，减少冲突。
