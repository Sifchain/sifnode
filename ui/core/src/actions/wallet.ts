import { Address, TxParams } from "../entities";
import { validateMnemonic } from "bip39";
import { Mnemonic } from "../entities/Wallet";
import { ActionContext } from ".";
import { effect } from "@vue/reactivity";

export default ({
  api,
  store,
}: ActionContext<
  "SifService" | "ClpService" | "EventBusService",
  "wallet"
>) => {
  const state = api.SifService.getState();

  const actions = {
    async getCosmosBalances(address: Address) {
      // TODO: validate sif prefix
      return await api.SifService.getBalance(address);
    },

    async connect(mnemonic: Mnemonic): Promise<string> {
      if (!mnemonic) throw "Mnemonic must be defined";
      if (!validateMnemonic(mnemonic)) throw "Invalid Mnemonic. Not sent.";
      return await api.SifService.setPhrase(mnemonic);
    },

    async sendCosmosTransaction(params: TxParams) {
      return await api.SifService.transfer(params);
    },

    async disconnect() {
      api.SifService.purgeClient();
    },

    async connectToWallet() {
      try {
        // TODO type
        await api.SifService.connect();
        store.wallet.sif.isConnected = true;
      } catch (error) {
        api.EventBusService.dispatch({
          type: "WalletConnectionErrorEvent",
          payload: {
            walletType: "sif",
            message: "Failed to connect to Keplr.",
          },
        });
      }
    },
  };

  effect(() => {
    if (store.wallet.sif.isConnected !== state.connected) {
      store.wallet.sif.isConnected = state.connected;
      if (store.wallet.sif.isConnected) {
        api.EventBusService.dispatch({
          type: "WalletConnectedEvent",
          payload: {
            walletType: "sif",
            address: store.wallet.sif.address,
          },
        });
      }
    }
  });

  effect(() => {
    store.wallet.sif.address = state.address;
  });

  effect(() => {
    store.wallet.sif.balances = state.balances;
  });

  return actions;
};
