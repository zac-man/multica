# Makefile 命令说明

## 快速开始

| 命令 | 说明 |
|---|---|
| `make dev` | 一键开发：自动配置环境、安装依赖、启动数据库、运行迁移、启动所有服务 |
| `make setup` | 首次设置：安装依赖、启动数据库、运行迁移 |
| `make start` | 启动后端 + 前端 |
| `make stop` | 停止当前项目的应用进程 |
| `make check` | 完整验证：类型检查 + 单元测试 + Go 测试 + E2E |

## 单独服务

| 命令 | 说明 |
|---|---|
| `make server` | 只启动 Go 后端 |
| `make daemon` | 启动本地 daemon（agent 运行时） |
| `make cli ARGS="..."` | 运行 multica CLI（如 `make cli ARGS="config"`） |
| `make build` | 编译 server + CLI + migrate 二进制文件到 `server/bin/` |
| `make test` | 运行 Go 单元测试 |
| `make clean` | 清理 `server/bin/` 和 `server/tmp/` |

## 数据库

| 命令 | 说明 |
|---|---|
| `make db-up` | 启动共享 PostgreSQL 容器 |
| `make db-down` | 停止共享 PostgreSQL 容器 |
| `make migrate-up` | 执行数据库迁移 |
| `make migrate-down` | 回滚数据库迁移 |
| `make sqlc` | 编辑 SQL 后重新生成 sqlc 代码 |

## 自托管（Docker 部署）

| 命令 | 说明 |
|---|---|
| `make selfhost` | 一键 Docker 部署：创建 env、构建镜像、启动服务 |
| `make selfhost-stop` | 停止所有 Docker Compose 服务 |

## Worktree

| 命令 | 说明 |
|---|---|
| `make worktree-env` | 生成带独立端口和数据库的 `.env.worktree` |
| `make setup-worktree` | 使用 `.env.worktree` 执行 setup |
| `make start-worktree` | 使用 `.env.worktree` 启动服务 |
| `make stop-worktree` | 使用 `.env.worktree` 停止服务 |
| `make check-worktree` | 使用 `.env.worktree` 运行完整验证 |

## 主分支

| 命令 | 说明 |
|---|---|
| `make setup-main` | 使用 `.env` 执行 setup |
| `make start-main` | 使用 `.env` 启动服务 |
| `make stop-main` | 使用 `.env` 停止服务 |
| `make check-main` | 使用 `.env` 运行完整验证 |
