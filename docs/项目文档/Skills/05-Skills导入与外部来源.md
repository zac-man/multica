# Skills 导入与外部来源

## 概述

Multica 支持从外部来源导入 Skills，避免从零开始编写。系统目前支持两种外部来源：

1. **ClawHub**（`clawhub.ai`）— Skills 市场/注册中心
2. **skills.sh**（GitHub 仓库）— 从 GitHub 仓库中解析 SKILL.md

## 导入入口

| 入口 | 方式 |
|------|------|
| Web UI | Skills 页面 → Create 按钮 → Import 标签页 |
| Desktop UI | 同 Web UI |
| CLI | `multica skill import --url <url>` |
| API | `POST /api/skills/import`，body: `{"url": "..."}` |

## URL 自动检测

**文件**：`server/internal/handler/skill.go`（`detectImportSource` 函数）

系统根据 URL 域名自动判断来源：

| URL 模式 | 来源 | 说明 |
|----------|------|------|
| `clawhub.ai/skills/...` 或 `www.clawhub.ai/...` | ClawHub | Skills 市场 |
| `skills.sh/...` 或包含 `github.com` | skills.sh | GitHub 仓库 |
| 其他 | 尝试按 skills.sh 处理 | 兼容模式 |

## ClawHub 导入流程

```
用户输入 URL: https://clawhub.ai/skills/code-review
         │
         ▼
┌─────────────────────────────────────────────────┐
│  1. 解析 Skill Slug                              │
│     URL → slug: "code-review"                    │
│                                                   │
│  2. 获取 Skill 元数据                              │
│     GET clawhub.ai/api/skills/{slug}              │
│     → 获取 displayName, summary, tags             │
│                                                   │
│  3. 获取最新版本                                   │
│     从响应中提取 latestVersion.version              │
│                                                   │
│  4. 获取版本详情（含文件列表）                       │
│     GET clawhub.ai/api/skills/{slug}/versions/{v}  │
│     → 获取 files[]（每个文件的 path 和 size）        │
│                                                   │
│  5. 逐个下载文件内容                                │
│     对每个 file entry:                              │
│     GET clawhub.ai/api/skills/{slug}/.../{path}    │
│     → 获取文件内容                                  │
│                                                   │
│  6. 创建 Skill 记录                                │
│     name = displayName (或 slug)                   │
│     description = summary                          │
│     content = SKILL.md 的文件内容                   │
│     files = 其余文件作为 skill_file                 │
└─────────────────────────────────────────────────┘
```

### ClawHub API 数据结构

```go
type clawhubGetSkillResponse struct {
    Skill         clawhubSkill          `json:"skill"`
    LatestVersion *clawhubLatestVersion `json:"latestVersion"`
}

type clawhubSkill struct {
    Slug        string            `json:"slug"`
    DisplayName string            `json:"displayName"`
    Summary     string            `json:"summary"`
    Tags        map[string]string `json:"tags"`
}

type clawhubVersionDetail struct {
    Version string             `json:"version"`
    Files   []clawhubFileEntry `json:"files"`
}

type clawhubFileEntry struct {
    Path string `json:"path"`
    Size int64  `json:"size"`
}
```

## skills.sh / GitHub 导入流程

```
用户输入 URL: https://skills.sh/owner/repo/skill-name
  或直接输入: https://github.com/owner/repo
         │
         ▼
┌─────────────────────────────────────────────────┐
│  1. 解析 GitHub Owner/Repo                       │
│     从 URL 中提取 owner 和 repo                    │
│                                                   │
│  2. 获取默认分支                                   │
│     GET api.github.com/repos/{owner}/{repo}       │
│     → 获取 default_branch（失败则回退 "main"）      │
│                                                   │
│  3. 搜索 SKILL.md 位置                             │
│     按优先级尝试多个路径约定：                       │
│     ┌───────────────────────────────────────┐     │
│     │ ① skills/{name}/SKILL.md             │     │
│     │ ② .claude/skills/{name}/SKILL.md     │     │
│     │ ③ plugin/skills/{name}/SKILL.md      │     │
│     │ ④ SKILL.md (仓库根目录)               │     │
│     └───────────────────────────────────────┘     │
│     GET api.github.com/repos/{owner}/{repo}/      │
│         contents/{path}?ref={branch}              │
│                                                   │
│  4. 下载 SKILL.md 内容                             │
│     使用 download_url 获取文件原文                  │
│                                                   │
│  5. 搜索并下载辅助文件                              │
│     列出 SKILL.md 所在目录下的所有文件               │
│     跳过 SKILL.md 本身                             │
│     对每个文件使用 download_url 获取内容             │
│                                                   │
│  6. 创建 Skill 记录                                │
│     name = 目录名（或仓库名）                       │
│     content = SKILL.md 内容                        │
│     files = 其余文件作为 skill_file                │
└─────────────────────────────────────────────────┘
```

### 支持的 URL 格式

| 格式 | 解析结果 |
|------|----------|
| `skills.sh/owner/repo/skill-name` | owner, repo, skill-name |
| `github.com/owner/repo` | owner, repo, 搜索所有 Skills |
| `github.com/owner/repo/tree/branch/path` | owner, repo, branch, path |

### SKILL.md 搜索路径优先级

系统会按以下顺序搜索 SKILL.md 文件：

1. `skills/{name}/SKILL.md` — Claude Code 社区约定
2. `.claude/skills/{name}/SKILL.md` — Claude Code 原生路径
3. `plugin/skills/{name}/SKILL.md` — 插件系统路径
4. `SKILL.md` — 仓库根目录（兜底）

这个优先级设计使得一个 GitHub 仓库可以包含多个 Skills，每个在各自的子目录下。

## skills-lock.json

**文件**：项目根目录 `skills-lock.json`

追踪本地安装的外部 Skills：

```json
[
  {
    "name": "frontend-design",
    "source": "anthropics/skills",
    "sourceType": "github",
    "computedHash": "abc123..."
  },
  {
    "name": "shadcn",
    "source": "shadcn/ui",
    "sourceType": "github",
    "computedHash": "def456..."
  }
]
```

作用：
- 记录 Skill 的来源仓库和哈希值
- 用于版本追踪和更新检测
- 不是核心业务逻辑，而是开发辅助工具

## 导入后的处理

无论从哪个来源导入，导入成功后：

1. 在 `skill` 表创建记录（workspace 级别）
2. 在 `skill_file` 表创建辅助文件记录
3. 发布 `skill:created` WebSocket 事件
4. 前端自动刷新 Skills 列表
5. 导入的 Skill 与手动创建的 Skill 完全等价

导入的 Skill 可以：
- 被分配给任何 Agent
- 被编辑和修改
- 被删除
- 被重新导入（会覆盖同名 Skill）

## 错误处理

| 错误场景 | 处理方式 |
|----------|----------|
| URL 无法识别 | 返回错误提示"无法识别的 URL" |
| ClawHub API 不可用 | 返回上游错误 |
| GitHub 仓库不存在 | 返回 404 错误 |
| 仓库中找不到 SKILL.md | 返回"未找到 Skill 文件"错误 |
| 同名 Skill 已存在 | 覆盖更新（Upsert 语义） |
| 网络超时 | 返回超时错误 |
