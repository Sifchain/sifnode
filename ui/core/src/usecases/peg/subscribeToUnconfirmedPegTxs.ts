import { UsecaseContext } from "..";
import { SubscribeToTx } from "./utils/subscribeToTx";
import { PegConfig } from "./index";

export const SubscribeToUnconfirmedPegTxs = ({
  services,
  store,
  config,
}: UsecaseContext<"ethbridge" | "bus", "tx" | "wallet"> & {
  config: PegConfig;
}) => () => {
  if (!store.wallet.eth.address) return () => {};

  // Update a tx state in the store
  const subscribeToTx = SubscribeToTx({ store, services });

  async function getSubscriptions() {
    const pendingTxs = await services.ethbridge.fetchUnconfirmedLockBurnTxs(
      store.wallet.eth.address,
      config.ethConfirmations,
    );
    return pendingTxs.map(subscribeToTx);
  }

  // Need to keep subscriptions syncronous so using promise
  const subscriptionsPromise = getSubscriptions();

  // Return unsubscribe synchronously
  return () => {
    subscriptionsPromise.then((subscriptions) =>
      subscriptions.forEach((unsubscribe) => unsubscribe()),
    );
  };
};
