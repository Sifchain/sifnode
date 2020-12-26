import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";

import Swap from "@/views/SwapPage.vue";
import Pool from "@/views/PoolPage.vue";
import CreatePool from "@/views/CreatePoolPage.vue";
import RemoveLiquidity from "@/views/RemoveLiquidityPage.vue";
import PegListingPage from "@/views/PegListingPage.vue";
import PegAssetPage from "@/views/PegAssetPage.vue";

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
  {
    path: "/peg",
    name: "PegListingPage",
    component: PegListingPage,
  },
  {
    path: "/peg/:assetFrom/:assetTo",
    name: "PegAssetPage",
    component: PegAssetPage,
  },
  {
    path: "/unpeg/:assetFrom/:assetTo",
    name: "UnpegAssetPage",
    component: PegAssetPage,
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
