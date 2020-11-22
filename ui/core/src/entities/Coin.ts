import { Asset } from "./Asset";
import { Network } from "./Network";

export function Coin(p: {
  decimals: number;
  imageUrl?: string;
  name: string;
  network: Network;
  symbol: string;
}) {
  Asset.set(p.symbol, p);
  return p;
}

export type Coin = ReturnType<typeof Coin>;
