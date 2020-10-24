import { ChainId } from "./ChainId";

export function Token(p: {
  symbol: string;
  decimals: number;
  name: string;
  chainId: ChainId;
  address: string;
}) {
  return p;
}

export type Token = ReturnType<typeof Token>;
