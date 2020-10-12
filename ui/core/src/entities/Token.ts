import { ChainId } from "./ChainId";
import { Asset } from "./Asset";

export type Token = Asset & {
  chainId: ChainId;
  address: string;
};

export function createToken(
  chainId: ChainId,
  address: string,
  decimals: number,
  symbol: string,
  name: string
): Token {
  return {
    chainId,
    address,
    decimals,
    symbol,
    name,
  };
}
