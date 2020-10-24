import { Network } from "./Network";

export function Token(p: {
  symbol: string;
  decimals: number;
  name: string;
  network: Network;
  address: string;
}) {
  return p;
}

export type Token = ReturnType<typeof Token>;
