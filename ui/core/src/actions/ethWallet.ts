import { effect } from "@vue/reactivity";
import { ActionContext } from "..";
import { Asset } from "../entities";
import B from "../entities/utils/B";

export default ({
  api,
  store,
}: ActionContext<"EthereumService", "wallet" | "asset">) => {
  const etheriumState = api.EthereumService.getState();

  const actions = {
    async disconnectWallet() {
      await api.EthereumService.disconnect();
    },
    async connectToWallet() {
      await api.EthereumService.connect();
    },
    async transferEthWallet(amount: number, recipient: string, asset: Asset) {
      const hash = await api.EthereumService.transfer({
        amount: B(amount, asset.decimals),
        recipient,
        asset,
      });
      return hash;
    },
  };

  effect(() => {
    store.wallet.eth.isConnected = etheriumState.connected;
  });

  effect(() => {
    store.wallet.eth.address = etheriumState.address;
  });

  effect(() => {
    store.wallet.eth.balances = etheriumState.balances;
  });

  effect(async () => {
    etheriumState.log; // trigger on log change
    await api.EthereumService.getBalance();
  });

  return actions;
};
