import { createApp } from "vue";
import "@/style.css";
import App from "@/App.vue";
import router from "@/router/index";
import pinia from "@/store/index";
import ArcoVue from "@arco-design/web-vue";
import { registerArcoIcons } from "@/plugins/arco-icons";
import "@arco-themes/vue-gi-demo/css/arco.css";

const app = createApp(App);

registerArcoIcons(app);

app
  .use(pinia)
  .use(router)
  .use(ArcoVue)
  .mount("#app");
