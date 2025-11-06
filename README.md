# NijiChat - 虚拟声优互动模拟器

基于AI技术的虚拟声优互动模拟器，为二次元爱好者提供与声优AI分身对话的沉浸式体验。

## 📋 项目概述

- **项目名称**: NijiChat
- **版本**: 1.0.0
- **技术栈**: Vue3 + TypeScript + Golang + Cloudflare Worker
- **部署**: Vercel (前端) + Cloudflare Workers (后端)

## ✨ 核心特性

- 🎭 **声优本人模拟**: 基于公开资料训练的AI角色，还原声优真实人格
- 💬 **多种互动模式**: 1v1聊天、双人剧场、群组剧场
- 🔐 **隐私安全**: 用户数据本地存储，AI配置不经过服务器
- 🤖 **智能调度**: AI自动选择最适合的声优参与对话

## 📁 项目结构

```
NijiChat/
├── .aidocs/              # 技术文档
│   ├── tdd.md           # 技术架构说明
│   └── ui-style-guide.md # UI设计规范
├── frontend/             # Vue3前端项目
│   ├── src/
│   │   ├── pages/       # 页面组件
│   │   ├── components/  # 公共组件
│   │   ├── stores/      # Pinia状态管理
│   │   ├── services/    # 服务层
│   │   ├── types/       # TypeScript类型
│   │   ├── utils/       # 工具函数
│   │   └── styles/      # 全局样式
│   ├── package.json
│   └── vite.config.ts
└── backend/              # Golang后端项目 (待创建)
```

## 🚀 快速开始

### 前端开发

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 类型检查
npm run type-check
```

访问 http://localhost:5173 查看应用

### 环境变量

复制 `.env.example` 为 `.env` 并配置：

```bash
# API基础URL
VITE_API_BASE_URL=http://localhost:8787
```

## 📖 技术文档

- [技术架构说明](./.aidocs/tdd.md) - 完整的技术架构设计文档
- [UI设计规范](./.aidocs/ui-style-guide.md) - UI组件和样式规范

## 🏗️ 架构设计

### 前端技术栈
- **框架**: Vue 3.4+ + TypeScript 5.0+
- **构建工具**: Vite 5.0+
- **状态管理**: Pinia 2.1+
- **路由**: Vue Router 4.2+
- **HTTP客户端**: Axios
- **本地存储**: LocalForage (IndexedDB)

### 后端技术栈（规划中）
- **语言**: Go 1.21+
- **框架**: Gin 1.9+
- **部署**: Cloudflare Workers
- **数据库**: Cloudflare D1 + KV

## 🎨 设计系统

采用紫色渐变主题，契合二次元虚拟声优概念：
- **主色调**: `#667eea` → `#764ba2`
- **设计风格**: 现代化、专业性、二次元风格
- **响应式**: 支持桌面端、平板和移动端

## 🔒 隐私与安全

- ✅ AI密钥完全本地存储，加密保护
- ✅ 聊天记录存储在本地IndexedDB
- ✅ 支持WebDAV云端备份（可选）
- ✅ 不经过服务器的AI调用

## 📝 开发进度

### Phase 1: 基础架构 ✅
- [x] 创建项目目录结构
- [x] 配置Vue3 + TypeScript + Vite
- [x] 实现类型系统
- [x] 创建样式系统
- [x] 实现工具函数层
- [x] 实现状态管理
- [x] 实现服务层
- [x] 创建路由和页面骨架

### Phase 2: 声优管理系统 (待开发)
- [ ] 萌娘百科API集成
- [ ] AI资料处理服务
- [ ] 声优资料管理界面
- [ ] 管理员权限控制

### Phase 3: 聊天核心功能 (待开发)
- [ ] vue-advanced-chat集成
- [ ] 1v1聊天功能
- [ ] AI模型调用
- [ ] 双声优剧场
- [ ] 群组剧场

### Phase 4: 高级功能 (待开发)
- [ ] WebDAV同步
- [ ] 主题切换
- [ ] 多语言支持

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙏 鸣谢

- Vue.js团队
- Vite团队
- 所有开源贡献者

---

**注意**: 本项目处于积极开发中，功能持续完善中。
