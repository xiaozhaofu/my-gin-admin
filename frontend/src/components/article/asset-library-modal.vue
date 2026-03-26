<template>
  <a-modal v-model:visible="innerVisible" title="资源库" width="960px" @cancel="close" :footer="false">
    <a-form :model="filters" layout="inline" class="filters">
      <a-form-item label="名称"><a-input v-model="filters.origin_name" allow-clear /></a-form-item>
      <a-form-item label="类型">
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
        <a-select v-model="filters.scene" allow-clear style="width: 180px">
          <a-option value="misc">misc</a-option>
          <a-option value="cover">cover</a-option>
          <a-option value="article-video">article-video</a-option>
          <a-option value="article-audio">article-audio</a-option>
          <a-option value="article-image">article-image</a-option>
        </a-select>
      </a-form-item>
      <a-space>
        <a-button type="primary" @click="load">搜索</a-button>
        <a-button @click="reset">重置</a-button>
      </a-space>
    </a-form>

    <a-table :data="list" :pagination="pagination" row-key="id" @page-change="pageChange">
      <template #columns>
        <a-table-column title="名称" data-index="origin_name" />
        <a-table-column title="存储类型">
          <template #cell="{ record }">
            <a-tag :color="providerColor(record.provider)">{{ providerLabel(record.provider) }}</a-tag>
          </template>
        </a-table-column>
        <a-table-column title="场景" data-index="scene" />
        <a-table-column title="预览">
          <template #cell="{ record }">
            <a-image v-if="record.type === 1" :src="record.path" width="72" />
            <video v-else-if="record.type === 2" :src="record.path" controls class="mini-video" />
            <audio v-else-if="record.type === 4" :src="record.path" controls class="mini-audio" />
          </template>
        </a-table-column>
        <a-table-column title="地址">
          <template #cell="{ record }">
            <a-typography-paragraph :ellipsis="{ rows: 1, expandable: true }">{{ record.path }}</a-typography-paragraph>
          </template>
        </a-table-column>
        <a-table-column title="操作">
          <template #cell="{ record }">
            <a-space>
              <a-button size="mini" @click="select(record)">选择</a-button>
              <a-button size="mini" @click="copy(record.path)">复制</a-button>
              <a-button size="mini" status="danger" @click="remove(record.id)">删除</a-button>
            </a-space>
          </template>
        </a-table-column>
      </template>
    </a-table>
  </a-modal>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue";
import { Message } from "@arco-design/web-vue";
import { uploadDeleteAPI, uploadListAPI, type UploadRecord } from "@/api/upload";

const props = defineProps<{ visible: boolean; type?: number; scene?: string; provider?: string }>();
const emit = defineEmits<{ "update:visible": [value: boolean]; select: [record: UploadRecord] }>();

const innerVisible = computed({
  get: () => props.visible,
  set: value => emit("update:visible", value)
});

const list = ref<UploadRecord[]>([]);
const pagination = reactive({ current: 1, pageSize: 10, total: 0 });
const filters = reactive<Record<string, any>>({ origin_name: "", type: undefined, provider: undefined, scene: "" });

watch(
  () => props.visible,
  visible => {
    if (visible) {
      pagination.current = 1;
      filters.type = props.type;
      filters.provider = props.provider;
      // 资源库默认不按 scene 过滤，确保在“资源上传”页上传到 misc 的历史资源也能被选出来。
      filters.scene = "";
      load();
    }
  }
);

const load = async () => {
  const res = await uploadListAPI({ ...filters, page: pagination.current, page_size: pagination.pageSize });
  list.value = res.data.list;
  pagination.total = res.data.total;
};

const pageChange = (page: number) => {
  pagination.current = page;
  load();
};

const reset = () => {
  filters.origin_name = "";
  filters.type = props.type;
  filters.provider = props.provider;
  filters.scene = "";
  pagination.current = 1;
  load();
};

const close = () => emit("update:visible", false);

const select = (record: UploadRecord) => {
  emit("select", record);
  close();
};

const copy = async (url: string) => {
  await navigator.clipboard.writeText(url);
  Message.success("地址已复制");
};

const remove = async (id: number) => {
  await uploadDeleteAPI(id);
  Message.success("资源已删除");
  load();
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
</script>

<style scoped>
.filters {
  margin-bottom: 16px;
}

.mini-video {
  width: 96px;
  max-height: 64px;
}

.mini-audio {
  width: 120px;
}
</style>
