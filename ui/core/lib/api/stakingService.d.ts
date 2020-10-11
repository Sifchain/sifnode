import { Delegation, Redelegation, StakingParameters, StakingPool, UnbondingDelegation, Validator } from "../entities/Staking";
export declare const stakingService: {
    getDelegatorDelegations(delegatorAddr: string): Promise<Delegation[]>;
    getDelegatorCurrentDelegation(delegatorAddr: string, validatorAddr: string): Promise<Delegation>;
    getDelegatorUnbondingDelegations(delegatorAddr: string): Promise<UnbondingDelegation[]>;
    getDelegatorUnbondingDelegationsBtwValidator(delegatorAddr: string, validatorAddr: string): Promise<UnbondingDelegation[]>;
    getRedelegations(delegator: String, validatorFrom: String, validatorTo: string): Promise<Redelegation[]>;
    getDelegatorValidators(delegatorAddr: string): Promise<Validator[]>;
    getDelegatorValidator(delegatorAddr: string, validatorAddr: string): Promise<Validator>;
    getValidators(status: string, page: number, limit: number): Promise<Validator[]>;
    getValidator(validatorAddr: string): Promise<Validator>;
    getValidatorDelegations(validatorAddr: string): Promise<Delegation[]>;
    getValidatorUnbondingDelegations(validatorAddr: string): Promise<UnbondingDelegation[]>;
    getPool(): Promise<StakingPool>;
    getParameters(): Promise<StakingParameters>;
};
