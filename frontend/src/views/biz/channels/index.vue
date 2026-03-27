<template>
  <div class="snow-page">
    <div class="channel-page">
      <section class="page-hero">
        <div>
          <div class="page-badge">Channel Center</div>
          <h1>渠道管理</h1>
          <p>统一维护内容渠道名称、渠道编码、状态与备注，文章发布和筛选会直接复用这里的数据。</p>
        </div>
        <div class="page-hero-side">
          <div class="page-hero-label">当前渠道总数</div>
          <div class="page-hero-value">{{ list.length }}</div>
        </div>
      </section>

      <a-card class="panel-card" :bordered="false">
        <template #title>渠道列表</template>
        <template #extra>
          <a-space wrap>
            <a-button v-if="canStatus" :disabled="!selectedRowKeys.length" @click="batchStatus(1)">批量启用</a-button>
            <a-button v-if="canStatus" :disabled="!selectedRowKeys.length" @click="batchStatus(2)">批量禁用</a-button>
            <a-button v-if="canCreate" type="primary" @click="open()">新增渠道</a-button>
          </a-space>
        </template>

        <a-table
          row-key="id"
          :data="list"
          :pagination="false"
          :row-selection="canStatus ? { type: 'checkbox' } : undefined"
          :selected-keys="selectedRowKeys"
          @selection-change="setSelected"
        >
          <template #columns>
            <a-table-column title="渠道名称" data-index="name" />
            <a-table-column title="渠道编码" data-index="code" />
            <a-table-column title="状态">
              <template #cell="{ record }">
                <a-tag :color="statusColor(record.status)">{{ statusLabel(record.status) }}</a-tag>
              </template>
            </a-table-column>
            <a-table-column title="备注">
              <template #cell="{ record }">
                <a-typography-paragraph :ellipsis="{ rows: 2, expandable: true }">{{ record.remark || "-" }}</a-typography-paragraph>
              </template>
            </a-table-column>
            <a-table-column title="操作">
              <template #cell="{ record }">
                <a-button v-if="canEdit" size="mini" @click="open(record)">编辑</a-button>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>

  <a-modal v-model:visible="visible" :title="current.id ? '编辑渠道' : '新增渠道'" :on-before-ok="submit">
    <a-form :model="current" layout="vertical">
      <a-form-item label="渠道名称"><a-input v-model="current.name" /></a-form-item>
      <a-form-item label="渠道编码"><a-input v-model="current.code" /></a-form-item>
      <a-form-item label="状态">
        <a-select v-model="current.status">
          <a-option :value="1">启用</a-option>
          <a-option :value="2">禁用</a-option>
        </a-select>
      </a-form-item>
      <a-form-item label="备注"><a-textarea v-model="current.remark" /></a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { channelListAPI, channelSaveAPI, channelStatusAPI, type ChannelItem } from "@/api/channel";
import { useConfirmAction } from "@/hooks/useConfirmAction";
import { usePageSelection } from "@/hooks/usePageSelection";
import { useSessionStore } from "@/store/modules/session";

const session = useSessionStore();
const { confirmBatchHide, confirmSave, runConfirmed } = useConfirmAction();
const list = ref<ChannelItem[]>([]);
const visible = ref(false);
const current = reactive<any>({ id: null, name: "", code: "", status: 1, remark: "" });
const {
  selectedKeys: selectedRowKeys,
  setSelected,
  clearSelected
} = usePageSelection(() => list.value, item => item.id);

const canCreate = computed(() => session.can("/api/v1/channels#POST"));
const canEdit = computed(() => session.can("/api/v1/channels/:id#PUT"));
const canStatus = computed(() => session.can("/api/v1/channels/status#PUT"));

const load = async () => {
  const res = await channelListAPI();
  list.value = res.data;
};

const open = (record?: ChannelItem) => {
  Object.assign(current, { id: null, name: "", code: "", status: 1, remark: "" }, record || {});
  visible.value = true;
};

const submit = async () =>
  confirmSave(async () => {
    await channelSaveAPI(current.id, {
      name: current.name,
      code: current.code,
      status: current.status,
      remark: current.remark
    });
    Message.success("保存成功");
    visible.value = false;
    load();
  }, current.id ? "当前渠道" : "新渠道");

const batchStatus = async (status: number) => {
  const action = async () => {
    await channelStatusAPI({ ids: selectedRowKeys.value, status });
    Message.success("状态已更新");
    clearSelected();
    load();
  };

  if (status === 2) {
    await confirmBatchHide(action, "选中渠道");
    return;
  }

  await runConfirmed(
    {
      title: "确认批量修改状态",
      content: "确认将选中渠道批量设置为启用吗？"
    },
    action
  );
};

const statusLabel = (status: number) => (status === 1 ? "启用" : "禁用");
const statusColor = (status: number) => (status === 1 ? "green" : "gray");

onMounted(load);
</script>

<style scoped lang="scss">
.channel-page {
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

@media (max-width: 1080px) {
  .page-hero {
    grid-template-columns: 1fr;
  }
}
</style>
