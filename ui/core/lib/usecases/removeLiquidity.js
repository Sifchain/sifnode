function renderRemoveLiquidityPageData(liquidityPool, token, tokenAmount) {
    // ...
    return {};
}
export default ({ api, store }) => ({
    intializeRemoveLiquidity() {
        // XXX: Need websocket listener https://docs.cosmos.network/master/core/events.html#subscribing-to-events
        //
        // const event$ = api.tendermintService.getSifchainEventStream('pool')
        //
        // store.clearPageValues()
        // event$.subscribe(store.updateMarketData)
        //
        // return () => {
        //   event$.unsubscribe()
        // }
    },
    // Render helpers that are business logic
    renderRemoveLiquidityPageData,
    // Commands
    async removeLiquidity(liquidityPool, token, tokenAmount) {
        // get wallet balances from store etc.
        //
        // store.setRemoveLiquidityTransactionInitiated()
        //
        // ...
        //
        // await api.transactionService.broadcast(tx)
        //
        // if error store.setRemoveLiquidityTransactionError(error)
        //
        // store.setRemoveLiquidityTransactionCompleted()
        // store.setNewPoolShare(poolShare)
    },
});
//# sourceMappingURL=removeLiquidity.js.map