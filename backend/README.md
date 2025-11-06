# NijiChat Backend

虚拟声优互动模拟器后端服务 - 基于 Golang + Cloudflare Workers

## 技术栈

- **语言**: Go 1.21+
- **框架**: Gin 1.9+ (HTTP路由)
- **Worker框架**: syumai/workers (Cloudflare Worker适配)
- **数据库**: Cloudflare D1 (SQLite兼容)
- **缓存**: Cloudflare KV
- **部署**: Cloudflare Workers

## 项目结构

```
backend/
├── main.go                    # 应用入口
├── handlers/                  # HTTP路由处理器
│   ├── seiyuu.go             # 声优相关API
│   ├── scheduler.go          # 智能调度器API
│   ├── admin.go              # 管理员功能API
│   ├── groups.go             # 群组管理API
│   └── moegirl.go            # 萌娘百科集成API
├── models/                    # 数据模型
│   ├── seiyuu.go             # 声优数据模型
│   ├── group.go              # 群组数据模型
│   └── admin.go              # 管理员数据模型
├── services/                  # 业务逻辑服务
│   ├── seiyuu_service.go     # 声优业务逻辑
│   ├── scheduler_service.go  # 调度算法实现
│   ├── moegirl_service.go    # 萌娘百科数据处理
│   └── ai_service.go         # AI API调用服务
├── middleware/                # 中间件
│   ├── auth.go               # 认证中间件
│   ├── cors.go               # CORS处理
│   └── logging.go            # 日志中间件
├── database/                  # 数据库操作
│   ├── d1.go                 # Cloudflare D1操作
│   ├── kv.go                 # Cloudflare KV操作
│   └── migrations/           # 数据库迁移脚本
│       ├── 001_create_seiyuu.sql
│       ├── 002_create_groups.sql
│       └── 003_create_admin.sql
├── utils/                     # 工具函数
│   ├── response.go           # API响应封装
│   ├── crypto.go             # 加密工具
│   └── validation.go         # 数据验证
├── go.mod                     # Go模块依赖
├── go.sum                     # 依赖校验和
└── wrangler.toml              # Cloudflare配置
```

## 开发指南

### 本地开发

1. 安装依赖
```bash
go mod download
npm install -g wrangler  # 安装 Wrangler CLI
```

2. 配置开发环境
```bash
# 复制环境变量示例文件
cp .dev.vars.example .dev.vars

# 编辑 .dev.vars，填入实际的密钥
vim .dev.vars
```

3. 运行开发服务器
```bash
wrangler dev

# 或使用远程数据库
wrangler dev --remote
```

### 数据库迁移

执行SQL迁移脚本:
```bash
# 创建D1数据库
wrangler d1 create seiyuu-chat-db

# 执行迁移
wrangler d1 execute seiyuu-chat-db --file=database/migrations/001_create_seiyuu.sql
wrangler d1 execute seiyuu-chat-db --file=database/migrations/002_create_groups.sql
wrangler d1 execute seiyuu-chat-db --file=database/migrations/003_create_admin.sql
```

### 部署

```bash
# 部署到生产环境
wrangler deploy --env production

# 部署到 staging 环境
wrangler deploy --env staging
```

详细的部署指南请参考: [DEPLOYMENT.md](./DEPLOYMENT.md)

## API文档

详细的API文档请参考: `docs/API.md`

### 公开接口

- `GET /api/seiyuu` - 获取已发布声优列表
- `GET /api/seiyuu/:id` - 获取声优详细资料
- `GET /api/groups` - 获取群组配置
- `POST /api/scheduler/select` - 智能选择声优

### 管理员接口

- `POST /api/admin/login` - 管理员登录
- `GET /api/admin/seiyuu` - 获取所有声优(含待审核)
- `POST /api/admin/seiyuu` - 创建声优
- `PUT /api/admin/seiyuu/:id` - 更新声优资料
- `DELETE /api/admin/seiyuu/:id` - 删除声优
- `GET /api/admin/moegirl/:name` - 获取萌娘百科原始数据
- `POST /api/admin/process-profile` - AI处理资料转Markdown

## 环境变量

### 开发环境
在 `.dev.vars` 文件中配置（参考 `.dev.vars.example`）：
- `ADMIN_JWT_SECRET` - JWT密钥
- `LIGHTWEIGHT_AI_API_KEY` - 轻量级AI API密钥
- `MOEGIRL_API_BASE_URL` - 萌娘百科API地址

### 生产环境
使用 `wrangler secret` 命令设置：
```bash
wrangler secret put ADMIN_JWT_SECRET --env production
wrangler secret put LIGHTWEIGHT_AI_API_KEY --env production
```

### 配置文件
- `wrangler.jsonc` - Cloudflare Workers 主配置文件（支持注释）
- `.dev.vars` - 开发环境变量（不提交到 Git）
- `.dev.vars.example` - 环境变量示例文件

## 贡献指南

1. Fork本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## License

MIT License
