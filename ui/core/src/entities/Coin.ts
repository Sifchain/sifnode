import { Network } from "./Network";

export function Coin(p: {
  symbol: string;
  decimals: number;
  name: string;
  network: Network;
}) {
  return p;
}

export type Coin = ReturnType<typeof Coin>;
