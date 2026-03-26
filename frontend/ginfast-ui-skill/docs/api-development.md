# API 开发指南

> 本指南介绍如何在 GinFast Tenant UI 中进行 API 接口开发和调用

## 概述

GinFast Tenant UI 使用 Axios 进行 HTTP 请求，提供了一套统一的 API 开发规范和工具函数。

## 目录结构

```
src/
├── api/
│   ├── types.d.ts       # 全局 API 类型定义
│   └── utils.ts         # API 工具函数
├── utils/
│   └── http/
│       ├── index.ts     # Axios 实例配置
│       └── types.d.ts   # HTTP 类型定义
```

## 基础类型定义

### API 响应类型

```typescript
// src/api/types.d.ts

/**
 * 基础 API 响应结构
 */
export interface BaseResult<T = any> {
  code: number;
  message: string;
  data: T;
}

/**
 * 分页列表响应结构
 */
export interface PageResult<T = any> {
  list: T[];
  total: number;
  pageNum: number;
  pageSize: number;
}

/**
 * 分页参数
 */
export interface PageParams {
  pageNum: number;
  pageSize: number;
}
```

## API 开发规范

### 框架 API

框架 API 放在 `src/api/` 目录下，按模块分类：

```
src/api/
├── types.d.ts           # 全局类型
├── utils.ts             # 工具函数
├── system/              # 系统模块
│   ├── user.ts
│   ├── role.ts
│   └── menu.ts
└── common/              # 通用模块
    ├── dict.ts
    └── upload.ts
```

### 插件 API

插件 API 放在 `src/plugins/{plugin}/api/` 目录下：

```
src/plugins/my-plugin/api/
└── my-api.ts
```

## API 开发模板

### 完整的 API 模块

```typescript
// src/api/system/user.ts 或 src/plugins/my-plugin/api/my-api.ts
import { http } from '@/utils/http';
import { baseUrlApi } from '@/api/utils';
import { BaseResult, PageResult, PageParams } from '@/api/types';

/**
 * 数据类型定义
 */
export interface UserData {
  id: number;
  username: string;
  nickname: string;
  email?: string;
  phone?: string;
  status: 'active' | 'inactive';
  createdAt: string;
}

/**
 * 列表查询参数
 */
export interface UserListParams extends PageParams {
  username?: string;
  status?: string;
}

/**
 * 列表响应类型
 */
export type UserListResult = BaseResult<PageResult<UserData>>;

/**
 * 单条数据响应类型
 */
export type UserResult = BaseResult<UserData>;

/**
 * 创建/更新参数
 */
export interface UserFormParams {
  username: string;
  nickname: string;
  email?: string;
  phone?: string;
  status: 'active' | 'inactive';
}

/**
 * 获取用户列表
 */
export const getUserList = (params: UserListParams) => {
  return http.request<UserListResult>(
    'get',
    baseUrlApi('system/user/list'),  // 框架 API
    // baseUrlApi('plugins/my-plugin/list'),  // 插件 API
    { params }
  );
};

/**
 * 获取用户详情
 */
export const getUserDetail = (id: number) => {
  return http.request<UserResult>(
    'get',
    baseUrlApi(`system/user/${id}`)
  );
};

/**
 * 创建用户
 */
export const createUser = (data: UserFormParams) => {
  return http.request<UserResult>(
    'post',
    baseUrlApi('system/user/add'),
    { data }
  );
};

/**
 * 更新用户
 */
export const updateUser = (data: Partial<UserFormParams> & { id: number }) => {
  return http.request<UserResult>(
    'put',
    baseUrlApi('system/user/edit'),
    { data }
  );
};

/**
 * 删除用户
 */
export const deleteUser = (id: number) => {
  return http.request<BaseResult>(
    'delete',
    baseUrlApi(`system/user/delete`),
    { data: { id } }
  );
};

/**
 * 批量删除用户
 */
export const batchDeleteUsers = (ids: number[]) => {
  return http.request<BaseResult>(
    'delete',
    baseUrlApi('system/user/batch'),
    { data: { ids } }
  );
};

/**
 * 导出用户数据
 */
export const exportUsers = (params: UserListParams) => {
  return http.request<Blob>(
    'get',
    baseUrlApi('system/user/export'),
    {
      params,
      responseType: 'blob'
    }
  );
};
```

## API 路径规范

### 框架 API 路径

```
/system/user/list
/system/role/list
/system/menu/list
/common/dict/list
```

### 插件 API 路径

```
/plugins/{plugin-name}/list
/plugins/{plugin-name}/add
/plugins/{plugin-name}/edit
/plugins/{plugin-name}/delete
/plugins/{plugin-name}/{id}
```

## HTTP 工具函数

### baseUrlApi

```typescript
// src/api/utils.ts

/**
 * 构造完整的 API URL
 * @param path API 路径
 * @returns 完整 URL
 */
export const baseUrlApi = (path: string) => {
  const base = import.meta.env.VITE_API_BASE_URL || '/api';
  return `${base}/${path}`;
};
```

### HTTP 实例

```typescript
// src/utils/http/index.ts

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { Message, Modal } from '@arco-design/web-vue';

// 创建 Axios 实例
const http: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器
http.interceptors.request.use(
  (config) => {
    // 添加 token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
http.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, message, data } = response.data;

    // 成功响应
    if (code === 200 || code === 0) {
      return data;
    }

    // 业务错误
    Message.error(message || '请求失败');
    return Promise.reject(new Error(message || '请求失败'));
  },
  (error) => {
    // HTTP 错误
    if (error.response) {
      const { status } = error.response;

      switch (status) {
        case 401:
          Message.error('未授权，请重新登录');
          // 跳转登录页
          break;
        case 403:
          Message.error('拒绝访问');
          break;
        case 404:
          Message.error('请求地址不存在');
          break;
        case 500:
          Message.error('服务器错误');
          break;
        default:
          Message.error(error.response.data?.message || '请求失败');
      }
    } else {
      Message.error('网络连接失败');
    }

    return Promise.reject(error);
  }
);

export { http };
```

## 在组件中使用 API

### 基本用法

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getUserList, createUser, updateUser, deleteUser } from '@/api/system/user';
import type { UserData, UserListParams } from '@/api/system/user';

// 数据状态
const dataList = ref<UserData[]>([]);
const loading = ref(false);
const total = ref(0);

// 获取列表
const fetchList = async (params: UserListParams) => {
  loading.value = true;
  try {
    const res = await getUserList(params);
    dataList.value = res.list;
    total.value = res.total;
  } catch (error) {
    console.error('获取列表失败:', error);
  } finally {
    loading.value = false;
  }
};

// 创建数据
const handleCreate = async (data: Partial<UserData>) => {
  try {
    await createUser(data);
    Message.success('创建成功');
    fetchList({ pageNum: 1, pageSize: 10 });
  } catch (error) {
    console.error('创建失败:', error);
  }
};

// 更新数据
const handleUpdate = async (data: Partial<UserData>) => {
  try {
    await updateUser(data);
    Message.success('更新成功');
    fetchList({ pageNum: 1, pageSize: 10 });
  } catch (error) {
    console.error('更新失败:', error);
  }
};

// 删除数据
const handleDelete = async (id: number) => {
  try {
    await deleteUser(id);
    Message.success('删除成功');
    fetchList({ pageNum: 1, pageSize: 10 });
  } catch (error) {
    console.error('删除失败:', error);
  }
};

onMounted(() => {
  fetchList({ pageNum: 1, pageSize: 10 });
});
</script>
```

### 在 Store 中使用 API

```typescript
// src/store/modules/user.ts 或 src/plugins/my-plugin/store/my-store.ts
import { defineStore } from 'pinia';
import { ref } from 'vue';
import { getUserList, createUser, updateUser, deleteUser } from '../api/user';

export const useUserStore = defineStore('user', () => {
  const dataList = ref([]);
  const loading = ref(false);

  const fetchDataList = async (params: any) => {
    loading.value = true;
    try {
      const res = await getUserList(params);
      dataList.value = res.list;
    } catch (error) {
      console.error('获取数据失败:', error);
    } finally {
      loading.value = false;
    }
  };

  const createData = async (data: any) => {
    try {
      await createUser(data);
      Message.success('创建成功');
    } catch (error) {
      console.error('创建失败:', error);
      throw error;
    }
  };

  return {
    dataList,
    loading,
    fetchDataList,
    createData
  };
});
```

## 文件上传 API

```typescript
// src/api/common/upload.ts
import { http } from '@/utils/http';
import { baseUrlApi } from '@/api/utils';
import { BaseResult } from '@/api/types';

export interface UploadResult {
  url: string;
  filename: string;
  size: number;
}

export type UploadResponse = BaseResult<UploadResult>;

/**
 * 上传文件
 */
export const uploadFile = (file: File, onProgress?: (percent: number) => void) => {
  const formData = new FormData();
  formData.append('file', file);

  return http.request<UploadResponse>(
    'post',
    baseUrlApi('common/upload'),
    {
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          onProgress(percent);
        }
      }
    }
  );
};

/**
 * 批量上传文件
 */
export const uploadFiles = (files: File[]) => {
  const formData = new FormData();
  files.forEach((file) => {
    formData.append('files', file);
  });

  return http.request<BaseResult<UploadResult[]>>(
    'post',
    baseUrlApi('common/upload/batch'),
    {
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    }
  );
};
```

## 错误处理

### 统一错误处理

```typescript
// 封装错误处理函数
export const handleApiError = (error: any, defaultMessage: string = '操作失败') => {
  console.error(error);

  let message = defaultMessage;

  if (error.response) {
    message = error.response.data?.message || message;
  } else if (error.message) {
    message = error.message;
  }

  Message.error(message);
  return Promise.reject(error);
};

// 使用示例
const fetchData = async () => {
  try {
    const res = await getUserList(params);
    // 处理响应
  } catch (error) {
    handleApiError(error, '获取数据失败');
  }
};
```

## API 最佳实践

1. **类型安全**：为所有 API 定义 TypeScript 类型
2. **统一规范**：遵循统一的命名和路径规范
3. **错误处理**：统一处理 API 错误
4. **请求取消**：在组件卸载时取消未完成的请求
5. **请求缓存**：对不常变化的数据进行缓存
6. **请求重试**：对失败请求进行适当重试

## 相关文档

- [框架应用开发指南](./framework-development.md)
- [插件应用开发指南](./plugin-development.md)
- [组件开发指南](./component-development.md)
