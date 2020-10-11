// No async means this cannot use the store or remote apis.
function renderCreatePoolData(amountA, amountB) {
    // ...
    return {};
}
export default ({ api, store }) => ({
    intializeCreatePoolUseCase( /*  */) {
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
    renderCreatePoolData,
    // Command and effect usecases
    async createPool(amountA, amountB) {
        // get wallet balances from store etc.
        //
        // store.setCreatePoolTransactionInitiated()
        //
        // ...
        //
        // await api.transactionService.broadcast(tx)
        //
        // if error store.seCreatePoolTransactionError(error)
        //
        // store.seCreatePoolTransactionCompleted()
    },
});
//# sourceMappingURL=createPool.js.map