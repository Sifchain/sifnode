import createActions from "./ethWallet";
import { IWalletService } from "../api/IWalletService";
import { Address, Asset, Network, Token, TxParams } from "../entities";
import { Msg } from "@cosmjs/launchpad";
let mockEthereumService: IWalletService;
let ethWalletActions: ReturnType<typeof createActions>;

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
  ethWalletActions = createActions({
    api: { EthereumService: mockEthereumService },
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
    Token({
      name: "Ethereum",
      network: Network.SIFCHAIN,
      address: "abcdefg",
      decimals: 18,
      symbol: "ceth",
    })
  );
  expect(mockEthereumService.transfer).toHaveBeenCalled();
});
