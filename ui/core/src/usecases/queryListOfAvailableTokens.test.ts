import queryListOfAvailableTokens from "./queryListOfAvailableTokens";
import { createStore, Store } from "../store";
import { AssetAmount, Token } from "../entities";

import JSBI from "jsbi";
import { USDC, USDT, BNB, CRO, FET } from "../constants/tokens";

const toBalance = (balance: number) => (tok: Token) =>
  AssetAmount.create(tok, JSBI.BigInt(balance));

const assetAmounts = [USDC, USDT].map(toBalance(100));

describe("queryListOfAvailableTokens", () => {
  describe("updateListOfAvailableTokens", () => {
    let store: Store;
    beforeEach(async () => {
      store = createStore();
      await queryListOfAvailableTokens({
        api: {
          walletService: {
            getAssetBalances: jest.fn(() => Promise.resolve(assetAmounts)),
          },
          tokenService: {
            getTopERC20Tokens: jest.fn(() =>
              Promise.resolve([BNB, CRO, USDC, USDT, FET])
            ),
          },
        },
        store,
        state: store.state,
      }).updateAvailableTokens();
    });

    it("should store the available tokens", () => {
      const others = [BNB, CRO, FET].map(toBalance(0));
      expect(store.state.tokenBalances).toEqual([...assetAmounts, ...others]);
    });
  });
});
