<template>
  <div class="chart-shell">
    <div class="chart-grid">
      <div v-for="n in 5" :key="n" class="chart-grid-line"></div>
    </div>
    <div class="chart-bars">
      <div v-for="item in normalizedData" :key="item.month" class="chart-bar-item">
        <div class="chart-bar-value">{{ item.label }}</div>
        <div class="chart-bar-track">
          <div class="chart-bar-fill" :style="{ height: `${item.height}%` }"></div>
        </div>
        <div class="chart-bar-month">{{ item.month }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{ data: Array<{ month: string; total: number }> }>();

const normalizedData = computed(() => {
  const max = Math.max(...props.data.map(item => item.total), 1);
  return props.data.map(item => ({
    ...item,
    height: Math.max(8, (item.total / max) * 100),
    label: (item.total / 100).toFixed(0)
  }));
});
</script>

<style scoped lang="scss">
.chart-shell {
  position: relative;
  height: 320px;
}

.chart-grid {
  position: absolute;
  inset: 0 0 28px 0;
  display: grid;
  grid-template-rows: repeat(5, 1fr);
}

.chart-grid-line {
  border-bottom: 1px dashed var(--color-border-2);
}

.chart-bars {
  position: absolute;
  inset: 0 0 0 0;
  display: grid;
  grid-template-columns: repeat(12, minmax(0, 1fr));
  gap: 10px;
  align-items: end;
}

.chart-bar-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  height: 100%;
}

.chart-bar-value {
  font-size: 11px;
  color: var(--color-text-3);
}

.chart-bar-track {
  display: flex;
  align-items: flex-end;
  width: 100%;
  height: 100%;
  min-height: 200px;
}

.chart-bar-fill {
  width: 100%;
  border-radius: 10px 10px 4px 4px;
  background: linear-gradient(180deg, #69b1ff 0%, #165dff 100%);
  box-shadow: 0 10px 18px rgb(22 93 255 / 16%);
}

.chart-bar-month {
  font-size: 12px;
  color: var(--color-text-2);
}
</style>
