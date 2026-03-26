# 插件应用开发指南

> 本指南介绍如何在 GinFast Tenant UI 中开发独立的插件模块

## 概述

插件应用开发是指在 `src/plugins/` 目录下创建独立的功能模块，适用于开发可插拔、独立维护的业务功能。

## 为什么选择插件开发

### 插件开发的优势

| 优势 | 说明 |
|------|------|
| **模块化** | 功能独立，职责清晰 |
| **可维护性** | 独立版本控制和维护 |
| **可插拔** | 按需启用/禁用功能 |
| **多租户支持** | 支持租户级别的功能定制 |
| **团队协作** | 不同团队独立开发不同插件 |

### 适用场景

- 独立的业务功能模块（如 CMS、CRM、数据分析）
- 需要独立维护的功能
- 可选的功能扩展
- 租户专属功能

## 插件目录结构

```
src/plugins/
└── your-plugin/          # 插件目录（小写，用连字符）
    ├── api/              # API 接口定义
    │   └── your-api.ts   # 具体接口文件
    ├── components/       # 插件专属组件
    │   └── your-component/
    │       ├── index.vue
    │       └── README.md
    ├── hooks/            # 插件专属 Hooks
    │   └── useYourHook.ts
    ├── store/            # 插件状态管理
    │   └── your-store.ts
    ├── views/            # 插件页面视图
    │   ├── your-page/
    │   │   ├── index.vue
    │   │   └── components/
    │   └── config/       # 配置页面
    │       └── components/
    ├── types/            # 插件类型定义
    │   └── index.ts
    └── README.md         # 插件说明文档
```

## 创建新插件

### 步骤 1：创建插件目录

```bash
# 在 src/plugins/ 下创建新插件目录
mkdir src/plugins/my-plugin
cd src/plugins/my-plugin

# 创建基本目录结构
mkdir api components hooks store views types
```

### 步骤 2：创建 API 接口

```typescript
// src/plugins/my-plugin/api/my-api.ts
import { http } from '@/utils/http';
import { baseUrlApi } from '@/api/utils';
import { BaseResult } from '@/api/types';

// 数据类型定义
export interface MyPluginData {
  id: number;
  name: string;
  description: string;
  status: 'active' | 'inactive';
}

// 列表返回类型
export type MyPluginListResult = BaseResult<{
  list: MyPluginData[];
  total: number;
}>;

// 列表查询参数
export interface MyPluginListParams {
  pageNum: number;
  pageSize: number;
  name?: string;
  status?: string;
}

// 单条数据返回类型
export type MyPluginResult = BaseResult<MyPluginData>;

// API 方法
export const getMyPluginList = (params: MyPluginListParams) => {
  return http.request<MyPluginListResult>(
    'get',
    baseUrlApi('plugins/my-plugin/list'),
    { params }
  );
};

export const getMyPluginDetail = (id: number) => {
  return http.request<MyPluginResult>(
    'get',
    baseUrlApi(`plugins/my-plugin/${id}`)
  );
};

export const createMyPlugin = (data: Omit<MyPluginData, 'id'>) => {
  return http.request<MyPluginResult>(
    'post',
    baseUrlApi('plugins/my-plugin/add'),
    { data }
  );
};

export const updateMyPlugin = (data: Partial<MyPluginData>) => {
  return http.request<MyPluginResult>(
    'put',
    baseUrlApi('plugins/my-plugin/edit'),
    { data }
  );
};

export const deleteMyPlugin = (id: number) => {
  return http.request<BaseResult>(
    'delete',
    baseUrlApi(`plugins/my-plugin/delete`),
    { data: { id } }
  );
};
```

### 步骤 3：创建 Store

```typescript
// src/plugins/my-plugin/store/my-plugin-store.ts
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type {
  MyPluginData,
  MyPluginListParams,
  MyPluginListResult
} from '../api/my-api';
import {
  getMyPluginList,
  createMyPlugin,
  updateMyPlugin,
  deleteMyPlugin
} from '../api/my-api';

export const useMyPluginStore = defineStore('my-plugin', () => {
  // State
  const dataList = ref<MyPluginData[]>([]);
  const loading = ref<boolean>(false);
  const total = ref<number>(0);
  const currentPage = ref<number>(1);
  const pageSize = ref<number>(10);
  const searchParams = ref<{
    name?: string;
    status?: string;
  }>({});

  // Getters
  const getDataList = computed(() => dataList.value);
  const isLoading = computed(() => loading.value);
  const getTotal = computed(() => total.value);
  const getCurrentPage = computed(() => currentPage.value);
  const getPageSize = computed(() => pageSize.value);
  const getSearchParams = computed(() => searchParams.value);

  // Actions
  const fetchDataList = async (params?: Partial<MyPluginListParams>) => {
    loading.value = true;
    try {
      if (params?.pageNum !== undefined) {
        currentPage.value = params.pageNum;
      }
      if (params?.pageSize !== undefined) {
        pageSize.value = params.pageSize;
      }

      const requestParams: MyPluginListParams = {
        pageNum: currentPage.value,
        pageSize: pageSize.value,
        ...searchParams.value
      };

      const res = await getMyPluginList(requestParams);
      if (res.data) {
        dataList.value = res.data.list || [];
        total.value = res.data.total || 0;
      }
    } catch (error) {
      console.error('获取数据失败:', error);
    } finally {
      loading.value = false;
    }
  };

  const handleSearch = (params: { name?: string; status?: string }) => {
    searchParams.value = { ...params };
    currentPage.value = 1;
    fetchDataList();
  };

  const handleReset = () => {
    searchParams.value = {};
    currentPage.value = 1;
    fetchDataList();
  };

  return {
    dataList,
    loading,
    total,
    currentPage,
    pageSize,
    searchParams,
    getDataList,
    isLoading,
    getTotal,
    getCurrentPage,
    getPageSize,
    getSearchParams,
    fetchDataList,
    handleSearch,
    handleReset
  };
});
```

### 步骤 4：创建视图页面

```vue
<!-- src/plugins/my-plugin/views/my-list/index.vue -->
<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-card title="我的插件列表" :loading="loading" :bordered="false">
        <!-- 搜索区域 -->
        <a-space wrap>
          <a-input-search
            v-model="searchForm.name"
            placeholder="请输入名称搜索"
            style="width: 240px;"
            @search="handleSearch"
            allow-clear
          />
          <a-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            style="width: 150px;"
            allow-clear
          >
            <a-option value="active">启用</a-option>
            <a-option value="inactive">禁用</a-option>
          </a-select>
          <a-button type="primary" @click="handleSearch">查询</a-button>
          <a-button @click="handleReset">重置</a-button>
          <a-button
            type="primary"
            @click="handleCreate"
            v-hasPerm="['plugins:my-plugin:add']"
          >
            <template #icon>
              <icon-plus />
            </template>
            <span>新增数据</span>
          </a-button>
        </a-space>

        <!-- 数据表格 -->
        <a-table
          :data="dataList"
          :loading="loading"
          :pagination="paginationConfig"
          :bordered="{ wrapper: true, cell: true }"
          @page-change="handlePageChange"
          @page-size-change="handlePageSizeChange"
        >
          <template #columns>
            <a-table-column title="ID" data-index="id" :width="70" align="center" />
            <a-table-column title="名称" data-index="name" :width="150" ellipsis tooltip />
            <a-table-column title="描述" data-index="description" :width="200" ellipsis tooltip />
            <a-table-column title="状态" data-index="status" :width="100">
              <template #cell="{ record }">
                <a-tag :color="record.status === 'active' ? 'green' : 'red'">
                  {{ record.status === 'active' ? '启用' : '禁用' }}
                </a-tag>
              </template>
            </a-table-column>
            <a-table-column title="操作" :width="200" fixed="right">
              <template #cell="{ record }">
                <a-space>
                  <a-button
                    size="small"
                    @click="handleEdit(record)"
                    v-hasPerm="['plugins:my-plugin:edit']"
                  >
                    编辑
                  </a-button>
                  <a-popconfirm
                    content="确定要删除这条数据吗？"
                    @ok="handleDelete(record.id)"
                  >
                    <a-button
                      size="small"
                      status="danger"
                      v-hasPerm="['plugins:my-plugin:delete']"
                    >
                      删除
                    </a-button>
                  </a-popconfirm>
                </a-space>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-card>

      <!-- 编辑/创建弹窗 -->
      <a-modal
        v-model:visible="modalVisible"
        :title="editingData.id ? '编辑数据' : '新增数据'"
        :on-before-ok="handleSave"
        @cancel="handleCancel"
      >
        <a-form :model="editingData" :rules="rules" ref="formRef">
          <a-form-item field="name" label="名称">
            <a-input v-model="editingData.name" placeholder="请输入名称" />
          </a-form-item>
          <a-form-item field="description" label="描述">
            <a-textarea
              v-model="editingData.description"
              placeholder="请输入描述"
              :max-length="200"
              show-word-limit
            />
          </a-form-item>
          <a-form-item field="status" label="状态">
            <a-radio-group v-model="editingData.status">
              <a-radio value="active">启用</a-radio>
              <a-radio value="inactive">禁用</a-radio>
            </a-radio-group>
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { useMyPluginStore } from '../../store/my-plugin-store';
import type { MyPluginData } from '../../api/my-api';
import { createMyPlugin, updateMyPlugin, deleteMyPlugin } from '../../api/my-api';
import { Message } from '@arco-design/web-vue';

const myPluginStore = useMyPluginStore();

// 搜索表单
const searchForm = reactive({
  name: '',
  status: ''
});

// 弹窗相关
const modalVisible = ref(false);
const formRef = ref();
const editingData = reactive<Partial<MyPluginData>>({
  name: '',
  description: '',
  status: 'active'
});

// 表单验证规则
const rules = {
  name: [{ required: true, message: '请输入名称' }],
  status: [{ required: true, message: '请选择状态' }]
};

// 数据列表
const dataList = computed(() => myPluginStore.getDataList);
const loading = computed(() => myPluginStore.isLoading);
const total = computed(() => myPluginStore.getTotal);
const currentPage = computed(() => myPluginStore.getCurrentPage);
const pageSize = computed(() => myPluginStore.getPageSize);

// 分页配置
const paginationConfig = computed(() => ({
  current: currentPage.value,
  pageSize: pageSize.value,
  total: total.value,
  showTotal: true,
  showPageSize: true
}));

// 搜索
const handleSearch = () => {
  myPluginStore.handleSearch(searchForm);
};

// 重置
const handleReset = () => {
  searchForm.name = '';
  searchForm.status = '';
  myPluginStore.handleReset();
};

// 新增
const handleCreate = () => {
  Object.assign(editingData, {
    id: undefined,
    name: '',
    description: '',
    status: 'active'
  });
  modalVisible.value = true;
};

// 编辑
const handleEdit = (record: MyPluginData) => {
  Object.assign(editingData, record);
  modalVisible.value = true;
};

// 保存
const handleSave = async () => {
  const valid = await formRef.value?.validate();
  if (!valid) {
    try {
      if (editingData.id) {
        await updateMyPlugin(editingData);
        Message.success('更新成功');
      } else {
        await createMyPlugin(editingData);
        Message.success('创建成功');
      }
      modalVisible.value = false;
      myPluginStore.fetchDataList();
    } catch (error) {
      console.error('保存失败:', error);
    }
  }
};

// 取消
const handleCancel = () => {
  modalVisible.value = false;
};

// 删除
const handleDelete = async (id: number) => {
  try {
    await deleteMyPlugin(id);
    Message.success('删除成功');
    myPluginStore.fetchDataList();
  } catch (error) {
    console.error('删除失败:', error);
  }
};

// 分页变化
const handlePageChange = (page: number) => {
  myPluginStore.fetchDataList({ pageNum: page });
};

const handlePageSizeChange = (size: number) => {
  myPluginStore.fetchDataList({ pageSize: size, pageNum: 1 });
};

// 初始化
onMounted(() => {
  myPluginStore.fetchDataList();
});
</script>

<style scoped lang="scss">
.snow-page {
  padding: 20px;
}
</style>
```

### 步骤 5：创建插件 README

```markdown
# My Plugin

## 描述

这是一个示例插件，展示如何开发 GinFast UI 插件。

## 功能特性

- 功能列表 1
- 功能列表 2
- 功能列表 3

## 目录结构

\`\`\`
my-plugin/
├── api/              # API 接口
├── components/       # 组件
├── hooks/            # Hooks
├── store/            # 状态管理
├── views/            # 页面视图
└── README.md         # 说明文档
\`\`\`

## 使用说明

1. 在后端配置插件路由和权限
2. 在前端菜单管理中添加插件菜单
3. 配置用户权限

## API 说明

| 接口 | 方法 | 说明 |
|------|------|------|
| /plugins/my-plugin/list | GET | 获取列表 |
| /plugins/my-plugin/{id} | GET | 获取详情 |
| /plugins/my-plugin/add | POST | 创建数据 |
| /plugins/my-plugin/edit | PUT | 更新数据 |
| /plugins/my-plugin/delete | DELETE | 删除数据 |

## 权限标识

- `plugins:my-plugin:list` - 查看列表
- `plugins:my-plugin:add` - 新增数据
- `plugins:my-plugin:edit` - 编辑数据
- `plugins:my-plugin:delete` - 删除数据
```

## 插件开发规范

### 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 插件目录 | 小写，连字符 | `my-plugin` |
| Store | `use{PluginName}Store` | `useMyPluginStore` |
| API 文件 | 小写，连字符 | `my-api.ts` |
| 权限标识 | `plugins:{plugin}:{action}` | `plugins:my-plugin:add` |

### API 路径规范

插件 API 路径统一使用 `/plugins/{plugin-name}/` 前缀：

```
/plugins/my-plugin/list
/plugins/my-plugin/add
/plugins/my-plugin/edit
/plugins/my-plugin/delete
/plugins/my-plugin/{id}
```

### 权限标识规范

插件权限标识统一使用 `plugins:{plugin-name}:{action}` 格式：

```
plugins:my-plugin:list
plugins:my-plugin:add
plugins:my-plugin:edit
plugins:my-plugin:delete
```

## 插件与框架的交互

### 使用框架工具函数

```typescript
import { http } from '@/utils/http';
import { baseUrlApi } from '@/api/utils';
import { Message } from '@arco-design/web-vue';
```

### 使用框架组件

```vue
<template>
  <s-user-select v-model="userId" />
  <s-department-select v-model="deptId" />
</template>
```

### 使用框架 Hooks

```typescript
import { useRoutingMethod } from '@/hooks/useRoutingMethod';
import { useThemeMethods } from '@/hooks/useThemeMethods';
```

## 注意事项

1. **独立性**：插件应尽量独立，减少对框架的依赖
2. **命名空间**：使用插件名作为命名空间前缀
3. **权限控制**：所有操作都应配置权限标识
4. **错误处理**：统一处理 API 错误
5. **国际化**：支持多语言的插件需要添加语言文件

## 相关文档

- [框架应用开发指南](./framework-development.md)
- [组件开发指南](./component-development.md)
- [API 开发指南](./api-development.md)
- [权限系统指南](./permission-guide.md)
