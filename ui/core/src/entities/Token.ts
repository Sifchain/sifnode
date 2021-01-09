import { Network } from "./Network";

export function Token(p: {
  address: string;
  decimals: number;
  imageUrl?: string;
  name: string;
  network: Network;
  symbol: string;
}) {
  return p;
}

export type Token = ReturnType<typeof Token>;
