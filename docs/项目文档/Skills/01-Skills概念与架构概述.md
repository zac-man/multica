# Skills 概念与架构概述

## 什么是 Skills

Skills 是 Multica 平台中 **可复用的 AI Agent 指令集**。每个 Skill 本质上是一个 Markdown 文档（SKILL.md），加上可选的辅助文件，定义了一组特定的工作流、编码规范、或领域知识。

Skills 的核心设计思想：**将 AI Agent 的能力模块化、可共享、可组合**。

类比理解：
- Skills 类似于人类团队中的"标准操作手册"（SOP）
- 每个 Skill 就像一个可插拔的能力模块，告诉 Agent 在特定场景下如何行动
- Agent 通过被分配多个 Skills 来获得复合能力

## Skills 在系统中的位置

```
┌─────────────────────────────────────────────────┐
│                   Workspace                      │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐       │
│  │  Agent A │   │  Agent B │   │  Agent C │       │
│  │  Skills: │   │  Skills: │   │  Skills: │       │
│  │  - S1    │   │  - S2    │   │  - S1    │       │
│  │  - S2    │   │  - S3    │   │  - S3    │       │
│  └─────────┘   └─────────┘   └─────────┘       │
│       │              │             │              │
│       ▼              ▼             ▼              │
│  ┌──────────────────────────────────────────┐    │
│  │           Skills 池（Workspace 级别）        │    │
│  │   S1: 代码审查 Skill                       │    │
│  │   S2: 提交规范 Skill                       │    │
│  │   S3: 文档生成 Skill                       │    │
│  └──────────────────────────────────────────┘    │
└─────────────────────────────────────────────────┘
```

关键关系：
- **Workspace → Skill**：一对多。Skills 在 Workspace 级别创建和管理，团队共享
- **Skill → SkillFile**：一对多。每个 Skill 可以有多个辅助文件（模板、示例代码等）
- **Agent ↔ Skill**：多对多。通过 `agent_skill` 关联表，一个 Agent 可以拥有多个 Skills，一个 Skill 可以分配给多个 Agents

## Skills 的生命周期

```
创建/导入 → 分配给 Agent → Agent 执行任务时注入 → Agent 运行时读取并遵循
    │            │                │                      │
    ▼            ▼                ▼                      ▼
  工作台管理   Agent 详情页    Daemon 自动写入         Provider 原生发现
  或 CLI       Skills Tab     工作目录                  .claude/skills/
  或导入URL                                           .config/opencode/skills/
```

1. **创建**：在工作台创建新的 Skill，填写名称、描述、SKILL.md 内容和辅助文件
2. **导入**：从外部来源（ClawHub、skills.sh / GitHub）导入已有 Skill
3. **分配**：在 Agent 详情页的 Skills 标签中，将 Workspace 内的 Skills 分配给 Agent
4. **执行时注入**：当 Agent 被分配任务并开始执行时，系统自动将 Skills 写入工作目录
5. **运行时发现**：AI Agent 运行时通过 Provider 原生机制自动发现并加载 Skills

## 演进历史

### Phase 1 — 文本字段（Migration 002，已废弃）

Agent 表上的简单 `skills TEXT` 列，存储纯文本。此列已被 Migration 008 删除。

### Phase 2 — 结构化 Skills（Migration 008，当前版本）

三个独立表实现完整的 Skills 体系：
- `skill` — Skill 实体（名称、描述、SKILL.md 内容、配置）
- `skill_file` — Skill 辅助文件
- `agent_skill` — Agent-Skill 多对多关联

## 与其他概念的区别

| 概念 | 作用 | 作用域 | 管理方式 |
|------|------|--------|----------|
| **Skills** | Agent 能力模块 | Workspace 级别 | Skills 页面、Agent 详情页 |
| **Agent Instructions** | Agent 的基本身份和行为指令 | Agent 级别 | Agent 设置页 |
| **Tools** | Agent 可调用的工具列表 | Agent 级别 | Agent 设置页 |
| **Triggers** | Agent 的触发条件 | Agent 级别 | Agent 设置页 |

Skills 与 Agent Instructions 的关系：
- Agent Instructions 是"你是谁"（Agent Identity），定义 Agent 的基本身份
- Skills 是"你会什么"（Agent Capabilities），定义 Agent 在特定场景下的具体操作流程
- 两者在运行时都会被注入到 Agent 的执行环境中

## 核心文件清单

| 层级 | 文件路径 | 职责 |
|------|----------|------|
| 数据库迁移 | `server/migrations/008_structured_skills.up.sql` | 创建 skills 三表结构 |
| SQL 查询 | `server/pkg/db/queries/skill.sql` | 所有 Skills 相关的 SQL |
| 生成代码 | `server/pkg/db/generated/skill.sql.go` | sqlc 生成的 Go 数据访问代码 |
| API Handler | `server/internal/handler/skill.go` | HTTP API 端点 |
| 任务服务 | `server/internal/service/task.go` | 加载 Agent Skills 到任务上下文 |
| 执行环境 | `server/internal/daemon/execenv/context.go` | 将 Skills 写入工作目录 |
| 运行时配置 | `server/internal/daemon/execenv/runtime_config.go` | 在 CLAUDE.md/AGENTS.md 中引用 Skills |
| 前端类型 | `packages/core/types/agent.ts` | TypeScript 类型定义 |
| API 客户端 | `packages/core/api/client.ts` | 前端 API 调用 |
| React Query | `packages/core/workspace/queries.ts` | 缓存和实时同步 |
| UI 页面 | `packages/views/skills/` | Skills 管理界面 |
| Agent Skills Tab | `packages/views/agents/components/tabs/skills-tab.tsx` | Agent 详情中的 Skills 管理 |
