import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import Home from "../views/Home.vue";
import List from "../views/List.vue";
import TransferEth from "../views/TransferEth.vue";
import SifWallet from "@/views/SifWallet.vue";
import Swap from "@/views/Swap.vue";
import Ui from "@/views/Ui.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/sifwallet",
    name: "SifWallet",
    component: SifWallet,
  },
  {
    path: "/swap",
    name: "Swap",
    component: Swap,
  },
  {
    path: "/list",
    name: "List",
    component: List,
  },
  {
    path: "/ethtransfer",
    name: "TransferEth",
    component: TransferEth,
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
