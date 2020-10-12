// Q:Is this a usecase or an api call?
export default ({ api, store }) => ({
    async broadcastTx(tx) {
        // IF "Set <Xn> Quantity of <X> Token"
        // POST: TX to Wallet	-> LocalStorage(App) -> Transaction -> <WALLETX>
        // RENDER: Loading.vue
        // RENDER: WatchWallet.vue (For progress, prompts)
        // xAddress, xQuantity
        // yAddress, yQuantity
    },
});
//# sourceMappingURL=broadcastTx.js.map