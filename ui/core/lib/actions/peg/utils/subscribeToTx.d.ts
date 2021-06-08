import { ActionContext } from "../..";
import { PegTxEventEmitter } from "../../../api/EthbridgeService/PegTxEventEmitter";
export declare function SubscribeToTx({ api, store, }: ActionContext<"EventBusService", "wallet" | "tx">): (tx: PegTxEventEmitter) => () => void;
