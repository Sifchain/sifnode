import { IAmount } from "../entities/Amount";
import { IAssetAmount } from "../entities/AssetAmount";
import { IAsset } from "../entities";
declare type IFormatOptionsBase = {
    exponent?: number;
    forceSign?: boolean;
    mode?: "number" | "percent";
    separator?: boolean;
    space?: boolean;
    prefix?: string;
    postfix?: string;
    zeroFormat?: string;
};
declare type IFormatOptionsMantissa<M = number | DynamicMantissa> = IFormatOptionsBase & {
    shorthand?: boolean;
    mantissa?: M;
    trimMantissa?: boolean | "integer";
};
declare type IFormatOptionsShorthandTotalLength = IFormatOptionsBase & {
    shorthand: true;
    totalLength?: number;
};
export declare type DynamicMantissa = Record<number | "infinity", number>;
export declare type IFormatOptions = IFormatOptionsMantissa | IFormatOptionsShorthandTotalLength;
/**
 * Takes an amount and a dynamic mantissa hash and returns the mantisaa value to use
 * @param amount amount given to format function
 * @param hash dynamic value hash to calculate mantissa from
 * @returns number of mantissa to send to formatter
 */
export declare function getMantissaFromDynamicMantissa(amount: IAmount, hash: DynamicMantissa): any;
export declare function round(decimal: string, places: number): string;
export declare type AmountNotAssetAmount<T extends IAmount> = T extends IAssetAmount ? never : T;
export declare function format<T extends IAmount>(amount: AmountNotAssetAmount<T>): string;
export declare function format<T extends IAmount>(amount: AmountNotAssetAmount<T>, asset: Exclude<IAsset, IAssetAmount>): string;
export declare function format<T extends IAmount>(amount: AmountNotAssetAmount<T>, options: IFormatOptions): string;
export declare function format<T extends IAmount>(amount: AmountNotAssetAmount<T>, asset: Exclude<IAsset, IAssetAmount>, options: IFormatOptions): string;
export declare function trimMantissa(decimal: string, integer?: boolean): string;
export {};
