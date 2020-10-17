import { reactive } from "@vue/reactivity";
import { Balance } from "../entities";
import {Ref, ref} from "@vue/reactivity"

import {SigningCosmosClient, Account} from "@cosmjs/launchpad";
export type WalletStore = {
  balances: Balance[];
  isConnected: boolean
};

export const wallet = reactive({
  balances: [],
  isConnected: false
}) as WalletStore;



export interface ICWalletStore extends  WalletStore {
  account?: Account,
  client?: SigningCosmosClient
} 

export function createWallet() {
  return 
}
export const CWalletStore = reactive({
  balances: [],
  isConnected: false,
  account: undefined,
  client: undefined
}) as ICWalletStore;

// toRefs?