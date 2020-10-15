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
  mnemonic?: string, // bip valudated string
  account?: Account,
  client?: SigningCosmosClient
} 

export const CWalletStore = reactive({
  balances: [],
  isConnected: false,
  mnemonic: undefined,
  account: undefined,
  client: undefined
}) as ICWalletStore;

// toRefs?