import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
// import Home from "../views/Home.vue";
// import List from "../views/List.vue";
// import TransferEth from "../views/TransferEth.vue";
// import SifWallet from "@/views/SifWallet.vue";
import Swap from "@/views/SwapPage.vue";
import Pool from "@/views/PoolPage.vue";
import AddLiquidity from "@/views/AddLiquidityPage.vue";
import CreatePair from "@/views/CreatePairPage.vue";
import TestModal from "@/views/TestModal.vue";
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
    path: "/testmodal",
    name: "TestModal",
    component: TestModal,
  },
  {
    path: "/pool/add-liquidity",
    name: "AddLiquidity",
    component: AddLiquidity,
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
