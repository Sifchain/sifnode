import { reactive } from "@vue/reactivity";
import { wallet, WalletStore } from "./wallet";
import { asset, AssetStore } from "./asset";
import { pools, PoolStore } from "./pools";
import { notifications, NotificationsStore } from "./notifications";
import { Pool } from "../entities";
export * from "./poolFinder";

export type Store = {
  wallet: WalletStore;
  asset: AssetStore;
  pools: PoolStore;
  accountpools: Pool[];
  notifications: NotificationsStore;
};

export function createStore() {
  const state = reactive<Store>({
    wallet,
    asset,
    pools,
    accountpools: [],
    notifications,
  }) as Store;

  return state;
}

export type WithStore<T extends keyof Store = keyof Store> = {
  store: Pick<Store, T>;
};
