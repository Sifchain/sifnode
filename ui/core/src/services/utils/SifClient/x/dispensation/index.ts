import { LcdClient, Msg } from "@cosmjs/launchpad";

type BaseReq = {
  from: string;
  chain_id: string;
  account_number?: string;
  sequence?: string;
};

type IClaimParams = {
  base_req: BaseReq;
  claim_type: "2" | "3";
  signer: string;
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

export interface DispensationExtension {
  readonly dispensation: {
    claim: (params: IClaimParams) => Promise<Msg>;
  };
}

export function setupDispensationExtension(
  base: LcdClient,
): DispensationExtension {
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
