import { reactive } from "@vue/reactivity";
import { wallet, WalletStore } from "./wallet";
import { asset, AssetStore } from "./asset";
import { pools, PoolStore } from "./pools";
import { notifications, NotificationsStore } from "./notifications";
import { LiquidityProvider, Pool } from "../entities";
export * from "./poolFinder";

// TODO: Add a tx lookup per blockchain so we have access to txs
// TODO: Consider storing local txs key in local storage as an effect

export type Store = {
  wallet: WalletStore;
  asset: AssetStore;
  pools: PoolStore;
  accountpools: { lp: LiquidityProvider; pool: Pool }[];
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
