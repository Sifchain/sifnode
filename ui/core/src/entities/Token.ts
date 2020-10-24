import { Asset } from "./Asset";
import { Network } from "./Network";

export function Token(p: {
  symbol: string;
  decimals: number;
  name: string;
  network: Network;
  address: string;
}) {
  Asset.set(p.symbol, p);
  return p;
}

export type Token = ReturnType<typeof Token>;
