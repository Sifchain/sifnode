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

export type IDispensationServiceContext = {
  nativeAsset: IAsset;
  sifApiUrl: string;
  sifRpcUrl: string;
  sifWsUrl: string;
  sifChainId: string;
  sifUnsignedClient?: SifUnSignedClient;
};

type IDispensationService = {
  claim: (params: { claimType: 2 | 3 }) => any;
};

// TS not null type guard
function notNull<T>(val: T | null): val is T {
  return val !== null;
}

export default function createDispensationService({
  sifApiUrl,
  nativeAsset,
  sifChainId,
  sifWsUrl,
  sifRpcUrl,
  sifUnsignedClient = new SifUnSignedClient(sifApiUrl, sifWsUrl, sifRpcUrl),
}: IDispensationServiceContext): IDispensationService {
  const client = sifUnsignedClient;

  const instance: IDispensationService = {
    async claim(params) {
      return await client.swap({
        base_req: { chain_id: sifChainId, from: params.fromAddress },
        claim_type: params.claimType,
        signer: params.fromAddress,
      });
    },
  };

  return instance;
}
