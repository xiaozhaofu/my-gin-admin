# 组件开发指南

> 本指南介绍如何在 GinFast Tenant UI 中开发和封装可复用组件

## 概述

GinFast Tenant UI 使用 Arco Design 作为基础 UI 组件库，同时提供了一套自定义的全局组件。本指南介绍如何开发符合项目规范的可复用组件。

## 组件分类

### 全局组件（框架组件）

放在 `src/components/` 目录下，使用 `s-` 前缀命名，全局注册，可在任何地方直接使用。

### 插件组件

放在 `src/plugins/{plugin}/components/` 目录下，仅在插件内部使用。

### 页面组件

放在 `src/views/` 或 `src/plugins/{plugin}/views/` 的 `components/` 子目录下，仅在当前页面使用。

## 全局组件开发规范

### 命名规范

- 使用 `s-` 前缀（snow 组件）
- 使用 kebab-case 命名
- 目录名与组件名一致

```
src/components/
└── s-my-component/
    ├── index.vue          # 组件主文件
    ├── README.md          # 组件文档
    └── types.ts           # 类型定义（可选）
```

### 组件模板

```vue
<!-- src/components/s-my-component/index.vue -->
<template>
  <div class="s-my-component" :class="{ 'is-disabled': disabled }">
    <!-- 组件内容 -->
    <slot />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

/**
 * 组件名称
 */
defineOptions({
  name: 'SMyComponent'
});

/**
 * 组件属性
 */
interface Props {
  /** 属性说明 */
  modelValue?: string;
  /** 是否禁用 */
  disabled?: boolean;
  /** 占位符 */
  placeholder?: string;
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  placeholder: '请输入'
});

/**
 * 组件事件
 */
interface Emits {
  (e: 'update:modelValue', value: string): void;
  (e: 'change', value: string): void;
  (e: 'focus', event: FocusEvent): void;
  (e: 'blur', event: FocusEvent): void;
}

const emit = defineEmits<Emits>();

/**
 * 组件状态
 */
const internalValue = ref(props.modelValue);

/**
 * 计算属性
 */
const isDisabled = computed(() => props.disabled);

/**
 * 监听属性变化
 */
watch(
  () => props.modelValue,
  (newVal) => {
    internalValue.value = newVal;
  }
);

/**
 * 方法
 */
const handleChange = (value: string) => {
  emit('update:modelValue', value);
  emit('change', value);
};

const handleFocus = (event: FocusEvent) => {
  emit('focus', event);
};

const handleBlur = (event: FocusEvent) => {
  emit('blur', event);
};

/**
 * 暴露给父组件的方法和属性
 */
defineExpose({
  focus: () => {
    // 聚焦逻辑
  },
  blur: () => {
    // 失焦逻辑
  }
});
</script>

<style scoped lang="scss">
.s-my-component {
  // 组件样式

  &.is-disabled {
    // 禁用状态样式
  }
}
</style>
```

### 组件文档

```markdown
<!-- src/components/s-my-component/README.md -->
# SMyComponent 我的组件

## 组件说明

这是一个示例组件的说明文档。

## 基本用法

\`\`\`vue
<template>
  <s-my-component v-model="value" />
</template>

<script setup lang="ts">
import { ref } from 'vue';

const value = ref('');
</script>
\`\`\`

## API

### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| modelValue | 绑定值 | `string` | - |
| disabled | 是否禁用 | `boolean` | `false` |
| placeholder | 占位符 | `string` | `'请输入'` |

### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| update:modelValue | 值变化时触发 | `(value: string)` |
| change | 值改变时触发 | `(value: string)` |
| focus | 获得焦点时触发 | `(event: FocusEvent)` |
| blur | 失去焦点时触发 | `(event: FocusEvent)` |

### Methods

| 方法名 | 说明 | 参数 |
|--------|------|------|
| focus | 使输入框获取焦点 | - |
| blur | 使输入框失去焦点 | - |

### Slots

| 插槽名 | 说明 |
|--------|------|
| default | 默认内容 |
```

## 常见组件类型

### 表单选择组件

```vue
<!-- src/components/s-user-select/index.vue -->
<template>
  <a-select
    v-model="internalValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    :filter-option="true"
    allow-search
    @change="handleChange"
  >
    <a-option
      v-for="item in userList"
      :key="item.id"
      :value="item.id"
      :label="item.name"
    >
      {{ item.name }}
    </a-option>
  </a-select>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';

interface User {
  id: number;
  name: string;
}

interface Props {
  modelValue?: number;
  disabled?: boolean;
  placeholder?: string;
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  placeholder: '请选择用户'
});

interface Emits {
  (e: 'update:modelValue', value: number | undefined): void;
  (e: 'change', value: number | undefined, user?: User): void;
}

const emit = defineEmits<Emits>();

const internalValue = ref<number | undefined>(props.modelValue);
const loading = ref(false);
const userList = ref<User[]>([]);

const fetchUserList = async () => {
  loading.value = true;
  try {
    // 获取用户列表
    // const res = await getUserList();
    // userList.value = res.data;
  } finally {
    loading.value = false;
  }
};

const handleChange = (value: number | undefined) => {
  const user = userList.value.find((u) => u.id === value);
  emit('update:modelValue', value);
  emit('change', value, user);
};

watch(
  () => props.modelValue,
  (newVal) => {
    internalValue.value = newVal;
  }
);

onMounted(() => {
  fetchUserList();
});
</script>
```

### 数据展示组件

```vue
<!-- src/components/s-data-card/index.vue -->
<template>
  <a-card :bordered="bordered" :loading="loading" class="s-data-card">
    <template #title>
      <slot name="title">
        <span class="card-title">{{ title }}</span>
      </slot>
    </template>
    <template #extra>
      <slot name="extra" />
    </template>
    <div class="card-content">
      <div class="data-value">
        <slot name="value">
          <span class="value-text">{{ value }}</span>
          <span v-if="unit" class="value-unit">{{ unit }}</span>
        </slot>
      </div>
      <div v-if="showTrend" class="data-trend">
        <a-tag :color="trendColor">
          <icon-arrow-up v-if="trend > 0" />
          <icon-arrow-down v-if="trend < 0" />
          {{ Math.abs(trend) }}%
        </a-tag>
        <span class="trend-label">较上周</span>
      </div>
    </div>
  </a-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  title?: string;
  value?: string | number;
  unit?: string;
  bordered?: boolean;
  loading?: boolean;
  trend?: number;
  showTrend?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  bordered: false,
  loading: false,
  trend: 0,
  showTrend: true
});

const trendColor = computed(() => {
  if (props.trend > 0) return 'green';
  if (props.trend < 0) return 'red';
  return 'gray';
});
</script>

<style scoped lang="scss">
.s-data-card {
  .card-title {
    font-size: 14px;
    color: var(--color-text-2);
  }

  .card-content {
    padding: 16px 0;

    .data-value {
      display: flex;
      align-items: baseline;
      margin-bottom: 8px;

      .value-text {
        font-size: 28px;
        font-weight: 600;
        color: var(--color-text-1);
      }

      .value-unit {
        margin-left: 8px;
        font-size: 14px;
        color: var(--color-text-3);
      }
    }

    .data-trend {
      display: flex;
      align-items: center;
      gap: 4px;

      .trend-label {
        font-size: 12px;
        color: var(--color-text-3);
      }
    }
  }
}
</style>
```

### 操作按钮组件

```vue
<!-- src/components/s-action-buttons/index.vue -->
<template>
  <a-space>
    <a-button
      v-if="showEdit"
      type="text"
      size="small"
      @click="handleEdit"
    >
      <template #icon>
        <icon-edit />
      </template>
      编辑
    </a-button>
    <a-button
      v-if="showDelete"
      type="text"
      size="small"
      status="danger"
      @click="handleDelete"
    >
      <template #icon>
        <icon-delete />
      </template>
      删除
    </a-button>
    <slot />
  </a-space>
</template>

<script setup lang="ts">
interface Props {
  showEdit?: boolean;
  showDelete?: boolean;
}

withDefaults(defineProps<Props>(), {
  showEdit: true,
  showDelete: true
});

interface Emits {
  (e: 'edit'): void;
  (e: 'delete'): void;
}

const emit = defineEmits<Emits>();

const handleEdit = () => {
  emit('edit');
};

const handleDelete = () => {
  emit('delete');
};
</script>
```

## 组件最佳实践

### 1. 使用 TypeScript 类型

```typescript
// 定义清晰的 Props 类型
interface Props {
  data: UserData[];
  loading?: boolean;
}

// 定义清晰的 Emits 类型
interface Emits {
  (e: 'update:modelValue', value: string): void;
  (e: 'submit', data: FormData): void;
}
```

### 2. 使用 Composition API

```vue
<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';

// 使用 ref 定义响应式数据
const count = ref(0);

// 使用 computed 定义计算属性
const doubleCount = computed(() => count.value * 2);

// 使用 watch 监听变化
watch(count, (newVal, oldVal) => {
  console.log(`count changed from ${oldVal} to ${newVal}`);
});

// 使用 onMounted 进行初始化
onMounted(() => {
  // 初始化逻辑
});
</script>
```

### 3. 合理使用插槽

```vue
<template>
  <div class="s-component">
    <!-- 默认插槽 -->
    <slot />

    <!-- 具名插槽 -->
    <slot name="header" />

    <!-- 作用域插槽 -->
    <slot name="item" :item="dataItem" :index="index" />
  </div>
</template>
```

### 4. 暴露方法和属性

```vue
<script setup lang="ts">
// 暴露方法给父组件
defineExpose({
  validate: () => {
    // 验证逻辑
    return true;
  },
  reset: () => {
    // 重置逻辑
  }
});
</script>
```

### 5. 样式隔离

```vue
<style scoped lang="scss">
// 使用 scoped 避免样式污染
.s-component {
  // 使用 CSS 变量
  color: var(--color-text-1);

  // 使用嵌套
  &.is-active {
    color: var(--color-primary-6);
  }

  // 使用 BEM 命名
  &__header {
    // ...
  }

  &__body {
    // ...
  }
}
</style>
```

## 组件测试

```vue
<!-- 测试组件的使用 -->
<template>
  <s-my-component
    v-model="value"
    :disabled="false"
    placeholder="请输入内容"
    @change="handleChange"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue';

const value = ref('');

const handleChange = (val: string) => {
  console.log('value changed:', val);
};
</script>
```

## 现有组件参考

项目中有许多现有组件可以作为参考：

| 组件 | 路径 | 说明 |
|------|------|------|
| s-user-select | src/components/select-user/ | 用户选择组件 |
| s-department-select | src/components/select-department/ | 部门选择组件 |
| s-excel-reader | src/components/s-excel-reader/ | Excel 读取组件 |
| s-code-view | src/components/s-code-view/ | 代码查看组件 |
| s-qrcode-draw | src/components/s-qrcode-draw/ | 二维码生成组件 |
| s-barcode-draw | src/components/s-barcode-draw/ | 条形码生成组件 |
| upload | src/components/upload/ | 文件上传组件 |
| wang-editor | src/components/wang-editor/ | 富文本编辑器 |

## 相关文档

- [框架应用开发指南](./framework-development.md)
- [插件应用开发指南](./plugin-development.md)
- [API 开发指南](./api-development.md)
- [Arco Design 文档](https://arco.design/vue/docs/start)
