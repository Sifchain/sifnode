import { AssetAmount } from "../entities";
declare type AssetAmountMap = Map<string, AssetAmount>;
export declare class State {
    constructor(o?: Partial<State>);
    marketcapTokenOrder: string[];
    userBalances: AssetAmountMap;
    get availableAssetAccounts(): AssetAmount[];
}
export declare class StoreActions {
    state: State;
    constructor(state: State);
    setUserBalances(balances: AssetAmount[]): void;
}
export declare const store: StoreActions;
export declare function createStore(state?: Partial<State>): StoreActions;
export {};
