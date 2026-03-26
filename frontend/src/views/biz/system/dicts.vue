<template>
  <a-card :bordered="false">
    <template #title>字典管理</template>
    <a-space direction="vertical" fill>
      <a-button v-if="session.can('/api/v1/dict-types#POST')" type="primary" @click="openType()">新增字典类型</a-button>
      <a-collapse>
        <a-collapse-item v-for="item in list" :key="item.id" :header="`${item.name} (${item.type_code})`">
          <a-space direction="vertical" fill>
            <a-space>
              <a-button v-if="session.can('/api/v1/dict-types/:id#PUT')" size="mini" @click="openType(item)">编辑类型</a-button>
              <a-button v-if="session.can('/api/v1/dict-items#POST')" size="mini" @click="openItem(item.id)">新增字典项</a-button>
              <a-button v-if="session.can('/api/v1/dict-types/:id#DELETE')" size="mini" status="danger" @click="removeType(item.id)">删除类型</a-button>
            </a-space>
            <a-table :data="item.items || []" :pagination="false" row-key="id">
              <template #columns>
                <a-table-column title="标签" data-index="label" />
                <a-table-column title="值" data-index="value" />
                <a-table-column title="排序" data-index="sort" />
                <a-table-column title="状态" data-index="status" />
                <a-table-column title="操作">
                  <template #cell="{ record }">
                    <a-space>
                      <a-button v-if="session.can('/api/v1/dict-items/:id#PUT')" size="mini" @click="openItem(item.id, record)">编辑</a-button>
                      <a-button v-if="session.can('/api/v1/dict-items/:id#DELETE')" size="mini" status="danger" @click="removeItem(record.id)">删除</a-button>
                    </a-space>
                  </template>
                </a-table-column>
              </template>
            </a-table>
          </a-space>
        </a-collapse-item>
      </a-collapse>
    </a-space>
  </a-card>

  <a-modal v-model:visible="typeVisible" :title="typeForm.id ? '编辑字典类型' : '新增字典类型'" :on-before-ok="submitType">
    <a-form :model="typeForm" layout="vertical">
      <a-form-item label="名称"><a-input v-model="typeForm.name" /></a-form-item>
      <a-form-item label="编码"><a-input v-model="typeForm.type_code" /></a-form-item>
      <a-form-item label="状态"><a-input-number v-model="typeForm.status" :min="1" :max="2" /></a-form-item>
      <a-form-item label="备注"><a-input v-model="typeForm.remark" /></a-form-item>
    </a-form>
  </a-modal>

  <a-modal v-model:visible="itemVisible" :title="itemForm.id ? '编辑字典项' : '新增字典项'" :on-before-ok="submitItem">
    <a-form :model="itemForm" layout="vertical">
      <a-form-item label="标签"><a-input v-model="itemForm.label" /></a-form-item>
      <a-form-item label="值"><a-input v-model="itemForm.value" /></a-form-item>
      <a-form-item label="排序"><a-input-number v-model="itemForm.sort" :min="0" /></a-form-item>
      <a-form-item label="状态"><a-input-number v-model="itemForm.status" :min="1" :max="2" /></a-form-item>
      <a-form-item label="备注"><a-input v-model="itemForm.remark" /></a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { dictItemDeleteAPI, dictItemSaveAPI, dictTypeDeleteAPI, dictTypeListAPI, dictTypeSaveAPI, type DictType } from "@/api/system";
import { useConfirmAction } from "@/hooks/useConfirmAction";
import { useSessionStore } from "@/store/modules/session";

const list = ref<DictType[]>([]);
const typeVisible = ref(false);
const itemVisible = ref(false);
const session = useSessionStore();
const { confirmDelete, confirmSave } = useConfirmAction();
const typeForm = reactive<any>({ id: null, name: "", type_code: "", status: 1, remark: "" });
const itemForm = reactive<any>({ id: null, type_id: 0, label: "", value: "", sort: 0, status: 1, remark: "" });

const load = async () => {
  const res = await dictTypeListAPI();
  list.value = res.data;
};

const openType = (record?: any) => {
  Object.assign(typeForm, { id: null, name: "", type_code: "", status: 1, remark: "" }, record || {});
  typeVisible.value = true;
};

const openItem = (typeID: number, record?: any) => {
  Object.assign(itemForm, { id: null, type_id: typeID, label: "", value: "", sort: 0, status: 1, remark: "" }, record || { type_id: typeID });
  itemVisible.value = true;
};

const submitType = async () => {
  return confirmSave(async () => {
    await dictTypeSaveAPI(typeForm.id, typeForm);
    Message.success("字典类型保存成功");
    typeVisible.value = false;
    load();
  }, typeForm.id ? "当前字典类型" : "新字典类型");
};

const submitItem = async () => {
  return confirmSave(async () => {
    await dictItemSaveAPI(itemForm.id, itemForm);
    Message.success("字典项保存成功");
    itemVisible.value = false;
    load();
  }, itemForm.id ? "当前字典项" : "新字典项");
};

const removeType = async (id: number) => {
  await confirmDelete(async () => {
    await dictTypeDeleteAPI(id);
    Message.success("字典类型已删除");
    load();
  }, "这个字典类型");
};

const removeItem = async (id: number) => {
  await confirmDelete(async () => {
    await dictItemDeleteAPI(id);
    Message.success("字典项已删除");
    load();
  }, "这个字典项");
};

onMounted(load);
</script>
