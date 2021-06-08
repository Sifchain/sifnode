import { WalletStore } from "./wallet";
import { AssetStore } from "./asset";
import { AccountPoolStore, PoolStore } from "./pools";
import { TxStore } from "./tx";
export * from "./poolFinder";
export declare type Store = {
    wallet: WalletStore;
    asset: AssetStore;
    pools: PoolStore;
    tx: TxStore;
    accountpools: AccountPoolStore;
};
export declare function createStore(): Store;
export declare type WithStore<T extends keyof Store = keyof Store> = {
    store: Pick<Store, T>;
};
