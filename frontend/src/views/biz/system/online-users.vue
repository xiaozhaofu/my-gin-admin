<template>
  <div class="snow-page">
    <div class="page-shell">
      <section class="page-hero cyan-hero">
        <div>
          <div class="page-badge">Online Session</div>
          <h1>在线用户</h1>
          <p>实时查看当前在线会话、最后活跃时间与过期时间，并支持后台强制下线。</p>
        </div>
        <div class="page-hero-side">
          <div class="page-hero-label">当前在线数</div>
          <div class="page-hero-value">{{ list.length }}</div>
        </div>
      </section>

      <a-card class="panel-card" :bordered="false">
        <template #title>在线会话列表</template>
        <a-table :data="list" :pagination="false" row-key="id">
          <template #columns>
            <a-table-column title="账号" data-index="username" />
            <a-table-column title="IP" data-index="ip" />
            <a-table-column title="最后活跃" data-index="last_active_at" />
            <a-table-column title="过期时间" data-index="expired_at" />
            <a-table-column title="User-Agent">
              <template #cell="{ record }">
                <a-typography-paragraph :ellipsis="{ rows: 2, expandable: true }">{{ record.user_agent }}</a-typography-paragraph>
              </template>
            </a-table-column>
            <a-table-column title="操作">
              <template #cell="{ record }">
                <a-button size="mini" status="danger" @click="forceOffline(record.id)">强制下线</a-button>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { onlineSessionForceOfflineAPI, onlineSessionListAPI } from "@/api/system";

const list = ref<any[]>([]);

const load = async () => {
  const res = await onlineSessionListAPI();
  list.value = res.data;
};

const forceOffline = async (id: number) => {
  await onlineSessionForceOfflineAPI(id);
  Message.success("用户已强制下线");
  load();
};

onMounted(load);
</script>

<style scoped lang="scss">
.page-shell {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-hero {
  display: grid;
  grid-template-columns: minmax(0, 1.65fr) minmax(220px, 0.55fr);
  gap: 18px;
  padding: 24px 28px;
  border-radius: 16px;
  color: #fff;
  box-shadow: 0 20px 38px rgb(15 23 42 / 14%);
}

.cyan-hero {
  background:
    linear-gradient(135deg, rgba(15, 198, 194, 0.92), rgba(45, 183, 245, 0.88)),
    radial-gradient(circle at top right, rgba(255, 255, 255, 0.14), transparent 38%);
}

.page-badge {
  display: inline-flex;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.14);
  font-size: 12px;
}

.page-hero h1 {
  margin: 14px 0 10px;
  font-family: AliFangYuanTi, "PingFang SC", sans-serif;
  font-size: 32px;
  line-height: 1.2;
}

.page-hero p {
  max-width: 720px;
  line-height: 1.8;
  color: rgba(255, 255, 255, 0.86);
}

.page-hero-side {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 18px 20px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(10px);
}

.page-hero-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.74);
}

.page-hero-value {
  margin-top: 8px;
  font-size: 28px;
  font-weight: 700;
}

.panel-card {
  border-radius: 16px;
  box-shadow: 0 12px 30px rgb(15 23 42 / 5%);
}

@media (max-width: 1080px) {
  .page-hero {
    grid-template-columns: 1fr;
  }
}
</style>
