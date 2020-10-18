import walletActions from "../actions/walletActions";
import { createStore, Store } from "../store";
import { Balance, Token } from "../entities";

import JSBI from "jsbi";
import { USDC, USDT, BNB, CRO, FET } from "../constants/tokens";
import { reactive } from "@vue/reactivity";

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
        accounts: [],
      });
      const actions = walletActions({
        api: {
          EtheriumService: {
            getState: () => etheriumState,
            transfer: async () => "",
            getBalance: jest.fn(async () => walletBalances),
            connect: jest.fn(async () => {
              etheriumState.connected = true;
            }),
            disconnect: jest.fn(async () => {
              etheriumState.connected = false;
            }),
            isConnected: () => true,
          },
        },
        store,
      });
      await actions.connectToWallet();
    });

    it("should store the available tokens", () => {
      expect(store.wallet.balances).toEqual([
        toBalance(100)(USDC),
        toBalance(100)(USDT),
      ]);
    });
  });
});
