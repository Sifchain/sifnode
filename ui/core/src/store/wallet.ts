import { reactive } from "@vue/reactivity";
import { Balance } from "../entities";

export type WalletStore = {
  balances: Balance[];
  isConnected: boolean
};

export const wallet = reactive({
  balances: [],
  isConnected: false
}) as WalletStore;
