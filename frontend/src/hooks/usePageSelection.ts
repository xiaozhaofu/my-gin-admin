import { computed, ref } from "vue";

export function usePageSelection<T, K extends string | number = number>(
  getRows: () => T[],
  getKey: (row: T) => K
) {
  const selectedKeys = ref<K[]>([]);

  const currentPageKeys = computed(() => getRows().map(getKey));

  const setSelected = (keys: (string | number)[]) => {
    selectedKeys.value = keys.map(key => key as K);
  };

  const selectCurrentPage = () => {
    selectedKeys.value = [...currentPageKeys.value];
  };

  const clearSelected = () => {
    selectedKeys.value = [];
  };

  const removeSelected = (key: K) => {
    selectedKeys.value = selectedKeys.value.filter(item => item !== key);
  };

  return {
    selectedKeys,
    currentPageKeys,
    setSelected,
    selectCurrentPage,
    clearSelected,
    removeSelected
  };
}
