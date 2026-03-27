<template>
  <div class="snow-page">
    <div class="batch-page">
      <section class="batch-hero">
        <div>
          <div class="batch-badge">Content Workspace</div>
          <h1>批量新增文章</h1>
          <p>按“公共配置 + 多条文章明细”一次性创建内容，适合同一菜单、同一渠道下的成组发布。</p>
        </div>
        <div class="batch-hero-side">
          <div class="batch-hero-label">当前待提交</div>
          <div class="batch-hero-value">{{ rows.length }}</div>
        </div>
      </section>

      <a-card class="panel-card" :bordered="false">
        <template #title>公共配置</template>
        <a-form :model="form" layout="vertical">
          <a-row :gutter="16">
            <a-col :xs="24" :lg="6">
              <a-form-item label="文章类型">
                <a-select v-model="form.type">
                  <a-option :value="1">纯文本</a-option>
                  <a-option :value="2">视频</a-option>
                  <a-option :value="4">音频</a-option>
                  <a-option :value="6">图片</a-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :xs="24" :lg="6">
              <a-form-item label="是否收费">
                <a-select v-model="form.is_paid">
                  <a-option :value="0">免费</a-option>
                  <a-option :value="1">收费</a-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :xs="24" :lg="12">
              <a-form-item label="三级菜单">
                <MenuMultiSelector v-model="form.menu_ids" :menus="menuTree" placeholder="点击选择第三级分类" />
                <div class="field-tip">直接展示三级菜单，只允许勾选第三级叶子节点，可多选，最终同步写入 `article_menus` 中间表。</div>
              </a-form-item>
            </a-col>
          </a-row>

          <a-row :gutter="16">
            <a-col :xs="24" :lg="6">
              <a-form-item label="渠道">
                <a-select v-model="form.channel_id" allow-clear placeholder="请选择渠道">
                  <a-option v-for="item in channels" :key="item.id" :value="item.id">
                    {{ item.name }}（{{ item.code }}）
                  </a-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :xs="24" :lg="6">
              <a-form-item label="状态">
                <a-input-number v-model="form.status" :min="1" :max="4" style="width: 100%" />
              </a-form-item>
            </a-col>
            <a-col :xs="24" :lg="6">
              <a-form-item label="排序">
                <a-input-number v-model="form.sort_order" :min="0" style="width: 100%" />
              </a-form-item>
            </a-col>
            <a-col :xs="24" :lg="6">
              <a-form-item label="封面类型">
                <a-input v-model="form.cover_type" />
              </a-form-item>
            </a-col>
          </a-row>

          <a-row :gutter="16">
            <a-col :xs="24" :lg="12">
              <a-form-item label="默认存储类型">
                <a-select v-model="assetProvider" placeholder="选择默认上传到哪个存储" style="width: 160px">
                  <a-option value="aliyun-oss">阿里云 OSS</a-option>
                  <a-option value="tencent-cos">腾讯云 COS</a-option>
                  <a-option value="huawei-obs">华为云 OBS</a-option>
                  <a-option value="local">本地存储</a-option>
                </a-select>
                <div class="field-tip">公共封面、每条封面覆盖和每条正文资源上传都会默认使用这里的存储类型。</div>
              </a-form-item>
            </a-col>
          </a-row>

          <a-row :gutter="16">
            <a-col :span="24">
              <a-form-item label="公共封面图">
                <div class="cover-fields-grid">
                  <div v-for="item in coverFields" :key="item.key" class="cover-field-block">
                    <a-input v-model="form[item.key]" :placeholder="item.placeholder" />
                    <a-upload :custom-request="uploadCoverByField(item.key)" :show-file-list="false" accept=".jpg,.jpeg,.png,.gif,.webp">
                      <template #upload-button>
                        <a-button type="outline" long>{{ item.label }}</a-button>
                      </template>
                    </a-upload>
                    <a-space v-if="form[item.key]" wrap class="row-cover-actions">
                      <a-button size="small" @click="openUrl(form[item.key])">打开</a-button>
                      <a-button size="small" @click="copyUrl(form[item.key])">复制地址</a-button>
                      <a-button size="small" @click="downloadUrl(form[item.key])">下载</a-button>
                      <a-button size="small" status="danger" @click="clearCoverField(form, item.key)">清空</a-button>
                    </a-space>
                    <a-image v-if="form[item.key]" :src="form[item.key]" width="140" />
                  </div>
                </div>
              </a-form-item>
            </a-col>
          </a-row>

          <div class="option-row">
            <a-checkbox :model-value="form.is_top === 1" @change="form.is_top = $event ? 1 : 0">置顶</a-checkbox>
            <a-checkbox :model-value="form.is_hot === 1" @change="form.is_hot = $event ? 1 : 0">热门</a-checkbox>
            <a-checkbox :model-value="form.is_recommend === 1" @change="form.is_recommend = $event ? 1 : 0">推荐</a-checkbox>
          </div>

        </a-form>
      </a-card>

      <a-card class="panel-card" :bordered="false">
        <template #title>文本导入</template>
        <template #extra>
          <a-space>
            <a-button type="primary" @click="importRows">导入到列表</a-button>
            <a-button @click="appendRow()">新增一行</a-button>
            <a-button status="danger" @click="clearRows">清空列表</a-button>
          </a-space>
        </template>

        <a-alert class="import-tip" type="info">
          {{ importTip }}
        </a-alert>
        <a-textarea
          v-model="draftText"
          :auto-size="{ minRows: 5, maxRows: 10 }"
          :placeholder="importPlaceholder"
        />
      </a-card>

      <a-card class="panel-card" :bordered="false">
        <template #title>文章明细</template>
        <template #extra>
          <a-space>
            <a-tag color="arcoblue">{{ articleTypeLabel }}</a-tag>
            <a-tag color="green">默认存储：{{ providerLabel }}</a-tag>
          </a-space>
        </template>

        <div v-if="rows.length === 0" class="empty-rows">
          <a-empty description="当前还没有待提交的文章明细" />
        </div>

        <div v-else class="row-list">
          <div v-for="(item, index) in rows" :key="item.key" class="row-card">
            <div class="row-card-head">
              <div class="row-card-title">第 {{ index + 1 }} 篇</div>
              <a-space>
                <a-button size="mini" @click="duplicateRow(index)">复制</a-button>
                <a-button size="mini" status="danger" @click="removeRow(index)">删除</a-button>
              </a-space>
            </div>

            <a-row :gutter="16">
              <a-col :xs="24" :lg="10">
                <a-form-item label="标题">
                  <a-input v-model="item.title" placeholder="请输入标题" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :lg="14">
                <a-form-item label="摘要">
                  <a-input v-model="item.summary" placeholder="用于列表概览，可留空" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-form-item v-if="form.type === 1" label="正文">
              <RichTextEditor v-model="item.content" />
            </a-form-item>

            <a-form-item v-else :label="bodyUploadLabel">
              <div class="row-resource-panel">
                <a-alert>{{ bodyUploadTip }}</a-alert>
                <div v-if="extractAssetUrl(item.content)" class="row-upload-preview">
                  <a-link :href="extractAssetUrl(item.content)" target="_blank">{{ extractAssetUrl(item.content) }}</a-link>
                  <video v-if="form.type === 2" :src="extractAssetUrl(item.content)" controls class="row-preview-video" />
                  <audio v-else-if="form.type === 4" :src="extractAssetUrl(item.content)" controls class="row-preview-audio" />
                  <a-image v-else-if="form.type === 6" :src="extractAssetUrl(item.content)" width="180" />
                </div>
                <a-space wrap class="row-upload-actions">
                  <a-upload :custom-request="uploadRowAssetByRow(item)" :show-file-list="false" :accept="uploadAccept">
                    <template #upload-button>
                      <a-button type="outline">{{ rowUploadButtonText }}</a-button>
                    </template>
                  </a-upload>
                  <a-button v-if="extractAssetUrl(item.content)" size="small" @click="openUrl(extractAssetUrl(item.content))">打开</a-button>
                  <a-button v-if="extractAssetUrl(item.content)" size="small" @click="copyUrl(extractAssetUrl(item.content))">复制地址</a-button>
                  <a-button v-if="extractAssetUrl(item.content)" size="small" @click="downloadUrl(extractAssetUrl(item.content))">下载</a-button>
                  <a-button v-if="extractAssetUrl(item.content)" size="small" status="danger" @click="clearRowAsset(item)">清空资源</a-button>
                </a-space>
              </div>
            </a-form-item>

            <a-row :gutter="16">
              <a-col :span="24">
                <a-form-item label="封面覆盖">
                  <div class="cover-fields-grid">
                    <div v-for="coverField in coverFields" :key="coverField.key" class="cover-field-block">
                      <a-input v-model="item[coverField.key]" :placeholder="coverField.overridePlaceholder" />
                      <a-upload :custom-request="uploadRowCoverByField(item, coverField.key)" :show-file-list="false" accept=".jpg,.jpeg,.png,.gif,.webp">
                        <template #upload-button>
                          <a-button type="outline" long>{{ coverField.overrideLabel }}</a-button>
                        </template>
                      </a-upload>
                      <div class="field-tip">留空则继承公共封面图。</div>
                      <a-space v-if="effectiveCover(item, coverField.key)" wrap class="row-cover-actions">
                        <a-button size="small" @click="openUrl(effectiveCover(item, coverField.key))">打开</a-button>
                        <a-button size="small" @click="copyUrl(effectiveCover(item, coverField.key))">复制地址</a-button>
                        <a-button size="small" @click="downloadUrl(effectiveCover(item, coverField.key))">下载</a-button>
                        <a-button v-if="item[coverField.key]" size="small" status="danger" @click="clearCoverField(item, coverField.key)">清空覆盖</a-button>
                      </a-space>
                      <a-image v-if="effectiveCover(item, coverField.key)" :src="effectiveCover(item, coverField.key)" width="140" />
                    </div>
                  </div>
                </a-form-item>
              </a-col>
            </a-row>
          </div>
        </div>
      </a-card>

      <div class="action-bar">
        <a-space>
          <a-button type="primary" :loading="submitting" @click="submit">批量保存</a-button>
          <a-button @click="router.back()">返回列表</a-button>
        </a-space>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { useRouter } from "vue-router";
import { articleBatchCreateAPI, type ArticleBatchCreateItem } from "@/api/article";
import { channelListAPI, type ChannelItem } from "@/api/channel";
import { useConfirmAction } from "@/hooks/useConfirmAction";
import { menuCascaderAPI } from "@/api/menu";
import { uploadFileAPI } from "@/api/upload";
import MenuMultiSelector from "@/components/article/menu-multi-selector.vue";
import RichTextEditor from "@/components/article/rich-text-editor.vue";
import { useUploadPreferenceStore } from "@/store/modules/upload-preference";

type BatchRow = ArticleBatchCreateItem & {
  key: number;
};

type CoverFieldKey = "cover_large" | "cover_medium" | "cover_small";

type CoverFieldConfig = {
  key: CoverFieldKey;
  label: string;
  placeholder: string;
  overrideLabel: string;
  overridePlaceholder: string;
};

const router = useRouter();
const { confirmDelete, confirmSave } = useConfirmAction();
const channels = ref<ChannelItem[]>([]);
const menuTree = ref<any[]>([]);
const draftText = ref("");
const uploadPreference = useUploadPreferenceStore();
const assetProvider = computed({
  get: () => uploadPreference.provider,
  set: value => uploadPreference.setProvider(value)
});
const submitting = ref(false);
const rowSeed = ref(1);
const rows = ref<BatchRow[]>([]);
const form = reactive<Record<string, any>>({
  type: 1,
  menu_ids: [],
  channel_id: undefined,
  status: 1,
  sort_order: 0,
  cover_type: "1",
  cover_large: "",
  cover_medium: "",
  cover_small: "",
  is_paid: 0,
  is_top: 0,
  is_hot: 0,
  is_recommend: 0
});

const articleTypeLabel = computed(() => {
  switch (form.type) {
    case 2:
      return "视频文章";
    case 4:
      return "音频文章";
    case 6:
      return "图片文章";
    default:
      return "纯文本文章";
  }
});

const providerLabel = computed(() => {
  switch (assetProvider.value) {
    case "aliyun-oss":
      return "阿里云 OSS";
    case "tencent-cos":
      return "腾讯云 COS";
    case "huawei-obs":
      return "华为云 OBS";
    default:
      return "本地存储";
  }
});

const uploadAccept = computed(() => (form.type === 2 ? ".mp4,.mov,.avi,.mkv" : form.type === 4 ? ".mp3,.wav,.aac,.m4a" : ".jpg,.jpeg,.png,.gif,.webp"));
const rowUploadButtonText = computed(() => (form.type === 2 ? "上传视频资源" : form.type === 4 ? "上传音频资源" : "上传图片资源"));
const bodyUploadLabel = computed(() => (form.type === 2 ? "视频正文" : form.type === 4 ? "音频正文" : "图片正文"));
const bodyUploadTip = computed(() => `${providerLabel.value}中的文件上传成功后，会自动写入当前这条文章的正文。`);
const uploadScene = computed(() => (form.type === 2 ? "article-video" : form.type === 4 ? "article-audio" : "article-image"));
const importTip = computed(() =>
  form.type === 1
    ? "每行一篇，支持 `标题|摘要|正文` 或 `标题[TAB]摘要[TAB]正文`。如果只有两段，会按 `标题|正文` 解析。"
    : "非纯文本类型下，文本导入只负责标题和摘要；正文资源请在每条记录里单独上传。支持 `标题|摘要` 或 `标题[TAB]摘要`。"
);
const importPlaceholder = computed(() =>
  form.type === 1
    ? "示例：\n春夜助眠|适合晚间放松的短文|正文内容...\n晨间呼吸练习\t3分钟带你进入专注状态\t正文内容..."
    : "示例：\n春夜助眠视频|适合晚间放松的短视频\n深睡白噪音合集\t夜间持续播放"
);

const coverFields = computed<CoverFieldConfig[]>(() => {
  if (form.type === 4) {
    return [
      {
        key: "cover_large",
        label: "上传大封面图（560x340）",
        placeholder: "大封面图上传后自动回填 URL",
        overrideLabel: "上传本条大封面图（560x340）",
        overridePlaceholder: "留空则继承公共大封面图"
      },
      {
        key: "cover_medium",
        label: "上传中封面图（300x300）",
        placeholder: "中封面图上传后自动回填 URL",
        overrideLabel: "上传本条中封面图（300x300）",
        overridePlaceholder: "留空则继承公共中封面图"
      },
      {
        key: "cover_small",
        label: "上传小封面图（104x104）",
        placeholder: "小封面图上传后自动回填 URL",
        overrideLabel: "上传本条小封面图（104x104）",
        overridePlaceholder: "留空则继承公共小封面图"
      }
    ];
  }

  return [
    {
      key: "cover_large",
      label: "上传封面图",
      placeholder: "封面图上传后自动回填 URL",
      overrideLabel: "上传本条封面图",
      overridePlaceholder: "留空则继承公共封面图"
    }
  ];
});

const createRow = (preset?: Partial<BatchRow>): BatchRow => ({
  title: "",
  summary: "",
  content: "",
  cover_large: "",
  cover_medium: "",
  cover_small: "",
  cover_type: "",
  ...preset,
  key: rowSeed.value++
});

const appendRow = (preset?: Partial<BatchRow>) => {
  rows.value.push(createRow(preset));
};

const clearRows = () => {
  rows.value = [];
};

const removeRow = (index: number) => {
  confirmDelete(async () => {
    rows.value.splice(index, 1);
  }, `第 ${index + 1} 条文章草稿`);
};

const duplicateRow = (index: number) => {
  appendRow({ ...rows.value[index] });
};

const importRows = () => {
  const lines = draftText.value
    .split(/\r?\n/)
    .map(line => line.trim())
    .filter(Boolean);

  if (lines.length === 0) {
    Message.warning("请先输入待导入内容");
    return;
  }

  const imported: BatchRow[] = [];
  let invalidCount = 0;
  for (const line of lines) {
    const separator = line.includes("|") ? "|" : "\t";
    const parts = line.split(separator).map(item => item.trim());
    const title = parts[0] || "";
    if (form.type === 1) {
      const summary = parts.length > 2 ? parts[1] || "" : "";
      const content = (parts.length > 2 ? parts.slice(2) : parts.slice(1)).join(separator).trim();
      if (!title || !content) {
        invalidCount += 1;
        continue;
      }
      imported.push(createRow({ title, summary, content }));
      continue;
    }

    const summary = parts[1] || "";
    if (!title) {
      invalidCount += 1;
      continue;
    }
    imported.push(createRow({ title, summary, content: "" }));
  }

  if (imported.length === 0) {
    Message.warning("没有成功解析出可用文章数据");
    return;
  }

  rows.value = [...rows.value, ...imported];
  draftText.value = "";
  Message.success(`已导入 ${imported.length} 篇文章${invalidCount ? `，忽略 ${invalidCount} 行格式错误` : ""}`);
};

const uploadCover = async (option: any) => {
  const fd = new FormData();
  fd.append("file", option.fileItem.file);
  fd.append("scene", "cover");
  fd.append("provider", assetProvider.value);
  const res = await uploadFileAPI(fd);
  form.cover_large = res.data.url;
  Message.success("公共封面上传成功");
  option.onSuccess(res);
};

const uploadCoverByField = (field: CoverFieldKey) => async (option: any) => {
  const fd = new FormData();
  fd.append("file", option.fileItem.file);
  fd.append("scene", "cover");
  fd.append("provider", assetProvider.value);
  const res = await uploadFileAPI(fd);
  form[field] = res.data.url;
  Message.success("公共封面上传成功");
  option.onSuccess(res);
};

const uploadRowCover = async (option: any, row: BatchRow) => {
  const fd = new FormData();
  fd.append("file", option.fileItem.file);
  fd.append("scene", "cover");
  fd.append("provider", assetProvider.value);
  const res = await uploadFileAPI(fd);
  row.cover_large = res.data.url;
  Message.success("本条封面上传成功");
  option.onSuccess(res);
};
const uploadRowCoverByRow = (row: BatchRow) => (option: any) => uploadRowCover(option, row);
const uploadRowCoverByField = (row: BatchRow, field: CoverFieldKey) => async (option: any) => {
  const fd = new FormData();
  fd.append("file", option.fileItem.file);
  fd.append("scene", "cover");
  fd.append("provider", assetProvider.value);
  const res = await uploadFileAPI(fd);
  row[field] = res.data.url;
  Message.success("本条封面上传成功");
  option.onSuccess(res);
};

const uploadRowAsset = async (option: any, row: BatchRow) => {
  const fd = new FormData();
  fd.append("file", option.fileItem.file);
  fd.append("scene", uploadScene.value);
  fd.append("provider", assetProvider.value);
  const res = await uploadFileAPI(fd);
  row.content = buildAssetHTML(res.data.url);
  Message.success(`${bodyUploadLabel.value}上传成功`);
  option.onSuccess(res);
};
const uploadRowAssetByRow = (row: BatchRow) => (option: any) => uploadRowAsset(option, row);

const buildAssetHTML = (url: string) => {
  if (form.type === 2) {
    return `<p><a href="${url}" target="_blank" rel="noreferrer">${url}</a></p><p><video controls style="max-width:100%" src="${url}"></video></p>`;
  }
  if (form.type === 4) {
    return `<p><a href="${url}" target="_blank" rel="noreferrer">${url}</a></p><p><audio controls src="${url}"></audio></p>`;
  }
  return `<p><a href="${url}" target="_blank" rel="noreferrer">${url}</a></p><p><img src="${url}" alt="" style="max-width:100%" /></p>`;
};

const extractAssetUrl = (html: string) => {
  const src = html.match(/src="([^"]+)"/i)?.[1];
  if (src) return src;
  return html.match(/href="([^"]+)"/i)?.[1] || "";
};

const clearRowAsset = (row: BatchRow) => {
  row.content = "";
};

const effectiveCover = (row: BatchRow, field: CoverFieldKey) => row[field] || form[field] || "";

const clearCoverField = (target: Record<string, any>, field: CoverFieldKey) => {
  target[field] = "";
};

const openUrl = (url: string) => {
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

const fetchMenus = async () => {
  const res = await menuCascaderAPI();
  menuTree.value = res.data;
};

const fetchChannels = async () => {
  const res = await channelListAPI();
  channels.value = res.data;
};

const submitAction = async () => {
  const menuIDs = Array.isArray(form.menu_ids) ? [...new Set(form.menu_ids.map(Number).filter(Boolean))] : [];
  if (!menuIDs.length) {
    Message.warning("请至少选择一个第三级菜单");
    return;
  }
  if (!form.channel_id) {
    Message.warning("请选择渠道");
    return;
  }
  if (rows.value.length === 0) {
    Message.warning("请至少维护一条文章明细");
    return;
  }
  const titleMissingIndex = rows.value.findIndex(item => !item.title.trim());
  if (titleMissingIndex >= 0) {
    Message.warning(`第 ${titleMissingIndex + 1} 条文章标题不能为空`);
    return;
  }
  const contentMissingIndex = rows.value.findIndex(item =>
    form.type === 1 ? !item.content.trim() : !extractAssetUrl(item.content)
  );
  if (contentMissingIndex >= 0) {
    Message.warning(form.type === 1 ? `第 ${contentMissingIndex + 1} 条文章正文不能为空` : `第 ${contentMissingIndex + 1} 条文章还没有上传正文资源`);
    return;
  }
  if (!form.cover_large && rows.value.some(item => !item.cover_large)) {
    Message.warning("请配置公共大封面图，或为每条文章单独填写封面图");
    return;
  }
  if (form.type === 4 && (
    (!form.cover_medium && rows.value.some(item => !item.cover_medium)) ||
    (!form.cover_small && rows.value.some(item => !item.cover_small))
  )) {
    Message.warning("音频文章需要配置大、中、小三种封面图");
    return;
  }

  submitting.value = true;
  try {
    const payload = {
      type: form.type,
      cover: form.cover_large,
      cover_large: form.cover_large,
      cover_medium: form.cover_medium,
      cover_small: form.cover_small,
      cover_type: form.cover_type,
      menu_id: menuIDs[0],
      menu_ids: menuIDs,
      channel_id: Number(form.channel_id),
      sort_order: Number(form.sort_order || 0),
      is_paid: form.is_paid,
      is_top: form.is_top,
      is_hot: form.is_hot,
      is_recommend: form.is_recommend,
      status: Number(form.status || 1),
      items: rows.value.map(item => ({
        title: item.title.trim(),
        summary: item.summary.trim(),
        content: item.content.trim(),
        cover: item.cover_large?.trim() || "",
        cover_large: item.cover_large?.trim() || "",
        cover_medium: item.cover_medium?.trim() || "",
        cover_small: item.cover_small?.trim() || "",
        cover_type: item.cover_type?.trim() || ""
      }))
    };
    const res = await articleBatchCreateAPI(payload);
    Message.success(`已成功新增 ${res.data.created} 篇文章`);
    router.push("/articles");
  } finally {
    submitting.value = false;
  }
};

const submit = async () => {
  await confirmSave(async () => {
    await submitAction();
  }, `这 ${rows.value.length} 篇文章`);
};

onMounted(async () => {
  appendRow();
  await Promise.all([fetchMenus(), fetchChannels()]);
});
</script>

<style scoped lang="scss">
.batch-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.batch-hero {
  display: grid;
  grid-template-columns: minmax(0, 1.7fr) minmax(220px, 0.55fr);
  gap: 18px;
  padding: 24px 28px;
  border-radius: 16px;
  background:
    linear-gradient(135deg, rgba(22, 93, 255, 0.96), rgba(64, 128, 255, 0.86)),
    radial-gradient(circle at top right, rgba(255, 255, 255, 0.16), transparent 38%);
  color: #fff;
  box-shadow: 0 20px 38px rgb(22 93 255 / 16%);
}

.batch-badge {
  display: inline-flex;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.14);
  font-size: 12px;
}

.batch-hero h1 {
  margin: 14px 0 10px;
  font-family: AliFangYuanTi, "PingFang SC", sans-serif;
  font-size: 32px;
  line-height: 1.2;
}

.batch-hero p {
  max-width: 760px;
  line-height: 1.8;
  color: rgba(255, 255, 255, 0.86);
}

.batch-hero-side {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 18px 20px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(10px);
}

.batch-hero-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.74);
}

.batch-hero-value {
  margin-top: 8px;
  font-size: 28px;
  font-weight: 700;
}

.panel-card {
  border-radius: 16px;
  box-shadow: 0 12px 30px rgb(15 23 42 / 5%);
}

.option-row {
  display: flex;
  flex-wrap: wrap;
  gap: 18px;
}

.field-tip {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.6;
  color: var(--color-text-3);
}

.cover-fields-grid {
  display: grid;
  gap: 14px;
}

.cover-field-block {
  display: grid;
  gap: 10px;
  width: 100%;
}

.cover-preview-block {
  margin-top: 12px;
}

.import-tip {
  margin-bottom: 12px;
}

.empty-rows {
  padding: 24px 0;
}

.row-list {
  display: grid;
  gap: 16px;
}

.row-card {
  padding: 18px;
  border: 1px solid rgba(22, 93, 255, 0.08);
  border-radius: 16px;
  background: linear-gradient(180deg, #fff 0%, #f8fbff 100%);
}

.row-card :deep(.surface) {
  min-height: 220px;
}

.row-card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.row-card-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--color-text-1);
}

.row-resource-panel {
  display: grid;
  gap: 12px;
}

.row-upload-preview {
  display: grid;
  gap: 10px;
}

.row-preview-video {
  width: 220px;
  max-width: 100%;
  border-radius: 10px;
}

.row-preview-audio {
  width: 260px;
  max-width: 100%;
}

.row-upload-actions,
.row-cover-actions {
  margin-top: 4px;
}

.action-bar {
  position: sticky;
  bottom: 0;
  z-index: 10;
  display: flex;
  justify-content: flex-end;
  padding: 14px 18px;
  border: 1px solid rgba(22, 93, 255, 0.08);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(10px);
  box-shadow: 0 -6px 20px rgb(15 23 42 / 5%);
}

@media (max-width: 1080px) {
  .batch-hero {
    grid-template-columns: 1fr;
  }
}
</style>
