import createActions from "./ethWallet";

import { Address, Asset, Network, TxParams } from "../entities";
import { Msg } from "@cosmjs/launchpad";
import { IWalletService } from "../api/IWalletService";
let mockEthereumService: IWalletService & {};
let mockEventBusService: any;
let ethWalletActions: ReturnType<typeof createActions>;
let dispatch = jest.fn();

beforeEach(() => {
  mockEthereumService = {
    getState: () => ({
      address: "",
      accounts: [],
      connected: true,
      balances: [],
      log: "",
    }),
    getSupportedTokens: () => [],
    isConnected: () => true,
    connect: async () => {},
    disconnect: jest.fn(async () => {}),
    transfer: jest.fn(async (params: TxParams) => ""),
    getBalance: async (address: Address, asset?: Asset) => [],
    signAndBroadcast: async (msg: Msg, memo?: string) => {},
    setPhrase: async (phrase: string) => "",
    purgeClient: () => {},
    onProviderNotFound: () => {},
  };

  mockEventBusService = {
    dispatch: dispatch,
  };

  ethWalletActions = createActions({
    api: {
      EthereumService: mockEthereumService,
      EventBusService: mockEventBusService,
    },
    store: {
      asset: { topTokens: [] },
      wallet: {
        eth: {
          balances: [],
          address: "",
          isConnected: true,
        },
        sif: {
          balances: [],
          address: "",
          isConnected: true,
        },
      },
    },
  });
});

test("Calls transfer correctly", async () => {
  await ethWalletActions.transferEthWallet(
    123,
    "abcdef",
    Asset({
      name: "Ethereum",
      label: "ETH",
      network: Network.SIFCHAIN,
      address: "abcdefg",
      decimals: 18,
      symbol: "ceth",
    }),
  );
  expect(mockEthereumService.transfer).toHaveBeenCalled();
});
