"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getConfig = void 0;
// TODO - Conditional load or build-time tree shake
const config_localnet_json_1 = __importDefault(require("./config.localnet.json"));
const config_devnet_json_1 = __importDefault(require("./config.devnet.json"));
const config_testnet_json_1 = __importDefault(require("./config.testnet.json"));
const config_mainnet_json_1 = __importDefault(require("./config.mainnet.json"));
const assets_ethereum_localnet_json_1 = __importDefault(require("./assets.ethereum.localnet.json"));
const assets_ethereum_sifchain_devnet_json_1 = __importDefault(require("./assets.ethereum.sifchain-devnet.json"));
const assets_ethereum_sifchain_testnet_json_1 = __importDefault(require("./assets.ethereum.sifchain-testnet.json"));
const assets_ethereum_mainnet_json_1 = __importDefault(require("./assets.ethereum.mainnet.json"));
const assets_sifchain_localnet_json_1 = __importDefault(require("./assets.sifchain.localnet.json"));
const assets_sifchain_mainnet_json_1 = __importDefault(require("./assets.sifchain.mainnet.json"));
const parseConfig_1 = require("./utils/parseConfig");
const entities_1 = require("./entities");
// Save assets for sync lookup throughout the app via Asset.get('symbol')
function cacheAsset(asset) {
    return entities_1.Asset(asset);
}
function getConfig(config = "localnet", sifchainAssetTag = "sifchain.localnet", ethereumAssetTag = "ethereum.localnet") {
    const assetMap = {
        "sifchain.localnet": parseConfig_1.parseAssets(assets_sifchain_localnet_json_1.default.assets),
        "sifchain.mainnet": parseConfig_1.parseAssets(assets_sifchain_mainnet_json_1.default.assets),
        "ethereum.localnet": parseConfig_1.parseAssets(assets_ethereum_localnet_json_1.default.assets),
        "ethereum.devnet": parseConfig_1.parseAssets(assets_ethereum_sifchain_devnet_json_1.default.assets),
        "ethereum.testnet": parseConfig_1.parseAssets(assets_ethereum_sifchain_testnet_json_1.default.assets),
        "ethereum.mainnet": parseConfig_1.parseAssets(assets_ethereum_mainnet_json_1.default.assets),
    };
    const sifchainAssets = assetMap[sifchainAssetTag];
    const ethereumAssets = assetMap[ethereumAssetTag];
    const allAssets = [...sifchainAssets, ...ethereumAssets].map(cacheAsset);
    const configMap = {
        localnet: parseConfig_1.parseConfig(config_localnet_json_1.default, allAssets),
        devnet: parseConfig_1.parseConfig(config_devnet_json_1.default, allAssets),
        testnet: parseConfig_1.parseConfig(config_testnet_json_1.default, allAssets),
        mainnet: parseConfig_1.parseConfig(config_mainnet_json_1.default, allAssets),
    };
    return configMap[config.toLowerCase()];
}
exports.getConfig = getConfig;
//# sourceMappingURL=config.js.map