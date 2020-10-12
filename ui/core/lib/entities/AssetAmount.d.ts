import { Asset } from "./Asset";
import { BigintIsh, Fraction, Rounding } from "./fraction/Fraction";
export declare class AssetAmount extends Fraction {
    asset: Asset;
    amount: BigintIsh;
    constructor(asset: Asset, amount: BigintIsh);
    toSignificant(significantDigits?: number, format?: object, rounding?: Rounding): string;
    toFixed(decimalPlaces?: number, format?: object, rounding?: Rounding): string;
    toExact(format?: object): string;
    static create(asset: Asset, amount: BigintIsh): AssetAmount;
}
export declare type AssetBalancesByAddress = {
    [address: string]: AssetAmount | undefined;
};
