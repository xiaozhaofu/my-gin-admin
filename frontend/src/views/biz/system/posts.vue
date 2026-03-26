<template>
  <a-card :bordered="false">
    <template #title>岗位管理</template>
    <a-space direction="vertical" fill>
      <a-button type="primary" @click="open()">新增岗位</a-button>
      <a-table :data="list" :pagination="false" row-key="id">
        <template #columns>
          <a-table-column title="名称" data-index="name" />
          <a-table-column title="编码" data-index="code" />
          <a-table-column title="排序" data-index="sort" />
          <a-table-column title="状态" data-index="status" />
          <a-table-column title="备注" data-index="remark" />
          <a-table-column title="操作">
            <template #cell="{ record }">
              <a-space>
                <a-button size="mini" @click="open(record)">编辑</a-button>
                <a-button size="mini" status="danger" @click="remove(record.id)">删除</a-button>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-space>
  </a-card>
  <a-modal v-model:visible="visible" :title="form.id ? '编辑岗位' : '新增岗位'" :on-before-ok="submit">
    <a-form :model="form" layout="vertical">
      <a-form-item label="名称"><a-input v-model="form.name" /></a-form-item>
      <a-form-item label="编码"><a-input v-model="form.code" /></a-form-item>
      <a-form-item label="排序"><a-input-number v-model="form.sort" :min="0" /></a-form-item>
      <a-form-item label="状态"><a-input-number v-model="form.status" :min="1" :max="2" /></a-form-item>
      <a-form-item label="备注"><a-input v-model="form.remark" /></a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { postDeleteAPI, postListAPI, postSaveAPI, type PostItem } from "@/api/system";
import { useConfirmAction } from "@/hooks/useConfirmAction";

const list = ref<PostItem[]>([]);
const visible = ref(false);
const { confirmDelete, confirmSave } = useConfirmAction();
const form = reactive<any>({ id: null, name: "", code: "", sort: 0, status: 1, remark: "" });

const load = async () => {
  const res = await postListAPI();
  list.value = res.data;
};

const open = (record?: any) => {
  Object.assign(form, { id: null, name: "", code: "", sort: 0, status: 1, remark: "" }, record || {});
  visible.value = true;
};

const submit = async () => {
  return confirmSave(async () => {
    await postSaveAPI(form.id, form);
    Message.success("岗位已保存");
    visible.value = false;
    load();
  }, form.id ? "当前岗位" : "新岗位");
};

const remove = async (id: number) => {
  await confirmDelete(async () => {
    await postDeleteAPI(id);
    Message.success("岗位已删除");
    load();
  }, "这个岗位");
};

onMounted(load);
</script>
