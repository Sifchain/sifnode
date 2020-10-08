import { CurrencyAmount, Currency, Token, TokenAmount } from "../entities";

type CurrencyBalances = {
  [address: string]: CurrencyAmount | undefined;
};

type TokenBalances = {
  [address: string]: TokenAmount | undefined;
};

export const walletService = {
  async getCurrencyBalances(
    account: string,
    currencies: Currency[]
  ): Promise<CurrencyBalances> {
    return {};
  },

  async getTokenBalances(
    account: string,
    currencies: Token[]
  ): Promise<TokenBalances> {
    return {};
  },
};
