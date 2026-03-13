# Backend API Tester

一个用 Go 编写的 Web 测试面板，覆盖后端项目的核心接口，支持创角、任务生成、NPC 对话、道具生成和公告查询等功能的可视化测试。

## 功能模块

| 模块 | 接口 | 功能描述 |
|------|------|----------|
| ⚙️ 配置 | — | 设置后端地址和 JWT Token，支持快速登录 |
| 🧙 创角 | `/character/create_npc_charac` `/character/gen_npc_charac` | 两阶段创角流程，展示角色数据和图片 |
| 📜 任务生成 | `/quest_v2/generate_outline_and_quest` | 提交任务生成，自动轮询，展示任务数据和关联道具 |
| 💬 对话 | `/gameagent/talk` | 多轮 NPC 对话，展示回复、动作和选项 |
| 🗡️ 道具生成 | `/itemGen/generate` | AI 全流程道具生成，展示图标和脚本 |
| 📢 公告 | `/info/announcement` `/info/agreement` | 公告和用户协议查询 |

## 快速开始

### 前置条件

- [Go 1.21+](https://go.dev/dl/)

### 运行

```bash
# 克隆仓库
git clone https://github.com/kunwl123456/backend-api-tester.git
cd backend-api-tester

# 方式一：通过环境变量配置（可选）
cp .env.example .env
# 编辑 .env 填入后端地址和 Token

# 方式二：直接运行，在 Web UI 中配置
go run .

# 打开浏览器
# http://localhost:8080
```

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `BACKEND_BASE_URL` | 后端 API 地址 | 空（启动后在 UI 填写） |
| `BACKEND_JWT_TOKEN` | JWT Token | 空（可通过 UI 登录获取） |
| `BACKEND_USERNAME` | 默认用户名 | `player_001` |
| `PORT` | 本工具监听端口 | `8080` |

> 所有配置均可在启动后通过首页 UI 动态修改，无需重启。

## 架构

```
浏览器 → Go HTTP Server（代理层） → Backend API
```

所有 API 请求由 Go 服务端转发，避免浏览器跨域问题，且 Token 不暴露在前端。

## 项目结构

```
.
├── main.go              # 入口，路由注册
├── config/              # 运行时配置管理
├── handlers/            # 各模块接口代理
│   ├── proxy.go         # 通用 HTTP 工具
│   ├── auth.go          # 配置 & 登录
│   ├── character.go     # 创角模块
│   ├── quest.go         # 任务 & 道具模块
│   ├── dialogue.go      # 对话模块
│   └── announcement.go  # 公告 & 信息模块
├── templates/           # HTML 页面模板
├── static/              # 静态资源（CSS）
├── .env.example         # 配置模板
└── go.mod
```

## 安全说明

- 本工具**仅供本地开发测试使用**，不建议部署到公网
- `.env` 文件已加入 `.gitignore`，不会被提交
- 代码中**无任何硬编码**的 URL、Token 或密钥
