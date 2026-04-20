# Multica CLI 命令参考

> 完整的 `multica` 命令行工具命令列表。

## 全局参数

| 参数 | 环境变量 | 说明 |
|------|---------|------|
| `--server-url` | `MULTICA_SERVER_URL` | 服务端地址 |
| `--workspace-id` | `MULTICA_WORKSPACE_ID` | 工作区 ID |
| `--profile` | — | 配置档案名（隔离配置、daemon 状态和工作区） |

---

## 快速开始

```bash
multica setup              # 云端：一键配置 + 登录 + 启动 daemon
multica setup self-host    # 自部署：配置本地服务

# 或分步执行
multica login              # 浏览器 OAuth 登录
multica daemon start       # 启动本地 daemon
```

---

## 认证

| 命令 | 说明 |
|------|------|
| `multica login` | 浏览器 OAuth 登录 |
| `multica login --token` | Token 登录（适用于无浏览器环境） |
| `multica auth status` | 查看认证状态 |
| `multica auth logout` | 退出登录 |

---

## Issue（任务）

### 列表与查询

| 命令 | 说明 |
|------|------|
| `multica issue list` | 列出任务 |
| `multica issue get <id>` | 查看任务详情 |
| `multica issue search <query>` | 搜索任务 |

**`issue list` 参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--status` | — | 按状态筛选 |
| `--priority` | — | 按优先级筛选 |
| `--assignee` | — | 按指派人筛选 |
| `--project` | — | 按项目 ID 筛选 |
| `--limit` | `50` | 最大返回数量 |
| `--output` | `table` | 输出格式：`table` / `json` |

### 创建与修改

| 命令 | 说明 |
|------|------|
| `multica issue create` | 创建任务 |
| `multica issue update <id>` | 更新任务 |
| `multica issue status <id> <status>` | 更改状态 |
| `multica issue assign <id> --to <name>` | 指派任务 |
| `multica issue assign <id> --unassign` | 取消指派 |

**`issue create` 参数：**

| 参数 | 说明 |
|------|------|
| `--title` | 标题（必填） |
| `--description` | 描述 |
| `--status` | 状态 |
| `--priority` | 优先级 |
| `--assignee` | 指派人 |
| `--parent` | 父任务 ID |
| `--project` | 项目 ID |
| `--due-date` | 截止日期（RFC3339） |
| `--attachment` | 附件路径（可多次指定） |

**有效状态值：** `backlog`、`todo`、`in_progress`、`in_review`、`done`、`blocked`、`cancelled`

### 评论

| 命令 | 说明 |
|------|------|
| `multica issue comment list <issue-id>` | 列出评论 |
| `multica issue comment add <issue-id>` | 添加评论 |
| `multica issue comment delete <comment-id>` | 删除评论 |

**`comment add` 参数：**

| 参数 | 说明 |
|------|------|
| `--content` | 评论内容 |
| `--content-stdin` | 从 stdin 读取内容 |
| `--parent` | 回复的父评论 ID |
| `--attachment` | 附件路径 |

### 执行记录

| 命令 | 说明 |
|------|------|
| `multica issue runs <issue-id>` | 查看执行历史 |
| `multica issue run-messages <task-id>` | 查看执行消息日志 |
| `multica issue run-messages <task-id> --since 42` | 增量获取（从序列号 42 之后） |

---

## Project（项目）

| 命令 | 说明 |
|------|------|
| `multica project list` | 列出项目 |
| `multica project get <id>` | 查看项目详情 |
| `multica project create` | 创建项目 |
| `multica project update <id>` | 更新项目 |
| `multica project delete <id>` | 删除项目 |
| `multica project status <id> <status>` | 更改项目状态 |

**`project create` 参数：**

| 参数 | 说明 |
|------|------|
| `--title` | 标题（必填） |
| `--description` | 描述 |
| `--status` | 状态 |
| `--icon` | 图标（emoji） |
| `--lead` | 负责人 |

**有效项目状态：** `planned`、`in_progress`、`paused`、`completed`、`cancelled`

---

## Agent（智能体）

| 命令 | 说明 |
|------|------|
| `multica agent list` | 列出智能体 |
| `multica agent get <id>` | 查看详情 |
| `multica agent create` | 创建智能体 |
| `multica agent update <id>` | 更新智能体 |
| `multica agent archive <id>` | 归档智能体 |
| `multica agent restore <id>` | 恢复智能体 |
| `multica agent tasks <id>` | 查看任务列表 |

**`agent create` 参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--name` | — | 名称（必填） |
| `--description` | — | 描述 |
| `--instructions` | — | 指令 |
| `--runtime-id` | — | 运行时 ID（必填） |
| `--runtime-config` | — | 运行时配置（JSON） |
| `--custom-args` | — | 自定义 CLI 参数（JSON 数组，如 `'["--model", "o3"]'`） |
| `--visibility` | `private` | 可见性：`private` / `workspace` |
| `--max-concurrent-tasks` | `6` | 最大并发任务数 |

### Agent Skills

| 命令 | 说明 |
|------|------|
| `multica agent skills list <agent-id>` | 查看智能体技能 |
| `multica agent skills set <agent-id> --skill-ids id1,id2` | 设置智能体技能 |

---

## Workspace（工作区）

| 命令 | 说明 |
|------|------|
| `multica workspace list` | 列出工作区 |
| `multica workspace get [id]` | 查看工作区详情 |
| `multica workspace members [id]` | 列出成员 |

---

## Skill（技能）

| 命令 | 说明 |
|------|------|
| `multica skill list` | 列出技能 |
| `multica skill get <id>` | 查看技能详情 |
| `multica skill create` | 创建技能 |
| `multica skill update <id>` | 更新技能 |
| `multica skill delete <id>` | 删除技能 |
| `multica skill import --url <url>` | 从 URL 导入技能（clawhub.ai 或 skills.sh） |

**`skill create` 参数：**

| 参数 | 说明 |
|------|------|
| `--name` | 名称（必填） |
| `--description` | 描述 |
| `--content` | 技能内容（SKILL.md 正文） |
| `--config` | 配置（JSON） |

### Skill Files

| 命令 | 说明 |
|------|------|
| `multica skill files list <skill-id>` | 列出技能文件 |
| `multica skill files upsert <skill-id>` | 创建或更新技能文件 |
| `multica skill files delete <skill-id> <file-id>` | 删除技能文件 |

---

## Daemon（守护进程）

| 命令 | 说明 |
|------|------|
| `multica daemon start` | 启动 daemon（后台运行） |
| `multica daemon start --foreground` | 前台运行（调试用） |
| `multica daemon stop` | 停止 daemon |
| `multica daemon restart` | 重启 daemon |
| `multica daemon status` | 查看状态 |
| `multica daemon logs` | 查看日志（默认 50 行） |
| `multica daemon logs -f` | 实时跟踪日志 |
| `multica daemon logs -n 100` | 查看最近 100 行 |

**`daemon start` 参数：**

| 参数 | 环境变量 | 默认值 | 说明 |
|------|---------|--------|------|
| `--daemon-id` | `MULTICA_DAEMON_ID` | 主机名 | 唯一标识 |
| `--device-name` | `MULTICA_DAEMON_DEVICE_NAME` | 主机名 | 设备显示名称 |
| `--runtime-name` | `MULTICA_AGENT_RUNTIME_NAME` | `Local Agent` | 运行时显示名称 |
| `--poll-interval` | `MULTICA_DAEMON_POLL_INTERVAL` | `3s` | 任务轮询间隔 |
| `--heartbeat-interval` | `MULTICA_DAEMON_HEARTBEAT_INTERVAL` | `15s` | 心跳间隔 |
| `--agent-timeout` | `MULTICA_AGENT_TIMEOUT` | `2h` | 单任务超时 |
| `--max-concurrent-tasks` | `MULTICA_DAEMON_MAX_CONCURRENT_TASKS` | `20` | 最大并发任务数 |

---

## Runtime（运行时）

| 命令 | 说明 |
|------|------|
| `multica runtime list` | 列出运行时 |
| `multica runtime usage <id>` | 查看 Token 用量 |
| `multica runtime activity <id>` | 查看每小时任务活动 |
| `multica runtime ping <id>` | Ping 运行时检查连通性 |
| `multica runtime update <id> --target-version <ver>` | 更新运行时 CLI 版本 |

**`runtime usage` 参数：** `--days 90`（默认 90 天，最大 365）

---

## 其他命令

| 命令 | 说明 |
|------|------|
| `multica version` | 查看版本信息 |
| `multica update` | 更新到最新版本（自动检测 Homebrew / 二进制安装） |
| `multica config show` | 查看配置 |
| `multica config set <key> <value>` | 设置配置项（支持 `server_url`、`app_url`、`workspace_id`） |
| `multica attachment download <id>` | 下载附件 |
| `multica repo checkout <url>` | 检出仓库（daemon 任务内使用） |
