import { effect } from "@vue/reactivity";
import { UsecaseContext } from "../..";
import { Asset, IAsset } from "../../entities";
import B from "../../entities/utils/B";
import { isSupportedEVMChain } from "../utils";

export default ({
  services,
  store,
}: UsecaseContext<"eth" | "bus", "wallet" | "asset">) => {
  services.eth.onProviderNotFound(() => {
    services.bus.dispatch({
      type: "WalletConnectionErrorEvent",
      payload: {
        walletType: "eth",
        message: "Metamask not found.",
      },
    });
  });

  services.eth.onChainIdDetected((chainId) => {
    store.wallet.eth.chainId = chainId;
  });

  const etheriumState = services.eth.getState();

  const actions = {
    isSupportedNetwork() {
      return isSupportedEVMChain(store.wallet.eth.chainId);
    },

    async disconnectWallet() {
      await services.eth.disconnect();
    },

    async connectToWallet() {
      try {
        await services.eth.connect();
      } catch (err) {
        services.bus.dispatch({
          type: "WalletConnectionErrorEvent",
          payload: {
            walletType: "eth",
            message: "Failed to connect to Metamask.",
          },
        });
      }
    },

    async transferEthWallet(amount: number, recipient: string, asset: Asset) {
      const hash = await services.eth.transfer({
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
        services.bus.dispatch({
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
    await services.eth.getBalance();
  });

  return actions;
};
