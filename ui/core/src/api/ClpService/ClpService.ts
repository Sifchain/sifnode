import {
  Amount,
  Asset,
  IAsset,
  IAssetAmount,
  LiquidityProvider,
  Network,
  Pool,
} from "../../entities";

import { SifUnSignedClient } from "../utils/SifClient";
import { toPool } from "../utils/SifClient/toPool";

export type ClpServiceContext = {
  nativeAsset: IAsset;
  sifApiUrl: string;
  sifRpcUrl: string;
  sifWsUrl: string;
  sifChainId: string;
  sifUnsignedClient?: SifUnSignedClient;
};

type IClpService = {
  getPools: () => Promise<Pool[]>;
  getPoolSymbolsByLiquidityProvider: (address: string) => Promise<string[]>;
  swap: (params: {
    fromAddress: string;
    sentAmount: IAssetAmount;
    receivedAsset: Asset;
    minimumReceived: IAssetAmount;
  }) => any;
  addLiquidity: (params: {
    fromAddress: string;
    nativeAssetAmount: IAssetAmount;
    externalAssetAmount: IAssetAmount;
  }) => any;
  createPool: (params: {
    fromAddress: string;
    nativeAssetAmount: IAssetAmount;
    externalAssetAmount: IAssetAmount;
  }) => any;
  getLiquidityProvider: (params: {
    symbol: string;
    lpAddress: string;
  }) => Promise<LiquidityProvider | null>;
  removeLiquidity: (params: {
    wBasisPoints: string;
    asymmetry: string;
    asset: IAsset;
    fromAddress: string;
  }) => any;
  claimRewards: (params: { address: string }) => any;
};

// TS not null type guard
function notNull<T>(val: T | null): val is T {
  return val !== null;
}

export default function createClpService({
  sifApiUrl,
  nativeAsset,
  sifChainId,
  sifWsUrl,
  sifRpcUrl,
  sifUnsignedClient = new SifUnSignedClient(sifApiUrl, sifWsUrl, sifRpcUrl),
}: ClpServiceContext): IClpService {
  const client = sifUnsignedClient;

  const instance: IClpService = {
    async getPools() {
      try {
        const rawPools = await client.getPools();
        return (
          rawPools
            .map(toPool(nativeAsset))
            // toPool can return a null pool for invalid pools lets filter them out
            .filter(notNull)
        );
      } catch (error) {
        return [];
      }
    },
    async getPoolSymbolsByLiquidityProvider(address: string) {
      // Unfortunately it is expensive for the backend to
      // filter pools so we need to annoyingly do this in two calls
      // First we get the metadata
      const poolMeta = await client.getAssets(address);
      if (!poolMeta) return [];
      return poolMeta.map(({ symbol }) => symbol);
    },

    async addLiquidity(params: {
      fromAddress: string;
      nativeAssetAmount: IAssetAmount;
      externalAssetAmount: IAssetAmount;
    }) {
      return await client.addLiquidity({
        base_req: { chain_id: sifChainId, from: params.fromAddress },
        external_asset: {
          source_chain: params.externalAssetAmount.asset.network as string,
          symbol: params.externalAssetAmount.asset.symbol,
          ticker: params.externalAssetAmount.asset.symbol,
        },
        external_asset_amount: params.externalAssetAmount.toBigInt().toString(),
        native_asset_amount: params.nativeAssetAmount.toBigInt().toString(),
        signer: params.fromAddress,
      });
    },

    async createPool(params) {
      return await client.createPool({
        base_req: { chain_id: sifChainId, from: params.fromAddress },
        external_asset: {
          source_chain: params.externalAssetAmount.asset.network as string,
          symbol: params.externalAssetAmount.asset.symbol,
          ticker: params.externalAssetAmount.asset.symbol,
        },
        external_asset_amount: params.externalAssetAmount.toBigInt().toString(),
        native_asset_amount: params.nativeAssetAmount.toBigInt().toString(),
        signer: params.fromAddress,
      });
    },

    async swap(params) {
      return await client.swap({
        base_req: { chain_id: sifChainId, from: params.fromAddress },
        received_asset: {
          source_chain: params.receivedAsset.network as string,
          symbol: params.receivedAsset.symbol,
          ticker: params.receivedAsset.symbol,
        },
        sent_amount: params.sentAmount.toBigInt().toString(),
        sent_asset: {
          source_chain: params.sentAmount.asset.network as string,
          symbol: params.sentAmount.asset.symbol,
          ticker: params.sentAmount.asset.symbol,
        },
        min_receiving_amount: params.minimumReceived.toBigInt().toString(),
        signer: params.fromAddress,
      });
    },
    async getLiquidityProvider(params) {
      const response = await client.getLiquidityProvider(params);
      let asset: IAsset;
      const {
        LiquidityProvider: liquidityProvider,
        native_asset_balance,
        external_asset_balance,
      } = response.result;
      const {
        asset: { symbol },
        liquidity_provider_units,
        liquidity_provider_address,
      } = liquidityProvider;

      try {
        asset = Asset(symbol);
      } catch (err) {
        asset = Asset({
          name: symbol,
          label: symbol,
          symbol,
          network: Network.SIFCHAIN,
          decimals: 18,
        });
      }

      return LiquidityProvider(
        asset,
        Amount(liquidity_provider_units),
        liquidity_provider_address,
        Amount(native_asset_balance),
        Amount(external_asset_balance),
      );
    },

    async removeLiquidity(params) {
      return await client.removeLiquidity({
        asymmetry: params.asymmetry,
        base_req: { chain_id: sifChainId, from: params.fromAddress },
        external_asset: {
          source_chain: params.asset.network as string,
          symbol: params.asset.symbol,
          ticker: params.asset.symbol,
        },
        signer: params.fromAddress,
        w_basis_points: params.wBasisPoints,
      });
    },

    async claimRewards(params) {
      const rewards = await client.claimRewards({ address: "" });
      console.log("rewards", rewards);
      return rewards;
    },
  };

  return instance;
}
