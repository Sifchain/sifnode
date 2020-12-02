import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";

import Swap from "@/views/SwapPage.vue";
import Pool from "@/views/PoolPage.vue";
import CreatePool from "@/views/CreatePoolPage.vue";
import Ui from "@/views/Ui.vue";
import RemoveLiquidity from "@/views/RemoveLiquidityPage.vue";

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
    component: CreatePool,
  },
  {
    path: "/pool/create-pool",
    name: "CreatePool",
    component: CreatePool,
  },
  {
    path: "/pool/remove-liquidity",
    name: "RemoveLiquidity",
    component: RemoveLiquidity,
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
