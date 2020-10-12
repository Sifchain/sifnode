// No async means this cannot use the store or remote apis.
function renderLiquidityData(amountA, amountB) {
    // ...
    return {};
}
export default ({ api, store }) => ({
    // Listener effects
    intializeAddLiquidityUseCase( /*  */) {
        // XXX: Need websocket listener https://docs.cosmos.network/master/core/events.html#subscribing-to-events
        //
        // const event$ = api.tendermintService.getSifchainEventStream()
        //
        // event$.subscribe(store.updateWithEvent)
        //
        // return () => {
        //   event$.unsubscribe()
        // }
    },
    // Render helpers that are business logic
    renderLiquidityData,
    // Command and effect usecases
    async addLiquidity(amountA, amountB) {
        // get wallet balances from store etc.
        //
        // store.setAddLiquidityTransactionInitiated()
        //
        // ...
        //
        // await api.transactionService.broadcast(tx)
        //
        // if error store.setAddLiquidityTransactionError(error)
        //
        // store.setAddLiquidityTransactionCompleted()
    },
});
//# sourceMappingURL=addLiquidity.js.map