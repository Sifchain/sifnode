import { UsecaseContext } from "..";

export default ({
  services,
  store,
}: UsecaseContext<
  "sif" | "clp" | "bus" | "dispensation",
  "pools" | "wallet" | "accountpools"
>) => {
  // Create the context for passing to commands, queries and subscriptions
  const ctx = { services, store };

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
  const commands = {
    async claim(params: { claimType: "2" | "3"; fromAddress: string }) {
      console.log(params);
      if (!store.wallet.sif.address) throw "No from address provided for swap";

      const tx = await services.dispensation.claim(params);
      console.log("=======");
      return await services.sif.signAndBroadcast(tx.value.msg);
    },
  };

  return commands;
};
