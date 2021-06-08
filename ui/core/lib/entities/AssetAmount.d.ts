import { IAmount } from "./Amount";
import { IAsset } from "./Asset";
import JSBI from "jsbi";
export declare type IAssetAmount = Readonly<IAsset> & {
    readonly asset: IAsset;
    readonly amount: IAmount;
    toBigInt(): JSBI;
    toString(detailed?: boolean): string;
    /**
     * Return the derived value for the AssetAmount based on the asset's decimals
     * For example lets say we have one eth:
     *
     * AssetAmount("eth", "100000000000000000").toDerived().equalTo(Amount("1")); // true
     * AssetAmount("usdc", "1000000").toDerived().equalTo(Amount("1")); // true
     *
     * All Math operators on AssetAmounts work on BaseUnits and return Amount objects. We have explored
     * creating a scheme for allowing mathematical operations to combine AssetAmounts but it is not clear as to
     * how to combine assets and which asset gets precedence. These rules would need to be internalized
     * by the team which may not be intuitive. Also performing mathematical operations on differing currencies
     * doesn't make conceptual sense.
     *
     *   Eg. What is 1 USDC x 1 ETH?
     *       - is it a value in ETH?
     *       - is it a value in USDC?
     *       - Mathematically it would be an ETHUSDC but we have no concept of this in our systems nor do we need one
     *
     * Instead when mixing AssetAmounts it is recommended to first extract the derived amounts and perform calculations on those
     *
     *   Eg.
     *   const ethAmount = oneEth.toDerived();
     *   const usdcAmount = oneUsdc.toDerived();
     *   const newAmount = ethAmount.multiply(usdcAmount);
     *   const newUsdcAmount = AssetAmount('usdc', newAmount);
     *
     * @returns IAmount
     */
    toDerived(): IAmount;
    add(other: IAmount | string): IAmount;
    divide(other: IAmount | string): IAmount;
    equalTo(other: IAmount | string): boolean;
    greaterThan(other: IAmount | string): boolean;
    greaterThanOrEqual(other: IAmount | string): boolean;
    lessThan(other: IAmount | string): boolean;
    lessThanOrEqual(other: IAmount | string): boolean;
    multiply(other: IAmount | string): IAmount;
    sqrt(): IAmount;
    subtract(other: IAmount | string): IAmount;
};
export declare function AssetAmount(asset: IAsset | string, amount: IAmount | string): IAssetAmount;
export declare function isAssetAmount(value: any): value is IAssetAmount;
