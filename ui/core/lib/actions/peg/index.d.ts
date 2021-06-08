import { ActionContext } from "..";
import { Address, IAsset, IAssetAmount, TransactionStatus } from "../../entities";
/**
 * Shared peg config for use throughout the peg feature
 */
export declare type PegConfig = {
    ethConfirmations: number;
};
declare const _default: ({ api, store, }: ActionContext<"SifService" | "EthbridgeService" | "EventBusService" | "EthereumService", // Select the services you want to access
// Select the services you want to access
"wallet" | "tx">) => {
    subscribeToUnconfirmedPegTxs: () => () => void;
    getSifTokens(): IAsset[];
    getEthTokens(): IAsset[];
    calculateUnpegFee(asset: IAsset): IAssetAmount;
    unpeg(assetAmount: IAssetAmount): Promise<TransactionStatus>;
    approve(address: Address, assetAmount: IAssetAmount): Promise<any>;
    peg(assetAmount: IAssetAmount): Promise<TransactionStatus>;
};
export default _default;
