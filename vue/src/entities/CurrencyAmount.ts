import { Currency } from "./Currency";
import { Fraction } from "./Fraction";

export type CurrencyAmount = Fraction & { currency: Currency; amount: BigInt };
