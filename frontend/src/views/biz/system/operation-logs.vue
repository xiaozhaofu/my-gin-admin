<template>
  <div class="snow-page">
    <div class="log-page">
      <section class="page-hero">
        <div>
          <div class="page-badge">Audit Timeline</div>
          <h1>操作日志</h1>
          <p>统一查看接口调用路径、方法、耗时、请求体与执行结果，用于排查后台行为和接口故障。</p>
        </div>
        <div class="page-hero-side">
          <div class="page-hero-label">当前日志总数</div>
          <div class="page-hero-value">{{ pagination.total }}</div>
        </div>
      </section>

      <a-card class="panel-card" :bordered="false">
        <template #title>筛选条件</template>
        <a-form :model="filters" layout="inline" class="filters">
          <a-form-item label="账号"><a-input v-model="filters.username" allow-clear /></a-form-item>
          <a-form-item label="方法">
            <a-select v-model="filters.method" allow-clear style="width: 120px">
              <a-option value="POST">POST</a-option>
              <a-option value="PUT">PUT</a-option>
              <a-option value="DELETE">DELETE</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="路径"><a-input v-model="filters.path" allow-clear /></a-form-item>
          <a-space wrap>
            <a-button type="primary" @click="load">搜索</a-button>
            <a-button @click="reset">重置</a-button>
            <a-button status="danger" @click="clearAll">清空日志</a-button>
          </a-space>
        </a-form>
      </a-card>

      <a-card class="panel-card" :bordered="false">
        <template #title>日志列表</template>
        <a-table :data="list" row-key="id" :pagination="pagination" @page-change="pageChange">
          <template #columns>
            <a-table-column title="账号" data-index="username" />
            <a-table-column title="方法" data-index="method" />
            <a-table-column title="路径" data-index="path" />
            <a-table-column title="状态码" data-index="status_code" />
            <a-table-column title="结果">
              <template #cell="{ record }">
                <a-tag :color="record.success ? 'green' : 'red'">{{ record.success ? "成功" : "失败" }}</a-tag>
              </template>
            </a-table-column>
            <a-table-column title="耗时(ms)" data-index="duration_ms" />
            <a-table-column title="IP" data-index="client_ip" />
            <a-table-column title="请求体">
              <template #cell="{ record }">
                <a-typography-paragraph :ellipsis="{ rows: 2, expandable: true }">{{ record.request_body }}</a-typography-paragraph>
              </template>
            </a-table-column>
            <a-table-column title="时间" data-index="created_at" />
            <a-table-column title="操作">
              <template #cell="{ record }">
                <a-button size="mini" status="danger" @click="remove(record.id)">删除</a-button>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { operationLogClearAPI, operationLogDeleteAPI, operationLogListAPI, type OperationLogItem } from "@/api/system";
import { useConfirmAction } from "@/hooks/useConfirmAction";

const list = ref<OperationLogItem[]>([]);
const { confirmDelete } = useConfirmAction();
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const filters = reactive<Record<string, any>>({
  username: "",
  method: undefined,
  path: ""
});

const load = async () => {
  const res = await operationLogListAPI({
    ...filters,
    page: pagination.current,
    page_size: pagination.pageSize
  });
  list.value = res.data.list;
  pagination.total = res.data.total;
};

const pageChange = (page: number) => {
  pagination.current = page;
  load();
};

const reset = () => {
  filters.username = "";
  filters.method = undefined;
  filters.path = "";
  pagination.current = 1;
  load();
};

const remove = async (id: number) => {
  await confirmDelete(async () => {
    await operationLogDeleteAPI(id);
    Message.success("日志已删除");
    load();
  }, "这条日志");
};

const clearAll = async () => {
  await confirmDelete(async () => {
    await operationLogClearAPI();
    Message.success("日志已清空");
    load();
  }, "全部操作日志");
};

onMounted(load);
</script>

<style scoped lang="scss">
.log-page {
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
  background:
    linear-gradient(135deg, rgba(22, 93, 255, 0.96), rgba(64, 128, 255, 0.86)),
    radial-gradient(circle at top right, rgba(255, 255, 255, 0.16), transparent 38%);
  color: #fff;
  box-shadow: 0 20px 38px rgb(22 93 255 / 16%);
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

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 16px;
}

@media (max-width: 1080px) {
  .page-hero {
    grid-template-columns: 1fr;
  }
}
</style>
