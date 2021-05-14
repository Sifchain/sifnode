import { ActionContext } from "..";

export default ({
  api,
  store,
}: ActionContext<
  // Once we have moved all interactors to their own files this can be
  // ActionContext<any,any> or renamed to InteractorContext<any,any>
  "SifService" | "EthbridgeService" | "EventBusService" | "EthereumService", // Select the services you want to access
  "wallet" | "tx" // Select the store keys you want to access
>) => {
  // Create the context for passing to commands, queries and subscriptions
  const ctx = { api, store };

  /* 
    TODO: suggestion externalize all interactors injecting ctx would look like the following

    const commands = {
      unpeg: Unpeg(ctx),
      peg: Peg(ctx),
    }

    const queries = {
      getSifTokens: GetSifTokens(ctx),
      getEthTokens: GetEthTokens(ctx),
      calculateUnpegFee: CalculateUnpegFee(ctx),
    }
    
    const subscriptions = {
      subscribeToUnconfirmedPegTxs: SubscribeToUnconfirmedPegTxs(ctx),
    }
  */

  // Rename and split this up to subscriptions, commands, queries
  const actions = {
    async claimRewards() {
      if (!store.wallet.sif.address) throw "No from address provided for swap";

      // const tx = await api.DispensationService.claim( {fromAddress: store.wallet.sif.address, });
      console.log("=======");

      return "signed";
    },
  };

  return actions;
};
