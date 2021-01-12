import { LcdClient, Msg } from "@cosmjs/launchpad";

type BaseReq = {
  from: string;
  chain_id: string;
  account_number?: string;
  sequence?: string;
};

type BurnOrLockReq = {
  base_req: BaseReq;
  ethereum_chain_id: string;
  token_contract_address: string;
  cosmos_sender: string;
  ethereum_receiver: string;
  amount: string;
  symbol: string;
};

export interface EthbridgeExtension {
  readonly ethbridge: {
    burn: (params: BurnOrLockReq) => Promise<Msg>;
  };
}

export function setupEthbridgeExtension(base: LcdClient): EthbridgeExtension {
  return {
    ethbridge: {
      burn: async (params) => {
        return await base.post(`/ethbridge/burn`, params);
      },
    },
  };
}
