<template>
  <div class="snow-page">
    <div class="menu-page">
      <section class="page-hero">
        <div>
          <div class="page-badge">Menu Tree</div>
          <h1>内容菜单</h1>
          <p>直接读取 `menus` 数据表，按父子关系展示三级结构。新增或编辑时仅允许选择合法父级，保持层级稳定。</p>
        </div>
        <div class="page-hero-side">
          <div class="page-hero-label">当前菜单总数</div>
          <div class="page-hero-value">{{ flatTree.length }}</div>
        </div>
      </section>

      <a-card class="panel-card" :bordered="false">
        <template #title>菜单树</template>
        <template #extra>
          <a-space wrap>
            <span class="toolbar-tip">读取 `menus` 数据表并按父子层级展示</span>
            <a-button type="primary" @click="openCreate">新增菜单</a-button>
          </a-space>
        </template>
        <a-table :data="flatTree" row-key="id" :pagination="false">
          <template #columns>
            <a-table-column title="菜单名称">
              <template #cell="{ record }">
                <div class="menu-name-cell" :style="{ paddingLeft: `${record._depth * 18}px` }">
                  <span class="menu-level-line" v-if="record._depth > 0"></span>
                  <a-space>
                    <s-svg-icon v-if="record.icon" :name="record.icon" :size="18" />
                    <span>{{ record.name }}</span>
                  </a-space>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="层级" data-index="level" />
            <a-table-column title="父子路径">
              <template #cell="{ record }">{{ record._chainLabel }}</template>
            </a-table-column>
            <a-table-column title="路径" data-index="path" />
            <a-table-column title="页面路径">
              <template #cell="{ record }">
                <div class="page-path-cell">
                  <a-typography-paragraph :ellipsis="{ rows: 1, expandable: true }">
                    {{ record.page_path || "-" }}
                  </a-typography-paragraph>
                  <a-space v-if="record.page_path" size="mini">
                    <a-button size="mini" @click="copyPagePath(record.page_path)">复制</a-button>
                  </a-space>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="图标">
              <template #cell="{ record }">
                <div class="icon-cell">
                  <a-image v-if="isImageIcon(record.icon)" :src="record.icon" width="40" height="40" fit="cover" />
                  <a-space v-else-if="record.icon">
                    <s-svg-icon :name="record.icon" :size="18" />
                    <span>{{ record.icon }}</span>
                  </a-space>
                  <span v-else>-</span>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="状态">
              <template #cell="{ record }">{{ record.is_active ? "启用" : "禁用" }}</template>
            </a-table-column>
            <a-table-column title="操作">
              <template #cell="{ record }">
                <a-space>
                  <a-button size="mini" @click="openEdit(record)">编辑</a-button>
                  <a-button size="mini" status="danger" @click="remove(record.id)">删除</a-button>
                </a-space>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>
  <a-modal v-model:visible="visible" :title="current.id ? '编辑菜单' : '新增菜单'" :on-before-ok="submit">
    <a-form :model="current" layout="vertical">
      <a-form-item label="名称"><a-input v-model="current.name" /></a-form-item>
      <a-form-item label="页面路径"><a-input v-model="current.page_path" placeholder="例如：/home/recommend" /></a-form-item>
      <a-form-item label="图标">
        <a-space direction="vertical" style="width: 100%">
          <a-input v-model="current.icon" placeholder="图标名或图片地址，例如：home / classify / /uploads/icon.png" />
          <a-space wrap>
            <a-upload :custom-request="uploadIcon" :show-file-list="false" accept=".jpg,.jpeg,.png,.gif,.webp">
              <template #upload-button>
                <a-button type="outline">上传图标图片</a-button>
              </template>
            </a-upload>
            <a-button type="outline" @click="iconLibraryVisible = true">从资源库选择</a-button>
          </a-space>
          <div v-if="current.icon" class="icon-preview-box">
            <a-image v-if="isImageIcon(current.icon)" :src="current.icon" width="48" height="48" fit="cover" />
            <a-space v-else>
              <s-svg-icon :name="current.icon" :size="18" />
              <span>{{ current.icon }}</span>
            </a-space>
          </div>
        </a-space>
      </a-form-item>
      <a-form-item label="父级菜单">
        <a-cascader v-model="current.parent_id" :options="parentOptions" allow-clear check-strictly />
      </a-form-item>
        <a-form-item label="排序"><a-input-number v-model="current.sort_order" :min="0" /></a-form-item>
        <a-form-item label="启用"><a-switch v-model="current.is_active" /></a-form-item>
      </a-form>
  </a-modal>
  <AssetLibraryModal v-model:visible="iconLibraryVisible" :type="1" :provider="uploadPreference.provider" @select="selectIconFromLibrary" />
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { menuCreateAPI, menuDeleteAPI, menuTreeAPI, menuUpdateAPI, type ContentMenu } from "@/api/menu";
import { useConfirmAction } from "@/hooks/useConfirmAction";
import { uploadFileAPI } from "@/api/upload";
import type { UploadRecord } from "@/api/upload";
import AssetLibraryModal from "@/components/article/asset-library-modal.vue";
import { useUploadPreferenceStore } from "@/store/modules/upload-preference";

const visible = ref(false);
const iconLibraryVisible = ref(false);
const items = ref<ContentMenu[]>([]);
const { confirmDelete, confirmSave } = useConfirmAction();
const current = reactive<any>({ id: 0, name: "", parent_id: undefined, sort_order: 0, is_active: true, page_path: "", icon: "" });
const uploadPreference = useUploadPreferenceStore();

const fetchTree = async () => {
  const res = await menuTreeAPI();
  items.value = res.data;
};

const findMenuPath = (nodes: ContentMenu[], targetID?: number, parents: number[] = []): number[] | undefined => {
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

const flatTree = computed(() => {
  const result: Array<ContentMenu & { _depth: number; _chainLabel: string }> = [];
  const visit = (nodes: ContentMenu[], parents: string[] = [], depth = 0) => {
    nodes.forEach(node => {
      const chain = [...parents, node.name];
      result.push({
        ...node,
        _depth: depth,
        _chainLabel: chain.join(" / ")
      });
      if (node.children?.length) visit(node.children, chain, depth + 1);
    });
  };
  visit(items.value);
  return result;
});

const options = computed(() => {
  const convert = (nodes: ContentMenu[]): any[] =>
    nodes.map(item => ({
      value: item.id,
      label: item.name,
      children: item.children?.length ? convert(item.children) : undefined
    }));
  return convert(items.value);
});

const descendantIDs = (nodes: ContentMenu[], targetID: number): number[] => {
  const result: number[] = [];
  const walk = (items: ContentMenu[]): boolean => {
    for (const item of items) {
      if (item.id === targetID) {
        collect(item.children || []);
        return true;
      }
      if (walk(item.children || [])) {
        return true;
      }
    }
    return false;
  };
  const collect = (items: ContentMenu[]) => {
    items.forEach(item => {
      result.push(item.id);
      collect(item.children || []);
    });
  };
  walk(nodes);
  return result;
};

const parentOptions = computed(() => {
  const disabledIDs = current.id ? new Set([current.id, ...descendantIDs(items.value, current.id)]) : new Set<number>();
  const convert = (nodes: ContentMenu[]): any[] =>
    nodes
      .filter(item => !disabledIDs.has(item.id))
      .map(item => ({
        value: item.id,
        label: item.name,
        disabled: item.level >= 3,
        children: item.children?.length ? convert(item.children) : undefined
      }));
  return convert(items.value);
});

const openCreate = () => {
  Object.assign(current, { id: 0, name: "", parent_id: undefined, sort_order: 0, is_active: true, page_path: "", icon: "" });
  visible.value = true;
};

const openEdit = (record: ContentMenu) => {
  Object.assign(current, {
    ...record,
    parent_id: record.parent_id ? findMenuPath(items.value, record.parent_id) : undefined
  });
  visible.value = true;
};

const submit = async () => {
  return confirmSave(async () => {
    const payload = { ...current, parent_id: Array.isArray(current.parent_id) ? current.parent_id.at(-1) : current.parent_id };
    if (current.id) {
      await menuUpdateAPI(current.id, payload);
    } else {
      await menuCreateAPI(payload);
    }
    Message.success("保存成功");
    visible.value = false;
    fetchTree();
  }, current.id ? "当前内容菜单" : "新内容菜单");
};

const remove = async (id: number) => {
  await confirmDelete(async () => {
    await menuDeleteAPI(id);
    Message.success("删除成功");
    fetchTree();
  }, "这个内容菜单");
};

const uploadIcon = async (option: any) => {
  const fd = new FormData();
  fd.append("file", option.fileItem.file);
  fd.append("scene", "misc");
  fd.append("provider", uploadPreference.provider);
  const res = await uploadFileAPI(fd);
  current.icon = res.data.url;
  Message.success("图标上传成功");
  option.onSuccess(res);
};

const selectIconFromLibrary = (record: UploadRecord) => {
  current.icon = record.path;
  Message.success("已选择图标资源");
};

const copyPagePath = async (pagePath: string) => {
  await navigator.clipboard.writeText(pagePath);
  Message.success("页面路径已复制");
};

const isImageIcon = (icon?: string) => !!icon && /(\.png|\.jpg|\.jpeg|\.gif|\.webp|^https?:\/\/|^\/uploads\/)/i.test(icon);

onMounted(fetchTree);
</script>

<style scoped lang="scss">
.menu-page {
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

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.toolbar-tip {
  font-size: 12px;
  color: var(--color-text-3);
}

.menu-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 32px;
}

.menu-level-line {
  width: 10px;
  height: 10px;
  border-left: 2px solid rgb(var(--primary-4));
  border-bottom: 2px solid rgb(var(--primary-4));
  opacity: 0.8;
}

.page-path-cell {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.icon-cell {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 40px;
}

.icon-preview-box {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 48px;
}

@media (max-width: 1080px) {
  .page-hero {
    grid-template-columns: 1fr;
  }
}
</style>
