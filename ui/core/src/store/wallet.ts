import { reactive } from "@vue/reactivity";

import { Address, Balance, SifBalance } from "../entities";

export type WalletStore = {
  eth: {
    balances: Balance[];
    isConnected: boolean;
    address: Address;
  };
  sif: {
    balances: readonly SifBalance[]; // 🤔 why readonly?
    isConnected: boolean;
    address: Address;
  };
};

export const wallet = reactive({
  eth: {
    isConnected: false,
    address: "",
    balances: [],
  },
  sif: {
    isConnected: false,
    address: "",
    balances: [],
  },
}) as WalletStore;
