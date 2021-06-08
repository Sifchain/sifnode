import { Asset } from "./Asset";
import { IAmount } from "./Amount";
export declare function LiquidityProvider(asset: Asset, units: IAmount, address: string, nativeAmount: IAmount, externalAmount: IAmount): {
    asset: import("./Asset").IAsset;
    units: IAmount;
    address: string;
    nativeAmount: IAmount;
    externalAmount: IAmount;
};
export declare type LiquidityProvider = {
    asset: Asset;
    units: IAmount;
    address: string;
    nativeAmount: IAmount;
    externalAmount: IAmount;
};
