import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";

import Swap from "@/views/SwapPage.vue";
import Pool from "@/views/PoolPage.vue";
import CreatePool from "@/views/CreatePoolPage.vue";
import RemoveLiquidity from "@/views/RemoveLiquidityPage.vue";
import SinglePool from "@/components/poolList/SinglePool.vue";
import PegListingPage from "@/views/PegListingPage.vue";
import PegAssetPage from "@/views/PegAssetPage.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    redirect: { name: "PegListingPage" },
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
    path: "/pool/:externalAsset",
    name: "SinglePool",
    component: SinglePool,
  },
  {
    path: "/pool/add-liquidity/:externalAsset?",
    name: "AddLiquidity",
    component: CreatePool,
    props: {
      title: "Add Liquidity",
    },
  },
  {
    path: "/pool/create-pool",
    name: "CreatePool",
    component: CreatePool,
    props: {
      title: "Create Pair",
    },
  },
  {
    path: "/pool/remove-liquidity/:externalAsset?",
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
    path: "/peg/reverse/:assetFrom/:assetTo",
    name: "UnpegAssetPage",
    component: PegAssetPage,
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
