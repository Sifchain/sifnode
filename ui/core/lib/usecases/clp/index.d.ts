import { IAsset, IAssetAmount } from "../../entities";
import { UsecaseContext } from "..";
declare const _default: ({ services, store, }: UsecaseContext<"sif" | "clp" | "bus", "pools" | "wallet" | "accountpools">) => {
    swap(sentAmount: IAssetAmount, receivedAsset: IAsset, minimumReceived: IAssetAmount): Promise<import("../../entities").TransactionStatus>;
    addLiquidity(nativeAssetAmount: IAssetAmount, externalAssetAmount: IAssetAmount): Promise<import("../../entities").TransactionStatus>;
    removeLiquidity(asset: IAsset, wBasisPoints: string, asymmetry: string): Promise<import("../../entities").TransactionStatus>;
    disconnect(): Promise<void>;
};
export default _default;
