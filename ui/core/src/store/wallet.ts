import { reactive } from "@vue/reactivity";
import { Balance } from "../entities";

export type WalletStore = {
  balances: Balance[];
  etheriumIsConnected: boolean;
};

export const wallet = reactive({
  balances: [],
  etheriumIsConnected: false,
}) as WalletStore;
