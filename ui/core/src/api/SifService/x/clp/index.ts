import { LcdClient, Msg } from "@cosmjs/launchpad";
import { AssetAmount } from "../../../../entities";
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

type ClpAddLiquidity = (params: {
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
}) => Promise<StdTx>;

export interface ClpExtension {
  readonly clp: {
    swap: ClpCmdSwap;
    getPools: ClpQueryPools;
    addLiquidity: ClpAddLiquidity;
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
      addLiquidity: async (params) => {
        return await base.post(`/clp/addLiquidity`, params);
      },
    },
  };
}
