import { IWalletService } from "src/api";
import { Balance } from "src/entities";

export function getMockWalletService(
  state: {
    address: string;
    accounts: string[];
    connected: boolean;
    log: string;
  },
  walletBalances: Balance[],
  service: Partial<IWalletService> = {}
): IWalletService {
  return {
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
