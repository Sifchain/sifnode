import { effect } from "@vue/reactivity";
import { ActionContext } from "..";
import { Asset } from "../entities";
import B from "../entities/utils/B";

export default ({
  api,
  store,
}: ActionContext<
  "EthereumService" | "EventBusService",
  "wallet" | "asset"
>) => {
  api.EthereumService.onProviderNotFound(() => {
    api.EventBusService.dispatch({
      type: "WalletConnectionErrorEvent",
      payload: {
        walletType: "eth",
        message: "Metamask not found.",
      },
    });
  });

  api.EthereumService.onChainIdDetected((chainId) => {
    store.wallet.eth.chainId = chainId;
  });

  const etheriumState = api.EthereumService.getState();

  const actions = {
    async disconnectWallet() {
      await api.EthereumService.disconnect();
    },
    async connectToWallet() {
      try {
        await api.EthereumService.connect();
      } catch (err) {
        api.EventBusService.dispatch({
          type: "WalletConnectionErrorEvent",
          payload: {
            walletType: "eth",
            message: "Failed to connect to Metamask.",
          },
        });
      }
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
    // Only show connected when we have an address
    if (store.wallet.eth.isConnected !== etheriumState.connected) {
      store.wallet.eth.isConnected =
        etheriumState.connected && !!etheriumState.address;

      if (store.wallet.eth.isConnected) {
        api.EventBusService.dispatch({
          type: "WalletConnectedEvent",
          payload: {
            walletType: "eth",
            address: store.wallet.eth.address,
          },
        });
      }
    }
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
