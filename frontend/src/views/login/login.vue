<template>
  <div class="login-page">
    <div class="hero">
      <div class="hero-badge">生产级睡眠内容后台</div>
      <h1>Sleep Admin Console</h1>
      <p>基于 Go + Gin + Casbin + Vue3 + Arco 的管理后台。</p>
    </div>
    <a-card class="login-card" :bordered="false">
      <div class="card-header">
        <img :src="logoURL" alt="logo" />
        <div>
          <div class="card-title">管理员登录</div>
          <div class="card-sub">默认账号 `admin / admin`</div>
        </div>
      </div>
      <a-form :model="form" layout="vertical" @submit.prevent="submit">
        <a-form-item field="username" label="账号">
          <a-input v-model="form.username" placeholder="请输入账号" />
        </a-form-item>
        <a-form-item field="password" label="密码">
          <a-input-password v-model="form.password" placeholder="请输入密码" />
        </a-form-item>
        <a-button type="primary" long html-type="submit" :loading="loading" @click="submit">登录</a-button>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { useRouter } from "vue-router";
import { useSessionStore } from "@/store/modules/session";
import logoURL from "@/assets/logo/snow.svg";

const router = useRouter();
const session = useSessionStore();
const loading = ref(false);
const form = reactive({
  username: "admin",
  password: "admin"
});

const submit = async () => {
  if (loading.value) {
    return;
  }

  const username = form.username.trim();
  const password = form.password;

  if (!username) {
    Message.warning("请输入账号");
    return;
  }

  if (!password) {
    Message.warning("请输入密码");
    return;
  }

  loading.value = true;
  try {
    await session.login({
      username,
      password
    });
    Message.success("登录成功");
    await router.push("/home");
  } catch (error: any) {
    console.error("登录失败:", error);
    Message.error(error?.response?.data?.message || error?.message || "登录失败");
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 1.2fr 420px;
  align-items: center;
  padding: 48px;
  gap: 48px;
  background:
    linear-gradient(135deg, rgba(8, 145, 178, 0.85), rgba(30, 64, 175, 0.9)),
    url("@/assets/img/login-bg.jpg") center/cover;
}

.hero {
  color: #fff;

  h1 {
    font-family: "AlimamaFangYuanTiVF-Thin", sans-serif;
    font-size: 64px;
    line-height: 1;
    margin: 12px 0 16px;
  }

  p {
    font-size: 18px;
    max-width: 560px;
    opacity: 0.88;
  }
}

.hero-badge {
  display: inline-flex;
  padding: 8px 14px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.12);
  backdrop-filter: blur(8px);
}

.login-card {
  padding: 8px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(16px);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;

  img {
    width: 42px;
    height: 42px;
  }
}

.card-title {
  font-size: 20px;
  font-weight: 700;
}

.card-sub {
  color: var(--color-text-3);
}

@media (max-width: 960px) {
  .login-page {
    grid-template-columns: 1fr;
    padding: 24px;
  }

  .hero h1 {
    font-size: 42px;
  }
}
</style>
