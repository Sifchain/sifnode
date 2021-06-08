import { Asset } from "./Asset";
import { IAssetAmount } from "./AssetAmount";
import { IAmount } from "./Amount";
export declare type Pool = ReturnType<typeof Pool>;
export declare type IPool = Omit<Pool, "poolUnits" | "calculatePoolUnits">;
export declare function Pool(a: IAssetAmount, // native asset
b: IAssetAmount, // external asset
poolUnits?: IAmount): {
    amounts: [IAssetAmount, IAssetAmount];
    otherAsset: (asset: import("./Asset").IAsset) => IAssetAmount;
    symbol: () => string;
    contains: (...assets: import("./Asset").IAsset[]) => boolean;
    toString: () => string;
    getAmount: (asset: string | import("./Asset").IAsset) => IAssetAmount;
    poolUnits: IAmount;
    priceAsset(asset: Asset): IAssetAmount;
    calcProviderFee(x: IAssetAmount): IAssetAmount;
    calcPriceImpact(x: IAssetAmount): IAmount;
    calcSwapResult(x: IAssetAmount): IAssetAmount;
    calcReverseSwapResult(Sa: IAssetAmount): IAssetAmount;
    calculatePoolUnits(nativeAssetAmount: IAssetAmount, externalAssetAmount: IAssetAmount): IAmount[];
};
export declare function CompositePool(pair1: IPool, pair2: IPool): IPool;
