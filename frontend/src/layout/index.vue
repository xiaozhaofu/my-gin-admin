<template>
  <a-layout class="shell">
    <a-layout-sider class="shell-sider" collapsible breakpoint="lg" :collapsed="collapsed" @collapse="collapsed = $event">
      <div class="sider-inner">
        <div class="brand">
          <div class="brand-mark">
            <img :src="logoURL" alt="logo" />
          </div>
          <div v-if="!collapsed">
            <div class="brand-title">Content Admin</div>
            <div class="brand-sub">内容管理控制台</div>
          </div>
        </div>

        <div v-if="!collapsed" class="sider-section">
          <div class="section-label">导航分组</div>
          <div class="section-value">按业务域访问页面与功能</div>
        </div>

        <a-scrollbar outer-class="sider-scrollbar">
          <a-menu
            class="sider-menu"
            :selected-keys="[route.path]"
            :default-open-keys="defaultOpenKeys"
            auto-open-selected
            @menu-item-click="goTo"
          >
            <template v-for="group in visibleMenuGroups" :key="group.key">
              <a-menu-item v-if="!group.children?.length" :key="group.path">
                <template #icon>
                  <component :is="group.icon" />
                </template>
                {{ group.title }}
              </a-menu-item>
              <a-sub-menu v-else :key="group.key">
                <template #icon>
                  <component :is="group.icon" />
                </template>
                <template #title>{{ group.title }}</template>
                <a-menu-item v-for="item in group.children" :key="item.path">
                  {{ item.title }}
                </a-menu-item>
              </a-sub-menu>
            </template>
          </a-menu>
        </a-scrollbar>
      </div>
    </a-layout-sider>
    <a-layout class="shell-main">
      <a-layout-header class="header">
        <div class="header-left">
          <div class="header-title-group">
            <div class="header-title">{{ currentTitle }}</div>
            <div class="header-sub">按业务域分组展示页面与能力入口</div>
          </div>
        </div>
        <a-space size="large" class="header-right">
          <a-space class="quick-links" size="small">
            <a-button size="small" type="text" @click="goTo('/home')">工作台</a-button>
            <a-button size="small" type="text" @click="goTo('/articles')">文章</a-button>
            <a-button size="small" type="text" @click="goTo('/uploads')">资源</a-button>
          </a-space>
          <a-dropdown @select="onQuickActionSelect">
            <a-button class="header-action" type="outline">
              <template #icon>
                <icon-plus />
              </template>
              快捷操作
            </a-button>
            <template #content>
              <a-doption value="/articles/new">新增文章</a-doption>
              <a-doption value="/articles/batch">批量新增</a-doption>
              <a-doption value="/content-menus">内容菜单</a-doption>
              <a-doption value="/uploads">上传资源</a-doption>
              <a-doption value="/system/online-users">在线用户</a-doption>
            </template>
          </a-dropdown>
          <a-popover trigger="click" position="bl">
            <div class="notice-card notice-trigger">
              <div class="notice-main">
                <div class="notice-title">系统提示</div>
                <div class="notice-text">当前界面已按业务域进行分组，可从左侧分区快速进入目标页面。</div>
              </div>
              <a-badge :count="noticeItems.length" :max-count="9">
                <icon-notification class="notice-icon" />
              </a-badge>
            </div>
            <template #content>
              <div class="notice-panel">
                <div class="notice-panel-head">
                  <span>待处理提醒</span>
                  <a-tag size="small" color="arcoblue">{{ noticeItems.length }} 条</a-tag>
                </div>
                <div v-for="item in noticeItems" :key="item.title" class="notice-panel-item" @click="handleNoticeClick(item.path)">
                  <div class="notice-panel-item-top">
                    <span>{{ item.title }}</span>
                    <a-tag size="small" :color="item.color">{{ item.label }}</a-tag>
                  </div>
                  <div class="notice-panel-item-text">{{ item.text }}</div>
                </div>
              </div>
            </template>
          </a-popover>
          <div class="welcome-card">
            <div class="welcome-label">当前登录</div>
            <div class="welcome-value">{{ session.profile?.nickname || session.profile?.username || "未登录" }}</div>
          </div>
          <a-tag color="arcoblue">{{ session.profile?.username }}</a-tag>
          <a-button type="outline" @click="logout">退出</a-button>
        </a-space>
      </a-layout-header>
      <a-layout-content class="content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { routes } from "@/router";
import { useSessionStore } from "@/store/modules/session";
import logoURL from "@/assets/logo/snow.svg";
import {
  IconApps,
  IconDashboard,
  IconList,
  IconNotification,
  IconPlus,
  IconSettings
} from "@arco-design/web-vue/es/icon";

const route = useRoute();
const router = useRouter();
const session = useSessionStore();
const collapsed = ref(false);

type FlatMenu = {
  path: string;
  title: string;
  permission?: string;
};

type MenuGroup = {
  key: string;
  title: string;
  icon: any;
  path?: string;
  children?: FlatMenu[];
};

type HeaderNotice = {
  title: string;
  text: string;
  label: string;
  color: "arcoblue" | "green" | "orangered";
  path: string;
};

const menus: FlatMenu[] = routes
  .find(item => item.path === "/")
  ?.children?.filter(item => !item.meta?.hidden)
  .map(item => ({ path: `/${item.path}`, title: item.meta?.title as string, permission: item.meta?.permission as string | undefined })) || [];

const noticeItems: HeaderNotice[] = [
  {
    title: "内容菜单",
    text: "客户端内容菜单和文章入口已经打通，新增或调整内容时优先维护这里。",
    label: "内容",
    color: "arcoblue",
    path: "/content-menus"
  },
  {
    title: "资源上传",
    text: "上传页和文章资源选择都支持按云存储类型筛选和上传。",
    label: "资源",
    color: "green",
    path: "/uploads"
  },
  {
    title: "在线会话",
    text: "上线前先检查在线用户和操作日志，避免会话和权限问题带着进生产。",
    label: "巡检",
    color: "orangered",
    path: "/system/online-users"
  }
];

const menuGroups = computed<MenuGroup[]>(() => {
  const home = menus.find(item => item.path === "/home");
  const contentMenus = menus.filter(item => ["/articles", "/articles/batch", "/uploads", "/content-menus"].includes(item.path));
  const systemMenus = menus.filter(item => item.path.startsWith("/system/"));
  const capabilityMenus = menus.filter(item => item.path === "/modules");

  return [
    { key: "workbench", title: "工作台", icon: IconDashboard, path: home?.path, children: [] },
    { key: "content", title: "内容管理", icon: IconApps, children: contentMenus },
    { key: "system", title: "系统管理", icon: IconSettings, children: systemMenus },
    { key: "capability", title: "通用能力", icon: IconList, children: capabilityMenus }
  ].filter(group => group.path || group.children?.length);
});

const visibleMenuGroups = computed(() =>
  menuGroups.value
    .map(group => ({
      ...group,
      children: group.children?.filter(item => !item.permission || session.hasPermission(item.permission))
    }))
    .filter(group => group.path || group.children?.length)
    .map(group => {
      if (group.path) {
        return group;
      }
      if ((group.children?.length || 0) === 1) {
        return {
          ...group,
          path: group.children?.[0]?.path,
          children: []
        };
      }
      return group;
    })
);

const defaultOpenKeys = computed(() =>
  visibleMenuGroups.value
    .filter(group => (group.children?.length || 0) > 0)
    .map(group => group.key)
);

const currentTitle = computed(() => (route.meta.title as string) || "Content Admin");

const goTo = (key: string) => {
  if (!key || key === route.path) {
    return;
  }
  router.push(key);
};
const onQuickActionSelect = (key: string | number) => goTo(String(key));
const handleNoticeClick = (path: string) => goTo(path);
const logout = () => {
  session.logout();
  router.push("/login");
};
</script>

<style scoped lang="scss">
.shell {
  min-height: 100vh;
  height: 100vh;
  background: linear-gradient(180deg, #f5f7fb 0%, #eef3fb 100%);
}

.shell :deep(.arco-layout) {
  min-height: 0;
}

.shell-sider {
  background: linear-gradient(180deg, #f7f9fd 0%, #edf2fb 100%);
  border-right: 1px solid var(--color-border-2);
}

.sider-inner {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.brand {
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 20px 16px 18px;
  color: #fff;
}

.brand-mark {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  border-radius: 14px;
  background: linear-gradient(135deg, #2f54eb 0%, #1890ff 100%);
  box-shadow: 0 10px 24px rgb(47 84 235 / 28%);

  img {
    width: 34px;
    height: 34px;
  }
}

.brand-title {
  font-size: 16px;
  font-weight: 700;
  color: #1d2129;
}

.brand-sub {
  font-size: 12px;
  color: #86909c;
}

.sider-section {
  margin: 0 16px 12px;
  padding: 14px 16px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.88);
  box-shadow: 0 8px 18px rgb(15 23 42 / 4%);
}

.section-label {
  font-size: 12px;
  color: #86909c;
}

.section-value {
  margin-top: 6px;
  font-size: 13px;
  color: #1d2129;
  font-weight: 600;
}

.sider-scrollbar {
  flex: 1;
  min-height: 0;
}

.sider-menu {
  padding: 0 10px 16px;
  background: transparent;
}

.sider-menu :deep(.arco-menu-item),
.sider-menu :deep(.arco-menu-inline-header) {
  border-radius: 12px;
}

.sider-menu :deep(.arco-menu-item) {
  margin: 2px 0;
  color: #4e5969;
  font-weight: 500;
  transition: all 0.2s ease;
}

.sider-menu :deep(.arco-menu-item:hover) {
  background: rgba(22, 93, 255, 0.06);
  color: #1d2129;
}

.sider-menu :deep(.arco-menu-item .arco-icon) {
  color: inherit;
}

.sider-menu :deep(.arco-menu-inline-header) {
  margin: 4px 0 8px;
  background: rgba(22, 93, 255, 0.08);
  color: #165dff;
  font-weight: 700;
}

.sider-menu :deep(.arco-menu-inline-header:hover) {
  background: rgba(22, 93, 255, 0.12);
  color: #165dff;
}

.sider-menu :deep(.arco-menu-inline-header .arco-icon) {
  color: #165dff;
}

.sider-menu :deep(.arco-menu-item.arco-menu-selected) {
  background: linear-gradient(90deg, rgba(22, 93, 255, 0.16) 0%, rgba(22, 93, 255, 0.1) 100%);
  box-shadow: inset 3px 0 0 #165dff;
  color: #165dff;
  font-weight: 700;
}

.sider-menu :deep(.arco-menu-item.arco-menu-selected:hover) {
  background: linear-gradient(90deg, rgba(22, 93, 255, 0.2) 0%, rgba(22, 93, 255, 0.12) 100%);
  color: #165dff;
}

.sider-menu :deep(.arco-menu-item.arco-menu-selected .arco-icon),
.sider-menu :deep(.arco-menu-item.arco-menu-selected .arco-menu-icon),
.sider-menu :deep(.arco-menu-item.arco-menu-selected .arco-menu-item-inner),
.sider-menu :deep(.arco-menu-item.arco-menu-selected .arco-menu-item-content) {
  color: #165dff;
}

.sider-menu :deep(.arco-menu-item.arco-menu-selected::after) {
  display: none;
}

.shell-main {
  min-width: 0;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px 0 28px;
  background: rgba(255, 255, 255, 0.9);
  border-bottom: 1px solid var(--color-border-2);
  backdrop-filter: blur(12px);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-title-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.header-title {
  font-size: 18px;
  font-weight: 700;
  color: #1d2129;
}

.header-sub {
  font-size: 12px;
  color: #86909c;
}

.header-right {
  flex-wrap: wrap;
}

.header-action {
  border-radius: 12px;
}

.quick-links {
  padding: 6px 8px;
  border-radius: 12px;
  background: rgba(15, 23, 42, 0.03);
}

.notice-card {
  display: flex;
  align-items: center;
  gap: 12px;
  max-width: 280px;
  padding: 8px 12px;
  border-radius: 12px;
  background: rgba(255, 125, 0, 0.08);
}

.notice-trigger {
  cursor: pointer;
}

.notice-main {
  min-width: 0;
}

.notice-title {
  font-size: 11px;
  color: #b75d00;
  font-weight: 700;
}

.notice-text {
  margin-top: 2px;
  font-size: 12px;
  line-height: 1.5;
  color: #7a4a0d;
}

.notice-icon {
  font-size: 18px;
  color: #b75d00;
}

.notice-panel {
  width: 320px;
  display: grid;
  gap: 10px;
}

.notice-panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
  font-weight: 700;
  color: #1d2129;
}

.notice-panel-item {
  padding: 12px;
  border-radius: 12px;
  background: #f7f9fc;
  transition: background 0.2s ease;
  cursor: pointer;
}

.notice-panel-item:hover {
  background: #edf3ff;
}

.notice-panel-item-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-size: 13px;
  font-weight: 600;
  color: #1d2129;
}

.notice-panel-item-text {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.6;
  color: #4e5969;
}

.welcome-card {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 8px 12px;
  border-radius: 12px;
  background: rgba(22, 93, 255, 0.06);
}

.welcome-label {
  font-size: 11px;
  color: #86909c;
}

.welcome-value {
  font-size: 13px;
  font-weight: 600;
  color: #1d2129;
}

.content {
  min-height: 0;
  overflow-y: auto;
  padding: 24px;
}

@media (max-width: 960px) {
  .header {
    padding: 12px 16px;
    align-items: flex-start;
  }

  .header-right {
    gap: 8px;
    justify-content: flex-end;
  }

  .notice-card {
    max-width: none;
  }
}
</style>
