import JSBI from "jsbi";
export declare type BigintIsh = JSBI | bigint | string;
export declare enum Rounding {
    ROUND_DOWN = 0,
    ROUND_HALF_UP = 1,
    ROUND_UP = 2
}
export declare const ZERO: JSBI;
export declare const ONE: JSBI;
export declare const TWO: JSBI;
export declare const THREE: JSBI;
export declare const FIVE: JSBI;
export declare const TEN: JSBI;
export declare const _100: JSBI;
export declare const _997: JSBI;
export declare const _1000: JSBI;
export declare function parseBigintIsh(bigintIsh: BigintIsh): JSBI;
export interface IFraction {
    readonly quotient: JSBI;
    readonly remainder: IFraction;
    readonly numerator: JSBI;
    readonly denominator: JSBI;
    invert(): IFraction;
    add(other: IFraction | BigintIsh): IFraction;
    subtract(other: IFraction | BigintIsh): IFraction;
    lessThan(other: IFraction | BigintIsh): boolean;
    lessThanOrEqual(other: IFraction | BigintIsh): boolean;
    equalTo(other: IFraction | BigintIsh): boolean;
    greaterThan(other: IFraction | BigintIsh): boolean;
    greaterThanOrEqual(other: IFraction | BigintIsh): boolean;
    multiply(other: IFraction | BigintIsh): IFraction;
    divide(other: IFraction | BigintIsh): IFraction;
    toSignificant(significantDigits: number, format?: object, rounding?: Rounding): string;
    toFixed(decimalPlaces: number, format?: object, rounding?: Rounding): string;
}
export declare function isFraction(value: unknown): value is IFraction;
export declare class Fraction implements IFraction {
    readonly numerator: JSBI;
    readonly denominator: JSBI;
    constructor(numerator: BigintIsh, denominator?: BigintIsh);
    get quotient(): JSBI;
    get remainder(): IFraction;
    invert(): IFraction;
    add(other: IFraction | BigintIsh): IFraction;
    subtract(other: IFraction | BigintIsh): IFraction;
    lessThan(other: IFraction | BigintIsh): boolean;
    lessThanOrEqual(other: IFraction | BigintIsh): boolean;
    equalTo(other: IFraction | BigintIsh): boolean;
    greaterThan(other: IFraction | BigintIsh): boolean;
    greaterThanOrEqual(other: IFraction | BigintIsh): boolean;
    multiply(other: IFraction | BigintIsh): IFraction;
    divide(other: IFraction | BigintIsh): IFraction;
    toSignificant(significantDigits: number, format?: object, rounding?: Rounding): string;
    toFixed(decimalPlaces: number, format?: object, rounding?: Rounding): string;
}
