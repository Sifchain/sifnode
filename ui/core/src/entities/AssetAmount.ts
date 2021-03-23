import { IAmount, Amount, _ExposeInternal } from "./Amount";
import { IAsset, Asset } from "./Asset";
import { IFraction } from "./fraction/Fraction";
import Big from "big.js";

export type IAssetAmount = Readonly<IAsset> &
  IAmount & {
    readonly asset: IAsset;
    readonly amount: IAmount;
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

    toBigInt() {
      return _amount.toBigInt();
    },

    toString(detailed?: boolean) {
      return `${_amount.toString(detailed)} ${_asset.symbol.toUpperCase()}`;
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
