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

type SwapResponse = {
  type: string;
  value: {
    msg: [
      {
        type: string;
        value: {
          Signer: string;
          SentAsset: {
            source_chain: string;
            symbol: string;
            ticker: string;
          };
          ReceivedAsset: {
            source_chain: string;
            symbol: string;
            ticker: string;
          };
          SentAmount: string;
        };
      }
    ];
    fee: { amount: []; gas: string };
    signatures: null;
    memo: string;
  };
};

type ClpCmdSwap = (params: SwapParams) => Promise<Msg>;

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
