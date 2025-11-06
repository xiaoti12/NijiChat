# NijiChat 前端项目

Vue3 + TypeScript + Vite 构建的现代化前端应用。

## 🚀 快速开始

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview

# TypeScript类型检查
npm run type-check
```

## 📦 依赖说明

### 核心依赖
- **vue**: ^3.4.0 - Vue 3框架
- **vue-router**: ^4.2.0 - 路由管理
- **pinia**: ^2.1.0 - 状态管理
- **axios**: ^1.6.0 - HTTP客户端
- **localforage**: ^1.10.0 - IndexedDB封装

### 开发依赖
- **vite**: ^5.0.0 - 构建工具
- **typescript**: ^5.3.0 - TypeScript支持
- **@vitejs/plugin-vue**: ^5.0.0 - Vue插件
- **vue-tsc**: ^1.8.0 - Vue TypeScript编译器

## 📁 目录结构

```
src/
├── main.ts           # 应用入口
├── App.vue           # 根组件
├── router.ts         # 路由配置
├── pages/            # 页面组件
│   ├── SeiyuuLibrary.vue  # 声优库主页
│   ├── ChatPage.vue       # 聊天页面
│   ├── Settings.vue       # 设置页面
│   └── Admin.vue          # 管理后台
├── components/       # 公共组件
├── stores/           # Pinia状态管理
│   ├── seiyuuStore.ts    # 声优数据
│   ├── chatStore.ts      # 聊天记录
│   ├── configStore.ts    # 用户配置
│   ├── aiModelStore.ts   # AI模型配置
│   └── adminStore.ts     # 管理员状态
├── services/         # 服务层
│   ├── apiService.ts     # 后端API
│   ├── aiService.ts      # AI模型调用
│   ├── storageService.ts # 本地存储
│   └── webdavService.ts  # WebDAV同步
├── types/            # TypeScript类型
│   ├── seiyuu.ts
│   ├── chat.ts
│   ├── ai.ts
│   └── api.ts
├── utils/            # 工具函数
│   ├── request.ts    # HTTP封装
│   ├── crypto.ts     # 加密工具
│   ├── date.ts       # 日期处理
│   └── validation.ts # 数据验证
└── styles/           # 全局样式
    ├── variables.css # CSS变量
    ├── reset.css     # 重置样式
    └── global.css    # 全局样式
```

## 🎨 样式系统

基于CSS变量的设计系统：
- **主色调**: 紫色渐变 `#667eea` → `#764ba2`
- **间距系统**: 4px基准单位
- **断点**: 480px / 768px / 1024px / 1400px

详见 [UI设计规范](../.aidocs/ui-style-guide.md)

## 🔧 开发规范

### 类型安全
- 所有组件使用TypeScript `<script setup lang="ts">`
- 严格类型检查，避免使用`any`
- 使用`@/types`中定义的类型

### 状态管理
- 使用Pinia进行状态管理
- Store采用Composition API风格
- 敏感数据加密存储

### API调用
- 统一使用`@/services/apiService`
- 错误处理在服务层统一处理
- 支持请求/响应拦截器

## 🚢 部署

### Vercel部署

1. 在Vercel中设置Root Directory为`frontend`
2. Build Command: `npm run build`
3. Output Directory: `dist`

或使用`vercel.json`配置文件自动部署。

### 环境变量

生产环境需要配置：
```
VITE_API_BASE_URL=https://your-worker.workers.dev
```

## 📝 待实现功能

- [ ] vue-advanced-chat集成
- [ ] 聊天界面完整实现
- [ ] AI模型配置UI
- [ ] WebDAV同步UI
- [ ] 主题切换功能
- [ ] 移动端优化

## 🐛 已知问题

目前无已知问题。

## 📄 许可证

MIT License
