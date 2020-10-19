import { effect } from "@vue/reactivity";
import { ActionContext } from "..";

export default ({
  api,
  store,
}: ActionContext<"EthereumService", "wallet" | "asset">) => {
  const actions = {
    async updateBalances(_?: string) {
      store.wallet.eth.balances = await api.EthereumService.getBalance();
    },
    async disconnectWallet() {
      await api.EthereumService.disconnect();
    },
    async connectToWallet() {
      await api.EthereumService.connect();
      actions.updateBalances();
    },
  };

  const etheriumState = api.EthereumService.getState();

  effect(async () => {
    store.wallet.eth.isConnected = etheriumState.connected;
    await actions.updateBalances();
  });

  effect(async () => {
    await actions.updateBalances(etheriumState.log);
  });

  return actions;
};
