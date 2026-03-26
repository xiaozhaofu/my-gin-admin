<template>
  <a-card :bordered="false">
    <template #title>后台菜单</template>
    <a-space direction="vertical" fill>
      <a-button v-if="session.can('/api/v1/admin-menus#POST')" type="primary" @click="open()">新增后台菜单</a-button>
      <a-table :data="flatList" :pagination="false" row-key="id">
        <template #columns>
          <a-table-column title="标题" data-index="title" />
          <a-table-column title="路由" data-index="path" />
          <a-table-column title="组件" data-index="component" />
          <a-table-column title="权限标识" data-index="permission" />
          <a-table-column title="接口">
            <template #cell="{ record }">{{ record.method }} {{ record.api_path }}</template>
          </a-table-column>
          <a-table-column title="类型" data-index="type" />
          <a-table-column title="操作">
            <template #cell="{ record }">
              <a-space>
                <a-button v-if="session.can('/api/v1/admin-menus/:id#PUT')" size="mini" @click="open(record)">编辑</a-button>
                <a-button v-if="session.can('/api/v1/admin-menus/:id#DELETE')" size="mini" status="danger" @click="remove(record.id)">删除</a-button>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-space>
  </a-card>
  <a-modal v-model:visible="visible" :title="current.id ? '编辑后台菜单' : '新增后台菜单'" :on-before-ok="submit">
    <a-form :model="current" layout="vertical">
      <a-form-item label="父级菜单">
        <a-tree-select
          v-model="current.parent_id"
          :data="tree"
          :field-names="{ key: 'id', title: 'title', children: 'children' }"
          allow-clear
        />
      </a-form-item>
      <a-form-item label="标题"><a-input v-model="current.title" /></a-form-item>
      <a-form-item label="名称"><a-input v-model="current.name" /></a-form-item>
      <a-form-item label="路由路径"><a-input v-model="current.path" /></a-form-item>
      <a-form-item label="组件">
        <a-auto-complete v-model="current.component" :data="componentOptions" placeholder="选择或输入组件路径" />
      </a-form-item>
      <a-form-item label="图标"><a-input v-model="current.icon" /></a-form-item>
      <a-form-item label="权限标识"><a-input v-model="current.permission" /></a-form-item>
      <a-form-item label="类型"><a-input-number v-model="current.type" :min="1" :max="3" /></a-form-item>
      <a-form-item label="排序"><a-input-number v-model="current.sort" :min="0" /></a-form-item>
      <a-form-item label="接口方法"><a-input v-model="current.method" placeholder="GET/POST/PUT/DELETE" /></a-form-item>
      <a-form-item label="接口路径"><a-input v-model="current.api_path" placeholder="/api/v1/..." /></a-form-item>
      <a-form-item label="隐藏"><a-switch v-model="current.hidden" /></a-form-item>
      <a-form-item label="缓存"><a-switch v-model="current.keep_alive" /></a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { adminMenuDeleteAPI, adminMenuSaveAPI, adminMenuTreeAPI, type AdminMenuNode } from "@/api/system";
import { useConfirmAction } from "@/hooks/useConfirmAction";
import { useSessionStore } from "@/store/modules/session";

const visible = ref(false);
const session = useSessionStore();
const { confirmDelete, confirmSave } = useConfirmAction();
const tree = ref<AdminMenuNode[]>([]);
const current = reactive<any>({
  id: null,
  parent_id: undefined,
  title: "",
  name: "",
  path: "",
  component: "",
  icon: "",
  permission: "",
  type: 1,
  sort: 0,
  hidden: false,
  keep_alive: true,
  method: "",
  api_path: ""
});

const componentOptions = [
  "home/home",
  "biz/articles/index",
  "biz/articles/form",
  "biz/menus/index",
  "biz/uploads/index",
  "biz/system/admin-menus",
  "biz/system/admins",
  "biz/system/roles",
  "biz/modules/index"
];

const load = async () => {
  const res = await adminMenuTreeAPI();
  tree.value = res.data;
};

const flatList = computed(() => {
  const out: AdminMenuNode[] = [];
  const walk = (nodes: AdminMenuNode[]) => {
    nodes.forEach(node => {
      out.push(node);
      if (node.children?.length) walk(node.children);
    });
  };
  walk(tree.value);
  return out;
});

const open = (record?: Partial<AdminMenuNode>) => {
  Object.assign(current, {
    id: null,
    parent_id: undefined,
    title: "",
    name: "",
    path: "",
    component: "",
    icon: "",
    permission: "",
    type: 1,
    sort: 0,
    hidden: false,
    keep_alive: true,
    method: "",
    api_path: ""
  }, record || {});
  visible.value = true;
};

const submit = async () => {
  return confirmSave(async () => {
    await adminMenuSaveAPI(current.id, current);
    Message.success("保存成功，角色权限策略已同步");
    visible.value = false;
    load();
  }, current.id ? "当前后台菜单" : "新后台菜单");
};

const remove = async (id: number) => {
  await confirmDelete(async () => {
    await adminMenuDeleteAPI(id);
    Message.success("删除成功");
    load();
  }, "这个后台菜单");
};

onMounted(load);
</script>
