import { CurrencyAmount } from "./CurrencyAmount";
import { Token } from "./Token";

export type TokenAmount = CurrencyAmount & { token: Token; amount: BigInt };
