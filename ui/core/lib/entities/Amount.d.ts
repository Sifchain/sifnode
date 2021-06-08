import JSBI from "jsbi";
import { IFraction } from "./fraction/Fraction";
export declare type IAmount = {
    toBigInt(): JSBI;
    toString(detailed?: boolean): string;
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
export declare function Amount(source: JSBI | bigint | string | IAmount): Readonly<IAmount>;
export declare type _ExposeInternal<T extends IAmount> = T & {
    _toInternal(): IFraction;
    _fromInternal(fraction: IFraction): IAmount;
};
