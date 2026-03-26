---
name: ginfast-ui-skill
description: GinFast UI frontend development skill. TRIGGER when writing/modifying Vue/TypeScript files, developing components/plugins, creating APIs, configuring routes, implementing permissions, or customizing themes. DO NOT TRIGGER for backend Go development or shell scripts.
---

# GinFast UI 开发技能

> GinFast Tenant UI 项目开发技能指南，帮助开发者快速掌握框架应用开发和插件应用开发

## 概述

GinFast Tenant UI 是一个基于 Vue3 + Vite6 + TypeScript + Arco Design 的现代化后台管理系统。本技能指南提供两种开发模式的完整指导：

- **框架应用开发**：在主框架基础上进行核心功能开发，修改 `src/` 目录下的核心代码
- **插件应用开发**：开发独立的插件模块，在 `src/plugins/` 目录下创建新功能

## 技术栈

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.5.15 | 前端框架，使用 Composition API |
| Vite | 6.x | 构建工具 |
| TypeScript | 5.2.2 | 编程语言 |
| Arco Design Vue | 2.57.0 | UI 组件库 |
| Pinia | 2.3.0 | 状态管理 |
| Vue Router | 4.3.0 | 路由管理 |
| Axios | 1.6.8 | HTTP 请求 |

## 开发模式选择

### 框架应用开发

适用于以下场景：
- 修改系统核心功能
- 添加全局组件、指令、工具函数
- 修改布局、主题、路由等基础配置
- 开发系统级功能（如用户管理、权限管理等）

**入口文档**: [framework-development.md](./docs/framework-development.md)

### 插件应用开发

适用于以下场景：
- 开发独立的业务功能模块
- 需要独立维护和更新的功能
- 可插拔的功能扩展
- 多租户场景下的租户专属功能

**入口文档**: [plugin-development.md](./docs/plugin-development.md)

## 目录结构

```
ginfast-tenant-front/
├── src/
│   ├── api/                 # 框架 API 接口
│   ├── components/          # 全局组件
│   ├── config/              # 全局配置
│   ├── directives/          # 全局指令
│   ├── hooks/               # Composition API Hooks
│   ├── lang/                # 国际化
│   ├── layout/              # 布局组件
│   ├── plugins/             # 插件目录 ⭐
│   ├── router/              # 路由配置
│   ├── store/               # 状态管理
│   ├── style/               # 全局样式
│   ├── typings/             # 类型声明
│   ├── utils/               # 工具函数
│   └── views/               # 框架视图页面
└── ginfast-ui-skill/        # 开发技能文档
    └── SKILL.md             # 入口文件（本文件）
```

## 快速导航

### 核心开发指南

| 文档 | 说明 |
|------|------|
| [framework-development.md](./docs/framework-development.md) | 框架应用开发完整指南 |
| [plugin-development.md](./docs/plugin-development.md) | 插件应用开发完整指南 |
| [component-development.md](./docs/component-development.md) | 组件开发指南 |
| [api-development.md](./docs/api-development.md) | API 接口开发指南 |
| [routing-guide.md](./docs/routing-guide.md) | 路由配置指南 |
| [permission-guide.md](./docs/permission-guide.md) | 权限系统指南 |
| [theme-guide.md](./docs/theme-guide.md) | 主题系统指南 |

### 开发规范

- 使用 TypeScript 进行类型安全开发
- 遵循 ESLint 和 Prettier 代码规范
- 使用 Composition API 编写组件
- API 接口统一放在 `api/` 目录
- 组件使用 `s-` 前缀命名（如 `s-user-select`）

## 开发环境

```bash
# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建生产版本
pnpm build:prod
```

## 获取帮助

根据你的开发需求，选择对应的文档：

1. **开发框架核心功能** → 阅读 [framework-development.md](./docs/framework-development.md)
2. **开发新插件** → 阅读 [plugin-development.md](./docs/plugin-development.md)
3. **开发可复用组件** → 阅读 [component-development.md](./docs/component-development.md)
4. **配置路由和权限** → 阅读 [routing-guide.md](./docs/routing-guide.md) 和 [permission-guide.md](./docs/permission-guide.md)

---

*本技能指南持续更新中，如有问题请参考项目 README.md 或联系开发团队*
