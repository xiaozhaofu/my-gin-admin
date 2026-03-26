import axios from "axios";
import { Message } from "@arco-design/web-vue";
import { useSessionStore } from "@/store/modules/session";
import router from "@/router";

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export const client = axios.create({
  baseURL: `${import.meta.env.VITE_APP_BASE_URL}/api/v1`,
  timeout: 30000
});

client.interceptors.request.use(config => {
  const session = useSessionStore();
  if (session.accessToken) {
    config.headers.Authorization = `Bearer ${session.accessToken}`;
  }
  return config;
});

client.interceptors.response.use(
  response => response.data,
  error => {
    if (error?.response?.status === 401) {
      const session = useSessionStore();
      session.logout();
      if (router.currentRoute.value.path !== "/login") {
        router.replace("/login");
      }
    }
    Message.error(error?.response?.data?.message || "请求失败");
    return Promise.reject(error);
  }
);
