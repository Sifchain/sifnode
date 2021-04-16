import { reactive } from "@vue/reactivity";

import { Address, IAssetAmount } from "../entities";

export type WalletStore = {
  eth: {
    chainId?: string; // 0x1 is mainnet
    balances: IAssetAmount[];
    isConnected: boolean;
    address: Address;
  };
  sif: {
    balances: IAssetAmount[];
    isConnected: boolean;
    address: Address;
  };
};

export const wallet = reactive<WalletStore>({
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
