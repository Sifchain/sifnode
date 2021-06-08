import { LiquidityProvider, Pool } from "../entities";
export declare type PoolStore = {
    [s: string]: Pool;
};
export declare type AccountPool = {
    lp: LiquidityProvider;
    pool: string;
};
export declare type AccountPoolStore = {
    [address: string]: {
        [pool: string]: AccountPool;
    };
};
export declare const pools: PoolStore;
export declare const accountpools: AccountPoolStore;
