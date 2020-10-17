import { initProvider } from "@metamask/inpage-provider";
import { ActionContext } from "..";

export default ({
  api,
  store,
}: ActionContext<"EtheriumService", "wallet" | "asset">) => ({
  // TODO: Move this out
  async init() {
    api.EtheriumService.onChange(this.handleChange);
  },

  async handleChange() {
    console.log("handleChange");
    const balances = await api.EtheriumService.getBalance();
    const isConnected = api.EtheriumService.isConnected();

    console.log({ balances, isConnected });
    store.wallet.etheriumIsConnected = isConnected;
    store.wallet.balances = balances;
  },
  async disconnectWallet() {
    await api.EtheriumService.disconnect();
  },
  async connectToWallet() {
    await api.EtheriumService.connect();
  },
});
