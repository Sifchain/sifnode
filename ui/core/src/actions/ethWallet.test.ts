import createActions from "./ethWallet";
import { IWalletService } from "../api/IWalletService";
import { Address, Asset, Network, Token, TxParams } from "../entities";
import { Msg } from "@cosmjs/launchpad";
let mockEthereumService: IWalletService;
let mockNotificationsService: any;
let ethWalletActions: ReturnType<typeof createActions>;
let notify = jest.fn();

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
    getBalance: async (address: Address, asset?: Asset | Token) => [],
    signAndBroadcast: async (msg: Msg, memo?: string) => {},
    setPhrase: async (phrase: string) => "",
    purgeClient: () => {},
  };

  mockNotificationsService = {
    notify: notify,
  };

  ethWalletActions = createActions({
    api: {
      EthereumService: mockEthereumService,
      NotificationService: mockNotificationsService,
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

test("Calls disconnect", async () => {
  await ethWalletActions.disconnectWallet();
  expect(mockEthereumService.disconnect).toHaveBeenCalled();
  expect(notify).toHaveBeenCalled();
});

test("Calls transfer correctly", async () => {
  await ethWalletActions.transferEthWallet(
    123,
    "abcdef",
    Token({
      name: "Ethereum",
      network: Network.SIFCHAIN,
      address: "abcdefg",
      decimals: 18,
      symbol: "ceth",
    })
  );
  expect(mockEthereumService.transfer).toHaveBeenCalled();
  expect(notify).toHaveBeenCalled();
});
