import { IWalletService } from "../../services/IWalletService";
import { IAssetAmount } from "../../entities";
export declare function getMockWalletService(state: {
    address: string;
    accounts: string[];
    connected: boolean;
    balances: IAssetAmount[];
    log: string;
}, walletBalances: IAssetAmount[], service?: Partial<IWalletService>): IWalletService;
