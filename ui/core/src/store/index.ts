import { reactive } from "@vue/reactivity";
import { wallet, WalletStore } from "./wallet";
import { asset, AssetStore } from "./asset";
import { accountpools, AccountPoolStore, pools, PoolStore } from "./pools";

export * from "./poolFinder";

// TODO: Add a tx lookup per blockchain so we have access to txs
// TODO: Consider storing local txs key in local storage as an effect

export type Store = {
  wallet: WalletStore;
  asset: AssetStore;
  pools: PoolStore;
  accountpools: AccountPoolStore;
};

export function createStore() {
  return reactive<Store>({
    wallet,
    asset,
    pools,
    accountpools,
  }) as Store;
}

export type WithStore<T extends keyof Store = keyof Store> = {
  store: Pick<Store, T>;
};
