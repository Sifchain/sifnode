import { spinner } from "zx/experimental";

const { DirectSecp256k1HdWallet } = require("@cosmjs/proto-signing");
const { SifSigningStargateClient } = require("@sifchain/stargate");

const { SIFNODE_NODE } = process.env;

export async function getSigningClient(mnemonic) {
  const { wallet, account, signingClient } = await spinner(
    "loading signing client",
    () =>
      within(async () => {
        const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
          prefix: "sif",
        });
        const [account] = await wallet.getAccounts();
        const signingClient = await SifSigningStargateClient.connectWithSigner(
          SIFNODE_NODE,
          wallet
        );
        return { wallet, account, signingClient };
      })
  );
  return { wallet, account, signingClient };
}
