import { effect } from "@vue/reactivity";
import { Asset } from "../entities";
import { ActionContext } from "..";
import B from "../entities/utils/B";
import { ETH } from "../constants";

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
    async transferEthWallet(
      amount: number,
      recipient: string,
      asset: Asset = ETH
    ) {
      const hash = await api.EthereumService.transfer({
        amount: B(amount, asset.decimals),
        recipient,
        asset,
      });
      return hash;
    },
  };

  const etheriumState = api.EthereumService.getState();

  effect(async () => {
    store.wallet.eth.isConnected = etheriumState.connected;
    await actions.updateBalances();
  });

  effect(async () => {
    etheriumState.log; // trigger on log change
    await actions.updateBalances();
  });

  return actions;
};
