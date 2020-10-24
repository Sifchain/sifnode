import { Asset } from "./Asset";
import { Network } from "./Network";

export function Coin(p: {
  symbol: string;
  decimals: number;
  name: string;
  network: Network;
}) {
  Asset.set(p.symbol, p);
  return p;
}

export type Coin = ReturnType<typeof Coin>;
