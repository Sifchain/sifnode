import { ApiContext } from "..";
import { Coin, Network, Token } from "../../entities";
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

type AssetConfig = CoinConfig | TokenConfig;
function isTokenConfig(a: any): a is TokenConfig {
  return typeof a?.address === "string";
}

function parseAsset(a: unknown) {
  if (isTokenConfig(a)) {
    return Token(a);
  }
  return Coin(a as CoinConfig);
}

export type ChainConfig = {
  sifAddrPrefix: string;
  sifApiUrl: string;
  sifWsUrl: string;
  web3Provider: "metamask" | string;
  assets: AssetConfig[];
  nativeAsset: string; // symbol
};
export function parseConfig(config: ChainConfig): ApiContext {
  const nativeAsset = config.assets.find(
    (a) => a.symbol === config.nativeAsset
  );

  if (!nativeAsset)
    throw new Error(
      "No nativeAsset defined for chain config:" + JSON.stringify(config)
    );

  return {
    sifAddrPrefix: config.sifAddrPrefix,
    sifApiUrl: config.sifApiUrl,
    sifWsUrl: config.sifWsUrl,
    getWeb3Provider:
      config.web3Provider === "metamask"
        ? getMetamaskProvider
        : async () => config.web3Provider,
    assets: (config.assets as AssetConfig[]).map(parseAsset),
    nativeAsset: parseAsset(nativeAsset),
  };
}
