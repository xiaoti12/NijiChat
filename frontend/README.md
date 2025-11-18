# NijiChat 前端项目

Vue 3 + TypeScript + Vite 构建的虚拟声优互动模拟器前端应用。

## 项目简介

NijiChat 是一个虚拟声优互动模拟器,支持以下核心功能:

- **单人聊天** - 1v1 声优互动对话
- **双人剧场** - 双声优互动剧场模式
- **多人群聊** - 群组声优互动聊天
- **声优管理** - 声优库浏览和管理
- **本地存储** - 基于 IndexedDB 的聊天记录持久化
- **数据同步** - 支持 WebDAV 同步功能

## 🚀 快速开始

### 安装依赖
```bash
npm install
```

### 开发命令
```bash
# 启动开发服务器 (端口 5173)
npm run dev

# 构建生产版本 (快速构建,不做类型检查)
npm run build

# 构建并进行类型检查 (vue-tsc + vite build)
npm run build-with-check

# 仅进行类型检查,不构建
npm run type-check

# 预览生产构建
npm run preview
```

### 环境配置

创建 `.env` 文件配置后端 API 地址:

```bash
# 开发环境 API 地址
VITE_API_BASE_URL=http://localhost:8787
```

## 📦 技术栈

### 核心依赖
- **Vue 3.4+** - 前端框架,使用 Composition API
- **TypeScript 5.3+** - 类型安全
- **Vite 5.0+** - 构建工具和开发服务器
- **Vue Router 4.2+** - 路由管理
- **Pinia 2.1+** - 状态管理
- **Axios 1.6+** - HTTP 客户端
- **localforage 1.10+** - IndexedDB 封装库

### 开发依赖
- **@vitejs/plugin-vue** - Vite 的 Vue 3 插件
- **vue-tsc** - Vue TypeScript 编译器
- **@types/node** - Node.js 类型定义

## 📁 目录结构

```
frontend/
├── src/
│   ├── main.ts              # 应用入口
│   ├── App.vue              # 根组件
│   ├── router.ts            # 路由配置
│   │
│   ├── pages/               # 页面组件
│   │   ├── ChatHome.vue         # 聊天首页
│   │   ├── SeiyuuLibrary.vue    # 声优库
│   │   ├── ChatPage.vue         # 1v1 聊天页面
│   │   ├── DualTheater.vue      # 双人剧场
│   │   ├── GroupTheater.vue     # 群组剧场
│   │   ├── Settings.vue         # 设置页面
│   │   ├── Admin.vue            # 管理后台
│   │   └── NotFound.vue         # 404 页面
│   │
│   ├── components/          # 公共组件
│   │   ├── SeiyuuCard.vue       # 声优卡片
│   │   ├── ChatBubble.vue       # 聊天气泡
│   │   ├── ModelSelector.vue    # AI 模型选择器
│   │   └── ...
│   │
│   ├── stores/              # Pinia 状态管理
│   │   ├── configStore.ts       # 全局配置(主题、语言等)
│   │   ├── aiModelStore.ts      # AI 模型配置(API Key、模型参数)
│   │   ├── seiyuuStore.ts       # 声优数据缓存
│   │   ├── chatStore.ts         # 聊天会话状态
│   │   └── adminStore.ts        # 管理员认证状态
│   │
│   ├── services/            # 服务层(业务逻辑封装)
│   │   ├── apiService.ts            # 后端 API 调用封装
│   │   ├── aiService.ts             # AI 模型调用(直接调用 AI API)
│   │   ├── storageService.ts        # IndexedDB 本地存储管理
│   │   ├── chatPersistenceService.ts # 聊天记录持久化
│   │   ├── schedulerService.ts      # 智能调度算法(前端实现)
│   │   ├── webdavService.ts         # WebDAV 同步服务
│   │
│   ├── types/               # TypeScript 类型定义
│   │   ├── index.ts             # 类型导出入口
│   │   ├── seiyuu.ts            # 声优相关类型
│   │   ├── chat.ts              # 聊天相关类型
│   │   ├── ai.ts                # AI 模型相关类型
│   │   └── api.ts               # API 接口类型
│   │
│   ├── utils/               # 工具函数
│   │   ├── request.ts           # HTTP 请求封装
│   │   ├── crypto.ts            # 加密工具
│   │   ├── date.ts              # 日期处理
│   │   └── validation.ts        # 数据验证
│   │
│   └── styles/              # 全局样式
│       ├── variables.css        # CSS 变量
│       ├── reset.css            # 重置样式
│       └── global.css           # 全局样式
│
├── public/                  # 静态资源
├── index.html               # HTML 入口
├── vite.config.ts           # Vite 配置
├── tsconfig.json            # TypeScript 配置
└── package.json             # 项目配置
```

## 🌐 路由结构

| 路由 | 页面 | 说明 |
|------|------|------|
| `/` | ChatHome.vue | 聊天首页 |
| `/seiyuu` | SeiyuuLibrary.vue | 声优库 |
| `/chat/:seiyuuId` | ChatPage.vue | 1v1 聊天 |
| `/dual-theater` | DualTheater.vue | 双人剧场 |
| `/group-theater/:groupId?` | GroupTheater.vue | 群组剧场 |
| `/settings` | Settings.vue | 设置页面 |
| `/admin` | Admin.vue | 管理后台(需要认证) |

注: 开发环境会跳过管理员权限校验

## 🏗️ 核心架构

### 数据管理层次
1. **配置存储** (configStore.ts) - 用户配置的响应式状态管理
2. **本地存储** (storageService.ts) - IndexedDB 操作封装,管理聊天记录和房间信息
3. **API 层** - 通过 Axios 与后端通信,支持代理配置

### 存储架构
- 使用两个 IndexedDB 实例:
  - `chatDB` - 用于聊天记录
  - `roomDB` - 用于房间信息
- 支持数据导入导出功能,便于数据迁移和备份

### AI 调用架构
- AI 模型调用在前端直接完成,不经过后端
- 支持多种 AI 模型配置和切换
- 智能调度算法在前端实现

## 🎨 样式系统

### 设计系统
- **基于 CSS 变量** 的主题系统
- **主色调**: 紫色渐变 `#667eea` → `#764ba2`
- **间距系统**: 4px 基准单位
- **响应式断点**: 480px / 768px / 1024px / 1400px

### 主题切换
- 通过 `data-theme` 属性控制主题
- 支持明暗主题切换
- 配置保存在 configStore 中

## 🔧 开发规范

### TypeScript 类型安全
- 所有组件使用 `<script setup lang="ts">` 语法
- 严格类型检查,避免使用 `any`
- 统一使用 `@/types` 中定义的类型
- 导入路径使用 `@` 别名指向 `src` 目录

### 状态管理规范
- 使用 Pinia 进行状态管理
- Store 采用 Composition API 风格
- 敏感数据(如 API Key)需加密存储
- 状态更新遵循单向数据流

### API 调用规范
- 统一使用 `@/services/apiService` 调用后端接口
- 错误处理在服务层统一处理
- 支持请求/响应拦截器
- 后端 API 会通过 Vite 代理到开发服务器

### 代码组织
- 组件按功能模块组织
- 复用逻辑提取为 Composables
- 工具函数统一放在 `utils` 目录
- 避免循环依赖

## 🚢 部署

### Vercel 部署配置

1. 在 Vercel 项目设置中配置:
   - **Root Directory**: `frontend`
   - **Build Command**: `npm run build`
   - **Output Directory**: `dist`

2. 配置环境变量:
   ```
   VITE_API_BASE_URL=https://your-worker.workers.dev
   ```

3. 或使用 `vercel.json` 配置文件自动部署

### 构建优化

生产构建会自动进行代码分割:
- `vue-vendor` - Vue 相关库
- `utils` - 工具库 (axios, localforage)

### 环境变量说明

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| VITE_API_BASE_URL | 后端 API 地址 | http://localhost:8787 |

## 🔗 与后端联调

### 本地开发环境
1. 启动后端服务 (端口 8787):
   ```bash
   cd ../backend
   npm run dev
   ```

2. 启动前端服务 (端口 5173):
   ```bash
   npm run dev
   ```

3. 前端会自动代理 `/api` 请求到后端 (vite.config.ts:16-22)

### API 接口分类

**公开接口** (无需认证):
- `GET /api/seiyuu` - 获取声优列表
- `GET /api/seiyuu/:id` - 获取声优详情
- `GET /api/groups` - 获取群组列表
- `GET /api/groups/:id` - 获取群组详情
- `GET /api/seiyuu/:id/relationships` - 获取声优关系

**管理员接口** (需要 JWT 认证):
- `POST /api/admin/login` - 管理员登录
- `POST /api/admin/seiyuu` - 创建声优
- `PUT /api/admin/seiyuu/:id` - 更新声优
- `DELETE /api/admin/seiyuu/:id` - 删除声优
- (更多接口详见后端文档)

认证方式: 在请求头添加 `Authorization: Bearer <JWT_TOKEN>`

## 🎯 核心功能说明

### 聊天记录持久化
- 采用本地优先策略
- 使用 IndexedDB 持久化存储
- 支持导入导出功能
- 可选 WebDAV 云端同步

### AI 模型配置
- 支持多种 AI 模型
- API Key 加密存储
- 可自定义模型参数
- 前端直接调用 AI API

### 智能调度
- 前端实现调度算法
- 优化多轮对话体验
- 支持群聊场景

## 📋 类型系统

项目使用完整的 TypeScript 类型定义 (types/):

**声优相关**:
- `Seiyuu` - 声优信息
- `SeiyuuGroup` - 声优群组
- `SeiyuuRelationship` - 声优关系

**聊天相关**:
- `Message` - 消息
- `Room` - 房间
- `Conversation` - 对话
- `Theater` - 剧场

**AI 模型**:
- `AIModelConfig` - AI 模型配置
- `AICallOptions` - AI 调用选项

**API 接口**:
- `ApiResponse` - API 响应
- `PaginatedResponse` - 分页响应

## 🛠️ 常见开发场景

### 添加新页面
1. 在 `pages/` 中创建页面组件
2. 在 `router.ts` 中注册路由
3. 更新导航菜单

### 添加新的 API 接口
1. 在 `services/apiService.ts` 中添加调用方法
2. 在 `types/api.ts` 中定义接口类型
3. 更新相关 Store 和组件

### 添加新的声优字段
1. 修改 `types/seiyuu.ts` 中的类型定义
2. 更新 `services/apiService.ts` 中的接口调用
3. 更新相关组件和 Store

## 🐛 故障排查

### 开发服务器无法启动
- 检查端口 5173 是否被占用
- 确认依赖是否正确安装 (`npm install`)
- 查看 `vite.config.ts` 配置是否正确

### API 调用失败
- 确认后端服务是否启动 (端口 8787)
- 检查 `.env` 中的 `VITE_API_BASE_URL` 配置
- 查看浏览器控制台和网络请求

### TypeScript 类型错误
- 运行 `npm run type-check` 检查类型错误
- 确保使用了正确的类型定义
- 避免使用 `any` 类型

## 📝 项目特殊约定

- 使用 `@` 作为 `src` 目录别名
- 主题通过 `data-theme` 属性控制
- 聊天记录本地优先,使用 IndexedDB 持久化
- 开发环境会跳过管理员权限校验

## 📄 许可证

MIT License
