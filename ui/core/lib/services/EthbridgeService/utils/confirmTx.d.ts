import Web3 from "web3";
export declare function getConfirmations(web3: Web3, txHash: string): Promise<number>;
export declare function confirmTx({ web3, txHash, confirmations, onSuccess, onCheckConfirmation, }: {
    web3: Web3;
    txHash: string;
    confirmations: number;
    onSuccess?: () => void;
    onCheckConfirmation?: (confirmations: number) => void;
}): void;
