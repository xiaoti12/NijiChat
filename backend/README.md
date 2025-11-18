# NijiChat Backend

虚拟声优互动模拟器后端服务 - 基于 Go + Cloudflare Workers

## 项目简介

NijiChat 后端是一个基于 Golang 和 Cloudflare Workers 的无服务器应用,提供声优资料管理、关系管理、用户认证等核心 API 服务。

### 核心特性
- **无服务器架构** - 基于 Cloudflare Workers,全球边缘节点部署
- **RESTful API** - 提供完整的声优管理和查询接口
- **JWT 认证** - 管理员权限控制
- **D1 数据库** - 使用 Cloudflare D1 (SQLite 兼容) 进行数据持久化
- **KV 缓存** - 使用 Cloudflare KV 进行高性能缓存

## 📦 技术栈

### 核心技术
- **Go 1.21+** - 主要编程语言
- **syumai/workers** - Cloudflare Worker 适配框架
- **自定义 Router** - 类 Gin 风格的路由系统,适配 Workers 环境
- **Cloudflare D1** - SQLite 兼容的边缘数据库
- **Cloudflare KV** - 分布式键值存储
- **TinyGo** - 用于生产构建,生成更小的 WASM 文件

### 构建工具
- **Go 编译器** - 用于本地开发构建
- **TinyGo** - 用于生产部署构建
- **Wrangler** - Cloudflare Workers CLI 工具

## 🏗️ 项目架构

### 分层设计

**应用入口**:
- `main.go` - 应用入口,配置中间件和路由

**路由层** (router/):
- `router.go` - 自定义路由系统,类似 Gin 但适配 Cloudflare Workers
- `context.go` - 上下文封装,提供请求/响应处理

**处理器层** (handlers/):
- `seiyuu.go` - 声优相关 API 处理器
- `relationships.go` - 声优关系管理 API 处理器
- `admin.go` - 管理员功能 API 处理器

**服务层** (services/):
- `seiyuu_service.go` - 声优业务逻辑
- `relationship_service.go` - 声优关系业务逻辑

**数据层** (database/):
- `d1.go` - Cloudflare D1 数据库客户端
- `kv.go` - Cloudflare KV 存储客户端
- `migrations/` - 数据库迁移脚本

**模型层** (models/):
- `seiyuu.go` - 声优数据模型
- `group.go` - 群组数据模型
- `relationship.go` - 声优关系数据模型
- `admin.go` - 管理员数据模型
- `scheduler.go` - 调度器相关模型
- `errors.go` - 统一错误处理

**中间件** (middleware/):
- `auth.go` - JWT 认证中间件
- `cors.go` - 跨域配置中间件
- `logging.go` - 请求日志中间件
- `ratelimit.go` - 频率限制中间件

**工具函数** (utils/):
- `response.go` - API 响应封装
- `jwt.go` - JWT 令牌生成和验证
- `crypto.go` - 加密工具
- `validation.go` - 数据验证
- `http_client.go` - HTTP 客户端封装
- `auth.go` - 认证相关工具
- `env.go` - 环境变量工具

## 📁 项目结构

```
backend/
├── main.go                     # 应用入口
│
├── router/                     # 自定义路由系统
│   ├── router.go                  # 路由注册和匹配
│   └── context.go                 # 请求上下文封装
│
├── handlers/                   # HTTP 路由处理器
│   ├── seiyuu.go                  # 声优相关 API
│   ├── relationships.go           # 声优关系管理 API
│   └── admin.go                   # 管理员功能 API
│
├── services/                   # 业务逻辑服务
│   ├── seiyuu_service.go          # 声优业务逻辑
│   └── relationship_service.go    # 声优关系业务逻辑
│
├── models/                     # 数据模型
│   ├── seiyuu.go                  # 声优数据模型
│   ├── group.go                   # 群组数据模型
│   ├── relationship.go            # 声优关系模型
│   ├── admin.go                   # 管理员数据模型
│   ├── scheduler.go               # 调度器模型
│   └── errors.go                  # 统一错误处理
│
├── database/                   # 数据库操作
│   ├── d1.go                      # Cloudflare D1 操作
│   ├── kv.go                      # Cloudflare KV 操作
│   └── migrations/                # 数据库迁移脚本
│       ├── 001_create_seiyuu.sql
│       ├── 002_create_groups.sql
│       └── 004_create_seiyuu_relationships.sql
│
├── middleware/                 # 中间件
│   ├── auth.go                    # JWT 认证
│   ├── cors.go                    # CORS 处理
│   ├── logging.go                 # 请求日志
│   └── ratelimit.go               # 频率限制
│
├── utils/                      # 工具函数
│   ├── response.go                # API 响应封装
│   ├── jwt.go                     # JWT 工具
│   ├── crypto.go                  # 加密工具
│   ├── validation.go              # 数据验证
│   ├── http_client.go             # HTTP 客户端
│   ├── auth.go                    # 认证工具
│   └── env.go                     # 环境变量工具
│
├── build/                      # 构建输出目录
│   └── app.wasm                   # 编译后的 WASM 文件
│
├── go.mod                      # Go 模块依赖
├── go.sum                      # 依赖校验和
├── package.json                # npm 脚本配置
├── wrangler.jsonc              # Cloudflare Workers 配置
├── .dev.vars.example           # 环境变量示例文件
└── Makefile                    # Make 构建脚本
```

## 🚀 快速开始

### 前置要求

1. **Go 1.21+** - [安装 Go](https://golang.org/doc/install)
2. **Node.js 16+** - 用于运行 Wrangler
3. **Wrangler CLI** - Cloudflare Workers 开发工具

### 安装依赖

```bash
# 安装 Go 依赖
go mod download

# 全局安装 Wrangler CLI (可选,也可以使用 npx)
npm install -g wrangler

# 安装项目 npm 依赖
npm install
```

### 环境配置

1. 创建开发环境变量文件:
```bash
cp .dev.vars.example .dev.vars
```

2. 编辑 `.dev.vars` 文件,填入必需的配置:
```bash
# JWT 密钥
ADMIN_JWT_SECRET=dev-jwt-secret-key-please-change-in-production

# 管理员账号配置
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin456

# 环境标识
ENVIRONMENT=development
```

## 🛠️ 开发命令

### 构建命令

```bash
# 本地开发构建 (使用标准 Go 编译器,支持完整调试)
npm run build:local

# 生产部署构建 (使用 TinyGo,生成更小的 WASM 文件)
npm run build:deploy
```

### 开发服务器

```bash
# 启动本地开发服务器 (默认端口 8787)
npm run dev

# 或者
npm run start

# 使用本地数据库预览
wrangler dev

# 使用远程 staging 环境数据库
wrangler dev --remote
```

### 部署命令

```bash
# 部署到生产环境
npm run deploy

# 部署到 staging 环境
npm run deploy:staging
```

## 💾 数据库管理

### 创建数据库

```bash
# 创建 D1 数据库
wrangler d1 create SEIYUU_DB
```

### 执行数据库迁移

**重要**: 按顺序执行迁移脚本

```bash
# 1. 创建声优表
wrangler d1 execute SEIYUU_DB --file=database/migrations/001_create_seiyuu.sql --remote

# 2. 创建群组表
wrangler d1 execute SEIYUU_DB --file=database/migrations/002_create_groups.sql --remote

# 3. 创建声优关系表
wrangler d1 execute SEIYUU_DB --file=database/migrations/004_create_seiyuu_relationships.sql --remote
```

### 查询数据库

```bash
# 查询所有声优
wrangler d1 execute SEIYUU_DB --command="SELECT * FROM seiyuu" --remote

# 查询特定声优
wrangler d1 execute SEIYUU_DB --command="SELECT * FROM seiyuu WHERE id = 1" --remote
```

### KV 存储管理

```bash
# 列出所有 KV 键
wrangler kv:key list --binding=SEIYUU_KV

# 获取 KV 值
wrangler kv:key get "key_name" --binding=SEIYUU_KV

# 设置 KV 值
wrangler kv:key put "key_name" "value" --binding=SEIYUU_KV

# 删除 KV 键
wrangler kv:key delete "key_name" --binding=SEIYUU_KV
```

## 🌐 API 接口

### 公开接口 (无需认证)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/seiyuu` | 获取声优列表 |
| GET | `/api/seiyuu/:id` | 获取声优详情 |
| GET | `/api/groups` | 获取群组列表 |
| GET | `/api/groups/:id` | 获取群组详情 |
| GET | `/api/seiyuu/:id/relationships` | 获取声优关系 |

### 管理员接口 (需要 JWT 认证)

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/admin/login` | 管理员登录 |
| POST | `/api/admin/seiyuu` | 创建声优 |
| PUT | `/api/admin/seiyuu/:id` | 更新声优 |
| DELETE | `/api/admin/seiyuu/:id` | 删除声优 |
| POST | `/api/admin/groups` | 创建群组 |
| PUT | `/api/admin/groups/:id` | 更新群组 |
| DELETE | `/api/admin/groups/:id` | 删除群组 |
| POST | `/api/admin/relationships` | 创建声优关系 |
| PUT | `/api/admin/relationships/:id` | 更新声优关系 |
| DELETE | `/api/admin/relationships/:id` | 删除声优关系 |

### 认证方式

管理员接口需要在请求头中添加 JWT 令牌:

```bash
Authorization: Bearer <JWT_TOKEN>
```

## 🔐 环境变量配置

### 开发环境变量

在 `.dev.vars` 文件中配置 (不提交到 git):

| 变量名 | 说明 | 必需 |
|--------|------|------|
| `ADMIN_JWT_SECRET` | JWT 密钥,用于签名和验证令牌 | 是 |
| `ADMIN_USERNAME` | 管理员用户名 (默认: admin) | 否 |
| `ADMIN_PASSWORD` | 管理员密码 (默认: admin123) | 否 |
| `ENVIRONMENT` | 环境标识 (development/production) | 否 |
| `MOEGIRL_API_BASE_URL` | 萌娘百科API地址 | 否 |

### 生产环境变量

使用 `wrangler secret` 命令设置:

```bash
# 设置 JWT 密钥
wrangler secret put ADMIN_JWT_SECRET --env production

# 设置管理员账号 (可选,不设置则使用默认值)
wrangler secret put ADMIN_USERNAME --env production
wrangler secret put ADMIN_PASSWORD --env production
```

### 配置文件说明

- **wrangler.jsonc** - Cloudflare Workers 主配置文件 (支持注释)
- **.dev.vars** - 开发环境变量 (不提交到 git)
- **.dev.vars.example** - 环境变量示例文件

## 🔧 开发注意事项

### 关于路由系统

⚠️ **重要说明**:

本项目使用**自定义 router 包**而非 Gin 框架,因为 Cloudflare Workers 环境的特殊性。API 风格与 Gin 类似,但底层基于 `syumai/workers`。

**自定义路由特点**:
- 类 Gin 的路由注册方式
- 支持路由分组和中间件
- 请求上下文封装 (Context)
- 完全适配 Workers 环境

### 构建差异

**本地开发构建** (`npm run build:local`):
- 使用标准 Go 编译器
- 支持完整的调试功能
- 生成的 WASM 文件较大
- 适合开发和调试

**生产部署构建** (`npm run build:deploy`):
- 使用 TinyGo 编译器
- 生成更小的 WASM 文件
- 启动速度更快,内存占用更小
- 适合生产部署

### CORS 配置

开发环境默认允许所有跨域请求。生产环境需要在 `middleware/cors.go` 中配置允许的域名。

### 数据库表结构

**seiyuu** (声优表):
- 声优基本信息
- 性格特征、说话习惯
- 发布状态控制

**seiyuu_groups** (群组表):
- 群组配置信息
- 成员关系

**seiyuu_relationships** (关系表):
- 声优之间的关系定义
- 关系程度和描述

## 🛠️ 常见开发场景

### 添加新的 API 接口

1. 在 `handlers/` 中创建或修改处理器函数
2. 在 `main.go` 中注册路由
3. 如需要,在 `services/` 中添加业务逻辑
4. 更新 `models/` 中的数据模型

示例:
```go
// 1. 在 handlers/seiyuu.go 中添加处理器
func GetSeiyuuList(c *router.Context) {
    // 处理逻辑
}

// 2. 在 main.go 中注册路由
r.GET("/api/seiyuu", handlers.GetSeiyuuList)
```

### 添加新的数据模型

1. 在 `models/` 中创建模型文件
2. 定义结构体和方法
3. 在 `database/migrations/` 中创建迁移脚本
4. 执行迁移更新数据库

### 添加新的中间件

1. 在 `middleware/` 中创建中间件文件
2. 实现中间件函数
3. 在 `main.go` 中注册使用

示例:
```go
// 1. 在 middleware/ 中创建
func MyMiddleware() router.HandlerFunc {
    return func(c *router.Context) {
        // 中间件逻辑
        c.Next()
    }
}

// 2. 在 main.go 中使用
r.Use(middleware.MyMiddleware())
```

## 🐛 故障排查

### 开发服务器无法启动

- 检查端口 8787 是否被占用
- 确认 `.dev.vars` 文件是否存在且配置正确
- 检查 Go 依赖是否正确安装 (`go mod download`)

### 构建失败

- 确认 Go 版本 >= 1.21
- 对于生产构建,确认 TinyGo 已正确安装
- 检查 `wrangler.jsonc` 配置是否正确

### 数据库连接失败

- 确认数据库迁移是否执行
- 检查 `wrangler.jsonc` 中的数据库绑定配置
- 使用 `wrangler d1 execute` 测试数据库连接

### 部署失败

- 确认 Wrangler 已登录 (`wrangler login`)
- 检查 `wrangler.jsonc` 中的环境配置
- 确认生产环境密钥已设置

## 🔗 与前端联调

### 本地开发环境

1. 启动后端服务 (端口 8787):
   ```bash
   npm run dev
   ```

2. 前端会通过 Vite 代理将 `/api` 请求转发到 `http://localhost:8787`

3. 确保 CORS 中间件已正确配置,允许前端域名访问

### API 测试

使用 curl 测试接口:

```bash
# 获取声优列表
curl http://localhost:8787/api/seiyuu

# 管理员登录
curl -X POST http://localhost:8787/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}'

# 使用 JWT 令牌访问管理员接口
curl http://localhost:8787/api/admin/seiyuu \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

## 📋 项目特殊约定

- 使用自定义 `router` 包而非 Gin 框架
- 所有 API 响应统一使用 `utils/response.go` 中的封装函数
- 错误处理统一使用 `models/errors.go` 中定义的错误类型
- 数据库迁移必须按顺序执行
- 生产部署使用 TinyGo 构建以优化性能

## 📚 相关文档

- [Cloudflare Workers 文档](https://developers.cloudflare.com/workers/)
- [Cloudflare D1 文档](https://developers.cloudflare.com/d1/)
- [Cloudflare KV 文档](https://developers.cloudflare.com/kv/)
- [syumai/workers 文档](https://github.com/syumai/workers)
- [TinyGo 文档](https://tinygo.org/docs/)

## 📄 许可证

MIT License
