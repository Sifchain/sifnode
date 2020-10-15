import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import Home from "../views/Home.vue";
import List from "../views/List.vue";
import CWallet from '@/views/CWallet.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: '/cosmos',
    name: 'CWallet',
    component: CWallet
  },
  {
    path: "/list",
    name: "List",
    component: List,
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
