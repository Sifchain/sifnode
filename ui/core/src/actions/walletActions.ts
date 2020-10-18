import { effect } from "@vue/reactivity";
import { ActionContext } from "..";

export default ({
  api,
  store,
}: ActionContext<"EtheriumService", "wallet" | "asset">) => {
  const actions = {
    async updateBalances(_?: string) {
      store.wallet.balances = await api.EtheriumService.getBalance();
    },
    async disconnectWallet() {
      await api.EtheriumService.disconnect();
    },
    async connectToWallet() {
      await api.EtheriumService.connect();
      actions.updateBalances();
    },
  };

  const etheriumState = api.EtheriumService.getState();

  effect(async () => {
    store.wallet.etheriumIsConnected = etheriumState.connected;
    await actions.updateBalances();
  });

  effect(async () => {
    await actions.updateBalances(etheriumState.log);
  });

  return actions;
};
