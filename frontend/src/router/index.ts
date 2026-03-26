import { createRouter, createWebHashHistory, type RouteRecordRaw } from "vue-router";
import { Message } from "@arco-design/web-vue";
import { useSessionStore } from "@/store/modules/session";

export const routes: RouteRecordRaw[] = [
  { path: "/login", component: () => import("@/views/login/login.vue"), meta: { guest: true } },
  {
    path: "/",
    component: () => import("@/layout/index.vue"),
    redirect: "/home",
    children: [
      { path: "home", name: "home", component: () => import("@/views/home/home.vue"), meta: { title: "工作台" } },
      { path: "articles", name: "articles", component: () => import("@/views/biz/articles/index.vue"), meta: { title: "文章管理", permission: "/api/v1/articles#GET" } },
      { path: "articles/batch", name: "article-batch", component: () => import("@/views/biz/articles/batch.vue"), meta: { title: "批量新增", permission: "/api/v1/articles/batch#POST" } },
      { path: "articles/new", name: "article-new", component: () => import("@/views/biz/articles/form.vue"), meta: { title: "新增文章", permission: "/api/v1/articles#POST", hidden: true } },
      { path: "articles/:id", name: "article-edit", component: () => import("@/views/biz/articles/form.vue"), meta: { title: "文章详情", permission: "/api/v1/articles/:id#PUT", hidden: true } },
      { path: "content-menus", name: "content-menus", component: () => import("@/views/biz/menus/index.vue"), meta: { title: "内容菜单", permission: "/api/v1/menus/tree#GET" } },
      { path: "uploads", name: "uploads", component: () => import("@/views/biz/uploads/index.vue"), meta: { title: "资源上传", permission: "/api/v1/uploads#GET" } },
      { path: "system/admin-menus", name: "system-admin-menus", component: () => import("@/views/biz/system/admin-menus.vue"), meta: { title: "后台菜单", permission: "/api/v1/admin-menus/tree#GET" } },
      { path: "system/admins", name: "system-admins", component: () => import("@/views/biz/system/admins.vue"), meta: { title: "管理员", permission: "/api/v1/admins#GET" } },
      { path: "system/roles", name: "system-roles", component: () => import("@/views/biz/system/roles.vue"), meta: { title: "角色权限", permission: "/api/v1/roles#GET" } },
      { path: "system/dicts", name: "system-dicts", component: () => import("@/views/biz/system/dicts.vue"), meta: { title: "字典管理", permission: "/api/v1/dict-types#GET" } },
      { path: "system/configs", name: "system-configs", component: () => import("@/views/biz/system/configs.vue"), meta: { title: "系统参数", permission: "/api/v1/sys-configs#GET" } },
      { path: "system/operation-logs", name: "system-operation-logs", component: () => import("@/views/biz/system/operation-logs.vue"), meta: { title: "操作日志", permission: "/api/v1/operation-logs#GET" } },
      { path: "system/login-logs", name: "system-login-logs", component: () => import("@/views/biz/system/login-logs.vue"), meta: { title: "登录日志", permission: "/api/v1/login-logs#GET" } },
      { path: "system/online-users", name: "system-online-users", component: () => import("@/views/biz/system/online-users.vue"), meta: { title: "在线用户", permission: "/api/v1/online-sessions#GET" } },
      { path: "system/apis", name: "system-apis", component: () => import("@/views/biz/system/apis.vue"), meta: { title: "API 管理", permission: "/api/v1/admin-menus/tree#GET" } },
      { path: "system/jobs", name: "system-jobs", component: () => import("@/views/biz/system/jobs.vue"), meta: { title: "定时任务", permission: "/api/v1/jobs#GET" } },
      { path: "system/depts", name: "system-depts", component: () => import("@/views/biz/system/depts.vue"), meta: { title: "部门管理", permission: "/api/v1/depts/tree#GET" } },
      { path: "system/posts", name: "system-posts", component: () => import("@/views/biz/system/posts.vue"), meta: { title: "岗位管理", permission: "/api/v1/posts#GET" } },
      { path: "modules", name: "modules", component: () => import("@/views/biz/modules/index.vue"), meta: { title: "通用模块" } }
    ]
  },
  { path: "/:pathMatch(.*)*", redirect: "/home" }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes
});

router.beforeEach(async (to, _, next) => {
  const session = useSessionStore();
  if (to.meta.guest) {
    if (session.isAuthenticated) {
      next("/home");
    } else {
      next();
    }
    return;
  }

  if (!session.isAuthenticated) {
    next("/login");
    return;
  }

  if (!session.profileLoaded) {
    try {
      await session.fetchProfile();
    } catch {
      session.logout();
      Message.error("登录状态已失效，请重新登录");
      next("/login");
      return;
    }
  }

  const permission = to.meta.permission as string | undefined;
  if (permission && !session.hasPermission(permission)) {
    Message.warning("当前账号没有访问该页面的权限");
    next("/home");
    return;
  }
  next();
});

export default router;
