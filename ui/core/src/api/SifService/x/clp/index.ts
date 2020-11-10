import { LcdClient, Msg } from "@cosmjs/launchpad";
import { StdTx } from "../../../../entities/noncore/Bank";

export type SwapParams = {
  sent_asset: {
    symbol: string;
    ticker: string;
    source_chain: string;
  };
  received_asset: {
    symbol: string;
    ticker: string;
    source_chain: string;
  };
  base_req: {
    from: string;
    chain_id: string;
  };
  signer: string;
  sent_amount: string;
};

export type LiquidityParams = {
  base_req: {
    from: string;
    chain_id: string;
  };
  external_asset: {
    source_chain: string;
    symbol: string;
    ticker: string;
  };
  native_asset_amount: string;
  external_asset_amount: string;
  signer: string;
};

export type RawPool = {
  external_asset: {
    source_chain: string;
    symbol: string;
    ticker: string;
  };
  native_asset_balance: string;
  external_asset_balance: string;
  pool_units: string;
  pool_address: string;
};

type ClpCmdSwap = (params: SwapParams) => Promise<Msg>;

type ClpQueryPools = () => Promise<RawPool[]>;

type ClpAddLiquidity = (params: LiquidityParams) => Promise<StdTx>;
type ClpCreatePool = (params: LiquidityParams) => Promise<StdTx>;

export interface ClpExtension {
  readonly clp: {
    swap: ClpCmdSwap;
    getPools: ClpQueryPools;
    addLiquidity: ClpAddLiquidity;
    createPool: ClpCreatePool;
  };
}

export function setupClpExtension(base: LcdClient): ClpExtension {
  return {
    clp: {
      getPools: async () => {
        return (await base.get(`/clp/getPools`)).result;
      },

      swap: async (params) => {
        return await base.post(`/clp/swap`, params);
      },

      addLiquidity: async (params) => {
        return await base.post(`/clp/addLiquidity`, params);
      },

      createPool: async (params) => {
        return await base.post(`/clp/createPool`, params);
      },
    },
  };
}
