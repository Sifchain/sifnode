import { SlashingParameters, ValidatorSignInfo } from "../entities/Slashing";
import { Transaction } from "../entities/Transaction";
export declare const slashingService: {
    getSigningInfos(page: number, limit: number): Promise<ValidatorSignInfo[]>;
    unjailJailedValidtor(validatorAddr: string, body: Transaction): Promise<any>;
    getParameters(): Promise<SlashingParameters>;
};
