<template>
  <a-card :bordered="false">
    <template #title>API 管理</template>
    <a-alert>这里复用后台菜单中配置了 <code>method + api_path</code> 的记录作为 API 权限清单。</a-alert>
    <a-table :data="apiRows" :pagination="false" row-key="id" class="table">
      <template #columns>
        <a-table-column title="标题" data-index="title" />
        <a-table-column title="权限标识" data-index="permission" />
        <a-table-column title="方法" data-index="method" />
        <a-table-column title="路径" data-index="api_path" />
        <a-table-column title="菜单类型" data-index="type" />
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { adminMenuTreeAPI, type AdminMenuNode } from "@/api/system";

const tree = ref<AdminMenuNode[]>([]);

const flatten = (nodes: AdminMenuNode[]): AdminMenuNode[] => {
  const out: AdminMenuNode[] = [];
  const walk = (items: AdminMenuNode[]) => {
    items.forEach(item => {
      out.push(item);
      if (item.children?.length) walk(item.children);
    });
  };
  walk(nodes);
  return out;
};

const apiRows = computed(() =>
  flatten(tree.value).filter(item => item.api_path && item.method)
);

onMounted(async () => {
  const res = await adminMenuTreeAPI();
  tree.value = res.data;
});
</script>

<style scoped>
.table {
  margin-top: 16px;
}
</style>
