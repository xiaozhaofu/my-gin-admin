<template>
  <div ref="rootRef" class="menu-selector">
    <div class="menu-selector-trigger" @click="toggleOpen">
      <div v-if="selectedTags.length === 0" class="menu-selector-placeholder">{{ placeholder }}</div>
      <div v-else class="menu-selector-tags">
        <div v-for="tag in selectedTags" :key="tag.id" class="menu-selector-tag" @click.stop>
          <span class="menu-selector-tag-label">{{ tag.label }}</span>
          <button class="menu-selector-tag-close" type="button" @click.stop="removeSelected(tag.id)">x</button>
        </div>
      </div>
      <div class="menu-selector-trigger-side">
        <span class="menu-selector-count">已选 {{ selectedTags.length }}</span>
        <span class="menu-selector-arrow">{{ open ? "▲" : "▼" }}</span>
      </div>
    </div>

    <div v-if="open" class="menu-selector-panel">
      <div class="menu-selector-panel-head">
        <div class="menu-selector-panel-title">具体分类</div>
        <div class="menu-selector-panel-sub">优先勾选第三级分类；如果某个二级分类下没有第三级，则该二级分类本身可直接勾选。</div>
      </div>

      <div v-if="menus.length" class="menu-selector-groups">
        <section v-for="root in menus" :key="root.id" class="menu-selector-group">
          <div class="menu-selector-group-title">{{ root.name }}</div>
          <div v-if="root.children?.length" class="menu-selector-subgroups">
            <div v-for="child in root.children" :key="child.id" class="menu-selector-subgroup">
              <div class="menu-selector-subtitle">{{ child.name }}</div>
              <div v-if="child.children?.length" class="menu-selector-leaf-grid">
                <label
                  v-for="leaf in child.children"
                  :key="leaf.id"
                  class="menu-selector-leaf"
                  :class="{ 'is-checked': checkedSet.has(leaf.id) }"
                >
                  <input
                    class="menu-selector-checkbox"
                    type="checkbox"
                    :checked="checkedSet.has(leaf.id)"
                    @change="toggleLeaf(leaf.id)"
                  />
                  <span class="menu-selector-leaf-label">{{ leaf.name }}</span>
                </label>
              </div>
              <label
                v-else
                class="menu-selector-leaf is-secondary-selectable"
                :class="{ 'is-checked': checkedSet.has(child.id) }"
              >
                <input
                  class="menu-selector-checkbox"
                  type="checkbox"
                  :checked="checkedSet.has(child.id)"
                  @change="toggleLeaf(child.id)"
                />
                <span class="menu-selector-leaf-label">{{ child.name }}</span>
                <span class="menu-selector-leaf-note">二级可选</span>
              </label>
            </div>
          </div>
          <div v-else class="menu-selector-empty">当前一级分类下暂无二级分类</div>
        </section>
      </div>

      <a-empty v-else description="当前没有可用菜单，请先到内容菜单页维护" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import type { ContentMenu } from "@/api/menu";

const props = withDefaults(
  defineProps<{
    menus: ContentMenu[];
    modelValue: number[];
    placeholder?: string;
  }>(),
  {
    placeholder: "点击选择第三级分类"
  }
);

const emit = defineEmits<{
  "update:modelValue": [value: number[]];
}>();

const rootRef = ref<HTMLElement | null>(null);
const open = ref(false);

const selectedValues = computed(() => Array.isArray(props.modelValue) ? props.modelValue.map(Number).filter(Boolean) : []);
const checkedSet = computed(() => new Set(selectedValues.value));

const leafMeta = computed(() => {
  const map = new Map<number, { id: number; label: string }>();
  const walk = (items: ContentMenu[], chain: string[] = []) => {
    items.forEach(item => {
      const nextChain = [...chain, item.name];
      if (item.children?.length) {
        walk(item.children, nextChain);
        return;
      }
      map.set(item.id, { id: item.id, label: nextChain.join(" / ") });
    });
  };
  walk(props.menus || []);
  return map;
});

const selectedTags = computed(() =>
  selectedValues.value
    .map(id => leafMeta.value.get(id))
    .filter((item): item is { id: number; label: string } => Boolean(item))
);

const emitValues = (values: number[]) => {
  emit("update:modelValue", [...new Set(values.map(Number).filter(Boolean))]);
};

const toggleLeaf = (id: number) => {
  const next = new Set(selectedValues.value);
  if (next.has(id)) {
    next.delete(id);
  } else {
    next.add(id);
  }
  emitValues([...next]);
};

const removeSelected = (id: number) => {
  emitValues(selectedValues.value.filter(item => item !== id));
};

const toggleOpen = () => {
  open.value = !open.value;
};

const handleDocumentClick = (event: MouseEvent) => {
  const target = event.target as Node | null;
  if (!rootRef.value || !target) {
    return;
  }
  if (!rootRef.value.contains(target)) {
    open.value = false;
  }
};

onMounted(() => {
  document.addEventListener("mousedown", handleDocumentClick);
});

onBeforeUnmount(() => {
  document.removeEventListener("mousedown", handleDocumentClick);
});
</script>

<style scoped lang="scss">
.menu-selector {
  position: relative;
  width: 100%;
}

.menu-selector-trigger {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  min-height: 32px;
  padding: 4px 12px;
  border: 1px solid rgba(22, 93, 255, 0.18);
  border-radius: 14px;
  background: #fff;
  cursor: pointer;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.menu-selector-trigger:hover,
.menu-selector-trigger:focus-within {
  border-color: rgba(22, 93, 255, 0.42);
  box-shadow: 0 0 0 3px rgba(22, 93, 255, 0.08);
}

.menu-selector-placeholder {
  line-height: 24px;
  color: var(--color-text-3);
}

.menu-selector-tags {
  display: flex;
  flex: 1;
  flex-wrap: wrap;
  gap: 8px;
}

.menu-selector-tag {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  max-width: 100%;
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(22, 93, 255, 0.08);
  color: #165dff;
}

.menu-selector-tag-label {
  max-width: 420px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.menu-selector-tag-close {
  width: 18px;
  height: 18px;
  border: 0;
  border-radius: 50%;
  background: rgba(22, 93, 255, 0.14);
  color: #165dff;
  line-height: 18px;
  cursor: pointer;
}

.menu-selector-trigger-side {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 24px;
  white-space: nowrap;
}

.menu-selector-count {
  font-size: 12px;
  color: var(--color-text-3);
}

.menu-selector-arrow {
  font-size: 12px;
  color: var(--color-text-3);
}

.menu-selector-panel {
  position: absolute;
  top: calc(100% + 10px);
  left: 0;
  z-index: 50;
  width: min(900px, 100%);
  max-height: 520px;
  overflow: auto;
  padding: 18px;
  border: 1px solid rgba(22, 93, 255, 0.16);
  border-radius: 18px;
  background: linear-gradient(180deg, #ffffff 0%, #f7faff 100%);
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.12);
}

.menu-selector-panel-head {
  margin-bottom: 14px;
}

.menu-selector-panel-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--color-text-1);
}

.menu-selector-panel-sub {
  margin-top: 4px;
  font-size: 12px;
  color: var(--color-text-3);
}

.menu-selector-groups {
  display: grid;
  gap: 14px;
}

.menu-selector-group {
  padding: 14px;
  border: 1px solid rgba(22, 93, 255, 0.1);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.88);
}

.menu-selector-group-title {
  margin-bottom: 12px;
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text-1);
}

.menu-selector-subgroups {
  display: grid;
  gap: 12px;
}

.menu-selector-subgroup {
  padding: 12px;
  border-radius: 14px;
  background: rgba(246, 248, 252, 0.92);
}

.menu-selector-subtitle {
  margin-bottom: 10px;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-1);
}

.menu-selector-leaf-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 10px;
}

.menu-selector-leaf {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 44px;
  padding: 10px 12px;
  border: 1px solid rgba(22, 93, 255, 0.08);
  border-radius: 12px;
  background: #fff;
  cursor: pointer;
  transition: border-color 0.2s ease, background 0.2s ease;
}

.menu-selector-leaf:hover {
  border-color: rgba(22, 93, 255, 0.24);
  background: rgba(22, 93, 255, 0.04);
}

.menu-selector-leaf.is-checked {
  border-color: rgba(22, 93, 255, 0.32);
  background: rgba(22, 93, 255, 0.08);
}

.menu-selector-leaf.is-secondary-selectable {
  justify-content: space-between;
}

.menu-selector-checkbox {
  width: 16px;
  height: 16px;
  margin: 0;
}

.menu-selector-leaf-label {
  flex: 1;
  line-height: 1.5;
  color: var(--color-text-1);
  word-break: break-word;
}

.menu-selector-leaf-note {
  flex-shrink: 0;
  margin-left: 10px;
  font-size: 12px;
  color: #165dff;
}

.menu-selector-empty {
  font-size: 12px;
  color: var(--color-text-3);
}

@media (max-width: 768px) {
  .menu-selector-panel {
    width: 100%;
    padding: 14px;
  }

  .menu-selector-leaf-grid {
    grid-template-columns: 1fr;
  }

  .menu-selector-tag-label {
    max-width: 220px;
  }
}
</style>
