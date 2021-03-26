// TODO - Conditional load or build-time tree shake
import localnetconfig from "./config.localnet.json";
import devnetconfig from "./config.devnet.json";
import testnetconfig from "./config.testnet.json";
import mainnnetconfig from "./config.mainnet.json";

import assetsEthereumLocalnet from "./assets.ethereum.localnet.json";
import assetsEthereumDevnet from "./assets.ethereum.devnet.json";
import assetsEthereumTestnet from "./assets.ethereum.testnet.json";
import assetsEthereumMainnet from "./assets.ethereum.mainnet.json";

import assetsSifchainLocalnet from "./assets.sifchain.localnet.json";
import assetsSifchainMainnet from "./assets.sifchain.mainnet.json";

import {
  parseConfig,
  parseAssets,
  ChainConfig,
  AssetConfig,
} from "./utils/parseConfig";
import { Asset } from "./entities";
import { ApiContext } from "./api";

type ConfigMap = { [s: string]: ApiContext };
type AssetMap = { [s: string]: Asset[] };

// Save assets for sync lookup throughout the app via Asset.get('symbol')
function cacheAsset(asset: Asset) {
  return Asset(asset);
}

export type AppConfig = ApiContext; // Will include other injectables

export function getConfig(
  config = "localnet",
  sifchainAssetTag = "sifchain.localnet",
  ethereumAssetTag = "ethereum.localnet",
): AppConfig {
  const assetMap: AssetMap = {
    "sifchain.localnet": parseAssets(
      assetsSifchainLocalnet.assets as AssetConfig[],
    ),
    "sifchain.mainnet": parseAssets(
      assetsSifchainMainnet.assets as AssetConfig[],
    ),
    "ethereum.localnet": parseAssets(
      assetsEthereumLocalnet.assets as AssetConfig[],
    ),
    "ethereum.devnet": parseAssets(
      assetsEthereumDevnet.assets as AssetConfig[],
    ),
    "ethereum.testnet": parseAssets(
      assetsEthereumTestnet.assets as AssetConfig[],
    ),
    "ethereum.mainnet": parseAssets(
      assetsEthereumMainnet.assets as AssetConfig[],
    ),
  };

  const sifchainAssets = assetMap[sifchainAssetTag];
  const ethereumAssets = assetMap[ethereumAssetTag];
  const allAssets = [...sifchainAssets, ...ethereumAssets].map(cacheAsset);

  const configMap: ConfigMap = {
    localnet: parseConfig(localnetconfig as ChainConfig, allAssets),
    devnet: parseConfig(devnetconfig as ChainConfig, allAssets),
    testnet: parseConfig(testnetconfig as ChainConfig, allAssets),
    mainnet: parseConfig(mainnnetconfig as ChainConfig, allAssets),
  };

  return configMap[config.toLowerCase()];
}
