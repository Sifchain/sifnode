// based on draft swagger spec https://raw.githubusercontent.com/Sifchain/sifnode/c1bb5a268da8b519d0fc90f81fa194d31c0f82b3/api/openapi/swagger.yml?token=AAJSXWM6CDXYAEETSC6BJ2S7Q2JLS

import {
  Delegation,
  Redelegation,
  StakingParameters,
  StakingPool,
  UnbondingDelegation,
  Validator,
} from "../entities/Staking";

export const stakingService = {
  async getDelegatorDelegations(delegatorAddr: string): Promise<Delegation[]> {
    return [];
  },

  async getDelegatorCurrentDelegation(
    delegatorAddr: string,
    validatorAddr: string
  ): Promise<Delegation> {},

  async getDelegatorUnbondingDelegations(
    delegatorAddr: string
  ): Promise<UnbondingDelegation[]> {
    return [];
  },

  async getDelegatorUnbondingDelegationsBtwValidator(
    delegatorAddr: string,
    validatorAddr: string
  ): Promise<UnbondingDelegation[]> {
    return [];
  },

  async getRedelegations(
    delegator: String,
    validatorFrom: String,
    validatorTo: string
  ): Promise<Redelegation[]> {
    return [];
  },

  async getDelegatorValidators(delegatorAddr: string): Promise<Validator[]> {
    return [];
  },

  async getDelegatorValidator(
    delegatorAddr: string,
    validatorAddr: string
  ): Promise<Validator> {},

  async getValidators(
    status: string,
    page: number,
    limit: number
  ): Promise<Validator[]> {
    return [];
  },

  async getValidator(validatorAddr: string): Promise<Validator> {},

  async getValidatorDelegations(validatorAddr: string): Promise<Delegation[]> {
    return [];
  },

  async getValidatorUnbondingDelegations(
    validatorAddr: string
  ): Promise<UnbondingDelegation[]> {
    return [];
  },

  async getPool(): Promise<StakingPool> {},

  async getParameters(): Promise<StakingParameters> {},
};
