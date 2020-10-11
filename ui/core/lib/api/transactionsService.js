// based on draft swagger spec https://raw.githubusercontent.com/Sifchain/sifnode/c1bb5a268da8b519d0fc90f81fa194d31c0f82b3/api/openapi/swagger.yml?token=AAJSXWM6CDXYAEETSC6BJ2S7Q2JLS
export const transactionService = {
    // GET /txs/{hash}
    async getByhash(hash) { },
    // GET /txs/{hash}
    async search(actions, sender, page, limit, txheight) {
        return [];
    },
    // POST /txs
    async broadcast(tx) { },
    // POST /txs/encode
    async encode(tx) {
        return { tx: "somehash" };
    },
    async decode(tx) {
        return {};
    },
};
//# sourceMappingURL=transactionsService.js.map