# Skills 运行时注入机制

## 概述

Skills 的核心价值在于运行时注入——当 Agent 执行任务时，系统自动将分配给该 Agent 的 Skills 写入工作目录，使 AI Agent 能够发现并遵循这些指令。

这是 Skills 系统最关键的数据流路径。

## 完整注入流程

```
┌─────────────────────────────────────────────────────────────────────────┐
│  1. 任务创建                                                            │
│     Issue 被分配给 Agent                                                │
│         │                                                               │
│         ▼                                                               │
│  2. Daemon 认领任务（Task Claim）                                        │
│     Handler 调用 TaskService.LoadAgentSkills()                          │
│     → 查询 agent_skill 表获取 Skill 列表                                │
│     → 查询 skill_file 表获取每个 Skill 的文件                            │
│     → 组装为 AgentSkillData[] 返回                                      │
│         │                                                               │
│         ▼                                                               │
│  3. Daemon 准备执行环境                                                  │
│     convertSkillsForEnv() 转换数据格式                                  │
│     → SkillContextForEnv{Name, Content, Files[]}                       │
│         │                                                               │
│         ▼                                                               │
│  4. 写入工作目录                                                        │
│     execenv.writeContextFiles()                                         │
│     ┌──────────────────────────────────────────────┐                    │
│     │  Claude:  .claude/skills/{name}/SKILL.md     │                    │
│     │  OpenCode: .config/opencode/skills/{name}/   │                    │
│     │  Codex:   {codex-home}/skills/{name}/        │                    │
│     │  Gemini:  .agent_context/skills/{name}/      │                    │
│     │  Default: .agent_context/skills/{name}/      │                    │
│     └──────────────────────────────────────────────┘                    │
│         │                                                               │
│         ▼                                                               │
│  5. 注入运行时配置                                                      │
│     InjectRuntimeConfig() 写入 Provider 原生配置文件                     │
│     ┌──────────────────────────────────────────────┐                    │
│     │  Claude:  CLAUDE.md                          │                    │
│     │  Codex:   AGENTS.md                          │                    │
│     │  OpenCode: AGENTS.md                         │                    │
│     │  OpenClaw: AGENTS.md                         │                    │
│     │  Gemini:  GEMINI.md                          │                    │
│     └──────────────────────────────────────────────┘                    │
│         │                                                               │
│         ▼                                                               │
│  6. Agent 运行时自动发现并加载 Skills                                    │
└─────────────────────────────────────────────────────────────────────────┘
```

## Step 1: 任务服务加载 Skills

**文件**：`server/internal/service/task.go`（约 line 422-452）

```go
func (s *TaskService) LoadAgentSkills(ctx context.Context, agentID uuid.UUID) ([]AgentSkillData, error) {
    // 1. 查询 agent_skill JOIN skill 获取 Agent 拥有的所有 Skills
    skills, err := s.queries.ListAgentSkills(ctx, agentID)

    // 2. 对每个 Skill，查询其所有辅助文件
    var result []AgentSkillData
    for _, skill := range skills {
        files, _ := s.queries.ListSkillFiles(ctx, skill.ID)
        result = append(result, AgentSkillData{
            ID:          skill.ID.String(),
            Name:        skill.Name,
            Description: skill.Description,
            Content:     skill.Content,
            Files:       convertFiles(files),
        })
    }
    return result, nil
}
```

关键点：
- 只有 **Agent 被明确分配的 Skills** 才会被加载
- 不存在隐式继承——Skills 不会因为 Agent 在某个 Workspace 就自动拥有

## Step 2: Daemon 认领任务时携带 Skills

**文件**：`server/internal/handler/daemon.go`（约 line 388）

当 Daemon 认领任务时，Handler 调用 `TaskService.LoadAgentSkills`，将 Skills 数据放入认领响应的 `TaskAgentData.Skills` 字段。

```go
// Daemon 认领响应中包含 Agent 信息
type TaskAgentData struct {
    ID          string
    Name        string
    Instructions string
    Skills      []AgentSkillData  // ← Skills 在这里
    // ...
}
```

## Step 3: Daemon 数据转换

**文件**：`server/internal/daemon/daemon.go`（约 line 942-960）

```go
// convertSkillsForEnv 将 API 层数据转换为执行环境格式
func convertSkillsForEnv(skills []SkillData) []SkillContextForEnv {
    var result []SkillContextForEnv
    for _, s := range skills {
        var files []SkillFileContextForEnv
        for _, f := range s.Files {
            files = append(files, SkillFileContextForEnv{
                Path:    f.Path,
                Content: f.Content,
            })
        }
        result = append(result, SkillContextForEnv{
            Name:    s.Name,
            Content: s.Content,
            Files:   files,
        })
    }
    return result
}
```

## Step 4: 写入工作目录

**文件**：`server/internal/daemon/execenv/context.go`

### 目录结构规则

不同 Provider 使用不同的目录写入 Skills：

| Provider | Skills 目录 | 配置文件 | 发现机制 |
|----------|-------------|----------|----------|
| Claude | `{workDir}/.claude/skills/{name}/` | `CLAUDE.md` | 原生自动发现 |
| OpenCode | `{workDir}/.config/opencode/skills/{name}/` | `AGENTS.md` | 原生自动发现 |
| Codex | `{codex-home}/skills/{name}/` | `AGENTS.md` | 原生自动发现 |
| Gemini | `{workDir}/.agent_context/skills/{name}/` | `GEMINI.md` | 配置文件引用 |
| 其他 | `{workDir}/.agent_context/skills/{name}/` | 无 | 配置文件引用 |

### 目录名称处理

```go
// sanitizeSkillName 将 Skill 名称转换为安全的目录名
// 规则：小写 → 非字母数字替换为 "-" → 去除首尾 "-"
func sanitizeSkillName(name string) string {
    s := strings.ToLower(strings.TrimSpace(name))
    s = nonAlphaNum.ReplaceAllString(s, "-")
    s = strings.Trim(s, "-")
    if s == "" { s = "skill" }
    return s
}
```

示例：`"Code Review Expert"` → `"code-review-expert"`

### 文件写入结构

```
.claude/skills/
├── code-review-expert/
│   ├── SKILL.md                          # Skill 主指令文件
│   ├── templates/
│   │   └── review-checklist.md           # 辅助文件
│   └── examples/
│       └── good-pr-pattern.go            # 辅助文件
└── commit-standards/
    ├── SKILL.md
    └── conventional-commits.md
```

每个 Skill 目录包含：
- `SKILL.md`：**必选**，Skill 的核心指令内容
- 辅助文件：**可选**，按 `skill_file.path` 的相对路径结构写入

### 写入逻辑

```go
func writeSkillFiles(skillsDir string, skills []SkillContextForEnv) error {
    for _, skill := range skills {
        // 1. 为每个 Skill 创建目录（名称经过 sanitize）
        dir := filepath.Join(skillsDir, sanitizeSkillName(skill.Name))

        // 2. 写入 SKILL.md（主文件）
        os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(skill.Content), 0o644)

        // 3. 写入所有辅助文件
        for _, f := range skill.Files {
            fpath := filepath.Join(dir, f.Path)
            os.MkdirAll(filepath.Dir(fpath), 0o755)
            os.WriteFile(fpath, []byte(f.Content), 0o644)
        }
    }
}
```

## Step 5: 运行时配置注入

**文件**：`server/internal/daemon/execenv/runtime_config.go`

除了写入 Skills 文件，系统还会在 Provider 原生的配置文件中列出可用的 Skills：

### Claude — CLAUDE.md

```markdown
## Skills

You have the following skills installed (discovered automatically):

- **Code Review Expert**
- **Commit Standards**
```

因为 Claude Code 会原生扫描 `.claude/skills/` 目录，所以只需列出名称作为提示。

### Codex / OpenCode / OpenClaw — AGENTS.md

与 Claude 类似，列出 Skill 名称。这些 Provider 也会原生发现各自目录下的 Skills。

### Gemini — GEMINI.md

```markdown
## Skills

Detailed skill instructions are in `.agent_context/skills/`. Each subdirectory contains a `SKILL.md`.

- **Code Review Expert**
- **Commit Standards**
```

Gemini 不自动扫描目录，所以需要在配置中明确指向路径。

## issue_context.md 中的 Skills 引用

**文件**：`server/internal/daemon/execenv/context.go`（`renderIssueContext`）

除了 Skills 文件和运行时配置，`issue_context.md` 也会列出 Skills：

```markdown
# Task Assignment

**Issue ID:** abc-123

## Quick Start

Run `multica issue get abc-123 --output json` to fetch the full issue details.

## Agent Skills

The following skills are available to you:

- **Code Review Expert**
- **Commit Standards**
```

这提供了三重注入保障：
1. **Skills 文件**：Provider 原生目录中的完整 SKILL.md
2. **运行时配置**：CLAUDE.md / AGENTS.md / GEMINI.md 中的 Skills 列表
3. **Issue 上下文**：issue_context.md 中的 Skills 摘要

## 数据流全景图

```
                     DB Layer
                        │
          ┌─────────────┼───────────────┐
          │             │               │
    agent_skill       skill         skill_file
    (Agent-Skill     (Skill 实体)   (辅助文件)
     关联)
          │             │               │
          └─────────────┼───────────────┘
                        │
                   TaskService
              .LoadAgentSkills()
                        │
                        ▼
                 AgentSkillData[]
                        │
                   Handler Layer
              (Daemon Claim Response)
                        │
                        ▼
                  Daemon Layer
            convertSkillsForEnv()
                        │
                        ▼
               SkillContextForEnv[]
                        │
              ┌─────────┼──────────┐
              │                    │
     writeContextFiles()   InjectRuntimeConfig()
              │                    │
     ┌────────┼────────┐    ┌─────┼─────┐
     │        │        │    │           │
  .claude/ .config/ .agent_  CLAUDE.md  AGENTS.md  GEMINI.md
  skills/ opencode/ context/
          skills/  skills/
```

## 测试覆盖

**文件**：`server/internal/daemon/execenv/execenv_test.go`

关键测试用例：
- `TestWriteContextFiles`：验证 Skill 内容和辅助文件正确写入
- `TestWriteContextFilesOmitsSkillsWhenEmpty`：无 Skills 时不写入
- `TestWriteContextFilesClaudeNativeSkills`：验证 Claude 的 `.claude/skills/` 路径
- `TestWriteContextFilesOpencodeNativeSkills`：验证 OpenCode 的 `.config/opencode/skills/` 路径
- `TestInjectRuntimeConfigClaude`：验证 CLAUDE.md 中列出 Skill 名称
- `TestInjectRuntimeConfigGemini`：验证 GEMINI.md 中引用 Skills 路径
- `TestInjectRuntimeConfigCodex`：验证 AGENTS.md 中列出 Skill 名称
- `TestInjectRuntimeConfigNoSkills`：无 Skills 时配置文件不含 Skills 段
