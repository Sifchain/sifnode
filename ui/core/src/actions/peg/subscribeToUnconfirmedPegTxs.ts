import { ActionContext } from "..";
import { SubscribeToTx } from "./utils/subscribeToTx";
import { PegConfig } from "./index";

export const SubscribeToUnconfirmedPegTxs = ({
  api,
  store,
  config,
}: ActionContext<"EthbridgeService", "tx" | "wallet"> & {
  config: PegConfig;
}) => () => {
  // Update a tx state in the store
  const subscribeToTx = SubscribeToTx({ store });

  async function getSubscriptions() {
    const pendingTxs = await api.EthbridgeService.fetchUnconfirmedLockBurnTxs(
      store.wallet.eth.address,
      config.ethConfirmations
    );

    return pendingTxs.map(subscribeToTx);
  }

  // Need to keep subscriptions syncronous so using promise
  const subscriptionsPromise = getSubscriptions();

  // Return unsubscribe synchronously
  return () => {
    subscriptionsPromise.then(subscriptions =>
      subscriptions.forEach(unsubscribe => unsubscribe())
    );
  };
};
