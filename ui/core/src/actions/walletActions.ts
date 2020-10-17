import { ActionContext } from "..";

export default ({
  api,
  store,
}: ActionContext<"EtheriumService", "wallet" | "asset">) => ({
  async disconnectWallet() {
    await api.EtheriumService.disconnect();
    store.wallet.etheriumIsConnected = false;
    store.wallet.balances = [];
  },
  async connectToWallet() {
    await api.EtheriumService.connect();
    const isConnected = api.EtheriumService.isConnected();
    store.wallet.etheriumIsConnected = isConnected;
    await this.refreshWalletBalances();
  },
  async refreshWalletBalances() {
    const balances = await api.EtheriumService.getBalance();
    store.wallet.balances = balances;
  },
});
