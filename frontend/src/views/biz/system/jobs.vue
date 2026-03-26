<template>
  <a-card :bordered="false">
    <template #title>定时任务</template>
    <a-space direction="vertical" fill>
      <a-alert>本轮先实现任务配置管理，不包含后台执行器。</a-alert>
      <a-button v-if="session.can('/api/v1/jobs#POST')" type="primary" @click="open()">新增任务</a-button>
      <a-table :data="list" :pagination="false" row-key="id">
        <template #columns>
          <a-table-column title="名称" data-index="name" />
          <a-table-column title="Key" data-index="job_key" />
          <a-table-column title="Cron" data-index="cron_expr" />
          <a-table-column title="Target" data-index="target" />
          <a-table-column title="状态" data-index="status" />
          <a-table-column title="并发" data-index="concurrent" />
          <a-table-column title="备注" data-index="remark" />
          <a-table-column title="操作">
            <template #cell="{ record }">
              <a-space>
                <a-button v-if="session.can('/api/v1/jobs/:id#PUT')" size="mini" @click="open(record)">编辑</a-button>
                <a-button v-if="session.can('/api/v1/jobs/:id#DELETE')" size="mini" status="danger" @click="remove(record.id)">删除</a-button>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-space>
  </a-card>
  <a-modal v-model:visible="visible" :title="form.id ? '编辑任务' : '新增任务'" :on-before-ok="submit">
    <a-form :model="form" layout="vertical">
      <a-form-item label="名称"><a-input v-model="form.name" /></a-form-item>
      <a-form-item label="Key"><a-input v-model="form.job_key" /></a-form-item>
      <a-form-item label="Cron"><a-input v-model="form.cron_expr" /></a-form-item>
      <a-form-item label="Target"><a-input v-model="form.target" /></a-form-item>
      <a-form-item label="状态"><a-input-number v-model="form.status" :min="1" :max="2" /></a-form-item>
      <a-form-item label="允许并发"><a-switch v-model="form.concurrent" /></a-form-item>
      <a-form-item label="备注"><a-input v-model="form.remark" /></a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { jobDeleteAPI, jobListAPI, jobSaveAPI } from "@/api/system";
import { useConfirmAction } from "@/hooks/useConfirmAction";
import { useSessionStore } from "@/store/modules/session";

const list = ref<any[]>([]);
const visible = ref(false);
const session = useSessionStore();
const { confirmDelete, confirmSave } = useConfirmAction();
const form = reactive<any>({ id: null, name: "", job_key: "", cron_expr: "", target: "", status: 1, concurrent: false, remark: "" });

const load = async () => {
  const res = await jobListAPI();
  list.value = res.data;
};

const open = (record?: any) => {
  Object.assign(form, { id: null, name: "", job_key: "", cron_expr: "", target: "", status: 1, concurrent: false, remark: "" }, record || {});
  visible.value = true;
};

const submit = async () => {
  return confirmSave(async () => {
    await jobSaveAPI(form.id, form);
    Message.success("任务已保存");
    visible.value = false;
    load();
  }, form.id ? "当前任务" : "新任务");
};

const remove = async (id: number) => {
  await confirmDelete(async () => {
    await jobDeleteAPI(id);
    Message.success("任务已删除");
    load();
  }, "这个任务");
};

onMounted(load);
</script>
