import { Token, TokenAmount } from "../entities";
import { Context } from ".";
import { BigintIsh } from "src/entities/fraction/Fraction";
import JSBI from "jsbi";

/*
Drop down list of tokens
appears with the following list
of tokens:

1. Top 10 tokens from
users wallet with
corresponding amounts
from their wallet

2. Top 20 ERC-20 tokens
User sees a search bar where
they can type their ERC-20
token if it’s not listed.

*/

function toTokenAmount(amount: BigintIsh) {
  return (token: Token) => {
    return TokenAmount.create(token, amount);
  };
}

export default ({ api, store }: Context<"walletService" | "tokenService">) => ({
  async updateAvailableTokens() {
    const walletBalances = await api.walletService.getAssetBalances({
      limit: 10,
    });

    const topERCTokens = await api.tokenService.getTopERC20Tokens({
      limit: 20,
    });

    const walletTokens = walletBalances.map((assetAmount) => assetAmount.asset);

    const availableEmptyTokenAccounts = topERCTokens
      .filter((token) => !walletTokens.includes(token))
      .map(toTokenAmount(JSBI.BigInt("0")));

    store.setTokenBalances([...walletBalances, ...availableEmptyTokenAccounts]);
  },
});
