// TODO - Conditional load or build-time tree shake
import localnetconfig from "./config.localnet.json";
import sandpitconfig from "./config.sandpit.json";
import monkeybarsconfig from "./config.monkeybars.json"
import assetsEthereumLocalnet from "./assets.ethereum.localnet.json";
import assetsEthereumMainnet from "./assets.ethereum.mainnet.json";
import assetsEthereumRopsten from "./assets.ethereum.ropsten.json";
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
  ethereumAssetTag = "ethereum.localnet"
): AppConfig {
  const assetMap: AssetMap = {
    "sifchain.localnet": parseAssets(
      assetsSifchainLocalnet.assets as AssetConfig[]
    ),
    "sifchain.mainnet": parseAssets(
      assetsSifchainMainnet.assets as AssetConfig[]
    ),
    "sifchain.sandpit": parseAssets(
      assetsSifchainSandpit.assets as AssetConfig[]
    ),
    "ethereum.localnet": parseAssets(
      assetsEthereumLocalnet.assets as AssetConfig[]
    ),
    "ethereum.ropsten": parseAssets(
      assetsEthereumRopsten.assets as AssetConfig[]
    ),
    "ethereum.mainnet": parseAssets(
      assetsEthereumMainnet.assets as AssetConfig[]
    ),
  };

  const sifchainAssets = assetMap[sifchainAssetTag];
  const ethereumAssets = assetMap[ethereumAssetTag];
  const allAssets = [...sifchainAssets, ...ethereumAssets].map(cacheAsset);

  const configMap: ConfigMap = {
    localnet: parseConfig(localnetconfig as ChainConfig, allAssets),
    monkeybars: parseConfig(monkeybarsconfig as ChainConfig, allAssets),
    sandpit: parseConfig(sandpitconfig as ChainConfig, allAssets),
  };

  return configMap[config.toLowerCase()];
}
