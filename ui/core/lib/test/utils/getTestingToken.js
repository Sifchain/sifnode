"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getBalance = exports.getTestingTokens = exports.getTestingToken = void 0;
const assets_ethereum_localnet_json_1 = __importDefault(require("../../assets.ethereum.localnet.json"));
const assets_sifchain_localnet_json_1 = __importDefault(require("../../assets.sifchain.localnet.json"));
const parseConfig_1 = require("../../utils/parseConfig");
const entities_1 = require("../../entities");
const assets = [...assets_ethereum_localnet_json_1.default.assets, ...assets_sifchain_localnet_json_1.default.assets];
function getTestingToken(tokenSymbol) {
    const supportedTokens = parseConfig_1.parseAssets(assets).map((asset) => {
        entities_1.Asset.set(asset.symbol, asset);
        return asset;
    });
    const asset = supportedTokens.find(({ symbol }) => symbol.toUpperCase() === tokenSymbol.toUpperCase());
    if (!asset)
        throw new Error(`${tokenSymbol} not returned`);
    return asset;
}
exports.getTestingToken = getTestingToken;
function getTestingTokens(tokens) {
    return tokens.map(getTestingToken);
}
exports.getTestingTokens = getTestingTokens;
function getBalance(balances, symbol) {
    const bal = balances.find(({ asset }) => asset.symbol.toUpperCase() === symbol.toUpperCase());
    if (!bal)
        throw new Error("Symbol not found in balances");
    return bal;
}
exports.getBalance = getBalance;
//# sourceMappingURL=getTestingToken.js.map