import { Asset } from "./Asset";
import invariant from "tiny-invariant";
import _Big from "big.js";
import toFormat from "toformat";

import {
  BigintIsh,
  Fraction,
  IFraction,
  isFraction,
  parseBigintIsh,
  Rounding,
  TEN,
} from "./fraction/Fraction";
import JSBI from "jsbi";
import B from "./utils/B";

const Big = toFormat(_Big);

export interface IAssetAmount extends IFraction {
  toFixed(decimalPlaces?: number, format?: object, rounding?: Rounding): string;
  asset: Asset;
  amount: JSBI;
  toBaseUnits: () => JSBI;
  toBaseUnitsFr: () => IFraction;
  toFormatted: (p?: {
    separator?: boolean;
    symbol?: boolean;
    decimals?: number;
  }) => string;
}

export class _AssetAmount implements IAssetAmount {
  protected fraction: IFraction;
  constructor(public asset: Asset, public amount: JSBI) {
    this.fraction = new Fraction(
      amount,
      JSBI.exponentiate(TEN, JSBI.BigInt(asset.decimals)),
    );
  }

  public toBaseUnits() {
    return this.multiply(
      JSBI.exponentiate(TEN, JSBI.BigInt(this.asset.decimals)),
    ).quotient;
  }

  public toBaseUnitsFr() {
    return new Fraction(this.toBaseUnits());
  }

  public toSignificant(
    significantDigits = 6,
    format?: object,
    rounding: Rounding = Rounding.ROUND_DOWN,
  ): string {
    return this.fraction.toSignificant(significantDigits, format, rounding);
  }

  public toFixed(
    decimalPlaces = this.asset.decimals,
    format?: object,
    rounding: Rounding = Rounding.ROUND_DOWN,
  ): string {
    // Provisional: This breaks app if falsy. N
    // Do not know why necessary if only for display
    // invariant(decimalPlaces <= this.asset.decimals, "DECIMALS");
    return this.fraction.toFixed(decimalPlaces, format, rounding);
  }

  public toExact(format: object = { groupSeparator: "" }): string {
    Big.DP = this.asset.decimals;
    return new Big(this.fraction.numerator.toString())
      .div(this.fraction.denominator.toString())
      .toFormat(format);
  }

  public get quotient() {
    return this.fraction.quotient;
  }

  public get remainder() {
    return this.fraction.remainder;
  }
  public get numerator() {
    return this.fraction.numerator;
  }
  public get denominator() {
    return this.fraction.denominator;
  }

  public invert() {
    return this.fraction.invert();
  }
  public add(other: IFraction | BigintIsh) {
    return this.fraction.add(other);
  }

  public subtract(other: IFraction | BigintIsh) {
    return this.fraction.subtract(other);
  }

  public lessThan(other: IFraction | BigintIsh) {
    return this.fraction.lessThan(other);
  }

  public lessThanOrEqual(other: IFraction | BigintIsh) {
    return this.fraction.greaterThanOrEqual(other);
  }

  public equalTo(other: IFraction | BigintIsh) {
    return this.fraction.equalTo(other);
  }

  public greaterThan(other: IFraction | BigintIsh) {
    return this.fraction.greaterThan(other);
  }

  public greaterThanOrEqual(other: IFraction | BigintIsh) {
    return this.fraction.greaterThanOrEqual(other);
  }

  public multiply(other: IFraction | BigintIsh) {
    return this.fraction.multiply(other);
  }

  public divide(other: IFraction | BigintIsh) {
    return this.fraction.divide(other);
  }

  // NOTE: This might eventually take a format string
  public toFormatted(params?: {
    decimals?: number;
    separator?: boolean;
    symbol?: boolean;
  }) {
    const { symbol = true } = params || {};
    // If decimals is too high fraction will bark
    const safeDecimals =
      typeof params?.decimals !== "undefined"
        ? this.asset.decimals < params.decimals
          ? this.asset.decimals
          : params.decimals
        : undefined;

    return [
      this.toFixed(safeDecimals, {
        groupSeparator: params?.separator ? "," : "",
      }),
      symbol ? this.asset.symbol.toUpperCase() : "",
    ]
      .filter(Boolean)
      .join(" ");
  }

  public toString() {
    return this.toFormatted();
  }
}

export type AssetAmount = IAssetAmount;

/**
 * Represents an amount of an Asset
 *
 * @param asset The Asset for the denomination
 * @param amount If amount is in JSBI then the amount this creates will be in base units (eg. Wei) otherwise the amount will be in natural units
 * @param options inBaseUnit boolean - if the asset amount given as string is in base units or not eg. 1000000000000000000 = 1 ether
 */
export function AssetAmount(
  asset: Asset,
  amount: string | number | JSBI | IFraction,
  options?: { inBaseUnit?: boolean },
): IAssetAmount {
  const { inBaseUnit = false } = options ?? {};
  if (inBaseUnit && typeof amount === "string") {
    return new _AssetAmount(asset, JSBI.BigInt(amount));
  }
  const unfractionedAmount = isFraction(amount)
    ? amount.toFixed(asset.decimals)
    : amount;

  const jsbiAmount =
    unfractionedAmount instanceof JSBI
      ? unfractionedAmount
      : B(unfractionedAmount, asset?.decimals);

  return new _AssetAmount(asset, jsbiAmount);
}
