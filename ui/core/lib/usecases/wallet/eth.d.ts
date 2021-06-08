import { UsecaseContext } from "../..";
import { Asset } from "../../entities";
declare const _default: ({ services, store, }: UsecaseContext<"eth" | "bus", "wallet" | "asset">) => {
    isSupportedNetwork(): boolean;
    disconnectWallet(): Promise<void>;
    connectToWallet(): Promise<void>;
    transferEthWallet(amount: number, recipient: string, asset: Asset): Promise<string>;
};
export default _default;
