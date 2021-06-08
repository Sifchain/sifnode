"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.isAssetAmount = exports.AssetAmount = void 0;
const Amount_1 = require("./Amount");
const Asset_1 = require("./Asset");
const decimalShift_1 = require("../utils/decimalShift");
function AssetAmount(asset, amount) {
    var _a, _b;
    const _asset = ((_a = asset) === null || _a === void 0 ? void 0 : _a.asset) || Asset_1.Asset(asset);
    const _amount = ((_b = amount) === null || _b === void 0 ? void 0 : _b.amount) || Amount_1.Amount(amount);
    // Proxy all methods because it is clearer and
    // more explicit than prototypal inheritence
    const instance = {
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
            return _amount.multiply(decimalShift_1.fromBaseUnits("1", _asset));
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
            return _amount._toInternal();
        },
        _fromInternal(internal) {
            return _amount._fromInternal(internal);
        },
    };
    return instance;
}
exports.AssetAmount = AssetAmount;
function isAssetAmount(value) {
    return (value === null || value === void 0 ? void 0 : value.asset) && (value === null || value === void 0 ? void 0 : value.amount);
}
exports.isAssetAmount = isAssetAmount;
//# sourceMappingURL=AssetAmount.js.map