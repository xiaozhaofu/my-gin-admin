<template>
  <div class="snow-page">
    <div class="article-list-page">
      <section class="page-hero">
        <div>
          <div class="page-badge">Content Workspace</div>
          <h1>文章管理</h1>
          <p>统一查看文章、资源预览、三级菜单筛选与批量状态处理，整体风格与工作台、编辑页保持一致。</p>
        </div>
        <div class="page-hero-side">
          <div class="page-hero-label">当前文章总数</div>
          <div class="page-hero-value">{{ pagination.total }}</div>
        </div>
      </section>

      <a-card class="panel-card" :bordered="false">
        <template #title>筛选条件</template>
        <a-form :model="filters" layout="inline" class="filters">
          <a-form-item field="title" label="标题">
            <a-input v-model="filters.title" placeholder="文章标题" allow-clear />
          </a-form-item>
          <a-form-item field="type" label="类型">
            <a-select v-model="filters.type" allow-clear placeholder="全部类型">
              <a-option v-for="item in articleTypes" :key="item.value" :value="item.value">{{ item.label }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item field="menuId" label="三级菜单">
            <a-cascader v-model="filters.menu_id" :options="menuOptions" allow-clear />
          </a-form-item>
          <a-form-item field="status" label="状态">
            <a-select v-model="filters.status" allow-clear placeholder="全部状态">
              <a-option :value="1">正常</a-option>
              <a-option :value="2">隐藏</a-option>
              <a-option :value="3">待审核</a-option>
              <a-option :value="4">已驳回</a-option>
            </a-select>
          </a-form-item>
          <a-space wrap>
            <a-button type="primary" @click="fetchList">搜索</a-button>
            <a-button @click="resetFilters">重置</a-button>
            <a-button v-if="session.can('/api/v1/articles/batch#POST')" type="outline" @click="router.push('/articles/batch')">批量新增</a-button>
            <a-button v-if="session.can('/api/v1/articles#POST')" type="outline" @click="router.push('/articles/new')">新增文章</a-button>
          </a-space>
        </a-form>
      </a-card>

      <a-card class="panel-card" :bordered="false">
        <template #title>文章列表</template>
        <template #extra>
          <a-space wrap>
            <a-button v-if="canSelectArticles" :disabled="!currentPageRowKeys.length" @click="selectCurrentPage">选中当前页</a-button>
            <a-button v-if="canSelectArticles" :disabled="!selectedRowKeys.length" @click="clearSelected">取消选中</a-button>
            <a-button v-if="session.can('/api/v1/articles/status#PUT')" :disabled="!selectedRowKeys.length" @click="batchChangeStatus(1)">批量正常</a-button>
            <a-button v-if="session.can('/api/v1/articles/status#PUT')" :disabled="!selectedRowKeys.length" @click="batchChangeStatus(2)">批量隐藏</a-button>
            <a-button v-if="session.can('/api/v1/articles#DELETE')" status="danger" :disabled="!selectedRowKeys.length" @click="batchDelete">批量删除</a-button>
          </a-space>
        </template>

        <a-table
          row-key="id"
          :data="list"
          :pagination="pagination"
          :row-selection="canSelectArticles ? { type: 'checkbox' } : undefined"
          :selected-keys="selectedRowKeys"
          @page-change="pageChange"
          @selection-change="setSelected"
        >
          <template #columns>
            <a-table-column title="标题">
              <template #cell="{ record }">
                <div class="article-title-cell">
                  <a-link class="article-title-link" @click="router.push(`/articles/${record.id}`)">{{ record.title }}</a-link>
                  <div class="article-title-meta">菜单：{{ menuLabels(record) }} / 渠道：{{ channelLabel(record.channel_id) }}</div>
                  <a-image v-if="coverUrl(record)" :src="coverUrl(record)" width="88" />
                </div>
              </template>
            </a-table-column>
            <a-table-column title="类型" data-index="type" />
            <a-table-column title="内容预览">
              <template #cell="{ record }">
                <div v-if="record.type === 2" class="preview-block">
                  <video v-if="extractSrc(record)" :src="extractSrc(record)" controls class="preview-media" />
                  <a-space class="preview-actions" size="mini">
                    <a-button size="mini" @click="openInNewTab(extractSrc(record))">打开</a-button>
                    <a-button size="mini" @click="copyUrl(extractSrc(record))">复制</a-button>
                    <a-button size="mini" @click="downloadUrl(extractSrc(record))">下载</a-button>
                  </a-space>
                </div>
                <div v-else-if="record.type === 4" class="preview-block">
                  <audio v-if="extractSrc(record)" :src="extractSrc(record)" controls class="preview-audio" />
                  <a-space class="preview-actions" size="mini">
                    <a-button size="mini" @click="openInNewTab(extractSrc(record))">打开</a-button>
                    <a-button size="mini" @click="copyUrl(extractSrc(record))">复制</a-button>
                    <a-button size="mini" @click="downloadUrl(extractSrc(record))">下载</a-button>
                  </a-space>
                </div>
                <div v-else-if="record.type === 6" class="preview-block">
                  <a-image v-if="extractSrc(record)" :src="extractSrc(record)" width="96" />
                  <a-space class="preview-actions" size="mini">
                    <a-button size="mini" @click="openInNewTab(extractSrc(record))">打开</a-button>
                    <a-button size="mini" @click="copyUrl(extractSrc(record))">复制</a-button>
                    <a-button size="mini" @click="downloadUrl(extractSrc(record))">下载</a-button>
                  </a-space>
                </div>
                <a-typography-paragraph v-else :ellipsis="{ rows: 2 }">
                  {{ extractText(record.content?.content || "") }}
                </a-typography-paragraph>
              </template>
            </a-table-column>
            <a-table-column title="状态">
              <template #cell="{ record }">
                <a-select
                  v-if="session.can('/api/v1/articles/status#PUT')"
                  :model-value="record.status"
                  size="mini"
                  style="width: 120px"
                  @change="(value: string | number | boolean) => changeSingleStatus(record, Number(value))"
                >
                  <a-option v-for="item in statusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
                </a-select>
                <a-tag v-else :color="statusColor(record.status)">{{ statusLabel(record.status) }}</a-tag>
              </template>
            </a-table-column>
            <a-table-column title="浏览量" data-index="view_num" />
            <a-table-column title="创建时间" data-index="created_at" />
            <a-table-column title="操作">
              <template #cell="{ record }">
                <a-space>
                  <a-button v-if="session.can('/api/v1/articles/:id#PUT')" size="mini" @click="router.push(`/articles/${record.id}`)">编辑</a-button>
                  <a-button v-if="session.can('/api/v1/articles#DELETE')" size="mini" status="danger" @click="removeOne(record.id)">删除</a-button>
                </a-space>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { Message } from "@arco-design/web-vue";
import { useRoute, useRouter } from "vue-router";
import { articleDeleteAPI, articleListAPI, articleStatusAPI, type ArticleItem } from "@/api/article";
import { channelListAPI, type ChannelItem } from "@/api/channel";
import { menuCascaderAPI } from "@/api/menu";
import { useConfirmAction } from "@/hooks/useConfirmAction";
import { usePageSelection } from "@/hooks/usePageSelection";
import { useSessionStore } from "@/store/modules/session";

const route = useRoute();
const router = useRouter();
const session = useSessionStore();
const { confirmBatchDelete, confirmBatchHide, confirmDelete, runConfirmed } = useConfirmAction();
const list = ref<ArticleItem[]>([]);
const channels = ref<ChannelItem[]>([]);
const menuOptions = ref<any[]>([]);
const menuTree = ref<any[]>([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const filters = reactive<Record<string, any>>({
  title: "",
  type: undefined,
  menu_id: undefined,
  status: undefined
});
const {
  selectedKeys: selectedRowKeys,
  currentPageKeys: currentPageRowKeys,
  setSelected,
  selectCurrentPage,
  clearSelected,
  removeSelected
} = usePageSelection(() => list.value, item => item.id);
const canSelectArticles = computed(() => session.can('/api/v1/articles/status#PUT') || session.can('/api/v1/articles#DELETE'));

const articleTypes = [
  { label: "纯文本", value: 1 },
  { label: "视频", value: 2 },
  { label: "视频集合", value: 3 },
  { label: "音频", value: 4 },
  { label: "音频集合", value: 5 },
  { label: "图片", value: 6 },
  { label: "图集", value: 7 }
];

const statusOptions = [
  { label: "正常", value: 1 },
  { label: "隐藏", value: 2 },
  { label: "待审核", value: 3 },
  { label: "已驳回", value: 4 }
];

const fetchList = async () => {
  const menuID = Array.isArray(filters.menu_id) ? filters.menu_id.at(-1) : filters.menu_id;
  const res = await articleListAPI({ ...filters, menu_id: menuID, page: pagination.current, page_size: pagination.pageSize });
  list.value = res.data.list;
  pagination.total = res.data.total;
};

const fetchMenus = async () => {
  const res = await menuCascaderAPI();
  menuTree.value = res.data;
  menuOptions.value = convertMenus(res.data);
  syncMenuFilterFromRoute();
};

const fetchChannels = async () => {
  const res = await channelListAPI();
  channels.value = res.data;
};

const convertMenus = (items: any[]): any[] =>
  items.map(item => ({
    value: item.id,
    label: item.name,
    children: item.children ? convertMenus(item.children) : undefined
  }));

const flattenMenuMap = (nodes: any[], out = new Map<number, string>()) => {
  nodes.forEach(node => {
    out.set(Number(node.id), node.name);
    if (node.children?.length) {
      flattenMenuMap(node.children, out);
    }
  });
  return out;
};

const menuNameMap = computed(() => flattenMenuMap(menuTree.value));

const findMenuPath = (nodes: any[], targetID?: number, parents: number[] = []): number[] | undefined => {
  if (!targetID) return undefined;
  for (const node of nodes) {
    const next = [...parents, node.id];
    if (node.id === targetID) {
      return next;
    }
    const childPath = findMenuPath(node.children || [], targetID, next);
    if (childPath) {
      return childPath;
    }
  }
  return undefined;
};

const syncMenuFilterFromRoute = () => {
  const rawMenuID = route.query.menu_id;
  const menuID = typeof rawMenuID === "string" ? Number(rawMenuID) : Number(Array.isArray(rawMenuID) ? rawMenuID[0] : 0);
  if (menuID > 0) {
    filters.menu_id = findMenuPath(menuTree.value, menuID) || [menuID];
  }
};

const resetFilters = () => {
  filters.title = "";
  filters.type = undefined;
  filters.menu_id = undefined;
  filters.status = undefined;
  pagination.current = 1;
  selectedRowKeys.value = [];
  fetchList();
};

const pageChange = (page: number) => {
  pagination.current = page;
  selectedRowKeys.value = [];
  fetchList();
};

const batchDelete = async () => {
  await confirmBatchDelete(async () => {
    await articleDeleteAPI(selectedRowKeys.value);
    Message.success("删除成功");
    selectedRowKeys.value = [];
    fetchList();
  }, "当前页选中的文章");
};

const removeOne = async (id: number) => {
  await confirmDelete(async () => {
    await articleDeleteAPI([id]);
    Message.success("删除成功");
    removeSelected(id);
    fetchList();
  }, "这篇文章");
};

const batchChangeStatus = async (status: number) => {
  const statusName = statusLabel(status);
  const action = async () => {
    await articleStatusAPI({ ids: selectedRowKeys.value, status });
    Message.success("状态已更新");
    selectedRowKeys.value = [];
    fetchList();
  };

  if (status === 2) {
    await confirmBatchHide(action);
    return;
  }

  await runConfirmed(
    {
      title: "确认批量修改状态",
      content: `确认将当前页选中的文章批量设置为${statusName}吗？`
    },
    action
  );
};

const changeSingleStatus = async (record: ArticleItem, status: number) => {
  if (record.status === status) {
    return;
  }
  await runConfirmed(
    {
      title: "确认修改状态",
      content: `确认将《${record.title}》的状态修改为${statusLabel(status)}吗？`
    },
    async () => {
      await articleStatusAPI({ ids: [record.id], status });
      Message.success("状态已更新");
      fetchList();
    }
  );
};

const extractText = (html: string) => html.replace(/<[^>]+>/g, " ").replace(/\s+/g, " ").trim();

const extractSrc = (record: ArticleItem) => {
  const html = record.content?.content || "";
  const matched = html.match(/src="([^"]+)"/i);
  if (matched?.[1]) {
    return matched[1];
  }
  const link = html.match(/href="([^"]+)"/i);
  return link?.[1] || "";
};

const channelLabel = (channelID: number) => {
  const channel = channels.value.find(item => item.id === channelID);
  if (!channel) return String(channelID);
  return `${channel.name}（${channel.code}）`;
};

const menuLabels = (record: ArticleItem) => {
  const ids = Array.isArray(record.menu_ids) && record.menu_ids.length ? record.menu_ids : [record.menu_id];
  return ids
    .map(id => menuNameMap.value.get(Number(id)) || String(id))
    .filter(Boolean)
    .join(" / ");
};

const coverUrl = (record: ArticleItem) => {
  if (record.type === 4) {
    return record.cover_small || record.cover_medium || record.cover_large || "";
  }
  return record.cover_large || record.cover_medium || record.cover_small || "";
};

const statusLabel = (status: number) => {
  switch (status) {
    case 1:
      return "正常";
    case 2:
      return "隐藏";
    case 3:
      return "待审核";
    case 4:
      return "已驳回";
    default:
      return String(status);
  }
};

const statusColor = (status: number) => {
  switch (status) {
    case 1:
      return "green";
    case 2:
      return "gray";
    case 3:
      return "arcoblue";
    case 4:
      return "red";
    default:
      return "gray";
  }
};

const openInNewTab = (url: string) => {
  if (!url) return;
  window.open(url, "_blank", "noopener,noreferrer");
};

const copyUrl = async (url: string) => {
  if (!url) return;
  await navigator.clipboard.writeText(url);
  Message.success("地址已复制");
};

const downloadUrl = (url: string) => {
  if (!url) return;
  const a = document.createElement("a");
  a.href = url;
  a.download = "";
  a.target = "_blank";
  document.body.appendChild(a);
  a.click();
  a.remove();
};

onMounted(async () => {
  await Promise.all([fetchList(), fetchMenus(), fetchChannels()]);
});

watch(
  () => route.query.menu_id,
  () => {
    syncMenuFilterFromRoute();
    fetchList();
  }
);
</script>

<style scoped lang="scss">
.article-list-page {
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

.article-title-cell {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.article-title-link {
  font-weight: 700;
}

.article-title-meta {
  font-size: 12px;
  color: var(--color-text-3);
}

.preview-block {
  max-width: 120px;
}

.preview-media {
  width: 120px;
  max-height: 72px;
  border-radius: 8px;
}

.preview-audio {
  width: 120px;
}

.preview-actions {
  margin-top: 8px;
}

@media (max-width: 1080px) {
  .page-hero {
    grid-template-columns: 1fr;
  }
}
</style>
