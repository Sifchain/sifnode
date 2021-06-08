import { LcdClient, Msg } from "@cosmjs/launchpad";
declare type BaseReq = {
    from: string;
    chain_id: string;
    account_number?: string;
    sequence?: string;
};
declare type BurnOrLockReq = {
    base_req: BaseReq;
    ethereum_chain_id: string;
    token_contract_address: string;
    cosmos_sender: string;
    ethereum_receiver: string;
    amount: string;
    symbol: string;
    ceth_amount: string;
};
export interface EthbridgeExtension {
    readonly ethbridge: {
        burn: (params: BurnOrLockReq) => Promise<Msg>;
        lock: (params: BurnOrLockReq) => Promise<Msg>;
    };
}
export declare function setupEthbridgeExtension(base: LcdClient): EthbridgeExtension;
export {};
