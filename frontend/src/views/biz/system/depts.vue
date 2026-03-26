<template>
  <a-card :bordered="false">
    <template #title>部门管理</template>
    <a-space direction="vertical" fill>
      <a-button type="primary" @click="open()">新增部门</a-button>
      <a-table :data="flatList" :pagination="false" row-key="id">
        <template #columns>
          <a-table-column title="名称" data-index="name" />
          <a-table-column title="编码" data-index="code" />
          <a-table-column title="负责人" data-index="leader" />
          <a-table-column title="状态" data-index="status" />
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
  <a-modal v-model:visible="visible" :title="form.id ? '编辑部门' : '新增部门'" :on-before-ok="submit">
    <a-form :model="form" layout="vertical">
      <a-form-item label="父级部门">
        <a-tree-select v-model="form.parent_id" :data="tree" :field-names="{ key: 'id', title: 'name', children: 'children' }" allow-clear />
      </a-form-item>
      <a-form-item label="名称"><a-input v-model="form.name" /></a-form-item>
      <a-form-item label="编码"><a-input v-model="form.code" /></a-form-item>
      <a-form-item label="负责人"><a-input v-model="form.leader" /></a-form-item>
      <a-form-item label="电话"><a-input v-model="form.phone" /></a-form-item>
      <a-form-item label="邮箱"><a-input v-model="form.email" /></a-form-item>
      <a-form-item label="排序"><a-input-number v-model="form.sort" :min="0" /></a-form-item>
      <a-form-item label="状态"><a-input-number v-model="form.status" :min="1" :max="2" /></a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { deptDeleteAPI, deptListAPI, deptSaveAPI, type DeptNode } from "@/api/system";
import { useConfirmAction } from "@/hooks/useConfirmAction";

const tree = ref<DeptNode[]>([]);
const visible = ref(false);
const { confirmDelete, confirmSave } = useConfirmAction();
const form = reactive<any>({ id: null, parent_id: undefined, name: "", code: "", leader: "", phone: "", email: "", sort: 0, status: 1 });

const load = async () => {
  const res = await deptListAPI();
  tree.value = res.data;
};

const flatList = computed(() => {
  const out: DeptNode[] = [];
  const walk = (nodes: DeptNode[]) => {
    nodes.forEach(node => {
      out.push(node);
      if (node.children?.length) walk(node.children);
    });
  };
  walk(tree.value);
  return out;
});

const open = (record?: any) => {
  Object.assign(form, { id: null, parent_id: undefined, name: "", code: "", leader: "", phone: "", email: "", sort: 0, status: 1 }, record || {});
  visible.value = true;
};

const submit = async () => {
  return confirmSave(async () => {
    await deptSaveAPI(form.id, form);
    Message.success("部门已保存");
    visible.value = false;
    load();
  }, form.id ? "当前部门" : "新部门");
};

const remove = async (id: number) => {
  await confirmDelete(async () => {
    await deptDeleteAPI(id);
    Message.success("部门已删除");
    load();
  }, "这个部门");
};

onMounted(load);
</script>
