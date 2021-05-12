import { LcdClient, Msg } from "@cosmjs/launchpad";

type BaseReq = {
  from: string;
  chain_id: string;
  account_number?: string;
  sequence?: string;
};

/*
    type DistributionType int64
    const Airdrop DistributionType = 1
    const LiquidityMining DistributionType = 2
    const ValidatorSubsidy DistributionType = 3

    type CreateClaimReq struct {
        BaseReq   rest.BaseReq     `json:"base_req"`
        Signer    string           `json:"signer"`
        ClaimType DistributionType `json:"claim_type"`   
    }
*/

type ICreateClaim = {
  base_req: BaseReq;
  claim_type: 2 | 3;
};

export interface IDispensationApi {
  readonly dispensation: {
    claim: (params: ICreateClaim) => Promise<Msg>;
  };
}

export function setupDispensationApi(base: LcdClient): IDispensationApi {
  return {
    dispensation: {
      claim: async (params) => {
        console.log(
          `/dispensation/createClaim`,
          JSON.stringify(params, null, 2),
        );
        return await base.post(`/dispensation/createClaim`, params);
      },
    },
  };
}
