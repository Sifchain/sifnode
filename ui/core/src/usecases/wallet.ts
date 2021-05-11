import { Address, TxParams } from "../entities";
import { validateMnemonic } from "bip39";
import { Mnemonic } from "../entities/Wallet";
import { UsecaseContext } from ".";
import { effect } from "@vue/reactivity";

export default ({
  services,
  store,
}: UsecaseContext<"sif" | "clp" | "bus", "wallet">) => {
  const state = services.sif.getState();

  const actions = {
    async getCosmosBalances(address: Address) {
      // TODO: validate sif prefix
      return await services.sif.getBalance(address);
    },

    async connect(mnemonic: Mnemonic): Promise<string> {
      if (!mnemonic) throw "Mnemonic must be defined";
      if (!validateMnemonic(mnemonic)) throw "Invalid Mnemonic. Not sent.";
      return await services.sif.setPhrase(mnemonic);
    },

    async sendCosmosTransaction(params: TxParams) {
      return await services.sif.transfer(params);
    },

    async disconnect() {
      services.sif.purgeClient();
    },

    async connectToWallet() {
      try {
        // TODO type
        await services.sif.connect();
        store.wallet.sif.isConnected = true;
      } catch (error) {
        services.bus.dispatch({
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
        services.bus.dispatch({
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
