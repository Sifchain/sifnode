import { ApiContext } from "..";
import { Asset, Coin, Network, Token } from "../../entities";
import { getMetamaskProvider } from "../EthereumService/utils/getMetamaskProvider";

type TokenConfig = {
  symbol: string;
  decimals: number;
  imageUrl?: string;
  name: string;
  address: string;
  network: Network;
};

type CoinConfig = {
  symbol: string;
  decimals: number;
  imageUrl?: string;
  name: string;
  network: Network;
};

export type AssetConfig = CoinConfig | TokenConfig;

function isTokenConfig(a: any): a is TokenConfig {
  return typeof a?.address === "string";
}

function parseAsset(a: unknown): Asset {
  if (isTokenConfig(a)) {
    return Token(a);
  }
  return Coin(a as CoinConfig);
}

export type ChainConfig = {
  sifAddrPrefix: string;
  sifApiUrl: string;
  sifWsUrl: string;
  sifChainId: string;
  web3Provider: "metamask" | string;
  assets: AssetConfig[];
  nativeAsset: string; // symbol
  bridgebankContractAddress: string;
};

export function parseAssets(configAssets: AssetConfig[]): Asset[] {
  return configAssets.map(parseAsset);
}

export function parseConfig(config: ChainConfig, assets: Asset[]): ApiContext {
  const nativeAsset = assets.find((a) => a.symbol === config.nativeAsset);

  if (!nativeAsset)
    throw new Error(
      "No nativeAsset defined for chain config:" + JSON.stringify(config)
    );

  // HACK: Filtering out our testing tokens if not in CI for the demo
  const HACK_blacklist = process.env.CI ? [] : ["atk", "btk", "catk", "cbtk"];
  const HACK_filteredAssets = assets.filter(
    (a) => !HACK_blacklist.includes(a.symbol)
  );

  return {
    sifAddrPrefix: config.sifAddrPrefix,
    sifApiUrl: config.sifApiUrl,
    sifWsUrl: config.sifWsUrl,
    sifChainId: config.sifChainId,
    getWeb3Provider:
      config.web3Provider === "metamask"
        ? getMetamaskProvider
        : async () => config.web3Provider,
    assets: HACK_filteredAssets,
    nativeAsset,
    bridgebankContractAddress: config.bridgebankContractAddress,
  };
}
