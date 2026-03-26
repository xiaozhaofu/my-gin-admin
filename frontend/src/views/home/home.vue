<template>
  <div class="snow-page">
    <div class="home-page">
      <section class="hero-panel">
        <div class="hero-main">
          <div class="hero-badge">Content Admin Dashboard</div>
          <div class="hero-title">内容管理后台工作台</div>
          <div class="hero-sub">
            对齐参考首页结构，集中展示常用功能、第三板指标、财务指标、销售额趋势和三级内容菜单。
          </div>
        </div>
        <div class="hero-side">
          <div class="hero-stat">
            <div class="hero-stat-label">当前登录账号</div>
            <div class="hero-stat-value">{{ session.profile?.nickname || session.profile?.username || "未登录" }}</div>
          </div>
          <div class="hero-stat">
            <div class="hero-stat-label">权限角色</div>
            <div class="hero-stat-value">{{ session.profile?.roles?.join(" / ") || "暂无" }}</div>
          </div>
        </div>
      </section>

      <section class="panel-box">
        <div class="box-title">
          <div>常用功能</div>
          <div class="box-sub">按业务域分组，整体结构参考 ginfast 首页</div>
        </div>
        <a-divider :margin="16" />
        <div class="group-list">
          <div v-for="group in combinedShortcutGroups" :key="group.key" class="group-card">
            <div class="group-card-title">
              <s-svg-icon :name="group.icon" :size="22" />
              <span>{{ group.title }}</span>
            </div>
            <a-grid class="shortcut-grid" :cols="{ xs: 1, sm: 2, lg: 2, xl: 3 }" :col-gap="16" :row-gap="16">
              <a-grid-item v-for="item in group.items" :key="item.path">
                <a-card hoverable class="shortcut-card" @click="goToShortcut(item)">
                  <div class="shortcut-icon">
                    <a-image v-if="isImageIcon(item.icon)" :src="item.icon" width="28" height="28" fit="cover" />
                    <s-svg-icon v-else :name="item.icon" :size="28" />
                  </div>
                  <div class="shortcut-title">{{ item.title }}</div>
                  <div class="shortcut-tip">{{ item.tip }}</div>
                </a-card>
              </a-grid-item>
            </a-grid>
          </div>
        </div>
      </section>

      <section class="panel-box">
        <div class="box-title">
          <div>第三版指标</div>
          <div class="box-sub">首页直接展示核心运营数据</div>
        </div>
        <a-divider :margin="16" />
        <div class="scoreboard">
          <div class="scoreboard-main">
            <div class="scoreboard-main-label">综合评分</div>
            <div class="scoreboard-main-value">{{ overviewScore }}</div>
            <div class="scoreboard-main-tip">基于文章、资源、管理员、角色、在线会话等维度计算</div>
          </div>
          <div class="scoreboard-items">
            <div v-for="item in metricCards" :key="item.key" class="scoreboard-item">
              <div class="scoreboard-item-top">
                <span>{{ item.title }}</span>
                <span :class="item.trend === 'up' ? 'trend-up' : 'trend-down'">
                  {{ item.trend === 'up' ? "↗" : "↘" }}
                </span>
              </div>
              <div class="scoreboard-item-value">{{ item.value }}</div>
              <div class="scoreboard-item-tip">{{ item.tip }}</div>
            </div>
          </div>
        </div>
      </section>

      <section class="panel-box">
        <div class="box-title">
          <div>财务指标</div>
          <div class="box-sub">基于订单与订单流水聚合展示</div>
        </div>
        <a-divider :margin="16" />
        <div class="finance-board">
          <div class="finance-highlight">
            <div class="finance-highlight-badge">经营快照</div>
            <div class="finance-highlight-value">{{ financeSummary.value }}</div>
            <div class="finance-highlight-title">{{ financeSummary.title }}</div>
            <div class="finance-highlight-sub">{{ financeSummary.tip }}</div>
            <div class="finance-highlight-tags">
              <a-tag color="arcoblue">渠道数 {{ pieData.length || 0 }}</a-tag>
              <a-tag color="green">趋势月份 {{ trendData.length || 0 }}</a-tag>
            </div>
          </div>
          <div class="finance-grid">
            <div v-for="item in financePanels" :key="item.title" class="finance-panel" :style="{ '--finance-accent': item.color }">
              <div class="finance-panel-top">
                <div class="finance-head">
                  <div class="finance-dot" :style="{ borderColor: item.color }"></div>
                  <span>{{ item.title }}</span>
                </div>
                <div class="finance-trend">{{ item.trendLabel }}</div>
              </div>
              <div class="finance-value">{{ item.value }}</div>
              <div class="finance-panel-sub">{{ item.tip }}</div>
              <div class="finance-progress">
                <span :style="{ width: `${item.progress}%`, background: item.color }"></span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="panel-box">
        <div class="box-title">
          <div>销售额趋势</div>
          <div class="box-sub">柱状图展示月销售额，饼图展示支付渠道结构</div>
        </div>
        <a-divider :margin="16" />
        <DataBox :trend="trendData" :pie="pieData" />
      </section>

      <section class="panel-box">
        <div class="box-title">
          <div>三级内容菜单</div>
          <div class="box-sub">直接读取 `menus` 数据表，按父子层级展示</div>
        </div>
        <a-divider :margin="16" />
        <div class="menu-tree-panel">
          <div v-if="menuTree.length === 0" class="menu-empty">
            <a-empty description="当前没有内容菜单数据" />
          </div>
          <div v-else class="menu-root-list">
            <div v-for="root in menuTree" :key="root.id" class="menu-root-card">
              <div class="menu-root-title">
                <span class="menu-root-marker"></span>
                <a-space>
                  <a-image v-if="isImageIcon(root.icon)" :src="root.icon" width="18" height="18" fit="cover" />
                  <s-svg-icon v-else-if="root.icon" :name="root.icon" :size="18" />
                  <span>{{ root.name }}</span>
                </a-space>
              </div>
              <div class="menu-root-meta">{{ levelLabel(root.level) }}菜单 · 页面路径：{{ root.page_path || "-" }}</div>
              <div v-if="root.children?.length" class="menu-children">
                <div v-for="child in root.children" :key="child.id" class="menu-child-card">
                  <div class="menu-child-title">
                    <a-space>
                      <a-image v-if="isImageIcon(child.icon)" :src="child.icon" width="16" height="16" fit="cover" />
                      <s-svg-icon v-else-if="child.icon" :name="child.icon" :size="16" />
                      <span>{{ child.name }}</span>
                    </a-space>
                    <a-tag size="small" color="arcoblue">{{ levelLabel(child.level) }}</a-tag>
                  </div>
                  <div class="menu-child-path">页面路径：{{ child.page_path || "-" }}</div>
                  <div v-if="child.children?.length" class="menu-grandchildren">
                    <div v-for="leaf in child.children" :key="leaf.id" class="menu-leaf-card">
                      <div class="menu-leaf-title">
                        <a-space>
                          <a-image v-if="isImageIcon(leaf.icon)" :src="leaf.icon" width="14" height="14" fit="cover" />
                          <s-svg-icon v-else-if="leaf.icon" :name="leaf.icon" :size="14" />
                          <span>{{ leaf.name }}</span>
                        </a-space>
                        <a-tag size="small" color="green">{{ levelLabel(leaf.level) }}</a-tag>
                      </div>
                      <div class="menu-leaf-meta">页面路径：{{ leaf.page_path || "-" }}</div>
                    </div>
                  </div>
                  <div v-else class="menu-child-empty">当前二级菜单下暂无三级菜单</div>
                </div>
              </div>
              <div v-else class="menu-child-empty">当前一级菜单下暂无子菜单</div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { Message } from "@arco-design/web-vue";
import { dashboardAPI, type DashboardFinanceItem, type DashboardMetric, type DashboardPieItem, type DashboardShortcutGroup, type DashboardShortcutItem, type DashboardTrendItem } from "@/api/dashboard";
import { menuTreeAPI, type ContentMenu } from "@/api/menu";
import DataBox from "@/views/home/components/data-box.vue";
import { useSessionStore } from "@/store/modules/session";

const router = useRouter();
const session = useSessionStore();
const menuTree = ref<ContentMenu[]>([]);
const metricCards = ref<DashboardMetric[]>([]);
const shortcutGroups = ref<DashboardShortcutGroup[]>([]);
const financeCards = ref<DashboardFinanceItem[]>([]);
const trendData = ref<DashboardTrendItem[]>([]);
const pieData = ref<DashboardPieItem[]>([]);

type HomeShortcut = DashboardShortcutItem & {
  menuID?: number;
  pagePath?: string;
};

type HomeShortcutGroup = {
  key: string;
  title: string;
  icon: string;
  items: HomeShortcut[];
};

type FinancePanel = DashboardFinanceItem & {
  progress: number;
  tip: string;
  trendLabel: string;
};

const visibleShortcutGroups = computed<HomeShortcutGroup[]>(() =>
  shortcutGroups.value
    .map(group => ({
      ...group,
      items: group.items.filter(item => !item.permission || session.can(item.permission))
    }))
    .filter(group => group.items.length > 0)
);

const combinedShortcutGroups = computed<HomeShortcutGroup[]>(() => {
  const groups = [...visibleShortcutGroups.value];
  const contentEntries = menuTree.value.slice(0, 6).map(item => ({
    path: "/articles",
    title: item.name,
    permission: "/api/v1/articles#GET",
    icon: item.icon || "folder-menu",
    tip: item.page_path ? `客户端页面：${item.page_path}` : "跳转到文章页并按菜单筛选",
    menuID: item.id,
    pagePath: item.page_path
  }));
  if (contentEntries.length > 0) {
    groups.unshift({
      key: "menu-entry",
      title: "内容入口",
      icon: "folder-menu",
      items: contentEntries
    });
  }
  return groups;
});

const overviewScore = computed(() => {
  if (metricCards.value.length === 0) return "0.0";
  const values = metricCards.value.map(item => Number(item.value) || 0);
  const normalized = values.reduce((sum, value) => sum + Math.min(value, 100), 0) / values.length;
  return (normalized / 10).toFixed(1);
});

const financePanels = computed<FinancePanel[]>(() => {
  const values = financeCards.value.map(item => Number(item.value) || 0);
  const maxValue = Math.max(...values, 1);
  return financeCards.value.map((item, index) => {
    const ratio = Math.round(((Number(item.value) || 0) / maxValue) * 100);
    const financeTips = ["支付侧核心金额汇总", "退款链路金额回收", "成功支付订单量", "退款订单量"];
    return {
      ...item,
      progress: Math.max(ratio, 12),
      tip: financeTips[index] || "经营指标概览",
      trendLabel: index < 2 ? "经营稳定" : "关注回收"
    };
  });
});

const financeSummary = computed(() => {
  const primary = financePanels.value[0];
  if (!primary) {
    return {
      title: "暂无财务数据",
      value: "--",
      tip: "订单与流水产生后会自动汇总到这里"
    };
  }
  return {
    title: primary.title,
    value: primary.value,
    tip: `当前财务概览共包含 ${financePanels.value.length} 项核心指标`
  };
});

const goTo = (path: string) => {
  router.push(path);
};

const goToShortcut = (item: HomeShortcut) => {
  if (item.menuID) {
    router.push({ path: "/articles", query: { menu_id: String(item.menuID) } });
    return;
  }
  goTo(item.path);
};

const isImageIcon = (icon?: string) => !!icon && /(\.png|\.jpg|\.jpeg|\.gif|\.webp|^https?:\/\/|^\/uploads\/)/i.test(icon);
const levelLabel = (level?: number) => {
  switch (level) {
    case 1:
      return "一级";
    case 2:
      return "二级";
    case 3:
      return "三级";
    default:
      return `第${level || 0}级`;
  }
};

const countMenus = (items: ContentMenu[]): number =>
  items.reduce((total, item) => total + 1 + countMenus(item.children || []), 0);

const loadDashboard = async () => {
  try {
    const [dashboardRes, menuRes] = await Promise.all([dashboardAPI(), menuTreeAPI()]);
    shortcutGroups.value = dashboardRes.data.shortcut_groups;
    metricCards.value = dashboardRes.data.metrics;
    financeCards.value = dashboardRes.data.finance;
    trendData.value = dashboardRes.data.order_trend;
    pieData.value = dashboardRes.data.order_pie;
    menuTree.value = menuRes.data;

    if (!metricCards.value.find(item => item.key === "menus")) {
      metricCards.value = [
        ...metricCards.value,
        { key: "menus", title: "内容菜单", value: String(countMenus(menuRes.data)), tip: "来自 menus 表", color: "#00b42a", trend: "up" }
      ];
    }
  } catch (error) {
    console.error("load dashboard failed", error);
    Message.error("首页数据加载失败");
  }
};

onMounted(loadDashboard);
</script>

<style scoped lang="scss">
.home-page {
  padding: $padding;
  background: $color-bg-1;
}

.hero-panel {
  display: grid;
  grid-template-columns: minmax(0, 1.8fr) minmax(280px, 0.8fr);
  gap: 20px;
  padding: 24px 28px;
  margin-bottom: calc($padding * 2);
  border-radius: 16px;
  background:
    linear-gradient(135deg, rgba(22, 93, 255, 0.96), rgba(64, 128, 255, 0.86)),
    radial-gradient(circle at top right, rgba(255, 255, 255, 0.16), transparent 36%);
  box-shadow: 0 20px 40px rgb(22 93 255 / 18%);
  color: #fff;
}

.hero-main {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.hero-badge {
  display: inline-flex;
  width: fit-content;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.14);
  font-size: $font-size-body-1;
}

.hero-title {
  font-family: AliFangYuanTi, "PingFang SC", sans-serif;
  font-size: 34px;
  font-weight: 700;
  line-height: 1.15;
}

.hero-sub {
  max-width: 720px;
  font-size: $font-size-body-3;
  line-height: 1.8;
  opacity: 0.9;
}

.hero-side {
  display: grid;
  gap: 16px;
}

.hero-stat {
  padding: 16px 18px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(10px);
}

.hero-stat-label {
  font-size: $font-size-body-1;
  opacity: 0.78;
}

.hero-stat-value {
  margin-top: 8px;
  font-size: $font-size-title-2;
  font-weight: 700;
  word-break: break-all;
}

.panel-box {
  padding: 20px 24px 24px;
  margin-bottom: calc($padding * 2);
  border-radius: 16px;
  background: $color-bg-white;
  box-shadow: 0 12px 30px rgb(15 23 42 / 6%);
}

.box-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-size: $font-size-body-3;
  color: $color-text-1;
}

.box-sub {
  font-size: $font-size-body-1;
  color: $color-text-3;
}

.group-list {
  display: grid;
  gap: 20px;
}

.group-card-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
  font-size: $font-size-title-1;
  font-weight: 700;
  color: $color-text-1;
}

.shortcut-card {
  height: 100%;
  min-height: 132px;
  border-radius: 14px;
  border: 1px solid transparent;
  background: linear-gradient(180deg, #ffffff 0%, #f8fbff 100%);
  transition: all 0.2s ease;
  cursor: pointer;
}

.shortcut-card:hover {
  border-color: rgb(var(--primary-3));
  box-shadow: 0 12px 24px rgb(22 93 255 / 10%);
}

.shortcut-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  margin-bottom: 14px;
  border-radius: 14px;
  background: linear-gradient(135deg, rgba(22, 93, 255, 0.12), rgba(64, 128, 255, 0.18));
  color: rgb(var(--primary-6));
}

.shortcut-title {
  font-size: $font-size-body-3;
  font-weight: 700;
  color: $color-text-1;
}

.shortcut-tip {
  margin-top: 8px;
  font-size: $font-size-body-1;
  line-height: 1.7;
  color: $color-text-3;
}

.metric-card,
.finance-card {
  border-radius: 14px;
  background: linear-gradient(180deg, #fff 0%, #f8fbff 100%);
}

.metric-head,
.finance-head {
  display: flex;
  align-items: center;
  gap: 8px;
  color: $color-text-2;
}

.metric-dot,
.finance-dot {
  box-sizing: border-box;
  width: 10px;
  height: 10px;
  border: 3px solid;
  border-radius: 50%;
}

.metric-value,
.finance-value {
  margin-top: 16px;
  font-family: AliFangYuanTi, "PingFang SC", sans-serif;
  font-size: 30px;
  font-weight: 700;
  color: $color-text-1;
}

.metric-tip {
  font-size: $font-size-body-1;
  color: $color-text-3;
}

.scoreboard {
  display: grid;
  grid-template-columns: minmax(240px, 0.9fr) minmax(0, 2.1fr);
  gap: 18px;
}

.scoreboard-main {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 212px;
  padding: 24px;
  border-radius: 18px;
  background:
    linear-gradient(145deg, rgba(255, 125, 0, 0.95), rgba(255, 154, 46, 0.92)),
    radial-gradient(circle at top right, rgba(255, 255, 255, 0.24), transparent 34%);
  box-shadow: 0 18px 34px rgb(255 125 0 / 16%);
  color: #fff;
}

.scoreboard-main-label {
  font-size: 13px;
  opacity: 0.9;
}

.scoreboard-main-value {
  margin-top: 16px;
  font-family: AliFangYuanTi, "PingFang SC", sans-serif;
  font-size: 42px;
  font-weight: 700;
  line-height: 1.1;
}

.scoreboard-main-tip {
  margin-top: 10px;
  line-height: 1.7;
  opacity: 0.92;
}

.scoreboard-items {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.scoreboard-item {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 98px;
  padding: 18px;
  border: 1px solid rgba(22, 93, 255, 0.08);
  border-radius: 16px;
  background: linear-gradient(180deg, #fff 0%, #f8fbff 100%);
  box-shadow: 0 12px 24px rgb(15 23 42 / 5%);
}

.scoreboard-item-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-size: 13px;
  font-weight: 600;
  color: $color-text-2;
}

.scoreboard-item-value {
  font-family: AliFangYuanTi, "PingFang SC", sans-serif;
  font-size: 28px;
  font-weight: 700;
  color: $color-text-1;
}

.scoreboard-item-tip {
  font-size: 12px;
  line-height: 1.6;
  color: $color-text-3;
}

.trend-up,
.trend-down {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 30px;
  height: 26px;
  padding: 0 8px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 700;
}

.trend-up {
  background: rgba(0, 180, 42, 0.12);
  color: #00b42a;
}

.trend-down {
  background: rgba(245, 63, 63, 0.12);
  color: #f53f3f;
}

.finance-board {
  display: grid;
  grid-template-columns: minmax(260px, 0.9fr) minmax(0, 2.1fr);
  gap: 18px;
}

.finance-highlight {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-height: 224px;
  padding: 22px;
  border-radius: 18px;
  background:
    linear-gradient(145deg, rgba(22, 93, 255, 0.96), rgba(99, 152, 255, 0.92)),
    radial-gradient(circle at top right, rgba(255, 255, 255, 0.22), transparent 36%);
  box-shadow: 0 16px 34px rgb(22 93 255 / 18%);
  color: #fff;
}

.finance-highlight-badge {
  display: inline-flex;
  width: fit-content;
  padding: 5px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.18);
  font-size: 12px;
}

.finance-highlight-value {
  margin-top: 18px;
  font-family: AliFangYuanTi, "PingFang SC", sans-serif;
  font-size: 34px;
  font-weight: 700;
  line-height: 1.15;
}

.finance-highlight-title {
  margin-top: 10px;
  font-size: 18px;
  font-weight: 700;
}

.finance-highlight-sub {
  margin-top: 8px;
  line-height: 1.7;
  opacity: 0.9;
}

.finance-highlight-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 20px;
}

.finance-highlight-tags :deep(.arco-tag) {
  background: rgba(255, 255, 255, 0.16);
  color: #fff;
  border-color: transparent;
}

.finance-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.finance-panel {
  padding: 18px;
  border-radius: 16px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 1) 0%, rgba(246, 249, 255, 0.95) 100%),
    linear-gradient(135deg, color-mix(in srgb, var(--finance-accent) 12%, white) 0%, white 100%);
  border: 1px solid rgba(22, 93, 255, 0.08);
  box-shadow: 0 12px 24px rgb(15 23 42 / 6%);
}

.finance-panel-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.finance-trend {
  padding: 4px 10px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--finance-accent) 12%, white);
  font-size: 12px;
  font-weight: 600;
  color: var(--finance-accent);
}

.finance-panel-sub {
  margin-top: 8px;
  font-size: 12px;
  color: $color-text-3;
}

.finance-progress {
  height: 7px;
  margin-top: 16px;
  overflow: hidden;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.06);
}

.finance-progress span {
  display: block;
  height: 100%;
  border-radius: 999px;
}

.menu-tree-panel {
  min-height: 160px;
}

.menu-root-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
}

.menu-root-card {
  padding: 18px;
  border: 1px solid $color-border-1;
  border-radius: 14px;
  background: linear-gradient(180deg, #fff 0%, #fbfdff 100%);
}

.menu-root-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: $font-size-title-1;
  font-weight: 700;
  color: $color-text-1;
}

.menu-root-marker {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: rgb(var(--primary-6));
}

.menu-root-meta {
  margin-top: 8px;
  font-size: $font-size-body-1;
  color: $color-text-3;
}

.menu-children {
  display: grid;
  gap: 12px;
  margin-top: 16px;
}

.menu-child-card {
  padding: 12px;
  border-radius: 12px;
  background: $color-fill-1;
}

.menu-child-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-size: $font-size-body-3;
  font-weight: 600;
}

.menu-grandchildren {
  display: grid;
  gap: 8px;
  margin-top: 12px;
}

.menu-leaf-card {
  padding: 10px 12px;
  border: 1px solid rgba(0, 180, 42, 0.12);
  border-radius: 10px;
  background: rgba(0, 180, 42, 0.04);
}

.menu-leaf-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-size: 13px;
  font-weight: 600;
  color: $color-text-1;
}

.menu-leaf-meta {
  margin-top: 6px;
  font-size: 12px;
  color: $color-text-3;
}

.menu-child-path {
  margin-top: 8px;
  font-size: 12px;
  color: $color-text-3;
}

.menu-child-empty {
  margin-top: 12px;
  font-size: $font-size-body-1;
  color: $color-text-3;
}

.menu-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 180px;
}

@media (max-width: 1080px) {
  .hero-panel {
    grid-template-columns: 1fr;
  }

  .scoreboard,
  .finance-board {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .scoreboard-items,
  .finance-grid {
    grid-template-columns: 1fr;
  }
}
</style>
