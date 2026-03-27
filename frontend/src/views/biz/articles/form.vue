<template>
  <div class="article-form-page">
    <div class="article-form-hero">
      <div>
        <div class="article-form-badge">{{ isEdit ? "内容维护" : "内容创建" }}</div>
        <h1>{{ isEdit ? "文章详情 / 编辑" : "新增文章" }}</h1>
        <p>采用后台 CMS 双栏布局，左侧专注编辑正文，右侧集中配置发布属性、三级菜单与封面资源。</p>
      </div>
      <div class="article-form-hero-side">
        <div class="hero-side-label">当前类型</div>
        <div class="hero-side-value">{{ articleTypeLabel }}</div>
      </div>
    </div>

    <div class="article-form-layout">
      <div class="article-form-main">
        <a-card class="article-panel" :bordered="false">
          <template #title>基础信息</template>
          <a-form :model="form" layout="vertical">
            <a-row :gutter="16">
              <a-col :xs="24" :lg="14">
                <a-form-item label="标题">
                  <a-input v-model="form.title" placeholder="请输入文章标题" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :lg="10">
                <a-form-item label="摘要">
                  <a-input v-model="form.summary" placeholder="请输入摘要，用于列表概览" />
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-card>

        <a-card class="article-panel" :bordered="false">
          <template #title>正文内容</template>
          <div class="editor-toolbar">
            <a-tag color="arcoblue">{{ articleTypeLabel }}</a-tag>
            <span class="editor-toolbar-tip">
              {{ form.type === 1 ? "纯文本模式下只需要维护富文本正文" : `${assetLabel}模式下可先上传资源，再继续编辑正文说明` }}
            </span>
          </div>

          <a-form :model="form" layout="vertical">
            <a-form-item v-if="form.type === 1" label="文章内容">
              <RichTextEditor ref="editorRef" v-model="form.content" />
            </a-form-item>

            <a-form-item v-else label="文章内容">
              <a-space direction="vertical" style="width: 100%">
                <a-alert>{{ assetLabel }}上传成功后会把访问地址自动插入到编辑器。</a-alert>

                <div v-if="currentAssetUrl" class="asset-preview">
                  <a-space direction="vertical" style="width: 100%">
                    <a-link :href="currentAssetUrl" target="_blank">{{ currentAssetUrl }}</a-link>
                    <video v-if="form.type === 2" :src="currentAssetUrl" controls class="asset-video" />
                    <audio v-else-if="form.type === 4" :src="currentAssetUrl" controls class="asset-audio" />
                    <a-image v-else-if="form.type === 6" :src="currentAssetUrl" width="180" />
                    <a-space wrap>
                      <a-button size="small" @click="openUrl(currentAssetUrl)">打开</a-button>
                      <a-button size="small" @click="copyUrl(currentAssetUrl)">复制地址</a-button>
                      <a-button size="small" @click="downloadUrl(currentAssetUrl)">下载</a-button>
                      <a-button size="small" type="outline">重新上传将替换当前资源</a-button>
                      <a-button size="small" status="danger" @click="clearAsset">清空当前资源</a-button>
                    </a-space>
                  </a-space>
                </div>

                <div class="asset-actions">
                  <a-select v-model="assetProvider" style="width: 220px" placeholder="设置全局默认上传类型">
                    <a-option value="aliyun-oss">上传到阿里云 OSS</a-option>
                    <a-option value="tencent-cos">上传到腾讯云 COS</a-option>
                    <a-option value="huawei-obs">上传到华为云 OBS</a-option>
                    <a-option value="local">上传到本地存储</a-option>
                  </a-select>
                  <a-upload :custom-request="uploadAsset" :show-file-list="false" :accept="uploadAccept">
                    <template #upload-button>
                      <a-button type="outline">{{ uploadButtonText }}</a-button>
                    </template>
                  </a-upload>
                  <a-button type="outline" @click="libraryVisible = true">从资源库选择</a-button>
                </div>

                <RichTextEditor ref="editorRef" v-model="form.content" />
              </a-space>
            </a-form-item>
          </a-form>
        </a-card>

      </div>

      <div class="article-form-side">
        <div class="article-form-side-sticky">
          <a-card class="article-panel" :bordered="false">
            <template #title>发布设置</template>
            <a-form :model="form" layout="vertical">
              <a-row :gutter="12">
                <a-col :span="8">
                  <a-form-item label="文章类型">
                    <a-select v-model="form.type">
                      <a-option :value="1">纯文本</a-option>
                      <a-option :value="2">视频</a-option>
                      <a-option :value="4">音频</a-option>
                      <a-option :value="6">图片</a-option>
                    </a-select>
                  </a-form-item>
                </a-col>
                <a-col :span="8">
                  <a-form-item label="渠道">
                    <a-select v-model="form.channel_id" allow-clear placeholder="请选择渠道">
                      <a-option v-for="item in channels" :key="item.id" :value="item.id">
                        {{ item.name }}（{{ item.code }}）
                      </a-option>
                    </a-select>
                  </a-form-item>
                </a-col>
                <a-col :span="8">
                  <a-form-item label="是否收费">
                    <a-select v-model="form.is_paid">
                      <a-option :value="0">免费</a-option>
                      <a-option :value="1">收费</a-option>
                    </a-select>
                  </a-form-item>
                </a-col>
              </a-row>

              <a-row :gutter="12">
                <a-col :span="12">
                  <a-form-item label="封面类型">
                    <a-input v-model="form.cover_type" />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item label="状态">
                    <a-input-number v-model="form.status" :min="1" :max="4" style="width: 100%" />
                  </a-form-item>
                </a-col>
              </a-row>
            </a-form>
          </a-card>

          <a-card class="article-panel article-menu-panel" :bordered="false">
            <template #title>三级菜单</template>
            <MenuMultiSelector v-model="form.menu_ids" :menus="menuTree" placeholder="点击选择第三级分类" />
            <div class="field-tip">直接读取 `menus` 表三级树结构，点击输入框展开三级分类菜单，只允许勾选第三级叶子节点，最终同步写入 `article_menus` 中间表。</div>
          </a-card>

          <a-card class="article-panel" :bordered="false">
            <template #title>封面资源</template>
            <a-space direction="vertical" style="width: 100%">
              <a-select v-model="assetProvider" placeholder="设置全局默认上传类型">
                <a-option value="aliyun-oss">阿里云</a-option>
                <a-option value="tencent-cos">腾讯云</a-option>
                <a-option value="huawei-obs">华为云</a-option>
                <a-option value="local">本地存储</a-option>
              </a-select>
              <div class="field-tip">这里设置的是全局默认上传类型。当前页面的封面上传、正文资源上传和资源库默认筛选都会跟随这个值。</div>
              <a-alert v-if="showLegacyAudioCoverHint" type="warning">
                当前音频文章是历史单封面数据，系统已保留大图，请补齐中图和小图后再保存。
              </a-alert>
              <div v-for="item in coverFields" :key="item.key" class="cover-field-block">
                <a-input v-model="form[item.key]" :placeholder="item.placeholder" />
                <a-upload :custom-request="uploadCoverByField(item.key)" :show-file-list="false" accept=".jpg,.jpeg,.png,.gif,.webp">
                  <template #upload-button>
                    <a-button type="outline" long>{{ item.label }}</a-button>
                  </template>
                </a-upload>
                <a-space v-if="form[item.key]" wrap>
                  <a-button size="small" @click="openUrl(form[item.key])">打开</a-button>
                  <a-button size="small" @click="copyUrl(form[item.key])">复制地址</a-button>
                  <a-button size="small" @click="downloadUrl(form[item.key])">下载</a-button>
                  <a-button size="small" status="danger" @click="clearCover(item.key)">清空</a-button>
                </a-space>
                <a-image v-if="form[item.key]" :src="form[item.key]" class="cover-preview" />
              </div>
            </a-space>
          </a-card>

          <div class="article-action-bar">
            <a-space direction="vertical" fill style="width: 100%">
              <a-button type="primary" long @click="submit">保存文章</a-button>
              <a-button long @click="router.back()">返回列表</a-button>
            </a-space>
          </div>
        </div>
      </div>
    </div>
  </div>

  <AssetLibraryModal v-model:visible="libraryVisible" :type="libraryType" :scene="uploadScene" :provider="assetProvider" @select="selectFromLibrary" />
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { useRoute, useRouter } from "vue-router";
import { articleCreateAPI, articleDetailAPI, articleUpdateAPI } from "@/api/article";
import { channelListAPI, type ChannelItem } from "@/api/channel";
import { useConfirmAction } from "@/hooks/useConfirmAction";
import { menuCascaderAPI } from "@/api/menu";
import { uploadFileAPI } from "@/api/upload";
import type { UploadRecord } from "@/api/upload";
import AssetLibraryModal from "@/components/article/asset-library-modal.vue";
import MenuMultiSelector from "@/components/article/menu-multi-selector.vue";
import RichTextEditor from "@/components/article/rich-text-editor.vue";
import { useUploadPreferenceStore } from "@/store/modules/upload-preference";

type EditorExpose = {
  insertHTML: (html: string) => Promise<void>;
  replaceHTML: (html: string) => void;
  focus: () => void;
};

const route = useRoute();
const router = useRouter();
const isEdit = computed(() => !!route.params.id);
const { confirmSave } = useConfirmAction();
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
const menuTree = ref<any[]>([]);
const channels = ref<ChannelItem[]>([]);
const editorRef = ref<EditorExpose | null>(null);
const libraryVisible = ref(false);
const uploadPreference = useUploadPreferenceStore();
const assetProvider = computed({
  get: () => uploadPreference.provider,
  set: value => uploadPreference.setProvider(value)
});
const form = reactive<Record<string, any>>({
  title: "",
  summary: "",
  type: 1,
  cover_large: "",
  cover_medium: "",
  cover_small: "",
  cover_type: "1",
  menu_ids: [],
  channel_id: undefined,
  is_paid: 0,
  status: 1,
  content: ""
});

const coverFields = computed(() => {
  if (form.type === 4) {
    return [
      { key: "cover_large", label: "上传大封面图（560x340）", placeholder: "大封面图上传后自动回填 URL" },
      { key: "cover_medium", label: "上传中封面图（300x300）", placeholder: "中封面图上传后自动回填 URL" },
      { key: "cover_small", label: "上传小封面图（104x104）", placeholder: "小封面图上传后自动回填 URL" }
    ];
  }
  return [{ key: "cover_large", label: "上传封面图", placeholder: "封面图上传后自动回填 URL" }];
});

const uploadAccept = computed(() => (form.type === 2 ? ".mp4,.mov,.avi,.mkv" : form.type === 4 ? ".mp3,.wav,.aac,.m4a" : ".jpg,.jpeg,.png,.gif,.webp"));
const assetLabel = computed(() => (form.type === 2 ? "视频资源" : form.type === 4 ? "音频资源" : "图片资源"));
const uploadButtonText = computed(() => (form.type === 2 ? "上传视频" : form.type === 4 ? "上传音频" : "上传图片"));
const uploadScene = computed(() => (form.type === 2 ? "article-video" : form.type === 4 ? "article-audio" : "article-image"));
const libraryType = computed(() => (form.type === 2 ? 2 : form.type === 4 ? 4 : 1));
const currentAssetUrl = computed(() => {
  const html = form.content || "";
  const src = html.match(/src="([^"]+)"/i)?.[1];
  if (src) return src;
  return html.match(/href="([^"]+)"/i)?.[1] || "";
});

const showLegacyAudioCoverHint = computed(
  () => isEdit.value && form.type === 4 && Boolean(form.cover_large) && (!form.cover_medium || !form.cover_small)
);

const uploadCover = async (option: any, field: string) => {
  const fd = new FormData();
  fd.append("file", option.fileItem.file);
  fd.append("scene", "cover");
  fd.append("provider", assetProvider.value);
  const res = await uploadFileAPI(fd);
  form[field] = res.data.url;
  Message.success("封面图上传成功");
  option.onSuccess(res);
};

const uploadCoverByField = (field: string) => (option: any) => uploadCover(option, field);

const uploadAsset = async (option: any) => {
  const fd = new FormData();
  fd.append("file", option.fileItem.file);
  fd.append("scene", uploadScene.value);
  fd.append("provider", assetProvider.value);
  const res = await uploadFileAPI(fd);
  const inserted = buildAssetHTML(res.data.url);
  if (editorRef.value) {
    editorRef.value.replaceHTML(inserted);
  } else {
    form.content = inserted;
  }
  Message.success(`${assetLabel.value}上传成功`);
  option.onSuccess(res);
};

const buildAssetHTML = (url: string) => {
  if (form.type === 2) {
    return `<p><a href="${url}" target="_blank" rel="noreferrer">${url}</a></p><p><video controls style="max-width:100%" src="${url}"></video></p>`;
  }
  if (form.type === 4) {
    return `<p><a href="${url}" target="_blank" rel="noreferrer">${url}</a></p><p><audio controls src="${url}"></audio></p>`;
  }
  return `<p><a href="${url}" target="_blank" rel="noreferrer">${url}</a></p><p><img src="${url}" alt="" style="max-width:100%" /></p>`;
};

const fetchMenus = async () => {
  const res = await menuCascaderAPI();
  menuTree.value = res.data;
};

const fetchChannels = async () => {
  const res = await channelListAPI();
  channels.value = res.data.filter(item => item.status === 1);
};

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

const clearCover = (field: string) => {
  form[field] = "";
};

const clearAsset = () => {
  form.content = "";
  editorRef.value?.replaceHTML("");
};

const selectFromLibrary = (record: UploadRecord) => {
  const inserted = buildAssetHTML(record.path);
  form.content = inserted;
  editorRef.value?.replaceHTML(inserted);
  Message.success("已从资源库选择资源");
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

const loadDetail = async () => {
  if (!isEdit.value) return;
  const res = await articleDetailAPI(route.params.id as string);
  Object.assign(form, res.data, {
    content: res.data.content?.content || "",
    menu_ids: Array.isArray(res.data.menu_ids) && res.data.menu_ids.length ? res.data.menu_ids : (res.data.menu_id ? [res.data.menu_id] : [])
  });
};

const submitAction = async () => {
  const menuIDs = Array.isArray(form.menu_ids) ? [...new Set(form.menu_ids.map(Number).filter(Boolean))] : [];
  const payload: Record<string, any> = { ...form, menu_ids: menuIDs, menu_id: menuIDs[0] };
  if (!payload.cover_large) {
    Message.error("请先上传大封面图");
    return false;
  }
  if (payload.type === 4 && (!payload.cover_medium || !payload.cover_small)) {
    Message.error("音频文章需要上传大、中、小三种封面图");
    return false;
  }
  if (!payload.menu_ids?.length) {
    Message.error("请至少选择一个内容菜单");
    return false;
  }
  if (isEdit.value) {
    await articleUpdateAPI(route.params.id as string, payload);
  } else {
    await articleCreateAPI(payload);
  }
  Message.success("保存成功");
  router.push("/articles");
  return true;
};

const submit = async () => {
  await confirmSave(async () => {
    await submitAction();
  }, isEdit.value ? "这篇文章" : "这篇新文章");
};

onMounted(async () => {
  await Promise.all([fetchMenus(), fetchChannels()]);
  await loadDetail();
});
</script>

<style scoped lang="scss">
.article-form-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 100%;
}

.article-form-hero {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(220px, 0.6fr);
  gap: 18px;
  padding: 24px 28px;
  border-radius: 16px;
  background:
    linear-gradient(135deg, rgba(15, 23, 42, 0.96), rgba(22, 93, 255, 0.9)),
    radial-gradient(circle at top right, rgba(255, 255, 255, 0.16), transparent 38%);
  color: #fff;
  box-shadow: 0 22px 40px rgb(15 23 42 / 14%);
}

.article-form-badge {
  display: inline-flex;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.14);
  font-size: 12px;
}

.article-form-hero h1 {
  margin: 14px 0 10px;
  font-family: AliFangYuanTi, "PingFang SC", sans-serif;
  font-size: 32px;
  line-height: 1.2;
}

.article-form-hero p {
  max-width: 760px;
  line-height: 1.8;
  color: rgba(255, 255, 255, 0.86);
}

.article-form-hero-side {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 18px 20px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(10px);
}

.hero-side-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.74);
}

.hero-side-value {
  margin-top: 8px;
  font-size: 22px;
  font-weight: 700;
}

.article-form-layout {
  display: grid;
  grid-template-columns: minmax(0, 1.22fr) minmax(520px, 1.04fr);
  gap: 16px;
  min-height: 0;
}

.article-form-main,
.article-form-side {
  min-width: 0;
}

.article-form-side-sticky {
  position: sticky;
  top: 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.article-action-bar {
  padding: 14px;
  border: 1px solid var(--color-border-2);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 12px 28px rgb(15 23 42 / 6%);
}

.article-panel {
  border-radius: 16px;
  box-shadow: 0 10px 30px rgb(15 23 42 / 5%);
}

.editor-toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 18px;
}

.editor-toolbar-tip,
.field-tip {
  font-size: 12px;
  color: var(--color-text-3);
}

.article-menu-panel {
  min-height: 0;
}

.asset-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.cover-field-block {
  display: grid;
  gap: 10px;
  width: 100%;
}

.asset-preview {
  padding: 12px;
  border: 1px solid var(--color-border-2);
  border-radius: 12px;
  background: var(--color-fill-1);
}

.asset-video {
  width: 280px;
  max-width: 100%;
  border-radius: 8px;
}

.asset-audio {
  width: 280px;
  max-width: 100%;
}

.cover-preview {
  width: 100%;
  max-width: 260px;
  border-radius: 12px;
}

@media (max-width: 1080px) {
  .article-form-layout,
  .article-form-hero {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .article-action-bar {
    padding: 12px;
  }

  .article-form-hero {
    padding: 20px;
  }

  .article-form-hero h1 {
    font-size: 26px;
  }
}
</style>
