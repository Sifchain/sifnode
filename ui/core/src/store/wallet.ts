import { reactive } from "@vue/reactivity";
import { SigningCosmosClient, Account } from "@cosmjs/launchpad";

import { Balance, SifAddress } from "../entities";

export type WalletStore = {
  balances: Balance[];
  etheriumIsConnected: boolean;
  isConnected: boolean
};

export const wallet = reactive({
  balances: [],
  etheriumIsConnected: false,
  isConnected: false
}) as WalletStore;

/*
 combining wallets premature without understand
  how we want them to work in /app
*/
export type SifWalletStore = {
  isConnected: boolean
  client?: SigningCosmosClient
  address?: SifAddress
  balances?: Account
} 
