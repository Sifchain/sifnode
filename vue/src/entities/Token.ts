import { ChainId } from "./ChainId";

// On chain tokens
export type Token = {
  chainId: ChainId;
  address: string;
  decimals: number;
  symbol?: string;
  name?: string;
};
