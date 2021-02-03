import { reactive } from "@vue/reactivity";
import { wallet, WalletStore } from "./wallet";
import { asset, AssetStore } from "./asset";
import { pools, PoolStore } from "./pools";
import { notifications, NotificationsStore } from "./notifications";
import { LiquidityProvider, Pool } from "../entities";
import { tx, TxStore } from "./tx";
export * from "./poolFinder";

export type Store = {
  wallet: WalletStore;
  asset: AssetStore;
  pools: PoolStore;
  accountpools: { lp: LiquidityProvider; pool: Pool }[];
  notifications: NotificationsStore;
  tx: TxStore;
};

export function createStore() {
  const state = reactive<Store>({
    wallet,
    asset,
    pools,
    accountpools: [],
    notifications,
    tx,
  }) as Store;

  return state;
}

export type WithStore<T extends keyof Store = keyof Store> = {
  store: Pick<Store, T>;
};
