import { initProvider } from "@metamask/inpage-provider";
import { ActionContext } from "..";

export default ({
  api,
  store,
}: ActionContext<"EtheriumService", "wallet" | "asset">) => ({
  // TODO: Move this out
  async init() {
    api.EtheriumService.onConnected(this.refreshWalletBalances);
  },
  async disconnectWallet() {
    await api.EtheriumService.disconnect();
    store.wallet.etheriumIsConnected = false;
    store.wallet.balances = [];
  },
  async connectToWallet() {
    await api.EtheriumService.connect();
    const isConnected = api.EtheriumService.isConnected();
    store.wallet.etheriumIsConnected = isConnected;
  },
  async refreshWalletBalances() {
    const balances = await api.EtheriumService.getBalance();
    store.wallet.balances = balances;
  },
});
