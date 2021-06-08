import { IAsset, IAssetAmount } from "../../entities";
import { ActionContext } from "..";
declare const _default: ({ api, store, }: ActionContext<"SifService" | "ClpService" | "EventBusService", "pools" | "wallet" | "accountpools">) => {
    swap(sentAmount: IAssetAmount, receivedAsset: IAsset, minimumReceived: IAssetAmount): Promise<import("../../entities").TransactionStatus>;
    addLiquidity(nativeAssetAmount: IAssetAmount, externalAssetAmount: IAssetAmount): Promise<import("../../entities").TransactionStatus>;
    removeLiquidity(asset: IAsset, wBasisPoints: string, asymmetry: string): Promise<import("../../entities").TransactionStatus>;
    disconnect(): Promise<void>;
};
export default _default;
