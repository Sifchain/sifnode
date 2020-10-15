// import { reactive } from "@vue/reactivity";
import { Balance } from "src/entities";

export type WalletStore = {
  balances: Balance[];
  isConnected: boolean;
};

export const wallet = {
  balances: [],
  isConnected: false,
} as WalletStore;
