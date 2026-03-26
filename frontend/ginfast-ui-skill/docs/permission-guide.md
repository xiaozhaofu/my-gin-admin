# 权限系统指南

> 本指南介绍如何在 GinFast Tenant UI 中实现和使用权限控制

## 概述

GinFast Tenant UI 采用 RBAC（基于角色的访问控制）权限模型，支持菜单权限和按钮权限两个维度的权限控制。

## 权限模型

### RBAC 模型结构

```
用户 (User)
  ↓
角色 (Role)
  ↓
权限 (Permission)
  ↓
资源 (Resource)
```

### 权限类型

| 类型 | 说明 | 示例 |
|------|------|------|
| 菜单权限 | 控制页面访问 | `system:user:list` |
| 按钮权限 | 控制操作按钮 | `system:user:add` |
| 数据权限 | 控制数据范围 | 本部门、本部门及子部门、全部数据 |

## 权限标识规范

### 框架权限标识

```
{module}:{resource}:{action}

示例：
system:user:list      # 系统模块-用户-列表
system:user:add       # 系统模块-用户-新增
system:user:edit      # 系统模块-用户-编辑
system:user:delete    # 系统模块-用户-删除
system:user:export    # 系统模块-用户-导出
```

### 插件权限标识

```
plugins:{plugin}:{resource}:{action}

示例：
plugins:my-plugin:data:list    # 我的插件-数据-列表
plugins:my-plugin:data:add     # 我的插件-数据-新增
plugins:my-plugin:data:edit    # 我的插件-数据-编辑
plugins:my-plugin:data:delete  # 我的插件-数据-删除
```

## 权限指令

### v-has-perm 指令

用于控制按钮级别的权限：

```vue
<template>
  <a-space>
    <!-- 单个权限 -->
    <a-button
      v-has-perm="['system:user:add']"
      type="primary"
      @click="handleAdd"
    >
      新增用户
    </a-button>

    <!-- 多个权限（满足其一即可） -->
    <a-button
      v-has-perm="['system:user:edit', 'system:user:delete']"
      @click="handleAction"
    >
      操作
    </a-button>

    <!-- 多个权限（必须全部满足） -->
    <a-button
      v-has-perm-all="['system:user:edit', 'system:user:export']"
      @click="handleBatchExport"
    >
      批量导出
    </a-button>
  </a-space>
</template>
```

### v-has-role 指令

用于控制角色级别的权限：

```vue
<template>
  <!-- 单个角色 -->
  <a-button v-has-role="['admin']">
    管理员功能
  </a-button>

  <!-- 多个角色（满足其一即可） -->
  <a-button v-has-role="['admin', 'super_admin']">
    高级功能
  </a-button>
</template>
```

## 权限指令实现

```typescript
// src/directives/modules/permission/has-perm.ts

import type { Directive, DirectiveBinding } from 'vue';
import { useUserStore } from '@/store/modules/user';

/**
 * 权限指令
 * 用法：v-has-perm="['system:user:add']"
 */
export const hasPerm: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string[]>) {
    const { value } = binding;

    if (value && value instanceof Array && value.length > 0) {
      const userStore = useUserStore();
      const perms = userStore.getPerms || [];

      const hasPermission = value.some((perm) => perms.includes(perm));

      if (!hasPermission) {
        el.parentNode?.removeChild(el);
      }
    } else {
      throw new Error('需要配置权限，如 v-has-perm="[\'system:user:add\']"');
    }
  }
};

/**
 * 权限指令（全部满足）
 * 用法：v-has-perm-all="['system:user:edit', 'system:user:export']"
 */
export const hasPermAll: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string[]>) {
    const { value } = binding;

    if (value && value instanceof Array && value.length > 0) {
      const userStore = useUserStore();
      const perms = userStore.getPerms || [];

      const hasAllPermission = value.every((perm) => perms.includes(perm));

      if (!hasAllPermission) {
        el.parentNode?.removeChild(el);
      }
    } else {
      throw new Error('需要配置权限，如 v-has-perm-all="[\'system:user:edit\', \'system:user:export\']"');
    }
  }
};
```

```typescript
// src/directives/modules/permission/has-role.ts

import type { Directive, DirectiveBinding } from 'vue';
import { useUserStore } from '@/store/modules/user';

/**
 * 角色指令
 * 用法：v-has-role="['admin']"
 */
export const hasRole: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string[]>) {
    const { value } = binding;

    if (value && value instanceof Array && value.length > 0) {
      const userStore = useUserStore();
      const roles = userStore.getRoles || [];

      const hasRole = value.some((role) => roles.includes(role));

      if (!hasRole) {
        el.parentNode?.removeChild(el);
      }
    } else {
      throw new Error('需要配置角色，如 v-has-role="[\'admin\']"');
    }
  }
};
```

## 权限检查函数

### 在 JavaScript 中检查权限

```typescript
// src/utils/permission.ts

import { useUserStore } from '@/store/modules/user';

/**
 * 检查是否有指定权限
 */
export const hasPermission = (perm: string | string[]): boolean => {
  const userStore = useUserStore();
  const perms = userStore.getPerms || [];

  if (Array.isArray(perm)) {
    return perm.some((p) => perms.includes(p));
  }

  return perms.includes(perm);
};

/**
 * 检查是否有全部权限
 */
export const hasPermissionAll = (perms: string[]): boolean => {
  const userStore = useUserStore();
  const userPerms = userStore.getPerms || [];

  return perms.every((perm) => userPerms.includes(perm));
};

/**
 * 检查是否有指定角色
 */
export const hasRole = (role: string | string[]): boolean => {
  const userStore = useUserStore();
  const roles = userStore.getRoles || [];

  if (Array.isArray(role)) {
    return role.some((r) => roles.includes(r));
  }

  return roles.includes(role);
};
```

### 在组件中使用

```vue
<script setup lang="ts">
import { computed } from 'vue';
import { hasPermission, hasRole } from '@/utils/permission';

// 检查权限
const canAdd = computed(() => hasPermission('system:user:add'));
const canEdit = computed(() => hasPermission('system:user:edit'));
const canDelete = computed(() => hasPermission('system:user:delete'));

// 检查角色
const isAdmin = computed(() => hasRole('admin'));

// 条件执行
const handleAction = () => {
  if (!hasPermission('system:user:edit')) {
    Message.error('没有权限');
    return;
  }

  // 执行操作
};
</script>

<template>
  <a-space>
    <a-button v-if="canAdd" type="primary" @click="handleAdd">
      新增
    </a-button>
    <a-button v-if="canEdit" @click="handleEdit">
      编辑
    </a-button>
    <a-button v-if="canDelete" status="danger" @click="handleDelete">
      删除
    </a-button>
  </a-space>
</template>
```

## 用户 Store 中的权限数据

```typescript
// src/store/modules/user.ts

import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export interface UserInfo {
  id: number;
  username: string;
  nickname: string;
  avatar?: string;
  email?: string;
  phone?: string;
  roles: string[];
  perms: string[];
}

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo | null>(null);
  const token = ref<string>('');

  // 获取角色列表
  const getRoles = computed(() => userInfo.value?.roles || []);

  // 获取权限列表
  const getPerms = computed(() => userInfo.value?.perms || []);

  // 获取用户信息
  const getUserInfo = async () => {
    try {
      const res = await getUserInfoApi();
      userInfo.value = res.data;
      return res;
    } catch (error) {
      console.error('获取用户信息失败:', error);
      throw error;
    }
  };

  // 登出
  const logOut = () => {
    userInfo.value = null;
    token.value = '';
    // 清除其他数据
  };

  return {
    userInfo,
    token,
    getRoles,
    getPerms,
    getUserInfo,
    logOut
  };
});
```

## 路由权限控制

### 路由守卫中的权限检查

```typescript
// src/router/index.ts

router.beforeEach(async (to, from, next) => {
  // 检查是否需要登录
  if (to.meta.requiresAuth !== false) {
    const token = hasRefreshToken();
    if (!token) {
      return next('/login');
    }
  }

  // 检查路由权限
  if (to.meta.perms) {
    const userStore = useUserStore();
    const perms = userStore.getPerms || [];
    const hasPermission = to.meta.perms.some((perm: string) => perms.includes(perm));

    if (!hasPermission) {
      Message.error('没有权限访问该页面');
      return next('/401');
    }
  }

  next();
});
```

### 动态路由权限过滤

```typescript
// 根据权限过滤路由
const filterRoutesByPermission = (routes: RouteRecordRaw[], perms: string[]): RouteRecordRaw[] => {
  return routes.filter((route) => {
    // 检查路由权限
    if (route.meta?.perms) {
      const hasPermission = route.meta.perms.some((perm: string) => perms.includes(perm));
      if (!hasPermission) {
        return false;
      }
    }

    // 递归过滤子路由
    if (route.children) {
      route.children = filterRoutesByPermission(route.children, perms);
    }

    return true;
  });
};
```

## 插件权限配置

### 插件权限定义

```typescript
// src/plugins/my-plugin/types/permissions.ts

/**
 * 插件权限标识
 */
export const MY_PLUGIN_PERMISSIONS = {
  LIST: 'plugins:my-plugin:list',
  ADD: 'plugins:my-plugin:add',
  EDIT: 'plugins:my-plugin:edit',
  DELETE: 'plugins:my-plugin:delete',
  EXPORT: 'plugins:my-plugin:export'
};

/**
 * 插件权限组
 */
export const MY_PLUGIN_PERMISSION_GROUPS = {
  VIEW: [MY_PLUGIN_PERMISSIONS.LIST],
  EDIT: [MY_PLUGIN_PERMISSIONS.ADD, MY_PLUGIN_PERMISSIONS.EDIT],
  DELETE: [MY_PLUGIN_PERMISSIONS.DELETE],
  ALL: [
    MY_PLUGIN_PERMISSIONS.LIST,
    MY_PLUGIN_PERMISSIONS.ADD,
    MY_PLUGIN_PERMISSIONS.EDIT,
    MY_PLUGIN_PERMISSIONS.DELETE,
    MY_PLUGIN_PERMISSIONS.EXPORT
  ]
};
```

### 在插件中使用权限

```vue
<!-- src/plugins/my-plugin/views/my-list/index.vue -->
<script setup lang="ts">
import { MY_PLUGIN_PERMISSIONS } from '../../types/permissions';
import { hasPermission } from '@/utils/permission';

const canAdd = computed(() => hasPermission(MY_PLUGIN_PERMISSIONS.ADD));
const canEdit = computed(() => hasPermission(MY_PLUGIN_PERMISSIONS.EDIT));
const canDelete = computed(() => hasPermission(MY_PLUGIN_PERMISSIONS.DELETE));
</script>

<template>
  <a-space>
    <a-button
      v-has-perm="[MY_PLUGIN_PERMISSIONS.ADD]"
      type="primary"
      @click="handleAdd"
    >
      新增
    </a-button>
    <a-button
      v-has-perm="[MY_PLUGIN_PERMISSIONS.EDIT]"
      @click="handleEdit"
    >
      编辑
    </a-button>
    <a-button
      v-has-perm="[MY_PLUGIN_PERMISSIONS.DELETE]"
      status="danger"
      @click="handleDelete"
    >
      删除
    </a-button>
  </a-space>
</template>
```

## 权限最佳实践

1. **统一规范**：遵循统一的权限标识命名规范
2. **最小权限**：默认不给权限，按需分配
3. **前后端一致**：前后端权限标识保持一致
4. **缓存优化**：权限数据缓存，减少重复请求
5. **友好提示**：无权限时给出友好提示
6. **日志记录**：记录权限验证失败日志

## 相关文档

- [框架应用开发指南](./framework-development.md)
- [插件应用开发指南](./plugin-development.md)
- [路由配置指南](./routing-guide.md)
