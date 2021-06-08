import { Address, TxParams } from "../../entities";
import { Mnemonic } from "../../entities/Wallet";
import { UsecaseContext } from "..";
declare const _default: ({ services, store, }: UsecaseContext<"sif" | "clp" | "bus", "wallet">) => {
    getCosmosBalances(address: Address): Promise<import("../../entities").IAssetAmount[]>;
    connect(mnemonic: Mnemonic): Promise<string>;
    sendCosmosTransaction(params: TxParams): Promise<any>;
    disconnect(): Promise<void>;
    connectToWallet(): Promise<void>;
};
export default _default;
