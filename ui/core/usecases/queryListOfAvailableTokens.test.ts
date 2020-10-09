import queryListOfAvailableTokens from "./queryListOfAvailableTokens";
import { createStore } from "../store";
import { createAssetAmount } from "../entities";
import { USDC, USDT } from "../constants/tokens";

const assetAmounts = [USDC, USDT].map((tok) => createAssetAmount(tok, 100n));

describe("queryListOfAvailableTokens", () => {
  describe("updateListOfAvailableTokens", () => {
    const store = createStore({
      marketcapTokenOrder: ["BNB", "USDT", "LINK", "CRO", "USDC"],
    });

    beforeEach(async () => {
      await queryListOfAvailableTokens({
        api: {
          walletService: {
            getAssetBalances: jest.fn(() => Promise.resolve(assetAmounts)),
          },
        },
        store,
        state: store.state,
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
          amount: 100n,
        },
        {
          symbol: "USDT",
          amount: 100n,
        },
        {
          symbol: "BNB",
          amount: 0n,
        },
        {
          symbol: "LINK",
          amount: 0n,
        },
        {
          symbol: "CRO",
          amount: 0n,
        },
      ]);
    });
  });
});
