// Q:Is this a usecase or an api call?

import { Context } from ".";

export default ({ api, store }: Context) => ({
  async broadcastTx(tx: any) {
    // IF "Set <Xn> Quantity of <X> Token"
    // POST: TX to Wallet	-> LocalStorage(App) -> Transaction -> <WALLETX>
    // RENDER: Loading.vue
    // RENDER: WatchWallet.vue (For progress, prompts)
    // xAddress, xQuantity
    // yAddress, yQuantity
  },
});
