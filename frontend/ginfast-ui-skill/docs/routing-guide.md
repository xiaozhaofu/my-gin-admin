# 路由配置指南

> 本指南介绍如何在 GinFast Tenant UI 中配置和管理路由

## 概述

GinFast Tenant UI 使用 Vue Router 4.x 进行路由管理，支持静态路由和动态路由两种模式。

## 路由目录结构

```
src/router/
├── index.ts             # 路由实例配置
├── route.ts             # 静态路由定义
└── route-output.ts      # 路由输出工具
```

## 路由模式

### 静态路由

静态路由在代码中直接定义，不需要后端支持：

```typescript
// src/router/route.ts

import type { RouteRecordRaw } from 'vue-router';

/**
 * 静态路由
 * 这些路由不需要权限验证，直接在代码中定义
 */
export const staticRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/login.vue'),
    meta: {
      title: '登录',
      requiresAuth: false
    }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/layout/index.vue'),
    redirect: '/home',
    children: [
      {
        path: 'home',
        name: 'Home',
        component: () => import('@/views/home/home.vue'),
        meta: {
          title: '首页',
          icon: 'icon-home',
          requiresAuth: true
        }
      }
    ]
  }
];

/**
 * 无权限路由
 * 401、404、500 等错误页面
 */
export const notFoundAndNoPower: RouteRecordRaw[] = [
  {
    path: '/401',
    name: '401',
    component: () => import('@/views/error/401.vue'),
    meta: {
      title: '401',
      requiresAuth: false
    }
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/error/404.vue'),
    meta: {
      title: '404',
      requiresAuth: false
    }
  },
  {
    path: '/500',
    name: '500',
    component: () => import('@/views/error/500.vue'),
    meta: {
      title: '500',
      requiresAuth: false
    }
  }
];
```

### 动态路由

动态路由由后端返回，根据用户权限动态添加：

```typescript
// src/store/modules/route-config.ts

import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { RouteRecordRaw } from 'vue-router';
import { useRouter } from 'vue-router';

export interface MenuData {
  id: number;
  parentId: number;
  name: string;
  path: string;
  component: string;
  icon?: string;
  title: string;
  type: 'menu' | 'button' | 'link';
  redirect?: string;
  hidden?: boolean;
  meta?: any;
  children?: MenuData[];
}

export const useRouteConfigStore = defineStore('route-config', () => {
  const router = useRouter();
  const routeTree = ref<RouteRecordRaw[]>([]);

  /**
   * 初始化路由
   * 从后端获取菜单数据并转换为路由
   */
  const initSetRouter = async () => {
    try {
      // 从后端获取菜单数据
      // const menuData = await getMenuList();

      // 转换为路由格式
      const routes = transformMenuToRoutes(menuData);

      // 动态添加路由
      routes.forEach((route) => {
        router.addRoute(route);
      });

      routeTree.value = routes;
    } catch (error) {
      console.error('初始化路由失败:', error);
    }
  };

  /**
   * 将菜单数据转换为路由配置
   */
  const transformMenuToRoutes = (menus: MenuData[]): RouteRecordRaw[] => {
    return menus
      .filter((menu) => menu.type === 'menu')
      .map((menu) => {
        const route: RouteRecordRaw = {
          path: menu.path,
          name: menu.name,
          component: loadComponent(menu.component),
          meta: {
            title: menu.title,
            icon: menu.icon,
            hidden: menu.hidden,
            ...menu.meta
          }
        };

        if (menu.redirect) {
          route.redirect = menu.redirect;
        }

        if (menu.children && menu.children.length > 0) {
          route.children = transformMenuToRoutes(menu.children);
        }

        return route;
      });
  };

  /**
   * 动态加载组件
   */
  const loadComponent = (componentPath: string) => {
    // 框架组件
    if (componentPath.startsWith('views/')) {
      return () => import(`@/${componentPath}`);
    }

    // 插件组件
    if (componentPath.startsWith('plugins/')) {
      const [_, plugin, ...rest] = componentPath.split('/');
      return () => import(`@/plugins/${plugin}/views/${rest.join('/')}`);
    }

    return () => import('@/views/error/404.vue');
  };

  return {
    routeTree,
    initSetRouter
  };
});
```

## 路由元信息

```typescript
declare module 'vue-router' {
  interface RouteMeta {
    /** 页面标题 */
    title?: string;
    /** 图标 */
    icon?: string;
    /** 是否需要登录 */
    requiresAuth?: boolean;
    /** 是否隐藏 */
    hidden?: boolean;
    /** 是否缓存 */
    keepAlive?: boolean;
    /** 权限标识 */
    perms?: string[];
    /** 角色标识 */
    roles?: string[];
    /** 是否外链 */
    isLink?: boolean;
    /** 外链地址 */
    linkUrl?: string;
  }
}
```

## 添加新路由

### 添加框架路由

```typescript
// src/router/route.ts

export const staticRoutes: RouteRecordRaw[] = [
  // ... 其他路由
  {
    path: '/my-page',
    name: 'MyPage',
    component: () => import('@/views/my-page/index.vue'),
    meta: {
      title: '我的页面',
      icon: 'icon-file',
      requiresAuth: true,
      keepAlive: true
    }
  }
];
```

### 添加插件路由

插件路由通常通过后端菜单配置动态添加，也可以在插件中定义：

```typescript
// src/plugins/my-plugin/router/index.ts

import type { RouteRecordRaw } from 'vue-router';

export const myPluginRoutes: RouteRecordRaw[] = [
  {
    path: '/plugins/my-plugin',
    name: 'MyPlugin',
    component: () => import('@/layout/index.vue'),
    redirect: '/plugins/my-plugin/list',
    meta: {
      title: '我的插件',
      icon: 'icon-apps',
      requiresAuth: true
    },
    children: [
      {
        path: 'list',
        name: 'MyPluginList',
        component: () => import('../views/my-list/index.vue'),
        meta: {
          title: '插件列表',
          icon: 'icon-list',
          requiresAuth: true,
          keepAlive: true
        }
      },
      {
        path: 'detail/:id',
        name: 'MyPluginDetail',
        component: () => import('../views/my-detail/index.vue'),
        meta: {
          title: '插件详情',
          hidden: true,
          requiresAuth: true
        }
      }
    ]
  }
];
```

## 路由导航

### 编程式导航

```typescript
import { useRouter } from 'vue-router';

const router = useRouter();

// 导航到指定路径
router.push('/home');

// 带参数导航
router.push({
  path: '/user',
  query: { id: 123 }
});

// 命名路由导航
router.push({
  name: 'UserDetail',
  params: { id: 123 }
});

// 替换当前路由
router.replace('/home');

// 后退
router.back();

// 前进
router.forward();
```

### 路由参数

```typescript
import { useRoute } from 'vue-router';

const route = useRoute();

// 获取查询参数
const id = route.query.id;

// 获取路径参数
const userId = route.params.id;

// 获取路由元信息
const title = route.meta.title;
```

## 路由守卫

### 全局前置守卫

```typescript
// src/router/index.ts

router.beforeEach(async (to, from, next) => {
  // 开启进度条
  NProgress.start();

  // 检查登录状态
  const token = hasRefreshToken();

  // 未登录跳转登录页
  if (!token && to.path !== '/login') {
    return next('/login');
  }

  // 已登录访问登录页，跳转首页
  if (token && to.path === '/login') {
    return next('/home');
  }

  // 检查路由是否已加载
  if (token && !routeStore.routeTree.length) {
    await routeStore.initSetRouter();
    return next({ ...to, replace: true });
  }

  next();
});
```

### 全局后置钩子

```typescript
router.afterEach((to, from) => {
  // 关闭进度条
  NProgress.done();

  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - GinFast` : 'GinFast';

  // 更新当前路由高亮
  currentlyRoute(to);
});
```

### 组件内守卫

```vue
<script setup lang="ts">
import { onBeforeRouteLeave, onBeforeRouteUpdate } from 'vue-router';

// 离开守卫
onBeforeRouteLeave((to, from, next) => {
  const answer = window.confirm('确定要离开吗？未保存的更改将丢失。');
  if (answer) {
    next();
  } else {
    next(false);
  }
});

// 更新守卫
onBeforeRouteUpdate((to, from, next) => {
  // 当路由改变但组件被复用时调用
  next();
});
</script>
```

## 路由工具函数

```typescript
// src/hooks/useRoutingMethod.ts

import { useRouter, useRoute } from 'vue-router';
import { Message } from '@arco-design/web-vue';

export function useRoutingMethod() {
  const router = useRouter();
  const route = useRoute();

  /**
   * 判断是否是动态路由
   */
  const isDynamicRoute = (path: string) => {
    return !staticRoutes.some((r) => r.path === path);
  };

  /**
   * 打开外链
   */
  const openExternalLinks = (to: any) => {
    if (to.meta?.isLink && to.meta?.linkUrl) {
      window.open(to.meta.linkUrl, '_blank');
      return false;
    }
    return true;
  };

  /**
   * 跳转到指定路由
   */
  const navigateTo = (path: string, query?: Record<string, any>) => {
    router.push({ path, query });
  };

  /**
   * 返回上一页
   */
  const goBack = () => {
    router.back();
  };

  /**
   * 刷新当前路由
   */
  const refreshRoute = () => {
    router.replace({ ...route, query: { ...route.query, t: Date.now() } });
  };

  return {
    isDynamicRoute,
    openExternalLinks,
    navigateTo,
    goBack,
    refreshRoute
  };
}
```

## 路由最佳实践

1. **懒加载**：使用动态 import 实现路由懒加载
2. **路由分组**：相关路由放在同一个模块中
3. **命名规范**：使用有意义的路由名称
4. **元信息**：充分利用路由元信息配置页面行为
5. **错误处理**：配置 404、401、500 等错误页面
6. **权限控制**：结合路由守卫实现权限控制

## 相关文档

- [框架应用开发指南](./framework-development.md)
- [插件应用开发指南](./plugin-development.md)
- [权限系统指南](./permission-guide.md)
- [Vue Router 官方文档](https://router.vuejs.org/)
