import { ApiContext } from "../api";
import { Asset, Coin, Network, Token } from "../entities";
import { getMetamaskProvider } from "../api/EthereumService/utils/getMetamaskProvider";

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

export type KeplrChainConfig = {
  rest: string;
  rpc: string;
  chainId: string;
  chainName: string;
  stakeCurrency: {
    coinDenom: string;
    coinMinimalDenom: string;
    coinDecimals: number;
  };
  bip44: {
    coinType: number;
  };
  bech32Config: {
    bech32PrefixAccAddr: string;
    bech32PrefixAccPub: string;
    bech32PrefixValAddr: string;
    bech32PrefixValPub: string;
    bech32PrefixConsAddr: string;
    bech32PrefixConsPub: string;
  };
  currencies: {
    coinDenom: string;
    coinMinimalDenom: string;
    coinDecimals: number;
  }[];
  feeCurrencies: {
    coinDenom: string;
    coinMinimalDenom: string;
    coinDecimals: number;
  }[];
  coinType: number;
  gasPriceStep: {
    low: number;
    average: number;
    high: number;
  };
};
export type ChainConfig = {
  sifAddrPrefix: string;
  sifApiUrl: string;
  sifWsUrl: string;
  sifRpcUrl: string;
  sifChainId: string;
  web3Provider: "metamask" | string;
  // assets: AssetConfig[];
  nativeAsset: string; // symbol
  bridgebankContractAddress: string;
  keplrChainConfig: KeplrChainConfig;
};

export function parseAssets(configAssets: AssetConfig[]): Asset[] {
  return configAssets.map(parseAsset);
}

export function parseConfig(config: ChainConfig, assets: Asset[]): ApiContext {
  const nativeAsset = assets.find((a) => a.symbol === config.nativeAsset);

  if (!nativeAsset)
    throw new Error(
      "No nativeAsset defined for chain config:" + JSON.stringify(config),
    );

  const bridgetokenContractAddress = (assets.find(
    (token) => token.symbol === "erowan",
  ) as Token).address;

  const sifAssets = assets
    .filter((asset) => asset.network === "sifchain")
    .map((sifAsset) => {
      return {
        coinDenom: sifAsset.symbol,
        coinDecimals: sifAsset.decimals,
        coinMinimalDenom: sifAsset.symbol,
      };
    });

  return {
    sifAddrPrefix: config.sifAddrPrefix,
    sifApiUrl: config.sifApiUrl,
    sifWsUrl: config.sifWsUrl,
    sifRpcUrl: config.sifRpcUrl,
    sifChainId: config.sifChainId,
    getWeb3Provider:
      config.web3Provider === "metamask"
        ? getMetamaskProvider
        : async () => config.web3Provider,
    assets,
    nativeAsset,
    bridgebankContractAddress: config.bridgebankContractAddress,
    bridgetokenContractAddress,
    keplrChainConfig: {
      ...config.keplrChainConfig,
      rest: config.sifApiUrl,
      rpc: config.sifRpcUrl,
      chainId: config.sifChainId,
      currencies: sifAssets,
    },
  };
}
