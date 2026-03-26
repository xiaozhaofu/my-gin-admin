<template>
  <div class="snow-page">
    <div class="role-page">
      <section class="page-hero">
        <div>
          <div class="page-badge">Permission Center</div>
          <h1>角色权限</h1>
          <p>统一管理角色、数据范围和后台菜单权限，并实时查看会授予哪些 API 权限。</p>
        </div>
        <div class="page-hero-side">
          <div class="page-hero-label">当前角色总数</div>
          <div class="page-hero-value">{{ list.length }}</div>
        </div>
      </section>

      <a-card class="panel-card" :bordered="false">
        <template #title>角色列表</template>
        <template #extra>
          <a-button v-if="session.can('/api/v1/roles#POST')" type="primary" @click="open()">新增角色</a-button>
        </template>
        <a-table :data="list" :pagination="false">
          <template #columns>
            <a-table-column title="名称" data-index="name" />
            <a-table-column title="编码" data-index="code" />
            <a-table-column title="说明" data-index="description" />
            <a-table-column title="数据范围" data-index="data_scope" />
            <a-table-column title="操作">
              <template #cell="{ record }">
                <a-button v-if="session.can('/api/v1/roles/:id#PUT')" size="mini" @click="open(record)">编辑</a-button>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>

  <a-modal v-model:visible="visible" title="角色" :on-before-ok="submit">
    <a-form :model="current" layout="vertical">
      <a-form-item label="名称"><a-input v-model="current.name" /></a-form-item>
      <a-form-item label="编码"><a-input v-model="current.code" /></a-form-item>
      <a-form-item label="说明"><a-input v-model="current.description" /></a-form-item>
      <a-form-item label="数据范围">
        <a-select v-model="current.data_scope">
          <a-option :value="1">全部数据</a-option>
          <a-option :value="2">本部门数据</a-option>
          <a-option :value="3">本部门及以下数据</a-option>
          <a-option :value="4">仅本人数据</a-option>
        </a-select>
      </a-form-item>
      <a-form-item label="菜单权限">
        <a-tree-select
          v-model="current.menu_ids"
          :data="menuTree"
          :field-names="{ key: 'id', title: 'title', children: 'children' }"
          multiple
          tree-checkable
          allow-clear
          placeholder="选择角色可访问的后台菜单"
        />
      </a-form-item>
      <a-form-item label="将授予的 API 权限">
        <a-space direction="vertical" fill>
          <a-alert v-if="grantedAPIs.length === 0">当前未选中任何可产生 API 权限的后台菜单。</a-alert>
          <a-tag v-for="item in grantedAPIs" :key="item" color="arcoblue">{{ item }}</a-tag>
        </a-space>
      </a-form-item>
      <a-form-item label="权限标识规范">
        <a-alert>
          页面展示权限由后台菜单的 <code>permission</code> 控制，API 鉴权由 <code>method + api_path</code> 生成 Casbin 策略。
        </a-alert>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { adminMenuTreeAPI, roleListAPI, roleSaveAPI, type AdminMenuNode, type RoleItem } from "@/api/system";
import { useConfirmAction } from "@/hooks/useConfirmAction";
import { useSessionStore } from "@/store/modules/session";

const list = ref<RoleItem[]>([]);
const menuTree = ref<AdminMenuNode[]>([]);
const visible = ref(false);
const session = useSessionStore();
const { confirmSave } = useConfirmAction();
const current = reactive<any>({ id: null, name: "", code: "", description: "", data_scope: 1, menu_ids: [] });

const load = async () => {
  const res = await roleListAPI();
  list.value = res.data;
};

const loadMenuTree = async () => {
  const res = await adminMenuTreeAPI();
  menuTree.value = res.data;
};

const flattenMenus = (nodes: AdminMenuNode[]): AdminMenuNode[] => {
  const out: AdminMenuNode[] = [];
  const walk = (items: AdminMenuNode[]) => {
    items.forEach(item => {
      out.push(item);
      if (item.children?.length) {
        walk(item.children);
      }
    });
  };
  walk(nodes);
  return out;
};

const grantedAPIs = computed(() => {
  const selected = new Set<number>(current.menu_ids || []);
  const menus = flattenMenus(menuTree.value);
  return menus
    .filter(item => selected.has(item.id) && item.api_path && item.method)
    .map(item => `${item.method} ${item.api_path}`);
});

const open = (record?: any) => {
  Object.assign(current, record || { id: null, name: "", code: "", description: "", data_scope: 1, menu_ids: [] });
  current.menu_ids = Array.isArray(current.menu_ids) ? [...current.menu_ids] : [];
  visible.value = true;
};

const submit = async () => {
  return confirmSave(async () => {
    await roleSaveAPI(current.id, {
      name: current.name,
      code: current.code,
      description: current.description,
      data_scope: current.data_scope,
      menu_ids: current.menu_ids
    });
    Message.success("保存成功");
    visible.value = false;
    load();
  }, current.id ? "当前角色" : "新角色");
};

onMounted(async () => {
  await Promise.all([load(), loadMenuTree()]);
});
</script>

<style scoped lang="scss">
.role-page {
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
    linear-gradient(135deg, rgba(15, 23, 42, 0.96), rgba(22, 93, 255, 0.9)),
    radial-gradient(circle at top right, rgba(255, 255, 255, 0.16), transparent 38%);
  color: #fff;
  box-shadow: 0 20px 38px rgb(15 23 42 / 14%);
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
