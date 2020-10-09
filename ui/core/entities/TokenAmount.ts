import { Token } from "./Token";

export type TokenAmount = { amount: BigInt; token: Token };

export function createTokenAmount(amount: BigInt, token: Token): TokenAmount {
  return {
    amount,
    token,
  };
}

export const amountToToken = (tokenAmount: TokenAmount) => tokenAmount.token;
