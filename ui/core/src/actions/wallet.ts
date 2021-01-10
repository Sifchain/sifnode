import { Address, TxParams } from "../entities";
import { validateMnemonic } from "bip39";
import { Mnemonic } from "../entities/Wallet";
import { ActionContext } from ".";
import { effect } from "@vue/reactivity";
import notify from "../api/utils/Notifications"
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
    async connectToWallet() {
      try  {
        await api.SifService.connect();
      } catch (error) {
        // to the ui??
        notify({type:"error", ...error})
      }
    },
    async disconnectWallet() {
      await api.SifService.disconnect();
    },
  };

  effect(() => {
    store.wallet.sif.isConnected = state.connected;
  });

  effect(() => {
    store.wallet.sif.address = state.address;
  });

  effect(() => {
    store.wallet.sif.balances = state.balances;
  });

  return actions;
};
