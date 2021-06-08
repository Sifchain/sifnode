import { UsecaseContext } from "../..";
import { PegTxEventEmitter } from "../../../services/EthbridgeService/PegTxEventEmitter";
export declare function SubscribeToTx({ services, store, }: UsecaseContext<"bus", "wallet" | "tx">): (tx: PegTxEventEmitter) => () => void;
