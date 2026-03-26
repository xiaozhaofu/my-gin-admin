<template>
  <div class="snow-page">
    <div class="admin-page">
      <section class="page-hero">
        <div>
          <div class="page-badge">Admin Console</div>
          <h1>管理员</h1>
          <p>统一管理后台账号、部门、岗位与角色归属，整体风格和首页、内容管理页保持一致。</p>
        </div>
        <div class="page-hero-side">
          <div class="page-hero-label">当前账号总数</div>
          <div class="page-hero-value">{{ list.length }}</div>
        </div>
      </section>

      <a-card class="panel-card" :bordered="false">
        <template #title>管理员列表</template>
        <template #extra>
          <a-button v-if="session.can('/api/v1/admins#POST')" type="primary" @click="open()">新增管理员</a-button>
        </template>
        <a-table :data="list" :pagination="false">
          <template #columns>
            <a-table-column title="账号">
              <template #cell="{ record }">
                <div class="cell-main">
                  <div class="cell-title">{{ record.username }}</div>
                  <div class="cell-sub">{{ record.nickname }}</div>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="联系方式">
              <template #cell="{ record }">
                <div class="cell-main">
                  <div class="cell-sub">手机：{{ record.phone || "未填写" }}</div>
                  <div class="cell-sub">邮箱：{{ record.email || "未填写" }}</div>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="状态" data-index="status" />
            <a-table-column title="角色">
              <template #cell="{ record }">{{ record.roles?.join(" / ") || "未分配角色" }}</template>
            </a-table-column>
            <a-table-column title="组织归属">
              <template #cell="{ record }">
                <div class="cell-main">
                  <div class="cell-sub">部门ID：{{ record.dept_id ?? "-" }}</div>
                  <div class="cell-sub">岗位ID：{{ record.post_id ?? "-" }}</div>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="操作">
              <template #cell="{ record }">
                <a-button v-if="session.can('/api/v1/admins/:id#PUT')" size="mini" @click="open(record)">编辑</a-button>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>

  <a-modal v-model:visible="visible" title="管理员" :on-before-ok="submit">
    <a-form :model="current" layout="vertical">
      <a-form-item label="账号"><a-input v-model="current.username" /></a-form-item>
      <a-form-item label="昵称"><a-input v-model="current.nickname" /></a-form-item>
      <a-form-item label="密码"><a-input-password v-model="current.password" placeholder="更新时留空则保持原密码" /></a-form-item>
      <a-form-item label="部门">
        <a-tree-select v-model="current.dept_id" :data="depts" :field-names="{ key: 'id', title: 'name', children: 'children' }" allow-clear />
      </a-form-item>
      <a-form-item label="岗位">
        <a-select v-model="current.post_id" allow-clear>
          <a-option v-for="item in posts" :key="item.id" :value="item.id">{{ item.name }}</a-option>
        </a-select>
      </a-form-item>
      <a-form-item label="角色">
        <a-select v-model="current.role_ids" multiple allow-clear placeholder="请选择角色">
          <a-option v-for="item in roles" :key="item.id" :value="item.id">{{ item.name }}</a-option>
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { adminListAPI, adminSaveAPI, deptListAPI, postListAPI, roleListAPI, type AdminItem, type DeptNode, type PostItem, type RoleItem } from "@/api/system";
import { useConfirmAction } from "@/hooks/useConfirmAction";
import { useSessionStore } from "@/store/modules/session";

const list = ref<AdminItem[]>([]);
const roles = ref<RoleItem[]>([]);
const depts = ref<DeptNode[]>([]);
const posts = ref<PostItem[]>([]);
const visible = ref(false);
const session = useSessionStore();
const { confirmSave } = useConfirmAction();
const current = reactive<any>({ id: null, username: "", nickname: "", password: "", dept_id: undefined, post_id: undefined, role_ids: [] });

const load = async () => {
  const res = await adminListAPI();
  list.value = res.data;
};

const loadRoles = async () => {
  const res = await roleListAPI();
  roles.value = res.data;
};

const loadDepts = async () => {
  const res = await deptListAPI();
  depts.value = res.data;
};

const loadPosts = async () => {
  const res = await postListAPI();
  posts.value = res.data;
};

const open = (record?: any) => {
  Object.assign(current, record || { id: null, username: "", nickname: "", password: "", dept_id: undefined, post_id: undefined, role_ids: [] });
  current.role_ids = Array.isArray(current.role_ids) ? [...current.role_ids] : [];
  visible.value = true;
};

const submit = async () => {
  return confirmSave(async () => {
    await adminSaveAPI(current.id, {
      username: current.username,
      nickname: current.nickname,
      password: current.password,
      dept_id: current.dept_id,
      post_id: current.post_id,
      role_ids: current.role_ids
    });
    Message.success("保存成功");
    visible.value = false;
    load();
  }, current.id ? "当前管理员" : "新管理员");
};

onMounted(async () => {
  await Promise.all([load(), loadRoles(), loadDepts(), loadPosts()]);
});
</script>

<style scoped lang="scss">
.admin-page {
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

.cell-main {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.cell-title {
  font-weight: 700;
}

.cell-sub {
  font-size: 12px;
  color: var(--color-text-3);
}

@media (max-width: 1080px) {
  .page-hero {
    grid-template-columns: 1fr;
  }
}
</style>
