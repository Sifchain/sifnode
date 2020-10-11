export default ({ api, store }) => ({
    async updateListOfAvailableTokens() {
        const walletBalances = await api.walletService.getAssetBalances();
        store.setUserBalances(walletBalances);
    },
});
//# sourceMappingURL=queryListOfAvailableTokens.js.map