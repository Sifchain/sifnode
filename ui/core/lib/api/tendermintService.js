// based on draft swagger spec https://raw.githubusercontent.com/Sifchain/sifnode/c1bb5a268da8b519d0fc90f81fa194d31c0f82b3/api/openapi/swagger.yml?token=AAJSXWM6CDXYAEETSC6BJ2S7Q2JLS
export const tendermintService = {
    // /node_info
    async getNodeInfo() { },
    // /syncing
    async getSyncing() {
        return false;
    },
    // GET /blocks/lastest
    async getBlockLatest() { },
    // GET /blocks/{height}
    async getBlockAtHeight(height) { },
    // /validatorsets/latest
    async getValidatorsetLatest() { },
    // /validatorsets/{height}
    async getValidatorsetAtHeight(height) { },
};
//# sourceMappingURL=tendermintService.js.map