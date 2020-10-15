import { ActionContext } from "..";

export default ({
  api,
  store,
}: ActionContext<"walletService", "wallet" | "asset">) => ({
  async disconnectWallet() {
    await api.walletService.disconnect();
    store.wallet.isConnected = false;
    store.wallet.balances = [];
  },
  async connectToWallet() {
    await api.walletService.connect();
    store.wallet.isConnected = api.walletService.isConnected();
    await this.refreshWalletBalances();
    // How do we listen for connection status?
    // What happens when the wallet drops off?
  },
  async refreshWalletBalances() {
    const balances = await api.walletService.getBalance();
    store.wallet.balances = balances;
  },
});
