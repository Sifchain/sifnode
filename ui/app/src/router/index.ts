import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";

import Swap from "@/views/SwapPage.vue";
import Pool from "@/views/PoolPage.vue";
import AddLiquidity from "@/views/AddLiquidityPage.vue";
import CreatePair from "@/views/CreatePoolPage.vue";
import Ui from "@/views/Ui.vue";

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

  {
    path: "/pool/add-liquidity",
    name: "AddLiquidity",
    component: CreatePair,
  },
  {
    path: "/pool/create-pair",
    name: "CreatePair",
    component: CreatePair,
  },

  // route for UI elements showcase - To Be Deleted
  {
    path: "/ui",
    name: "Ui",
    component: Ui,
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
