import { IAmount, Amount, _ExposeInternal } from "./Amount";
import { IAsset, Asset } from "./Asset";
import { IFraction } from "./fraction/Fraction";
import { fromBaseUnits } from "../utils/decimalShift";

import JSBI from "jsbi";

export type IAssetAmount = Readonly<IAsset> & {
  readonly asset: IAsset;
  readonly amount: IAmount;

  // COMPATABILITY OPERATORS
  // For use by display lib and in testing

  toBigInt(): JSBI;
  toString(detailed?: boolean): string;

  // CONVENIENCE UTILITIES
  // Utilty operators common enough to live on this class

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

  // MATH OPERATORS

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

export function AssetAmount(
  asset: IAsset | string,
  amount: IAmount | string,
): IAssetAmount {
  type _IAssetAmount = _ExposeInternal<IAssetAmount>;

  const _asset = (asset as IAssetAmount)?.asset || Asset(asset);
  const _amount = (amount as IAssetAmount)?.amount || Amount(amount);

  // Proxy all methods because it is clearer and
  // more explicit than prototypal inheritence
  const instance: _IAssetAmount = {
    get asset() {
      return _asset;
    },

    get amount() {
      return _amount;
    },

    get address() {
      return _asset.address;
    },

    get decimals() {
      return _asset.decimals;
    },

    get imageUrl() {
      return _asset.imageUrl;
    },

    get name() {
      return _asset.name;
    },

    get network() {
      return _asset.network;
    },

    get symbol() {
      return _asset.symbol;
    },

    get label() {
      return _asset.label;
    },

    toDerived() {
      return _amount.multiply(fromBaseUnits("1", _asset));
    },

    toBigInt() {
      return _amount.toBigInt();
    },

    toString() {
      return `${_amount.toString(false)} ${_asset.symbol.toUpperCase()}`;
    },

    add(other) {
      return _amount.add(other);
    },

    divide(other) {
      return _amount.divide(other);
    },

    equalTo(other) {
      return _amount.equalTo(other);
    },

    greaterThan(other) {
      return _amount.greaterThan(other);
    },

    greaterThanOrEqual(other) {
      return _amount.greaterThanOrEqual(other);
    },

    lessThan(other) {
      return _amount.lessThan(other);
    },

    lessThanOrEqual(other) {
      return _amount.lessThanOrEqual(other);
    },

    multiply(other) {
      return _amount.multiply(other);
    },

    sqrt() {
      return _amount.sqrt();
    },

    subtract(other) {
      return _amount.subtract(other);
    },

    // Internal methods need to be proxied here so this can be used as an Amount
    _toInternal() {
      return (_amount as _ExposeInternal<IAmount>)._toInternal();
    },

    _fromInternal(internal: IFraction) {
      return (_amount as _ExposeInternal<IAmount>)._fromInternal(internal);
    },
  };
  return instance;
}

export function isAssetAmount(value: any): value is IAssetAmount {
  return value?.asset && value?.amount;
}
