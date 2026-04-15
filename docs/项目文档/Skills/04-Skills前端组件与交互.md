# Skills 前端组件与交互

## 页面路由

Skills 管理页面在两个应用中均可访问：

| 应用 | 路由 | 入口文件 |
|------|------|----------|
| Web | `/skills` | `apps/web/app/(dashboard)/skills/page.tsx` |
| Desktop | `/skills` | `apps/desktop/src/renderer/src/routes.tsx` |

共享页面组件位于 `packages/views/skills/`，遵循项目跨平台架构规范。

侧边栏导航入口在 `packages/views/layout/app-sidebar.tsx` 中注册：
```typescript
{ href: "/skills", label: "Skills", icon: BookOpenText }
```

搜索索引在 `packages/views/search/search-command.tsx` 中包含 Skills 页面。

## 组件架构

```
packages/views/skills/
├── index.ts                              # barrel export
└── components/
    ├── index.ts                          # 组件索引
    ├── skills-page.tsx                   # 主页面（左右分栏布局）
    ├── create-skill-dialog.tsx           # 创建/导入 Skill 对话框
    └── file-tree.tsx                     # 文件树组件

packages/views/agents/components/tabs/
└── skills-tab.tsx                        # Agent 详情中的 Skills 标签页
```

## SkillsPage — 主管理页面

**文件**：`packages/views/skills/components/skills-page.tsx`

### 布局结构

采用 **可调整大小的左右双面板布局**：

```
┌──────────────────────────────────────────────────────────────┐
│  Skills                                           [+ Create] │
├─────────────────────┬────────────────────────────────────────┤
│                     │                                        │
│  Skill List         │  Skill Detail                          │
│  ┌───────────────┐  │  ┌──────────────────────────────────┐  │
│  │ Name          │  │  │ Name: [editable]                  │  │
│  │ Description   │  │  │ Description: [editable]           │  │
│  │ [files badge] │  │  │                                   │  │
│  └───────────────┘  │  │ ┌─ File Browser ───────────────┐ │  │
│  ┌───────────────┐  │  │ │ ▼ SKILL.md  (always first)   │ │  │
│  │ Name          │  │  │ │ ▶ templates/                 │ │  │
│  │ ...           │  │  │ │ ▶ examples/                  │ │  │
│  └───────────────┘  │  │ └──────────────────────────────┘ │  │
│  ...                │  │                                   │  │
│                     │  │ ┌─ File Viewer ────────────────┐ │  │
│                     │  │ │ [Markdown/Code editor]       │ │  │
│                     │  │ └──────────────────────────────┘ │  │
│                     │  │                                   │  │
│                     │  │ [Add File]  [Save]  [Delete]      │  │
│                     │  └──────────────────────────────────┘  │
└─────────────────────┴────────────────────────────────────────┘
```

### 左面板 — Skill 列表

- 展示 Workspace 内所有 Skills
- 每个 Skill 卡片显示：名称、描述、文件数量 badge
- 点击切换选中状态，右面板显示详情
- 数据来源：`useQuery(skillListOptions(wsId))`

### 右面板 — Skill 详情编辑器

功能：
- **名称/描述内联编辑**：直接在详情区编辑
- **文件浏览器**：树状结构浏览 Skill 文件
  - `SKILL.md` 始终排在第一位
  - 目录排在文件前面
  - 支持展开/折叠
- **文件查看器**：选中文件后查看内容
- **添加/删除辅助文件**：管理 Skill 的附加资源
- **保存按钮**：脏检测（dirty detection），仅在有变更时激活
- **删除确认对话框**：防止误删

### 数据流

```
SkillsPage
  │
  ├── useQuery(skillListOptions(wsId))  ← 获取 Skill 列表
  │       │
  │       └── 缓存键: ["workspaces", wsId, "skills"]
  │
  ├── useMutation(updateSkill)           ← 更新 Skill
  │
  ├── useMutation(deleteSkill)           ← 删除 Skill
  │
  └── useMutation(upsertSkillFile)       ← 更新辅助文件
```

## CreateSkillDialog — 创建/导入对话框

**文件**：`packages/views/skills/components/create-skill-dialog.tsx`

两个标签页：

### Create 标签页

```
┌────────────────────────────────────┐
│  Create                            │
│                                    │
│  Name:        [________________]   │
│  Description: [________________]   │
│               [________________]   │
│                                    │
│              [Create]              │
└────────────────────────────────────┘
```

创建后自动在详情页打开，用户可以继续编辑 SKILL.md 内容和添加辅助文件。

### Import 标签页

```
┌────────────────────────────────────┐
│  Import                            │
│                                    │
│  URL:  [https://_______________]   │
│                                    │
│  Supported: ClawHub, skills.sh     │
│                                    │
│              [Import]              │
└────────────────────────────────────┘
```

- 输入 URL 后系统自动检测来源（ClawHub 或 skills.sh）
- 导入完成后自动跳转到该 Skill 的详情页

## FileTree — 文件树组件

**文件**：`packages/views/skills/components/file-tree.tsx`

从文件路径构建虚拟树结构：

```
SKILL.md                  ← 始终排在首位
templates/
├── react-component.tsx
└── api-handler.go
examples/
├── good-pattern.go
└── bad-pattern.go
configs/
└── eslint-rules.json
```

排序规则：
1. `SKILL.md` 永远排第一
2. 目录排在文件前面
3. 同级按名称字母排序

## SkillsTab — Agent 详情中的 Skills 管理

**文件**：`packages/views/agents/components/tabs/skills-tab.tsx`

作为 Agent 详情页的一个标签页（与 Instructions、Tasks、Env、Settings 并列）：

```
┌──────────────────────────────────────────────────┐
│  Agent: Code Reviewer                            │
│  [Instructions] [Tasks] [Skills] [Env] [Settings]│
├──────────────────────────────────────────────────┤
│                                                  │
│  ℹ️ Local runtime skills are automatically       │
│     available and don't need to be assigned.     │
│                                                  │
│  Assigned Skills:                                │
│  ┌──────────────────────────────┐               │
│  │ ✅ Code Review Expert    [×] │               │
│  │ ✅ Commit Standards      [×] │               │
│  └──────────────────────────────┘               │
│                                                  │
│  [+ Add Skill]                                   │
│                                                  │
└──────────────────────────────────────────────────┘
```

功能：
- 显示已分配给该 Agent 的 Skills 列表
- **"Add Skill" 按钮**：弹出选择器对话框，展示 Workspace 内未被分配的 Skills
- **移除按钮**（`[×]`）：取消 Agent 与 Skill 的关联
- **提示横幅**：说明本地运行时 Skills 自动可用，无需手动分配

### 数据流

```
SkillsTab
  │
  ├── useQuery(agentSkills)  ← 获取 Agent 的 Skills
  │
  ├── useQuery(workspaceSkills)  ← 获取 Workspace 所有 Skills（用于选择器）
  │
  └── useMutation(setAgentSkills)  ← 全量设置 Agent 的 skill_ids
```

## 实时同步

**文件**：`packages/core/realtime/use-realtime-sync.ts`

Skills 通过 WebSocket 事件实时同步：

```typescript
// 监听 "skill" 前缀的 WS 事件
// skill:created, skill:updated, skill:deleted

// 100ms debounce 后自动 invalidate:
queryClient.invalidateQueries({
    queryKey: workspaceKeys.skills(wsId)
});
```

效果：
- 团队成员创建新 Skill → 其他成员的列表自动更新
- 编辑 Skill 内容 → 其他成员看到最新内容
- 删除 Skill → 列表中自动移除

## API 客户端方法

**文件**：`packages/core/api/client.ts`

```typescript
class ApiClient {
    // Skill CRUD
    listSkills(): Promise<Skill[]>
    getSkill(id: string): Promise<Skill>
    createSkill(data: CreateSkillRequest): Promise<Skill>
    updateSkill(id: string, data: UpdateSkillRequest): Promise<Skill>
    deleteSkill(id: string): Promise<void>

    // Skill 导入
    importSkill(data: { url: string }): Promise<Skill>

    // Agent-Skill 关联
    listAgentSkills(agentId: string): Promise<Skill[]>
    setAgentSkills(agentId: string, data: SetAgentSkillsRequest): Promise<void>
}
```

## 跨平台规范遵循

Skills 前端组件完全遵循项目的跨平台架构规范：

| 规范 | 遵循情况 |
|------|----------|
| 页面组件在 `packages/views/` | SkillsPage 在 `packages/views/skills/` |
| 不导入 `next/*` 或 `react-router-dom` | 使用 `NavigationAdapter` |
| Web 和 Desktop 路由注册 | 两端均已注册 |
| 状态管理 | TanStack Query 管理服务端状态 |
| 组件复用 | SkillsPage 被两个 App 共用 |
