# Skills 数据模型与 API

## 数据库 Schema

### `skill` 表 — Skill 实体

```sql
CREATE TABLE skill (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspace(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    content     TEXT NOT NULL DEFAULT '',       -- SKILL.md 正文内容
    config      JSONB NOT NULL DEFAULT '{}',    -- 扩展配置
    created_by  UUID REFERENCES "user"(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(workspace_id, name)                  -- 同一 Workspace 内名称唯一
);

CREATE INDEX idx_skill_workspace ON skill(workspace_id);
```

字段说明：
- `name`：Skill 名称，在 Workspace 内唯一，用于展示和目录名
- `description`：简短描述，用于列表展示
- `content`：SKILL.md 的完整 Markdown 内容，是 Skill 的核心指令
- `config`：JSONB 格式的扩展配置（当前为 `{}`，预留扩展）
- `created_by`：创建者，用于权限控制（创建者和管理员可编辑/删除）

### `skill_file` 表 — Skill 辅助文件

```sql
CREATE TABLE skill_file (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    skill_id   UUID NOT NULL REFERENCES skill(id) ON DELETE CASCADE,
    path       TEXT NOT NULL,                    -- 相对路径，如 "templates/example.go"
    content    TEXT NOT NULL,                    -- 文件内容
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(skill_id, path)                       -- 同一 Skill 内路径唯一
);

CREATE INDEX idx_skill_file_skill ON skill_file(skill_id);
```

典型用途：
- 代码模板文件（如 `templates/react-component.tsx`）
- 示例代码（如 `examples/api-handler.go`）
- 配置文件（如 `configs/eslint-rules.json`）

### `agent_skill` 表 — Agent-Skill 多对多关联

```sql
CREATE TABLE agent_skill (
    agent_id   UUID NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    skill_id   UUID NOT NULL REFERENCES skill(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (agent_id, skill_id)
);

CREATE INDEX idx_agent_skill_skill ON agent_skill(skill_id);
CREATE INDEX idx_agent_skill_agent ON agent_skill(agent_id);
```

设计要点：
- 一个 Agent 可以有多个 Skills
- 一个 Skill 可以分配给多个 Agents
- 通过 `SetAgentSkills` API 进行全量替换（而非增量操作）

### 实体关系图

```
workspace ──1:N──> skill ──1:N──> skill_file
                    │
                    │ M:N
                    ▼
                  agent
                    │
                    │ (通过 agent_skill 关联表)
```

## Go 数据模型

### 数据库生成模型（sqlc）

```go
// server/pkg/db/generated/models.go

type Skill struct {
    ID          uuid.UUID
    WorkspaceID uuid.UUID
    Name        string
    Description string
    Content     string
    Config      []byte          // JSONB
    CreatedBy   uuid.NullUUID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type SkillFile struct {
    ID        uuid.UUID
    SkillID   uuid.UUID
    Path      string
    Content   string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type AgentSkill struct {
    AgentID   uuid.UUID
    SkillID   uuid.UUID
    CreatedAt time.Time
}
```

### API 层模型（Handler）

```go
// server/internal/handler/skill.go

type SkillResponse struct {
    ID          string            `json:"id"`
    WorkspaceID string            `json:"workspace_id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Content     string            `json:"content"`
    Config      json.RawMessage   `json:"config"`
    CreatedBy   *string           `json:"created_by"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}

type SkillFileResponse struct {
    ID        string    `json:"id"`
    SkillID   string    `json:"skill_id"`
    Path      string    `json:"path"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type SkillWithFilesResponse struct {
    SkillResponse
    Files []SkillFileResponse `json:"files"`
}
```

### 任务执行层模型

```go
// server/internal/service/task.go
type AgentSkillData struct {
    ID          string
    Name        string
    Description string
    Content     string
    Files       []AgentSkillFileData
}

type AgentSkillFileData struct {
    Path    string
    Content string
}

// server/internal/daemon/execenv/execenv.go
type SkillContextForEnv struct {
    Name    string
    Content string
    Files   []SkillFileContextForEnv
}

type SkillFileContextForEnv struct {
    Path    string
    Content string
}
```

## TypeScript 前端类型

```typescript
// packages/core/types/agent.ts

interface Skill {
  id: string;
  workspace_id: string;
  name: string;
  description: string;
  content: string;                    // SKILL.md 内容
  config: Record<string, unknown>;    // 扩展配置
  files: SkillFile[];                 // 包含辅助文件列表
  created_by: string | null;
  created_at: string;
  updated_at: string;
}

interface SkillFile {
  id: string;
  skill_id: string;
  path: string;                       // 相对路径
  content: string;                    // 文件内容
  created_at: string;
  updated_at: string;
}

interface CreateSkillRequest {
  name: string;
  description?: string;
  content?: string;
  config?: Record<string, unknown>;
  files?: { path: string; content: string }[];
}

interface UpdateSkillRequest {
  name?: string;
  description?: string;
  content?: string;
  config?: Record<string, unknown>;
  files?: { path: string; content: string }[];
}

interface SetAgentSkillsRequest {
  skill_ids: string[];                // 全量替换
}
```

## HTTP API 端点

所有端点需要 Workspace 成员权限，通过 `X-Workspace-ID` Header 路由。

### Skill CRUD

| 方法 | 路径 | Handler | 说明 | 权限 |
|------|------|---------|------|------|
| GET | `/api/skills` | `ListSkills` | 列出 Workspace 所有 Skills | 成员 |
| POST | `/api/skills` | `CreateSkill` | 创建新 Skill | 成员 |
| GET | `/api/skills/{id}` | `GetSkill` | 获取 Skill 详情（含 files） | 成员 |
| PUT | `/api/skills/{id}` | `UpdateSkill` | 更新 Skill（支持部分更新） | 创建者或管理员 |
| DELETE | `/api/skills/{id}` | `DeleteSkill` | 删除 Skill | 创建者或管理员 |

### Skill 文件管理

| 方法 | 路径 | Handler | 说明 |
|------|------|---------|------|
| GET | `/api/skills/{id}/files` | `ListSkillFiles` | 列出 Skill 所有文件 |
| PUT | `/api/skills/{id}/files` | `UpsertSkillFile` | 创建或更新文件（按 path 去重） |
| DELETE | `/api/skills/{id}/files/{fileId}` | `DeleteSkillFile` | 删除单个文件 |

### Skill 导入

| 方法 | 路径 | Handler | 说明 |
|------|------|---------|------|
| POST | `/api/skills/import` | `ImportSkill` | 从外部来源导入 Skill |

请求体：`{ "url": "https://..." }`

支持的外部来源：
- **ClawHub**（`clawhub.ai`）：技能市场
- **skills.sh**（GitHub 仓库）：从 GitHub 仓库解析 SKILL.md

### Agent-Skill 关联

| 方法 | 路径 | Handler | 说明 |
|------|------|---------|------|
| GET | `/api/agents/{id}/skills` | `ListAgentSkills` | 获取 Agent 的所有 Skills |
| PUT | `/api/agents/{id}/skills` | `SetAgentSkills` | 全量设置 Agent 的 Skills |

## WebSocket 事件

```go
// server/pkg/protocol/events.go
EventSkillCreated = "skill:created"    // 创建或导入 Skill 时
EventSkillUpdated = "skill:updated"    // 更新 Skill 时
EventSkillDeleted = "skill:deleted"    // 删除 Skill 时
```

前端通过 `use-realtime-sync.ts` 监听这些事件，自动 invalidate Skills Query 缓存：
```typescript
// packages/core/realtime/use-realtime-sync.ts (line 104-108)
// "skill" 前缀的事件 → invalidate workspaceKeys.skills(wsId)
// 100ms debounce 防止频繁刷新
```

## CLI 命令

```bash
# Skill CRUD
multica skill list                                    # 列出 Workspace Skills
multica skill get <id>                                # 获取 Skill 详情（含文件）
multica skill create --name <name>                    # 创建 Skill
multica skill update <id> [--name ...] [--desc ...]   # 更新 Skill
multica skill delete <id>                             # 删除 Skill

# Skill 导入
multica skill import --url <url>                      # 从外部导入

# Skill 文件管理
multica skill files list <skill-id>                   # 列出文件
multica skill files upsert <skill-id>                 # 创建/更新文件
multica skill files delete <skill-id> <file-id>       # 删除文件
```

## 权限模型

- **查看 Skills**：Workspace 所有成员
- **创建 Skill**：Workspace 所有成员
- **编辑/删除 Skill**：仅 Skill 创建者或 Workspace Admin/Owner

权限检查逻辑（`canManageSkill` in `server/internal/handler/skill.go`）：
```go
func (h *Handler) canManageSkill(r *http.Request, skill *generated.Skill) bool {
    // 1. 获取当前请求用户 ID
    // 2. 获取 Workspace 成员角色
    // 3. 允许条件：请求者是 skill.CreatedBy 或角色为 admin/owner
}
```
