import walletActions from "../actions/ethWalletActions";
import { createStore, Store } from "../store";
import { Balance, Token } from "../entities";

import JSBI from "jsbi";
import { USDC, USDT, BNB, CRO, FET } from "../constants/tokens";
import { reactive } from "@vue/reactivity";
import { getMockWalletService } from "./utils/getMockWalletService";

const toBalance = (balance: number) => (tok: Token) =>
  Balance.create(tok, JSBI.BigInt(balance));

describe("queryListOfAvailableTokens", () => {
  describe("updateAvailableTokens", () => {
    let store: Store;
    let walletBalances = [USDC, USDT].map(toBalance(100));

    beforeEach(async () => {
      store = createStore();
      const etheriumState = reactive({
        connected: false,
        address: "",
        log: "",
        balances: [],
        accounts: [],
      });
      const actions = walletActions({
        api: {
          EthereumService: getMockWalletService(etheriumState, walletBalances),
        },
        store,
      });
      await actions.connectToWallet();
    });

    it("should store the available tokens", () => {
      expect(store.wallet.eth.balances).toEqual([
        toBalance(100)(USDC),
        toBalance(100)(USDT),
      ]);
    });

    // TODO: This test should also conver the list of tokens built for the
    // swap dropdowns including the busness logic of what swaps are allowed
  });
});
