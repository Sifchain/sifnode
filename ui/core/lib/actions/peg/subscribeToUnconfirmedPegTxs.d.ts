import { ActionContext } from "..";
import { PegConfig } from "./index";
export declare const SubscribeToUnconfirmedPegTxs: ({ api, store, config, }: import("../..").WithApi<"EthbridgeService" | "EventBusService"> & import("../..").WithStore<"wallet" | "tx"> & {
    config: PegConfig;
}) => () => () => void;
