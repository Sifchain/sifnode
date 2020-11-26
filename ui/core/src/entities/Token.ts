import { Asset } from "./Asset";
import { Network } from "./Network";

export function Token(p: {
  address: string;
  decimals: number;
  imageUrl?: string;
  name: string;
  network: Network;
  symbol: string;
}) {
  Asset.set(p.symbol, p);
  return p;
}

export type Token = ReturnType<typeof Token>;
