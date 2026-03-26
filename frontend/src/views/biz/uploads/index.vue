<template>
  <div class="snow-page">
    <div class="upload-page">
      <section class="page-hero">
        <div>
          <div class="page-badge">Asset Center</div>
          <h1>资源上传</h1>
          <p>统一管理图片、视频、音频资源，并按云存储类型、场景和文件类型进行筛选。页面风格与首页、文章编辑页保持一致。</p>
        </div>
        <div class="page-hero-side">
          <div class="page-hero-label">当前资源总数</div>
          <div class="page-hero-value">{{ pagination.total }}</div>
        </div>
      </section>

      <a-card class="panel-card" :bordered="false">
        <template #title>筛选条件</template>
        <a-form :model="filters" layout="inline" class="filters">
          <a-form-item label="名称">
            <a-input v-model="filters.origin_name" allow-clear placeholder="按原始文件名筛选" />
          </a-form-item>
          <a-form-item label="文件类型">
            <a-select v-model="filters.type" allow-clear style="width: 120px">
              <a-option :value="1">图片</a-option>
              <a-option :value="2">视频</a-option>
              <a-option :value="4">音频</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="云存储类型">
            <a-select v-model="filters.provider" allow-clear style="width: 160px" placeholder="全部类型">
              <a-option value="aliyun-oss">阿里云</a-option>
              <a-option value="tencent-cos">腾讯云</a-option>
              <a-option value="huawei-obs">华为云</a-option>
              <a-option value="local">本地存储</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="场景">
            <a-input v-model="filters.scene" allow-clear placeholder="按业务场景筛选" />
          </a-form-item>
          <a-space wrap>
            <a-button type="primary" @click="load">搜索</a-button>
            <a-button @click="reset">重置</a-button>
          </a-space>
        </a-form>
      </a-card>

      <a-card class="panel-card" :bordered="false">
        <template #title>上传入口</template>
        <div class="upload-actions">
          <a-select v-model="uploadPreference.provider" style="width: 220px" placeholder="设置全局默认上传类型">
            <a-option value="aliyun-oss">上传到阿里云 OSS</a-option>
            <a-option value="tencent-cos">上传到腾讯云 COS</a-option>
            <a-option value="huawei-obs">上传到华为云 OBS</a-option>
            <a-option value="local">上传到本地存储</a-option>
          </a-select>
          <div class="upload-default-tip">这里设置的是全局默认上传类型，文章新增、批量新增、菜单图标上传都会跟随这个默认值。</div>
          <a-upload :custom-request="upload" :show-file-list="false" multiple :accept="accept">
            <template #upload-button>
              <a-button type="primary">上传图片 / 音频 / 视频</a-button>
            </template>
          </a-upload>
          <a-alert>允许图片、音频、视频，音频明确支持 wav。上传时会按当前全局默认上传类型写入对应云存储或本地存储。</a-alert>
        </div>
      </a-card>

      <a-card class="panel-card" :bordered="false">
        <template #title>资源列表</template>
        <template #extra>
          <a-space wrap>
            <a-button v-if="canDeleteUploads" :disabled="!currentPageRowKeys.length" @click="selectCurrentPage">选中当前页</a-button>
            <a-button v-if="canDeleteUploads" :disabled="!selectedRowKeys.length" @click="clearSelected">取消选中</a-button>
            <a-button v-if="canDeleteUploads" status="danger" :disabled="!selectedRowKeys.length" @click="batchRemove">批量删除</a-button>
          </a-space>
        </template>
        <a-table
          row-key="id"
          :data="list"
          :pagination="pagination"
          :row-selection="canDeleteUploads ? { type: 'checkbox' } : undefined"
          :selected-keys="selectedRowKeys"
          @page-change="pageChange"
          @selection-change="setSelected"
        >
          <template #columns>
            <a-table-column title="原始文件名">
              <template #cell="{ record }">
                <div class="file-name-cell">
                  <div class="file-name">{{ record.origin_name }}</div>
                  <div class="file-meta">MD5：{{ record.md5 }}</div>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="文件类型">
              <template #cell="{ record }">
                <a-tag :color="typeColor(record.type)">{{ typeLabel(record.type) }}</a-tag>
              </template>
            </a-table-column>
            <a-table-column title="存储类型">
              <template #cell="{ record }">
                <a-tag :color="providerColor(record.provider)">{{ providerLabel(record.provider) }}</a-tag>
              </template>
            </a-table-column>
            <a-table-column title="场景" data-index="scene" />
            <a-table-column title="资源地址">
              <template #cell="{ record }">
                <a-link :href="record.path" target="_blank">{{ record.path }}</a-link>
              </template>
            </a-table-column>
            <a-table-column title="上传时间" data-index="created_at" />
            <a-table-column title="操作">
              <template #cell="{ record }">
                <a-space wrap>
                  <a-button size="mini" @click="copy(record.path)">复制</a-button>
                  <a-button size="mini" @click="open(record.path)">打开</a-button>
                  <a-button v-if="canDeleteUploads" size="mini" status="danger" @click="remove(record.id)">删除</a-button>
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
import { computed, onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { uploadBatchDeleteAPI, uploadDeleteAPI, uploadFileAPI, uploadListAPI, type UploadRecord } from "@/api/upload";
import { useConfirmAction } from "@/hooks/useConfirmAction";
import { usePageSelection } from "@/hooks/usePageSelection";
import { useUploadPreferenceStore } from "@/store/modules/upload-preference";
import { useSessionStore } from "@/store/modules/session";

const accept = ".jpg,.jpeg,.png,.gif,.webp,.mp3,.wav,.aac,.m4a,.mp4,.mov,.avi,.mkv";
const list = ref<UploadRecord[]>([]);
const { confirmBatchDelete, confirmDelete } = useConfirmAction();
const uploadPreference = useUploadPreferenceStore();
const session = useSessionStore();
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const filters = reactive<Record<string, any>>({
  origin_name: "",
  type: undefined,
  provider: undefined,
  scene: ""
});
const {
  selectedKeys: selectedRowKeys,
  currentPageKeys: currentPageRowKeys,
  setSelected,
  selectCurrentPage,
  clearSelected,
  removeSelected
} = usePageSelection(() => list.value, item => item.id);
const canDeleteUploads = computed(() => session.can("/api/v1/uploads#DELETE"));

const load = async () => {
  const res = await uploadListAPI({ ...filters, page: pagination.current, page_size: pagination.pageSize });
  list.value = res.data.list;
  pagination.total = res.data.total;
};

const upload = async (option: any) => {
  const fd = new FormData();
  fd.append("file", option.fileItem.file);
  fd.append("provider", uploadPreference.provider);
  await uploadFileAPI(fd);
  Message.success("上传成功");
  option.onSuccess();
  load();
};

const pageChange = (page: number) => {
  pagination.current = page;
  selectedRowKeys.value = [];
  load();
};

const reset = () => {
  filters.origin_name = "";
  filters.type = undefined;
  filters.provider = undefined;
  filters.scene = "";
  pagination.current = 1;
  selectedRowKeys.value = [];
  load();
};

const copy = async (url: string) => {
  await navigator.clipboard.writeText(url);
  Message.success("地址已复制");
};

const open = (url: string) => {
  window.open(url, "_blank", "noopener,noreferrer");
};

const remove = async (id: number) => {
  await confirmDelete(async () => {
    await uploadDeleteAPI(id);
    Message.success("资源已删除");
    removeSelected(id);
    load();
  }, "这个资源");
};

const batchRemove = async () => {
  await confirmBatchDelete(async () => {
    await uploadBatchDeleteAPI(selectedRowKeys.value);
    Message.success("资源已批量删除");
    selectedRowKeys.value = [];
    load();
  }, "当前页选中的资源");
};

const providerLabel = (provider: string) => {
  switch (provider) {
    case "aliyun-oss":
      return "阿里云";
    case "tencent-cos":
      return "腾讯云";
    case "huawei-obs":
      return "华为云";
    case "local":
      return "本地";
    default:
      return provider || "未标记";
  }
};

const providerColor = (provider: string) => {
  switch (provider) {
    case "aliyun-oss":
      return "orangered";
    case "tencent-cos":
      return "arcoblue";
    case "huawei-obs":
      return "green";
    case "local":
      return "gray";
    default:
      return "gray";
  }
};

const typeLabel = (type: number) => {
  switch (type) {
    case 1:
      return "图片";
    case 2:
      return "视频";
    case 4:
      return "音频";
    default:
      return String(type);
  }
};

const typeColor = (type: number) => {
  switch (type) {
    case 1:
      return "green";
    case 2:
      return "arcoblue";
    case 4:
      return "purple";
    default:
      return "gray";
  }
};

onMounted(load);
</script>

<style scoped lang="scss">
.upload-page {
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

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 16px;
}

.upload-actions {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.upload-default-tip {
  font-size: 12px;
  line-height: 1.6;
  color: var(--color-text-3);
}

.file-name-cell {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.file-name {
  font-weight: 700;
}

.file-meta {
  font-size: 12px;
  color: var(--color-text-3);
}

@media (max-width: 1080px) {
  .page-hero {
    grid-template-columns: 1fr;
  }
}
</style>
