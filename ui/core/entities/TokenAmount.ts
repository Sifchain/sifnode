import { Token } from "./Token";
import JSBI from "jsbi";

export type TokenAmount = { amount: JSBI; token: Token };

export function createTokenAmount(amount: JSBI, token: Token): TokenAmount {
  return {
    amount,
    token,
  };
}

export const amountToToken = (tokenAmount: TokenAmount) => tokenAmount.token;
