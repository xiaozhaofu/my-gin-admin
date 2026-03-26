# 主题系统指南

> 本指南介绍如何在 GinFast Tenant UI 中进行主题定制和样式管理

## 概述

GinFast Tenant UI 基于 Arco Design Vue 2.57.0，支持暗黑模式、自定义主题色等多种主题定制方式。

## 主题配置

### 主题 Store

```typescript
// src/store/modules/theme-config.ts

import { defineStore } from 'pinia';
import { ref } from 'vue';

export interface ThemeConfig {
  /** 主题模式 */
  mode: 'light' | 'dark';
  /** 主题色 */
  primaryColor: string;
  /** 侧边栏宽度 */
  sideBarWidth: string;
  /** 侧边栏暗色 */
  sideBarDark: boolean;
  /** 顶栏暗色 */
  headerDark: boolean;
  /** 页脚显示 */
  showFooter: boolean;
  /** 标签页显示 */
  showTabs: boolean;
  /** 面包屑显示 */
  showBreadcrumb: boolean;
}

export const useThemeConfigStore = defineStore('theme-config', () => {
  // 默认主题配置
  const themeConfig = ref<ThemeConfig>({
    mode: 'light',
    primaryColor: '#165DFF',
    sideBarWidth: '220px',
    sideBarDark: false,
    headerDark: false,
    showFooter: true,
    showTabs: true,
    showBreadcrumb: true
  });

  /**
   * 设置主题模式
   */
  const setMode = (mode: 'light' | 'dark') => {
    themeConfig.value.mode = mode;
    document.body.setAttribute('arco-theme', mode);
  };

  /**
   * 设置主题色
   */
  const setPrimaryColor = (color: string) => {
    themeConfig.value.primaryColor = color;
    document.body.style.setProperty('--primary-color', color);
  };

  /**
   * 设置侧边栏宽度
   */
  const setSideBarWidth = (width: string) => {
    themeConfig.value.sideBarWidth = width;
  };

  /**
   * 切换侧边栏暗色
   */
  const toggleSideBarDark = () => {
    themeConfig.value.sideBarDark = !themeConfig.value.sideBarDark;
  };

  /**
   * 切换顶栏暗色
   */
  const toggleHeaderDark = () => {
    themeConfig.value.headerDark = !themeConfig.value.headerDark;
  };

  /**
   * 切换页脚显示
   */
  const toggleShowFooter = () => {
    themeConfig.value.showFooter = !themeConfig.value.showFooter;
  };

  /**
   * 切换标签页显示
   */
  const toggleShowTabs = () => {
    themeConfig.value.showTabs = !themeConfig.value.showTabs;
  };

  /**
   * 切换面包屑显示
   */
  const toggleShowBreadcrumb = () => {
    themeConfig.value.showBreadcrumb = !themeConfig.value.showBreadcrumb;
  };

  return {
    themeConfig,
    setMode,
    setPrimaryColor,
    setSideBarWidth,
    toggleSideBarDark,
    toggleHeaderDark,
    toggleShowFooter,
    toggleShowTabs,
    toggleShowBreadcrumb
  };
});
```

## 主题变量

### Arco Design 主题变量

```scss
// src/style/var/global-theme.scss

// 主题色
$primary-color: #165DFF;

// 成功色
$success-color: #00B42A;

// 警告色
$warning-color: #FF7D00;

// 错误色
$danger-color: #F53F3F;

// 链接色
$link-color: #165DFF;

// 文本色
$text-color-primary: #1D2129;
$text-color-secondary: #4E5969;
$text-color-tertiary: #86909C;
$text-color-quaternary: #C9CDD4;

// 背景色
$bg-color-primary: #FFFFFF;
$bg-color-secondary: #F2F3F5;
$bg-color-tertiary: #E5E6EB;
$bg-color-quaternary: #C9CDD4;

// 边框色
$border-color-primary: #E5E6EB;
$border-color-secondary: #F2F3F5;

// 圆角
$border-radius-small: 2px;
$border-radius-medium: 4px;
$border-radius-large: 8px;
$border-radius-circle: 50%;

// 阴影
$shadow-light: 0 0 1px rgba(0, 0, 0, 0.3);
$shadow-medium: 0 2px 8px rgba(0, 0, 0, 0.08);
$shadow-heavy: 0 8px 24px rgba(0, 0, 0, 0.12);
```

### 自定义主题变量

```scss
// src/style/var/index.scss

// 布局
$layout-header-height: 60px;
$layout-footer-height: 48px;
$layout-side-bar-width: 220px;
$layout-side-bar-collapsed-width: 48px;

// 标签页
$tabs-height: 40px;
$tabs-bar-height: 36px;

// 过渡动画
$transition-duration: 0.3s;
$transition-timing-function: cubic-bezier(0.34, 0.69, 0.1, 1);

// Z-index
$z-index-dropdown: 1050;
$z-index-sticky: 1020;
$z-index-fixed: 1030;
$z-index-modal-backdrop: 1040;
$z-index-modal: 1050;
$z-index-popover: 1060;
$z-index-tooltip: 1070;
```

## 暗黑模式

### 暗黑模式样式

```scss
// src/style/var/dark-theme.scss

// 暗黑模式主题色
$dark-primary-color: #4080FF;

// 暗黑模式文本色
$dark-text-color-primary: #E5E6EB;
$dark-text-color-secondary: #C9CDD4;
$dark-text-color-tertiary: #86909C;

// 暗黑模式背景色
$dark-bg-color-primary: #17171A;
$dark-bg-color-secondary: #232324;
$dark-bg-color-tertiary: #2C2C2E;

// 暗黑模式边框色
$dark-border-color-primary: #2C2C2E;
$dark-border-color-secondary: #3A3A3C;
```

### 切换暗黑模式

```vue
<!-- src/layout/components/Header/components/theme-settings/index.vue -->
<script setup lang="ts">
import { useThemeConfigStore } from '@/store/modules/theme-config';

const themeStore = useThemeConfigStore();

const handleToggleTheme = () => {
  const newMode = themeStore.themeConfig.mode === 'light' ? 'dark' : 'light';
  themeStore.setMode(newMode);
};
</script>

<template>
  <a-switch
    :model-value="themeStore.themeConfig.mode === 'dark'"
    @change="handleToggleTheme"
  >
    <template #checked-icon>
      <icon-moon />
    </template>
    <template #unchecked-icon>
      <icon-sun />
    </template>
  </a-switch>
</template>
```

## 自定义主题色

### 使用 CSS 变量

```scss
// src/style/index.scss

:root {
  // 主题色
  --primary-color: #165DFF;
  --success-color: #00B42A;
  --warning-color: #FF7D00;
  --danger-color: #F53F3F;

  // 文本色
  --text-color-primary: #1D2129;
  --text-color-secondary: #4E5969;
  --text-color-tertiary: #86909C;

  // 背景色
  --bg-color-primary: #FFFFFF;
  --bg-color-secondary: #F2F3F5;
  --bg-color-tertiary: #E5E6EB;

  // 边框色
  --border-color-primary: #E5E6EB;
  --border-color-secondary: #F2F3F5;
}

// 暗黑模式
[arco-theme='dark'] {
  --primary-color: #4080FF;
  --text-color-primary: #E5E6EB;
  --text-color-secondary: #C9CDD4;
  --bg-color-primary: #17171A;
  --bg-color-secondary: #232324;
  --border-color-primary: #2C2C2E;
}
```

### 动态修改主题色

```typescript
// src/hooks/useThemeMethods.ts

import { watch } from 'vue';
import { useThemeConfigStore } from '@/store/modules/theme-config';

export function useThemeMethods() {
  const themeStore = useThemeConfigStore();

  /**
   * 初始化主题
   */
  const initTheme = () => {
    // 设置主题模式
    document.body.setAttribute('arco-theme', themeStore.themeConfig.mode);

    // 设置主题色
    document.body.style.setProperty('--primary-color', themeStore.themeConfig.primaryColor);
  };

  /**
   * 监听主题配置变化
   */
  watch(
    () => themeStore.themeConfig,
    (config) => {
      document.body.setAttribute('arco-theme', config.mode);
      document.body.style.setProperty('--primary-color', config.primaryColor);
    },
    { deep: true }
  );

  return {
    initTheme
  };
}
```

## 布局主题

### 侧边栏主题

```vue
<!-- src/layout/components/Aside/index.vue -->
<script setup lang="ts">
import { computed } from 'vue';
import { useThemeConfigStore } from '@/store/modules/theme-config';

const themeStore = useThemeConfigStore();

const asideTheme = computed(() => {
  return themeStore.themeConfig.sideBarDark ? 'dark' : 'light';
});
</script>

<template>
  <a-menu
    :theme="asideTheme"
    :style="{ width: themeStore.themeConfig.sideBarWidth }"
  >
    <!-- 菜单项 -->
  </a-menu>
</template>
```

### 顶栏主题

```vue
<!-- src/layout/components/Header/index.vue -->
<script setup lang="ts">
import { computed } from 'vue';
import { useThemeConfigStore } from '@/store/modules/theme-config';

const themeStore = useThemeConfigStore();

const headerClass = computed(() => {
  return {
    'header-dark': themeStore.themeConfig.headerDark
  };
});
</script>

<template>
  <div class="layout-header" :class="headerClass">
    <!-- 顶栏内容 -->
  </div>
</template>

<style scoped lang="scss">
.layout-header {
  background: var(--bg-color-primary);

  &.header-dark {
    background: var(--dark-bg-color-secondary);
  }
}
</style>
```

## 响应式设计

### 断点变量

```scss
// src/style/var/media.scss

// 断点
$breakpoint-xs: 480px;
$breakpoint-sm: 576px;
$breakpoint-md: 768px;
$breakpoint-lg: 992px;
$breakpoint-xl: 1200px;
$breakpoint-xxl: 1600px;

// 媒体查询 Mixin
@mixin respond-to($breakpoint) {
  @if $breakpoint == 'xs' {
    @media (max-width: $breakpoint-xs) { @content; }
  }
  @else if $breakpoint == 'sm' {
    @media (max-width: $breakpoint-sm) { @content; }
  }
  @else if $breakpoint == 'md' {
    @media (max-width: $breakpoint-md) { @content; }
  }
  @else if $breakpoint == 'lg' {
    @media (max-width: $breakpoint-lg) { @content; }
  }
  @else if $breakpoint == 'xl' {
    @media (max-width: $breakpoint-xl) { @content; }
  }
  @else if $breakpoint == 'xxl' {
    @media (max-width: $breakpoint-xxl) { @content; }
  }
}
```

### 响应式布局

```scss
// src/style/media/layout.scss

.sidebar {
  width: 220px;

  @include respond-to('md') {
    width: 180px;
  }

  @include respond-to('sm') {
    width: 100%;
    position: fixed;
    left: -100%;
    transition: left 0.3s;

    &.active {
      left: 0;
    }
  }
}
```

## 主题最佳实践

1. **使用 CSS 变量**：优先使用 CSS 变量定义主题色
2. **统一命名规范**：使用统一的命名规范定义变量
3. **暗黑模式支持**：确保所有组件都支持暗黑模式
4. **响应式设计**：考虑不同屏幕尺寸的显示效果
5. **性能优化**：避免频繁修改 DOM 样式
6. **持久化配置**：将用户主题配置保存到本地

## 相关文档

- [框架应用开发指南](./framework-development.md)
- [组件开发指南](./component-development.md)
- [Arco Design 主题定制](https://arco.design/vue/docs/theme)
