import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { loginAPI, meAPI, type LoginPayload } from "@/api/auth";

interface SessionProfile {
  id: number;
  username: string;
  roles: string[];
  permissions: string[];
  nickname?: string;
  avatar?: string;
}

export const useSessionStore = defineStore(
  "session",
  () => {
    const accessToken = ref("");
    const refreshToken = ref("");
    const profile = ref<SessionProfile | null>(null);
    const profileLoaded = ref(false);

    const isAuthenticated = computed(() => !!accessToken.value);

    const login = async (payload: LoginPayload) => {
      const res = await loginAPI(payload);
      if (res.code !== 0 || !res.data?.access_token || !res.data?.refresh_token) {
        throw new Error(res.message || "登录失败");
      }
      accessToken.value = res.data.access_token;
      refreshToken.value = res.data.refresh_token;
      profile.value = res.data.user;
      profileLoaded.value = true;
    };

    const fetchProfile = async () => {
      const res = await meAPI();
      profile.value = {
        ...profile.value,
        ...res.data,
        permissions: profile.value?.permissions || []
      };
      profileLoaded.value = true;
      return profile.value;
    };

    const logout = () => {
      accessToken.value = "";
      refreshToken.value = "";
      profile.value = null;
      profileLoaded.value = false;
    };

    const hasPermission = (permission: string) => {
      if (profile.value?.roles?.includes("admin")) {
        return true;
      }
      return profile.value?.permissions?.includes(permission) ?? false;
    };

    const can = (permission: string) => hasPermission(permission);

    return {
      accessToken,
      refreshToken,
      profile,
      profileLoaded,
      isAuthenticated,
      login,
      fetchProfile,
      logout,
      hasPermission,
      can
    };
  },
  {
    persist: {
      key: "sleep-admin-session-v2",
      paths: ["accessToken", "refreshToken", "profile"]
    }
  }
);
