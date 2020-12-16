import { reactive, Ref, toRefs } from "@vue/reactivity";
import { wallet, WalletStore } from "./wallet";
import { asset, AssetStore } from "./asset";
import { pools, PoolStore } from "./pools";
import { notifications, NotificationsStore } from "./notifications"
export * from "./poolFinder";

export type Store = {
  wallet: WalletStore;
  asset: AssetStore;
  pools: PoolStore;
  notifications: NotificationsStore;
};

export function createStore() {
  const state = reactive<Store>({
    wallet,
    asset,
    pools,
    notifications
  }) as Store;

  return state;
}

export type WithStore<T extends keyof Store = keyof Store> = {
  store: Pick<Store, T>;
};
