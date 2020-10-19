import { SifTransaction } from "../entities";
import { validateMnemonic } from "bip39";
import { Mnemonic, SifAddress } from "../entities/Wallet";
import { ActionContext } from ".";
import { effect } from "@vue/reactivity";

export default ({ api, store }: ActionContext<"SifService", "wallet">) => {
  const state = api.SifService.getState();

  const actions = {
    async getCosmosBalances(address: SifAddress) {
      // TODO: validate sif prefix
      return await api.SifService.getBalance(address);
    },

    async signInCosmosWallet(mnemonic: Mnemonic): Promise<string> {
      if (!mnemonic) throw "Mnemonic must be defined";
      if (!validateMnemonic(mnemonic)) throw "Invalid Mnemonic. Not sent.";
      return await api.SifService.setPhrase(mnemonic);
    },

    async sendCosmosTransaction(sifTransaction: SifTransaction) {
      return await api.SifService.transfer(sifTransaction);
    },

    async signOutCosmosWallet() {
      api.SifService.purgeClient();
    },
  };

  effect(() => {
    store.wallet.sif.isConnected = state.connected;
  });

  effect(() => {
    store.wallet.sif.address = state.address;
    // if (state.address)
    //   store.wallet.sif.balances = await actions.getCosmosBalances(
    //     state.address
    //   );
  });

  effect(() => {
    store.wallet.sif.balances = state.balances;
  });

  return actions;
};
