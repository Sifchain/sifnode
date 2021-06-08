import { ActionContext } from "..";
import { Asset } from "../entities";
declare const _default: ({ api, store, }: ActionContext<"EthereumService" | "EventBusService", "wallet" | "asset">) => {
    isSupportedNetwork(): boolean;
    disconnectWallet(): Promise<void>;
    connectToWallet(): Promise<void>;
    transferEthWallet(amount: number, recipient: string, asset: Asset): Promise<string>;
};
export default _default;
