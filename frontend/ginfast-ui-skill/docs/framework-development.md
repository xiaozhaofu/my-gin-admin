# 框架应用开发指南

> 本指南介绍如何在 GinFast Tenant UI 主框架基础上进行核心功能开发

## 概述

框架应用开发是指直接修改 `src/` 目录下的核心代码，适用于开发系统级功能、全局组件、布局修改等场景。

## 开发场景

### 适合框架开发的场景

| 场景 | 说明 | 示例 |
|------|------|------|
| 系统核心功能 | 修改或扩展系统核心功能 | 用户认证、权限控制、路由系统 |
| 全局组件 | 开发全局可复用的组件 | 通用表单组件、数据展示组件 |
| 布局修改 | 修改系统整体布局 | 侧边栏、顶部导航、标签页 |
| 全局配置 | 修改系统全局配置 | 主题配置、国际化配置 |
| 工具函数 | 添加全局工具函数 | 日期处理、文件处理、验证函数 |
| 全局指令 | 添加自定义 Vue 指令 | 权限指令、防抖节流指令 |

## 目录结构详解

```
src/
├── api/                 # API 接口定义
│   └── types.d.ts       # API 类型定义
├── assets/              # 静态资源
├── components/          # 全局组件（s-前缀）
├── config/              # 全局配置
├── directives/          # 自定义指令
│   ├── global/          # 全局指令
│   └── permission/      # 权限指令
├── globals/             # 全局挂载函数
├── hooks/               # Composition API Hooks
├── lang/                # 国际化
├── layout/              # 布局组件
├── plugins/             # 插件目录（独立开发）
├── router/              # 路由配置
├── store/               # Pinia 状态管理
├── style/               # 全局样式
├── typings/             # 类型声明
├── utils/               # 工具函数
└── views/               # 框架视图页面
```

## 核心开发指南

### 1. 添加全局组件

全局组件放在 `src/components/` 目录下，使用 `s-` 前缀命名：

```typescript
// src/components/s-my-component/index.vue
<template>
  <div class="s-my-component">
    <!-- 组件内容 -->
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

// 组件逻辑
</script>

<style scoped lang="scss">
.s-my-component {
  // 样式
}
</style>
```

### 2. 添加自定义指令

指令放在 `src/directives/` 目录下：

```typescript
// src/directives/modules/global/my-directive.ts
import type { Directive } from 'vue';

export const myDirective: Directive = {
  mounted(el, binding) {
    // 指令逻辑
  }
};
```

在 `src/directives/index.ts` 中注册：

```typescript
export const setupDirectives = (app: App) => {
  app.directive('my-directive', myDirective);
};
```

### 3. 添加全局工具函数

工具函数放在 `src/utils/` 目录下：

```typescript
// src/utils/my-tools.ts
/**
 * 我的工具函数
 */
export function myUtilFunction(param: string): boolean {
  // 实现逻辑
  return true;
}
```

### 4. 添加全局 Hooks

Hooks 放在 `src/hooks/` 目录下：

```typescript
// src/hooks/useMyHook.ts
import { ref } from 'vue';

export function useMyHook() {
  const data = ref<string>('');

  const fetchData = async () => {
    // 获取数据逻辑
  };

  return {
    data,
    fetchData
  };
}
```

### 5. 添加全局类型定义

类型定义放在 `src/typings/` 目录下：

```typescript
// src/typings/my-types.d.ts
export interface MyData {
  id: number;
  name: string;
}

export type MyStatus = 'active' | 'inactive';
```

### 6. 修改布局组件

布局组件放在 `src/layout/` 目录下：

```
src/layout/
├── index.vue              # 主布局入口
├── components/
│   ├── Aside/            # 侧边栏
│   ├── Header/           # 顶部导航
│   ├── Footer/           # 页脚
│   ├── Main/             # 主内容区
│   └── Tabs/             # 标签页
└── layout-*.vue          # 不同布局模式
```

### 7. 添加框架级 API

框架 API 放在 `src/api/` 目录下：

```typescript
// src/api/system/my-api.ts
import { http } from '@/utils/http';
import { baseUrlApi } from '@/api/utils';
import { BaseResult } from '@/api/types';

export interface MyData {
  id: number;
  name: string;
}

export type MyListResult = BaseResult<{
  list: MyData[];
  total: number;
}>;

export const getMyList = (params: any) => {
  return http.request<MyListResult>('get', baseUrlApi('system/my/list'), { params });
};
```

### 8. 添加框架级 Store

Store 放在 `src/store/modules/` 目录下：

```typescript
// src/store/modules/my-module.ts
import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useMyModuleStore = defineStore('my-module', () => {
  // State
  const data = ref<any[]>([]);

  // Getters
  const getData = computed(() => data.value);

  // Actions
  const fetchData = async () => {
    // 获取数据逻辑
  };

  return {
    data,
    getData,
    fetchData
  };
});
```

## 开发规范

### 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 组件 | `s-` 前缀，kebab-case | `s-user-select` |
| 指令 | `v-` 前缀，kebab-case | `v-has-perm` |
| Hooks | `use` 前缀，camelCase | `useMyHook` |
| 工具函数 | camelCase | `formatDate` |
| 类型 | PascalCase | `UserData` |
| 常量 | UPPER_SNAKE_CASE | `MAX_COUNT` |

### 文件组织

- 每个组件一个独立文件夹，包含 `index.vue` 和 `README.md`
- 相关文件放在同一目录下
- 使用 `index.ts` 作为模块导出入口

### 代码风格

- 使用 TypeScript 进行类型安全开发
- 使用 Composition API 编写组件
- 使用 `<script setup>` 语法
- 使用 SCSS 编写样式
- 遵循 ESLint 和 Prettier 规范

## 常见任务

### 修改主题配置

编辑 `src/store/modules/theme-config.ts`：

```typescript
export const useThemeConfigStore = defineStore('theme-config', () => {
  // 修改主题配置
});
```

### 添加国际化

在 `src/lang/modules/` 下添加语言文件：

```typescript
// src/lang/modules/zhCN.ts
export default {
  myModule: {
    title: '我的模块',
    // ...
  }
};
```

### 修改路由配置

编辑 `src/router/route.ts` 添加静态路由：

```typescript
export const staticRoutes: RouteRecordRaw[] = [
  {
    path: '/my-page',
    name: 'MyPage',
    component: () => import('@/views/my-page/index.vue')
  }
];
```

## 注意事项

1. **版本控制**：框架修改会影响所有插件，谨慎修改核心代码
2. **向后兼容**：修改公共 API 时保持向后兼容
3. **测试覆盖**：框架级修改需要充分测试
4. **文档更新**：修改后及时更新相关文档
5. **代码审查**：框架代码需要经过严格审查

## 相关文档

- [组件开发指南](./component-development.md)
- [API 开发指南](./api-development.md)
- [路由配置指南](./routing-guide.md)
- [权限系统指南](./permission-guide.md)
- [主题系统指南](./theme-guide.md)
