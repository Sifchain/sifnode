import { UsecaseContext } from "..";
import { PegConfig } from "./index";
export declare const SubscribeToUnconfirmedPegTxs: ({ services, store, config, }: import("../..").WithService<"ethbridge" | "bus"> & import("../..").WithStore<"wallet" | "tx"> & {
    config: PegConfig;
}) => () => () => void;
