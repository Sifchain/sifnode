import { TendermintBlock, TendermintNodeInfo, TendermintValidatorSet } from "../entities/Tendermint";
export declare const tendermintService: {
    getNodeInfo(): Promise<TendermintNodeInfo>;
    getSyncing(): Promise<boolean>;
    getBlockLatest(): Promise<TendermintBlock>;
    getBlockAtHeight(height: number): Promise<TendermintBlock>;
    getValidatorsetLatest(): Promise<TendermintValidatorSet>;
    getValidatorsetAtHeight(height: number): Promise<TendermintValidatorSet>;
};
