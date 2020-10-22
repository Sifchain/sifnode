import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
// import Home from "../views/Home.vue";
// import List from "../views/List.vue";
// import TransferEth from "../views/TransferEth.vue";
// import SifWallet from "@/views/SifWallet.vue";
import Swap from "@/views/SwapPage.vue";
import Pool from "@/views/PoolPage.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    redirect: { name: "Swap" },
  },
  {
    path: "/swap",
    name: "Swap",
    component: Swap,
  },
  {
    path: "/pool",
    name: "Pool",
    component: Pool,
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
