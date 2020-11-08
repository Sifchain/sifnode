import { LcdClient, Msg } from "@cosmjs/launchpad";

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

type ClpCmdSwap = (params: SwapParams) => Promise<Msg>;
type ClpQueryPools = () => Promise<
  {
    external_asset: {
      source_chain: string;
      symbol: string;
      ticker: string;
    };
    native_asset_balance: string;
    external_asset_balance: string;
    pool_units: string;
    pool_address: string;
  }[]
>;

export interface ClpExtension {
  readonly clp: {
    swap: ClpCmdSwap;
    getPools: ClpQueryPools;
  };
}

export function setupClpExtension(base: LcdClient): ClpExtension {
  return {
    clp: {
      swap: async (params) => {
        return await base.post(`/clp/swap`, params);
      },
      getPools: async () => {
        const response = await base.get(`/clp/getPools`);
        return response.result;
      },
    },
  };
}
