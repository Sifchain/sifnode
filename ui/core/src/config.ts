// TODO - Conditional load or build-time tree shake
import mainnnetconfig from "./config.mainnet.json";
import localnetconfig from "./config.localnet.json";
import sandpitconfig from "./config.sandpit.json";
import testnetconfig from "./config.testnet.json";
import assetsEthereumLocalnet from "./assets.ethereum.localnet.json";
import assetsEthereumMainnet from "./assets.ethereum.mainnet.json";
import assetsEthereumSandpit from "./assets.ethereum.sandpit.json";
import assetsSifchainLocalnet from "./assets.sifchain.localnet.json";
import assetsSifchainMainnet from "./assets.sifchain.mainnet.json";
import assetsSifchainSandpit from "./assets.sifchain.sandpit.json";

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
  Asset.set(asset.symbol, asset);
  return asset;
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
    "sifchain.sandpit": parseAssets(
      assetsSifchainSandpit.assets as AssetConfig[],
    ),
    "ethereum.localnet": parseAssets(
      assetsEthereumLocalnet.assets as AssetConfig[],
    ),
    "ethereum.ropsten": parseAssets(
      assetsEthereumSandpit.assets as AssetConfig[],
    ),
    "ethereum.mainnet": parseAssets(
      assetsEthereumMainnet.assets as AssetConfig[],
    ),
  };

  const sifchainAssets = assetMap[sifchainAssetTag];
  const ethereumAssets = assetMap[ethereumAssetTag];
  const allAssets = [...sifchainAssets, ...ethereumAssets].map(cacheAsset);

  const configMap: ConfigMap = {
    mainnet: parseConfig(mainnnetconfig as ChainConfig, allAssets),
    localnet: parseConfig(localnetconfig as ChainConfig, allAssets),
    testnet: parseConfig(testnetconfig as ChainConfig, allAssets),
    sandpit: parseConfig(sandpitconfig as ChainConfig, allAssets),
  };

  return configMap[config.toLowerCase()];
}
