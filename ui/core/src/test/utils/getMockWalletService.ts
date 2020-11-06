import { IWalletService } from "../../api/IWalletService";
import { Balance } from "../../entities";

export function getMockWalletService(
  state: {
    address: string;
    accounts: string[];
    connected: boolean;
    balances: Balance[];
    log: string;
  },
  walletBalances: Balance[],
  service: Partial<IWalletService> = {}
): IWalletService {
  return {
    setPhrase: async () => "",
    purgeClient: () => {},
    getState: () => state,
    transfer: async () => "",
    getBalance: jest.fn(async () => walletBalances),
    connect: jest.fn(async () => {
      state.connected = true;
    }),
    disconnect: jest.fn(async () => {
      state.connected = false;
    }),
    isConnected: () => true,
    ...service,
  };
}
