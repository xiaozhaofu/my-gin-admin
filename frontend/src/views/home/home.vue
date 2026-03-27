<template>
  <div class="snow-page">
    <div class="home-page">
      <section class="hero-panel">
        <div class="hero-main">
          <div class="hero-badge">Content Admin Dashboard</div>
          <div class="hero-title">内容管理后台工作台</div>
          <div class="hero-sub">
            首页聚焦核心经营数据，不再展示内容菜单和通用快捷区，直接保留订单、财务和趋势看板。
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

      <section v-if="orderOpsEntries.length" class="panel-box">
        <div class="box-title">
          <div>订单运营</div>
          <div class="box-sub">围绕订单查询、支付核单和退款巡检的快捷入口</div>
        </div>
        <a-divider :margin="16" />
        <div class="order-ops-grid">
          <a-card v-for="entry in orderOpsEntries" :key="entry.key" hoverable class="order-ops-card" @click="goToOrderShortcut(entry)">
            <div class="order-ops-head">
              <div class="order-ops-icon">
                <s-svg-icon :name="entry.icon" :size="24" />
              </div>
              <a-tag :color="entry.tagColor">{{ entry.tag }}</a-tag>
            </div>
            <div class="order-ops-title">{{ entry.title }}</div>
            <div class="order-ops-tip">{{ entry.tip }}</div>
          </a-card>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { Message } from "@arco-design/web-vue";
import { dashboardAPI, type DashboardFinanceItem, type DashboardMetric, type DashboardPieItem, type DashboardTrendItem } from "@/api/dashboard";
import DataBox from "@/views/home/components/data-box.vue";
import { useSessionStore } from "@/store/modules/session";

const router = useRouter();
const session = useSessionStore();
const metricCards = ref<DashboardMetric[]>([]);
const financeCards = ref<DashboardFinanceItem[]>([]);
const trendData = ref<DashboardTrendItem[]>([]);
const pieData = ref<DashboardPieItem[]>([]);

type FinancePanel = DashboardFinanceItem & {
  progress: number;
  tip: string;
  trendLabel: string;
};

type OrderOpsEntry = {
  key: string;
  title: string;
  tip: string;
  icon: string;
  tag: string;
  tagColor: "arcoblue" | "green" | "orangered";
  query?: Record<string, string>;
};

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

const orderOpsEntries = computed<OrderOpsEntry[]>(() => {
  if (!session.can("/api/v1/orders#GET")) return [];
  return [
    {
      key: "all-orders",
      title: "全部订单",
      tip: "查看全部支付订单、详情和导出数据",
      icon: "list",
      tag: "总览",
      tagColor: "arcoblue"
    },
    {
      key: "paid-orders",
      title: "已支付订单",
      tip: "直接筛到已支付订单，适合运营核单和交易排查",
      icon: "check-circle",
      tag: "支付",
      tagColor: "green",
      query: { status: "10" }
    },
    {
      key: "refund-orders",
      title: "退款中订单",
      tip: "直接筛到退款中订单，便于退款链路巡检和处理",
      icon: "warning-circle",
      tag: "退款",
      tagColor: "orangered",
      query: { status: "60" }
    }
  ];
});

const goToOrderShortcut = (entry: OrderOpsEntry) => {
  router.push({
    path: "/orders",
    query: entry.query || {}
  });
};

const loadDashboard = async () => {
  try {
    const dashboardRes = await dashboardAPI();
    metricCards.value = dashboardRes.data.metrics;
    financeCards.value = dashboardRes.data.finance;
    trendData.value = dashboardRes.data.order_trend;
    pieData.value = dashboardRes.data.order_pie;
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

.finance-head {
  display: flex;
  align-items: center;
  gap: 8px;
  color: $color-text-2;
}

.finance-dot {
  box-sizing: border-box;
  width: 10px;
  height: 10px;
  border: 3px solid;
  border-radius: 50%;
}

.finance-trend {
  padding: 4px 10px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--finance-accent) 12%, white);
  font-size: 12px;
  font-weight: 600;
  color: var(--finance-accent);
}

.finance-value {
  margin-top: 16px;
  font-family: AliFangYuanTi, "PingFang SC", sans-serif;
  font-size: 30px;
  font-weight: 700;
  color: $color-text-1;
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

.order-ops-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.order-ops-card {
  min-height: 158px;
  border-radius: 16px;
  border: 1px solid rgba(22, 93, 255, 0.08);
  background: linear-gradient(180deg, #fff 0%, #f8fbff 100%);
  box-shadow: 0 12px 24px rgb(15 23 42 / 5%);
  cursor: pointer;
  transition: all 0.2s ease;
}

.order-ops-card:hover {
  border-color: rgb(var(--primary-3));
  box-shadow: 0 14px 30px rgb(22 93 255 / 10%);
}

.order-ops-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.order-ops-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 14px;
  background: linear-gradient(135deg, rgba(22, 93, 255, 0.12), rgba(64, 128, 255, 0.18));
  color: rgb(var(--primary-6));
}

.order-ops-title {
  margin-top: 16px;
  font-size: 18px;
  font-weight: 700;
  color: $color-text-1;
}

.order-ops-tip {
  margin-top: 10px;
  font-size: 12px;
  line-height: 1.8;
  color: $color-text-3;
}

@media (max-width: 1080px) {
  .hero-panel,
  .scoreboard,
  .finance-board,
  .order-ops-grid {
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
