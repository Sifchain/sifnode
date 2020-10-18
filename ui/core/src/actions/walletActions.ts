import { effect } from "@vue/reactivity";
import { ActionContext } from "..";

export default ({
  api,
  store,
}: ActionContext<"EtheriumService", "wallet" | "asset">) => {
  const actions = {
    async updateBalances() {
      const balances = await api.EtheriumService.getBalance();
      store.wallet.balances = balances;
    },
    async disconnectWallet() {
      await api.EtheriumService.disconnect();
    },
    async connectToWallet() {
      await api.EtheriumService.connect();
    },
  };

  const etheriumState = api.EtheriumService.getReactive();

  effect(() => {
    console.log("connected");
    store.wallet.etheriumIsConnected = etheriumState.connected;
    actions.updateBalances();
  });

  effect(() => {
    const latestLog = etheriumState.log;
    actions.updateBalances();
  });

  return actions;
};
