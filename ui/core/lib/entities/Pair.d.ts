import { Asset } from "./Asset";
import { IAssetAmount } from "./AssetAmount";
export declare type Pair = ReturnType<typeof Pair>;
export declare function Pair(a: IAssetAmount, b: IAssetAmount): {
    amounts: [IAssetAmount, IAssetAmount];
    otherAsset(asset: Asset): IAssetAmount;
    symbol(): string;
    contains(...assets: Asset[]): boolean;
    getAmount(asset: Asset | string): IAssetAmount;
    toString(): string;
};
