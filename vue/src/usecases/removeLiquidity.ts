import { Context } from ".";
import { Pair, Token, TokenAmount } from "../entities";

function renderRemoveLiquidityPageData(
  liquidityPool: Pair,
  token: Token,
  tokenAmount: TokenAmount
): {
  canRemoveLiquidity: boolean;
  amount: TokenAmount;
  gasFees: TokenAmount; // TokenAmount or something else?
  shareOfPool: number;
  amountToRemoveIsTooHigh: boolean;
} {
  // ...
  return {} as any;
}

export default ({ api, store }: Context) => ({
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
  async removeLiquidity(
    liquidityPool: Pair,
    token: Token,
    tokenAmount: TokenAmount
  ) {
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
