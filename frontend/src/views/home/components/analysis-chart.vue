<template>
  <div class="pie-shell">
    <svg viewBox="0 0 220 220" class="pie-svg">
      <circle cx="110" cy="110" r="70" fill="none" stroke="#edf2fb" stroke-width="26" />
      <circle
        v-for="segment in segments"
        :key="segment.type"
        cx="110"
        cy="110"
        r="70"
        fill="none"
        :stroke="segment.color"
        stroke-width="26"
        stroke-linecap="round"
        :stroke-dasharray="segment.dashArray"
        :stroke-dashoffset="segment.dashOffset"
        transform="rotate(-90 110 110)"
      />
    </svg>
    <div class="pie-center">
      <div class="pie-center-label">总销售额</div>
      <div class="pie-center-value">{{ totalLabel }}</div>
    </div>
    <div class="pie-legend">
      <div v-for="segment in segments" :key="segment.type" class="pie-legend-item">
        <span class="pie-legend-dot" :style="{ background: segment.color }"></span>
        <span>{{ segment.type }}</span>
        <span class="pie-legend-value">{{ (segment.value / 100).toFixed(0) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{ data: Array<{ type: string; value: number }> }>();

const colors = ["#165dff", "#36cfc9", "#ff7d00", "#722ed1", "#eb0aa4", "#00b42a", "#86909c", "#ff4d4f"];
const circumference = 2 * Math.PI * 70;

const total = computed(() => props.data.reduce((sum, item) => sum + item.value, 0));
const totalLabel = computed(() => (total.value / 100).toFixed(0));

const segments = computed(() => {
  let offset = 0;
  return props.data.map((item, index) => {
    const ratio = total.value > 0 ? item.value / total.value : 0;
    const length = Math.max(6, ratio * circumference);
    const segment = {
      ...item,
      color: colors[index % colors.length],
      dashArray: `${length} ${circumference - length}`,
      dashOffset: -offset
    };
    offset += length;
    return segment;
  });
});
</script>

<style scoped lang="scss">
.pie-shell {
  position: relative;
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr);
  align-items: center;
  gap: 20px;
}

.pie-svg {
  width: 220px;
  height: 220px;
}

.pie-center {
  position: absolute;
  top: 50%;
  left: 110px;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.pie-center-label {
  font-size: 12px;
  color: var(--color-text-3);
}

.pie-center-value {
  margin-top: 4px;
  font-family: AliFangYuanTi, "PingFang SC", sans-serif;
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-1);
}

.pie-legend {
  display: grid;
  gap: 10px;
}

.pie-legend-item {
  display: grid;
  grid-template-columns: 10px minmax(0, 1fr) auto;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: var(--color-text-2);
}

.pie-legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.pie-legend-value {
  color: var(--color-text-1);
  font-weight: 700;
}

@media (max-width: 768px) {
  .pie-shell {
    grid-template-columns: 1fr;
    justify-items: center;
  }

  .pie-legend {
    width: 100%;
  }
}
</style>
