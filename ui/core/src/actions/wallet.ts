import { Address, TxParams } from "../entities";
import { validateMnemonic } from "bip39";
import { Mnemonic } from "../entities/Wallet";
import { ActionContext } from ".";
import { effect } from "@vue/reactivity";

export default ({
  api,
  store,
}: ActionContext<"SifService" | "ClpService", "wallet">) => {
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
  };

  effect(() => {
    store.wallet.sif.isConnected = state.connected;
  });

  effect(() => {
    console.log("(akasha: sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd)");
    console.log("(shadowfiend: sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5)");
    console.log("sifAddress:", state.address);
    store.wallet.sif.address = state.address;
  });

  effect(() => {
    store.wallet.sif.balances = state.balances;
  });

  return actions;
};
