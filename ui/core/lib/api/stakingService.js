// based on draft swagger spec https://raw.githubusercontent.com/Sifchain/sifnode/c1bb5a268da8b519d0fc90f81fa194d31c0f82b3/api/openapi/swagger.yml?token=AAJSXWM6CDXYAEETSC6BJ2S7Q2JLS
export const stakingService = {
    async getDelegatorDelegations(delegatorAddr) {
        return [];
    },
    async getDelegatorCurrentDelegation(delegatorAddr, validatorAddr) { },
    async getDelegatorUnbondingDelegations(delegatorAddr) {
        return [];
    },
    async getDelegatorUnbondingDelegationsBtwValidator(delegatorAddr, validatorAddr) {
        return [];
    },
    async getRedelegations(delegator, validatorFrom, validatorTo) {
        return [];
    },
    async getDelegatorValidators(delegatorAddr) {
        return [];
    },
    async getDelegatorValidator(delegatorAddr, validatorAddr) { },
    async getValidators(status, page, limit) {
        return [];
    },
    async getValidator(validatorAddr) { },
    async getValidatorDelegations(validatorAddr) {
        return [];
    },
    async getValidatorUnbondingDelegations(validatorAddr) {
        return [];
    },
    async getPool() { },
    async getParameters() { },
};
//# sourceMappingURL=stakingService.js.map