import { ChainId } from "./ChainId";

export function Coin(p: {
  symbol: string;
  decimals: number;
  name: string;
  chainId: ChainId;
}) {
  return p;
}

export type Coin = ReturnType<typeof Coin>;
