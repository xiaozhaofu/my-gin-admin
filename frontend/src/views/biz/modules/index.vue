<template>
  <a-row :gutter="[16, 16]">
    <a-col :span="8" v-for="item in list" :key="item.slug">
      <a-card :title="item.title" :bordered="false">
        <template #extra>{{ item.slug }}</template>
        <a-typography-paragraph v-for="field in item.fields.slice(0, 8)" :key="field" :ellipsis="{ rows: 1 }">
          {{ field }}
        </a-typography-paragraph>
      </a-card>
    </a-col>
  </a-row>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { moduleMetadataAPI } from "@/api/system";

const list = ref<Array<{ slug: string; title: string; fields: string[] }>>([]);

onMounted(async () => {
  const res = await moduleMetadataAPI();
  list.value = res.data;
});
</script>
