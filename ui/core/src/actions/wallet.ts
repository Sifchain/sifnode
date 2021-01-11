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

    // initialize() {
    // something like this on load
    // or maybe on createApi ??
    // but where notification
    // },
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
        const address: any = await api.SifService.connect();
        console.log(address)
        // notify connected
        notify({type: "success", message: "Sif Account connected", 
        detail: address })
        store.wallet.sif.isConnected = true
        store.wallet.sif.address = address
        // get balance in this context, not service.connect in case you connect but there's not balance
        const balances: any = await api.SifService.getBalance(address)
        // if no address found on chain therefore no balance throws from getBalance... 
        store.wallet.sif.balances = balances;
      } catch (error) {
        // to the ui??
        console.log(error)
        notify({type:"error", ...error})
      }
    },
    async disconnectWallet() {
      await api.SifService.disconnect();
    },
  };

  // effect(() => {
  //   store.wallet.sif.isConnected = state.connected;
  // });

  // effect(() => {
  //   store.wallet.sif.address = state.address;
  // });

  // effect(() => {
  //   store.wallet.sif.balances = state.balances;
  // });

  return actions;
};
