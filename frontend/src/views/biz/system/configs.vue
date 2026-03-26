<template>
  <div class="snow-page">
    <div class="page-shell">
      <section class="page-hero amber-hero">
        <div>
          <div class="page-badge">System Config</div>
          <h1>系统参数</h1>
          <p>统一维护站点级配置项，支持 Key / Value / 类型 / 备注管理，便于后台参数集中治理。</p>
        </div>
        <div class="page-hero-side">
          <div class="page-hero-label">当前配置项</div>
          <div class="page-hero-value">{{ list.length }}</div>
        </div>
      </section>

      <a-card class="panel-card" :bordered="false">
        <template #title>参数列表</template>
        <template #extra>
          <a-button v-if="session.can('/api/v1/sys-configs#POST')" type="primary" @click="open()">新增参数</a-button>
        </template>
        <a-table :data="list" :pagination="false" row-key="id">
          <template #columns>
            <a-table-column title="名称" data-index="config_name" />
            <a-table-column title="Key" data-index="config_key" />
            <a-table-column title="Value">
              <template #cell="{ record }">
                <a-typography-paragraph :ellipsis="{ rows: 2, expandable: true }">{{ record.config_value }}</a-typography-paragraph>
              </template>
            </a-table-column>
            <a-table-column title="类型" data-index="config_type" />
            <a-table-column title="备注" data-index="remark" />
            <a-table-column title="操作">
              <template #cell="{ record }">
                <a-space>
                  <a-button v-if="session.can('/api/v1/sys-configs/:id#PUT')" size="mini" @click="open(record)">编辑</a-button>
                  <a-button v-if="session.can('/api/v1/sys-configs/:id#DELETE')" size="mini" status="danger" @click="remove(record.id)">删除</a-button>
                </a-space>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>

  <a-modal v-model:visible="visible" :title="form.id ? '编辑参数' : '新增参数'" :on-before-ok="submit">
    <a-form :model="form" layout="vertical">
      <a-form-item label="名称"><a-input v-model="form.config_name" /></a-form-item>
      <a-form-item label="Key"><a-input v-model="form.config_key" /></a-form-item>
      <a-form-item label="Value"><a-textarea v-model="form.config_value" /></a-form-item>
      <a-form-item label="类型"><a-input-number v-model="form.config_type" :min="0" :max="1" /></a-form-item>
      <a-form-item label="备注"><a-input v-model="form.remark" /></a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { sysConfigDeleteAPI, sysConfigListAPI, sysConfigSaveAPI, type SysConfigItem } from "@/api/system";
import { useConfirmAction } from "@/hooks/useConfirmAction";
import { useSessionStore } from "@/store/modules/session";

const list = ref<SysConfigItem[]>([]);
const visible = ref(false);
const session = useSessionStore();
const { confirmDelete, confirmSave } = useConfirmAction();
const form = reactive<any>({ id: null, config_name: "", config_key: "", config_value: "", config_type: 0, remark: "" });

const load = async () => {
  const res = await sysConfigListAPI();
  list.value = res.data;
};

const open = (record?: any) => {
  Object.assign(form, { id: null, config_name: "", config_key: "", config_value: "", config_type: 0, remark: "" }, record || {});
  visible.value = true;
};

const submit = async () => {
  return confirmSave(async () => {
    await sysConfigSaveAPI(form.id, form);
    Message.success("系统参数保存成功");
    visible.value = false;
    load();
  }, form.id ? "当前系统参数" : "新系统参数");
};

const remove = async (id: number) => {
  await confirmDelete(async () => {
    await sysConfigDeleteAPI(id);
    Message.success("系统参数已删除");
    load();
  }, "这个系统参数");
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

.amber-hero {
  background:
    linear-gradient(135deg, rgba(255, 125, 0, 0.94), rgba(255, 159, 28, 0.88)),
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
