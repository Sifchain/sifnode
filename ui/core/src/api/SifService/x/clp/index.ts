import { LcdClient } from "@cosmjs/launchpad";

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

type ClpCmdSwap = (params: SwapParams) => any;

export interface ClpExtension {
  readonly clp: {
    swap: ClpCmdSwap;
  };
}

export function setupClpExtension(base: LcdClient): ClpExtension {
  return {
    clp: {
      swap: async (params) => {
        return await base.post(`/clp/swap`, params);
      },
    },
  };
}
