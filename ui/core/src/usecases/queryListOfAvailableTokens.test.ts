import queryListOfAvailableTokens from "./queryListOfAvailableTokens";
import { createStore, Store } from "../store";
import { AssetAmount, ChainId, createToken, Token } from "../entities";

import JSBI from "jsbi";
import { USDC, USDT, BNB, CRO, FET } from "../constants/tokens";

const toBalance = (balance: number) => (tok: Token) =>
  AssetAmount.create(tok, JSBI.BigInt(balance));

describe("queryListOfAvailableTokens", () => {
  describe("updateAvailableTokens", () => {
    let store: Store;
    let walletBalances = [USDC, USDT].map(toBalance(100));

    beforeEach(async () => {
      store = createStore();
      await queryListOfAvailableTokens({
        api: {
          walletService: {
            getAssetBalances: jest.fn(() => Promise.resolve(walletBalances)),
          },
          tokenService: {
            getTopERC20Tokens: jest.fn(() =>
              Promise.resolve([
                BNB,
                CRO,
                // test diff token instance with same symbol
                createToken(
                  ChainId.ETH_MAINNET,
                  "some address",
                  6,
                  "USDC",
                  "USDC"
                ),
                USDT,
                FET,
              ])
            ),
          },
        },
        store,
        state: store.state,
      }).updateAvailableTokens();
    });

    it("should store the available tokens", () => {
      // const others = [BNB, CRO, FET].map(toBalance(0));
      expect(store.state.tokenBalances).toEqual([
        toBalance(100)(USDC),
        toBalance(100)(USDT),
        toBalance(0)(BNB),
        toBalance(0)(CRO),
        toBalance(0)(FET),
      ]);
    });
  });
});
