import { computed, ref } from "vue";
import { defineStore } from "pinia";

export type UploadProvider = "aliyun-oss" | "tencent-cos" | "huawei-obs" | "local";

export const useUploadPreferenceStore = defineStore(
  "upload-preference",
  () => {
    const provider = ref<UploadProvider>("tencent-cos");

    const providerLabel = computed(() => {
      switch (provider.value) {
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

    const setProvider = (value: UploadProvider) => {
      provider.value = value;
    };

    return {
      provider,
      providerLabel,
      setProvider
    };
  },
  {
    persist: {
      key: "sleep-admin-upload-provider",
      paths: ["provider"]
    }
  }
);
