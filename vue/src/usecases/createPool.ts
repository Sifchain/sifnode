import { Context } from ".";
import { Token, TokenAmount } from "../entities";

// No async means this cannot use the store or remote apis.
function renderCreatePoolData(
  amountA: TokenAmount,
  amountB: TokenAmount
): {
  tokenAPerBRatio: number; // XXX: Fraction?
  tokenBPerARatio: number;
  tokenAAmountOwned: TokenAmount;
  tokenBAmountOwned: TokenAmount;
  shareOfPool: number;
  canCreatePool: boolean;
  isInsufficientFunds: boolean;
} {
  // ...
  return {} as any;
}

export default ({ api, store }: Context) => ({
  intializeCreatePoolUseCase(/*  */) {
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
  async createPool(amountA: TokenAmount, amountB: TokenAmount) {
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
