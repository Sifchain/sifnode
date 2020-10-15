import { reactive } from "@vue/reactivity";
import { Balance } from "../entities";
import {Ref, ref} from "@vue/reactivity"

export type WalletStore = {
  balances: Balance[];
  isConnected: boolean
};

export const wallet = reactive({
  balances: [],
  isConnected: false
}) as WalletStore;



export interface ICWalletStore extends  WalletStore {
  mnemonic?: string
} 

export const CWalletStore = reactive({
  balances: [],
  isConnected: false,
  mnemonic: ""
}) as ICWalletStore;

// toRefs?