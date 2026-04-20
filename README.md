<div align="center">

# 🚀 AI Router

**智能 AI 模型路由网关 — 统一接入、智能路由、成本可控**

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE.txt)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.5+-4FC08D?logo=vue.js)](https://vuejs.org/)

</div>

---

## 📖 项目简介

AI Router 是一个开源的 AI 模型路由网关，帮助您统一管理多个 AI 服务商（OpenAI、智谱 AI 等），提供智能路由、负载均衡、成本控制、用量统计等能力。兼容 OpenAI API 格式，可无缝对接现有应用。

## ✨ 核心特性

- 🔄 **智能路由** — 支持自动路由、成本优先、延迟优先、固定路由等多种策略
- 🔌 **多服务商接入** — 统一管理 OpenAI、智谱 AI 等 AI 服务商，兼容 OpenAI API 格式
- 💰 **成本控制** — 实时计费、余额管理、用量统计，精准掌控 AI 调用成本
- 🔑 **BYOK 支持** — 用户可自带 API Key（Bring Your Own Key），灵活使用自有额度
- 🖼️ **图像生成** — 支持 AI 图像生成接口，统一管理图像模型调用
- 🧩 **插件系统** — 可扩展的插件架构，支持联网搜索等增强能力
- 🛡️ **安全管控** — IP 黑名单、速率限制、用户权限分级管理
- 💳 **在线充值** — 集成 Stripe 支付，支持用户自助充值
- 📊 **数据看板** — 调用日志、Token 统计、成本分析，一目了然
- 🏥 **健康检查** — 自动检测服务商可用性，故障自动切换

## 🏗️ 技术架构

### 后端

| 技术 | 说明 |
|------|------|
| Go | 主语言，高性能后端 |
| Gin | HTTP Web 框架 |
| GORM | ORM 框架，支持自动建表 |
| PostgreSQL | 主数据库 |
| Redis | 缓存 / Session / 限流 |

### 前端

| 技术 | 说明 |
|------|------|
| Vue 3 | 前端框架 |
| Ant Design Vue | UI 组件库 |
| Pinia | 状态管理 |
| ECharts | 数据可视化 |
| Vite | 构建工具 |

### 项目结构

```
ai-router/
├── cmd/server/            # 程序入口
├── internal/
│   ├── adapter/           # 模型适配器（OpenAI、智谱等）
│   ├── common/            # 通用工具（分页、响应）
│   ├── config/            # 配置加载
│   ├── constant/          # 常量定义
│   ├── controller/        # 控制器层
│   ├── errno/             # 错误码定义
│   ├── middleware/         # 中间件（认证、CORS、限流、黑名单）
│   ├── model/
│   │   ├── dto/           # 数据传输对象
│   │   ├── entity/        # 数据库实体
│   │   └── vo/            # 视图对象
│   ├── repository/        # 数据访问层
│   ├── router/            # 路由注册
│   ├── service/           # 业务逻辑层
│   ├── strategy/          # 路由策略
│   └── task/              # 定时任务
└── web/                   # Vue 3 前端
    ├── src/
    │   ├── api/           # API 接口
    │   ├── components/    # 公共组件
    │   ├── pages/         # 页面
    │   │   ├── admin/     # 管理后台
    │   │   └── user/      # 用户页面
    │   ├── stores/        # 状态管理
    │   └── router/        # 前端路由
    └── Dockerfile
```

## 🚀 快速开始

### 环境要求

- Go 1.26+
- Node.js 22+
- PostgreSQL 14+
- Redis 6+

### 1. 克隆项目

```bash
git clone https://github.com/SuperJinggg/ai-router.git
cd ai-router
```

### 2. 配置环境变量

复制 `.env` 文件并修改配置：

```bash
cp .env .env.dev
```

关键配置项：

```env
# 数据库
POSTGRES_DSN=postgres://postgres:password@127.0.0.1:5432/ai_router?sslmode=disable

# Redis
REDIS_ADDR=127.0.0.1:6379

# Session 密钥（必须修改）
SESSION_SECRET=your-session-secret

# AI 服务商
AI_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode
AI_API_KEY=your-api-key
AI_MODEL=qwen-plus
```

### 3. 启动后端

项目使用 GORM AutoMigrate，启动时自动创建数据库表并插入初始测试数据：

```bash
go run ./cmd/server/
```

默认启动在 `http://localhost:8123`，首次启动会自动：
- 创建所有数据库表（user、api_key、model、model_provider 等）
- 插入两条测试用户数据

### 4. 启动前端

```bash
cd web
npm install
npm run dev
```

前端默认运行在 `http://localhost:5173`。

### 5. 预置测试账号

| 账号 | 密码 | 角色 |
|------|------|------|
| admin | 12345678 | 管理员 |
| user | 12345678 | 普通用户 |

> 密码采用 MD5 + 盐值（yupi）加密存储

## 🐳 Docker 部署

### 前端

```bash
cd web
docker build -t ai-router-web .
docker run -d -p 80:80 ai-router-web
```

### 后端

```bash
docker build -t ai-router-server .
docker run -d -p 8123:8123 ai-router-server
```

## 📡 API 概览

所有接口前缀为 `/api`，主要接口分组如下：

| 模块 | 路径 | 说明 |
|------|------|------|
| 健康检查 | `/api/health` | 服务健康状态 |
| 用户 | `/api/user` | 注册、登录、用户管理 |
| API Key | `/api/api/key` | API Key 创建与管理 |
| 服务商 | `/api/provider` | AI 服务商管理 |
| 模型 | `/api/model` | AI 模型管理 |
| 聊天 | `/api/v1/chat` | OpenAI 兼容聊天接口 |
| 内部聊天 | `/api/internal/chat` | 内部聊天接口（需登录） |
| 图像生成 | `/api/v1/images` | AI 图像生成接口 |
| 插件 | `/api/plugin` | 插件管理与执行 |
| BYOK | `/api/byok` | 自带 API Key 管理 |
| 充值 | `/api/recharge` | Stripe 充值 |
| 余额 | `/api/balance` | 余额与账单查询 |
| 统计 | `/api/stats` | 用量统计与分析 |
| 黑名单 | `/api/admin/blacklist` | IP 黑名单管理 |

### OpenAI 兼容接口

聊天补全接口兼容 OpenAI 格式，可直接替换 `base_url` 使用：

```
POST /api/v1/chat/completions
```

支持流式（SSE）和非流式响应。

## 🔀 路由策略

| 策略 | 说明 |
|------|------|
| `auto` | 自动路由 — 综合评分选择最优模型 |
| `cost_first` | 成本优先 — 优先选择价格最低的模型 |
| `latency_first` | 延迟优先 — 优先选择响应最快的模型 |
| `fixed` | 固定路由 — 使用指定模型，故障时自动降级 |

## 🤝 参与贡献

欢迎提交 Issue 和 Pull Request 参与贡献！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 发起 Pull Request

## 📄 许可证

本项目基于 [MIT License](LICENSE.txt) 开源。
