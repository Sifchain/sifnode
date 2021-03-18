import { IAmount, Amount, _ExposeInternal } from "./Amount";
import { IAsset, Asset } from "./Asset";

export type IAssetAmount = Readonly<IAsset> & IAmount;

export function AssetAmount(
  asset: IAsset | string,
  amount: IAmount | string,
): IAssetAmount {
  type _IAssetAmount = _ExposeInternal<IAssetAmount>;

  const _asset = Asset(asset);
  const _amount = Amount(amount);

  const instance: _IAssetAmount = {
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

    toString() {
      return `${_amount.toString()} ${_asset.symbol.toUpperCase()}`;
    },

    add(other) {
      return AssetAmount(_asset, _amount.add(other));
    },

    divide(other) {
      return AssetAmount(_asset, _amount.divide(other));
    },

    equalTo(other) {
      return _amount.equalTo(other); // TODO: do we care about assets? Suspect not.
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
      return AssetAmount(_asset, _amount.multiply(other));
    },

    sqrt() {
      return AssetAmount(_asset, _amount.sqrt());
    },

    subtract(other) {
      return AssetAmount(_asset, _amount.subtract(other));
    },

    _toInternal() {
      return (_amount as _ExposeInternal<IAmount>)._toInternal();
    },
  };
  return instance;
}
