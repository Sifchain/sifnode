import { reactive } from "@vue/reactivity";
import { wallet, WalletStore } from "./wallet";
import { asset, AssetStore } from "./asset";
import { accountpools, AccountPoolStore, pools, PoolStore } from "./pools";
import { notifications, NotificationsStore } from "./notifications";
import { LiquidityProvider, Pool } from "../entities";
import { tx, TxStore } from "./tx";
export * from "./poolFinder";

export * from "./poolFinder";

// TODO: Consider storing tx key in local storage as an optimization?
export type Store = {
  wallet: WalletStore;
  asset: AssetStore;
  pools: PoolStore;
  tx: TxStore;
  accountpools: AccountPoolStore;
  notifications: NotificationsStore;
};

export function createStore() {
  return reactive<Store>({
    wallet,
    asset,
    pools,
    tx,
    accountpools,
    notifications,
  }) as Store;
}

export type WithStore<T extends keyof Store = keyof Store> = {
  store: Pick<Store, T>;
};
