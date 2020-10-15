import walletActions from "./walletActions";
import { createStore, Store } from "../store";
import { Balance, Token } from "../entities";

import JSBI from "jsbi";
import { USDC, USDT, BNB, CRO, FET } from "../constants/tokens";

const toBalance = (balance: number) => (tok: Token) =>
  Balance.create(tok, JSBI.BigInt(balance));

describe("queryListOfAvailableTokens", () => {
  describe("updateAvailableTokens", () => {
    let store: Store;
    let walletBalances = [USDC, USDT].map(toBalance(100));

    beforeEach(async () => {
      store = createStore();
      await walletActions({
        api: {
          walletService: {
            getBalance: jest.fn(() => Promise.resolve(walletBalances)),
            connect: jest.fn(() => Promise.resolve(true)),
            disconnect: jest.fn(() => Promise.resolve()),
            isConnected: () => true,
          },
        },
        store,
      }).refreshWalletBalances();
    });

    it("should store the available tokens", () => {
      expect(store.wallet.balances).toEqual([
        toBalance(100)(USDC),
        toBalance(100)(USDT),
      ]);
    });
  });
});
