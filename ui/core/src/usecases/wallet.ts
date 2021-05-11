import { Address, TxParams } from "../entities";
import { validateMnemonic } from "bip39";
import { Mnemonic } from "../entities/Wallet";
import { UsecaseContext } from ".";
import { effect } from "@vue/reactivity";

export default ({
  services,
  store,
}: UsecaseContext<
  "SifService" | "ClpService" | "EventBusService",
  "wallet"
>) => {
  const state = services.SifService.getState();

  const actions = {
    async getCosmosBalances(address: Address) {
      // TODO: validate sif prefix
      return await services.SifService.getBalance(address);
    },

    async connect(mnemonic: Mnemonic): Promise<string> {
      if (!mnemonic) throw "Mnemonic must be defined";
      if (!validateMnemonic(mnemonic)) throw "Invalid Mnemonic. Not sent.";
      return await services.SifService.setPhrase(mnemonic);
    },

    async sendCosmosTransaction(params: TxParams) {
      return await services.SifService.transfer(params);
    },

    async disconnect() {
      services.SifService.purgeClient();
    },

    async connectToWallet() {
      try {
        // TODO type
        await services.SifService.connect();
        store.wallet.sif.isConnected = true;
      } catch (error) {
        services.EventBusService.dispatch({
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
        services.EventBusService.dispatch({
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
