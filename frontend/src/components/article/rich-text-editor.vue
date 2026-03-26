<template>
  <div class="editor">
    <div class="toolbar">
      <a-space wrap>
        <a-button size="mini" @click="run('bold')"><b>B</b></a-button>
        <a-button size="mini" @click="run('italic')"><i>I</i></a-button>
        <a-button size="mini" @click="run('underline')"><u>U</u></a-button>
        <a-button size="mini" @click="run('formatBlock', '<h1>')">H1</a-button>
        <a-button size="mini" @click="run('formatBlock', '<h2>')">H2</a-button>
        <a-button size="mini" @click="run('insertUnorderedList')">无序列表</a-button>
        <a-button size="mini" @click="run('insertOrderedList')">有序列表</a-button>
        <a-button size="mini" @click="createLink">链接</a-button>
        <a-button size="mini" @click="insertImage">图片</a-button>
        <a-button size="mini" @click="insertDivider">分割线</a-button>
        <a-button size="mini" @click="insertCodeBlock">代码块</a-button>
        <a-button size="mini" @click="run('removeFormat')">清除格式</a-button>
      </a-space>
    </div>
    <div
      ref="editorRef"
      class="surface"
      contenteditable="true"
      @input="syncValue"
      @blur="syncValue"
    />
  </div>
</template>

<script setup lang="ts">
import { Message } from "@arco-design/web-vue";
import { nextTick, ref, watch } from "vue";

const props = defineProps<{
  modelValue: string;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: string];
}>();

const editorRef = ref<HTMLDivElement | null>(null);

const syncValue = () => {
  emit("update:modelValue", editorRef.value?.innerHTML || "");
};

const run = async (command: string, value?: string) => {
  await nextTick();
  editorRef.value?.focus();
  document.execCommand(command, false, value);
  syncValue();
};

const createLink = async () => {
  const url = window.prompt("请输入链接地址");
  if (!url) {
    Message.info("已取消插入链接");
    return;
  }
  await run("createLink", url);
};

const insertImage = async () => {
  const url = window.prompt("请输入图片地址");
  if (!url) {
    Message.info("已取消插入图片");
    return;
  }
  await insertHTML(`<p><img src="${url}" alt="" style="max-width:100%" /></p>`);
};

const insertDivider = async () => {
  await insertHTML("<hr />");
};

const insertCodeBlock = async () => {
  await insertHTML('<pre><code>// 在这里输入代码</code></pre>');
};

const insertHTML = async (html: string) => {
  await nextTick();
  editorRef.value?.focus();
  document.execCommand("insertHTML", false, html);
  syncValue();
};

const replaceHTML = (html: string) => {
  if (!editorRef.value) {
    return;
  }
  editorRef.value.innerHTML = html;
  syncValue();
};

watch(
  () => props.modelValue,
  value => {
    if (editorRef.value && editorRef.value.innerHTML !== value) {
      editorRef.value.innerHTML = value || "";
    }
  },
  { immediate: true }
);

defineExpose({
  insertHTML,
  replaceHTML,
  focus: () => editorRef.value?.focus()
});
</script>

<style scoped lang="scss">
.editor {
  border: 1px solid var(--color-border-2);
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
}

.toolbar {
  padding: 12px;
  border-bottom: 1px solid var(--color-border-2);
  background: var(--color-fill-1);
}

.surface {
  min-height: 320px;
  padding: 16px;
  outline: none;
  line-height: 1.7;

  :deep(h1) {
    font-size: 30px;
    margin: 12px 0;
  }

  :deep(h2) {
    font-size: 24px;
    margin: 10px 0;
  }

  :deep(pre) {
    padding: 12px;
    border-radius: 10px;
    overflow: auto;
    background: #0f172a;
    color: #e2e8f0;
  }
}
</style>
