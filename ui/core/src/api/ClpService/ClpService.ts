import {
  Asset,
  AssetAmount,
  Coin,
  LiquidityProvider,
  Network,
  Pool,
} from "../../entities";
import { Fraction } from "../../entities/fraction/Fraction";

import { SifUnSignedClient } from "../utils/SifClient";
import { toPool } from "../utils/SifClient/toPool";

export type ClpServiceContext = {
  nativeAsset: Asset;
  sifApiUrl: string;
  sifWsUrl: string;
  sifUnsignedClient?: SifUnSignedClient;
};

type IClpService = {
  getPools: () => Promise<Pool[]>;
  getPoolsByLiquidityProvider: (address: string) => Promise<Pool[]>;
  swap: (params: {
    fromAddress: string;
    receivedAsset: Asset;
    sentAmount: AssetAmount;
  }) => any;
  addLiquidity: (params: {
    fromAddress: string;
    nativeAssetAmount: AssetAmount;
    externalAssetAmount: AssetAmount;
  }) => any;
  createPool: (params: {
    fromAddress: string;
    nativeAssetAmount: AssetAmount;
    externalAssetAmount: AssetAmount;
  }) => any;
  getLiquidityProvider: (params: {
    symbol: string;
    lpAddress: string;
  }) => Promise<LiquidityProvider | null>;
  removeLiquidity: (params: {
    wBasisPoints: string;
    asymmetry: string;
    asset: Asset;
    fromAddress: string;
  }) => any;
};

export default function createClpService({
  sifApiUrl,
  nativeAsset,
  sifWsUrl,
  sifUnsignedClient = new SifUnSignedClient(sifApiUrl, sifWsUrl),
}: ClpServiceContext): IClpService {
  const client = sifUnsignedClient;

  const instance: IClpService = {
    async getPools() {
      try {
        const rawPools = await client.getPools();
        return rawPools.map(toPool(nativeAsset));
      } catch (error) {
        return [];
      }
    },
    async getPoolsByLiquidityProvider(address: string) {
      // Unfortunately it is expensive for the backend to
      // filter pools so we need to annoyingly do this in two calls
      // First we get the metadata
      const poolMeta = await client.getAssets(address);
      if (!poolMeta) return [];
      const poolSymbols = poolMeta.map(({ symbol }) => symbol);

      // Then we get the pools and filter for the metadata
      const rawPools = await client.getPools();
      return rawPools
        .filter((rawPool) =>
          poolSymbols.includes(rawPool.external_asset.symbol)
        )
        .map(toPool(nativeAsset));
    },

    async addLiquidity(params: {
      fromAddress: string;
      nativeAssetAmount: AssetAmount;
      externalAssetAmount: AssetAmount;
    }) {
      return await client.addLiquidity({
        base_req: { chain_id: "sifnode", from: params.fromAddress },
        external_asset: {
          source_chain: params.externalAssetAmount.asset.network as string,
          symbol: params.externalAssetAmount.asset.symbol,
          ticker: params.externalAssetAmount.asset.symbol,
        },
        external_asset_amount: params.externalAssetAmount.toFixed(0),
        native_asset_amount: params.nativeAssetAmount.toFixed(0),
        signer: params.fromAddress,
      });
    },

    async createPool(params) {
      return await client.createPool({
        base_req: { chain_id: "sifnode", from: params.fromAddress },
        external_asset: {
          source_chain: params.externalAssetAmount.asset.network as string,
          symbol: params.externalAssetAmount.asset.symbol,
          ticker: params.externalAssetAmount.asset.symbol,
        },
        external_asset_amount: params.externalAssetAmount.toFixed(0),
        native_asset_amount: params.nativeAssetAmount.toFixed(0),
        signer: params.fromAddress,
      });
    },

    async swap(params) {
      return await client.swap({
        base_req: { chain_id: "sifchain", from: params.fromAddress },
        received_asset: {
          source_chain: params.receivedAsset.network as string,
          symbol: params.receivedAsset.symbol,
          ticker: params.receivedAsset.symbol,
        },
        sent_amount: params.sentAmount.numerator.toString(),
        sent_asset: {
          source_chain: params.sentAmount.asset.network as string,
          symbol: params.sentAmount.asset.symbol,
          ticker: params.sentAmount.asset.symbol,
        },
        signer: params.fromAddress,
      });
    },
    async getLiquidityProvider(params) {
      const response = await client.getLiquidityProvider(params);
      return LiquidityProvider(
        Coin({
          name: response.result.LiquidityProvider.asset.symbol,
          symbol: response.result.LiquidityProvider.asset.symbol,
          network: Network.SIFCHAIN,
          decimals: 18,
        }),
        new Fraction(
          response.result.LiquidityProvider.liquidity_provider_units
        ),
        response.result.LiquidityProvider.liquidity_provider_address
      );
    },

    async removeLiquidity(params) {
      return await client.removeLiquidity({
        asymmetry: params.asymmetry,
        base_req: { chain_id: "sifchain", from: params.fromAddress },
        external_asset: {
          source_chain: params.asset.network as string,
          symbol: params.asset.symbol,
          ticker: params.asset.symbol,
        },
        signer: params.fromAddress,
        w_basis_points: params.wBasisPoints,
      });
    },
  };

  return instance;
}
