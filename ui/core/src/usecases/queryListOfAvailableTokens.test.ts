import queryListOfAvailableTokens from "./queryListOfAvailableTokens";
import { createStore } from "../store";
import { AssetAmount } from "../entities";

import JSBI from "jsbi";
import { USDC, USDT } from "../constants/tokens";

const assetAmounts = [USDC, USDT].map((tok) =>
  AssetAmount.create(tok, JSBI.BigInt(100))
);

describe("queryListOfAvailableTokens", () => {
  describe("updateListOfAvailableTokens", () => {
    const store = createStore({
      marketcapTokenOrder: ["BNB", "USDT", "LINK", "CRO", "USDC"],
    });

    const { state } = store;

    beforeEach(async () => {
      await queryListOfAvailableTokens({
        api: {
          walletService: {
            getAssetBalances: jest.fn(() => Promise.resolve(assetAmounts)),
          },
        },
        store,
        state,
      }).updateListOfAvailableTokens();
    });

    it("should store the tokens in the wallet", () => {
      expect(store.state.userBalances.get("USDC")?.asset).toEqual(USDC);
    });

    it("should not contain other tokens", () => {
      expect(store.state.userBalances.get("ETH")).toBeUndefined();
    });

    it("should deliver tokens in order", () => {
      expect(
        store.state.availableAssetAccounts.map(
          ({ asset: { symbol }, amount }) => ({ symbol, amount })
        )
      ).toEqual([
        {
          symbol: "USDC",
          amount: JSBI.BigInt(100),
        },
        {
          symbol: "USDT",
          amount: JSBI.BigInt(100),
        },
        {
          symbol: "BNB",
          amount: JSBI.BigInt(0),
        },
        {
          symbol: "LINK",
          amount: JSBI.BigInt(0),
        },
        {
          symbol: "CRO",
          amount: JSBI.BigInt(0),
        },
      ]);
    });
  });
});
