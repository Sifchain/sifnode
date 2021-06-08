import { provider } from "web3-core";
import { IAssetAmount } from "../../entities";
import { PegTxEventEmitter } from "./PegTxEventEmitter";
import { SifUnSignedClient } from "../utils/SifClient";
export declare type EthbridgeServiceContext = {
    sifApiUrl: string;
    sifWsUrl: string;
    sifRpcUrl: string;
    sifChainId: string;
    bridgebankContractAddress: string;
    bridgetokenContractAddress: string;
    getWeb3Provider: () => Promise<provider>;
    sifUnsignedClient?: SifUnSignedClient;
};
export default function createEthbridgeService({ sifApiUrl, sifWsUrl, sifRpcUrl, sifChainId, bridgebankContractAddress, getWeb3Provider, sifUnsignedClient, }: EthbridgeServiceContext): {
    approveBridgeBankSpend(account: string, amount: IAssetAmount): Promise<any>;
    burnToEthereum(params: {
        fromAddress: string;
        ethereumRecipient: string;
        assetAmount: IAssetAmount;
        feeAmount: IAssetAmount;
    }): Promise<import("@cosmjs/launchpad").Msg>;
    lockToSifchain(sifRecipient: string, assetAmount: IAssetAmount, confirmations: number): {
        readonly hash: string | undefined;
        readonly symbol: string | undefined;
        setTxHash(hash: string): void;
        emit(e: import("./types").TxEventPrepopulated<import("./types").TxEvent>): void;
        onTxEvent(handler: (e: import("./types").TxEvent) => void): any;
        onTxHash(handler: (e: import("./types").TxEventHashReceived) => void): any;
        onEthConfCountChanged(handler: (e: import("./types").TxEventEthConfCountChanged) => void): any;
        onEthTxConfirmed(handler: (e: import("./types").TxEventEthTxConfirmed) => void): any;
        onSifTxConfirmed(handler: (e: import("./types").TxEventSifTxConfirmed) => void): any;
        onEthTxInitiated(handler: (e: import("./types").TxEventEthTxInitiated) => void): any;
        onSifTxInitiated(handler: (e: import("./types").TxEventSifTxInitiated) => void): any;
        onComplete(handler: (e: import("./types").TxEventComplete) => void): any;
        removeListeners(): any;
        onError(handler: (e: import("./types").TxEventError) => void): any;
    };
    lockToEthereum(params: {
        fromAddress: string;
        ethereumRecipient: string;
        assetAmount: IAssetAmount;
        feeAmount: IAssetAmount;
    }): Promise<import("@cosmjs/launchpad").Msg>;
    /**
     * Get a list of unconfirmed transaction hashes associated with
     * a particular address and return pegTxs associated with that hash
     * @param address contract address
     * @param confirmations number of confirmations required
     */
    fetchUnconfirmedLockBurnTxs(address: string, confirmations: number): Promise<PegTxEventEmitter[]>;
    burnToSifchain(sifRecipient: string, assetAmount: IAssetAmount, confirmations: number, account?: string | undefined): {
        readonly hash: string | undefined;
        readonly symbol: string | undefined;
        setTxHash(hash: string): void;
        emit(e: import("./types").TxEventPrepopulated<import("./types").TxEvent>): void;
        onTxEvent(handler: (e: import("./types").TxEvent) => void): any;
        onTxHash(handler: (e: import("./types").TxEventHashReceived) => void): any;
        onEthConfCountChanged(handler: (e: import("./types").TxEventEthConfCountChanged) => void): any;
        onEthTxConfirmed(handler: (e: import("./types").TxEventEthTxConfirmed) => void): any;
        onSifTxConfirmed(handler: (e: import("./types").TxEventSifTxConfirmed) => void): any;
        onEthTxInitiated(handler: (e: import("./types").TxEventEthTxInitiated) => void): any;
        onSifTxInitiated(handler: (e: import("./types").TxEventSifTxInitiated) => void): any;
        onComplete(handler: (e: import("./types").TxEventComplete) => void): any;
        removeListeners(): any;
        onError(handler: (e: import("./types").TxEventError) => void): any;
    };
};
