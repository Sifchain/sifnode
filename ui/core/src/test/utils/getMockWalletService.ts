import { IWalletService } from "../../api/IWalletService";
import { AssetAmount } from "../../entities";

export function getMockWalletService(
  state: {
    address: string;
    accounts: string[];
    connected: boolean;
    balances: AssetAmount[];
    log: string;
  },
  walletBalances: AssetAmount[],
  service: Partial<IWalletService> = {}
): IWalletService {
  return {
    setPhrase: async () => "",
    purgeClient: () => {},
    getState: () => state,
    transfer: async () => "",
    getBalance: jest.fn(async () => walletBalances),
    getSupportedTokens: () => [],
    connect: jest.fn(async () => {
      state.connected = true;
      state.balances = walletBalances;
    }),
    disconnect: jest.fn(async () => {
      state.connected = false;
    }),
    isConnected: () => true,
    ...service,
    signAndBroadcast: async (
      msg: { type: string; value: any },
      memo?: string
    ) => {},
    onProviderNotFound: () => {},
  };
}
