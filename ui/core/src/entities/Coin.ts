import { Network } from "./Network";

export function Coin(p: {
  decimals: number;
  imageUrl?: string;
  name: string;
  network: Network;
  symbol: string;
}) {
  return p;
}

export type Coin = ReturnType<typeof Coin>;
