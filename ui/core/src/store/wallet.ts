import { reactive } from "@vue/reactivity";

import { Address, IAssetAmount } from "../entities";

export type WalletStore = {
  eth: {
    chainId?: string;
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
    chainId: "",
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
