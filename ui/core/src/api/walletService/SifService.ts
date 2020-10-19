import {
  Secp256k1HdWallet,
  SigningCosmosClient,
  makeCosmoshubPath,
  CosmosClient,
  Account,
  coins,
} from "@cosmjs/launchpad";
import { SifWalletStore } from "../../store/wallet";
import { ADDR_PREFIX, API } from "../../constants";
import { Mnemonic, SifAddress } from "../../entities/Wallet";
import { SifTransaction } from "../../entities/Transaction";

// Warning This creates a client object used to *sign* TX
export async function cosmosSignin(
  mnemonic: Mnemonic
): Promise<SigningCosmosClient> {
  try {
    if (!mnemonic) {
      throw "No mnemonic. Can't generate wallet.";
    }
    const wallet = await Secp256k1HdWallet.fromMnemonic(
      mnemonic,
      makeCosmoshubPath(0),
      ADDR_PREFIX
    );
    const [{ address }] = await wallet.getAccounts();
    return new SigningCosmosClient(API, address, wallet);
  } catch (error) {
    throw error;
  }
}

export async function getCosmosBalance(address: SifAddress): Promise<Account> {
  if (!address) throw "Address undefined. Fail";
  if (address.length !== 42) throw "Address not valid (length). Fail"; // this is simple check, limited to default address type (check bech32)
  const client = new CosmosClient(API);
  try {
    const account = await client.getAccount(address);
    if (!account) throw "No Address found on chain";
    return account;
  } catch (error) {
    throw error;
  }
}

export async function signAndBroadcast(
  sifWalletClient: SifWalletStore["client"],
  sifTransaction: SifTransaction
): Promise<any> {
  if (!sifWalletClient) throw "No signed in client. Sign in with mnemonic.";
  if (!sifTransaction)
    throw "No user input data. Define who, what, and for how much.";
  // this seems like anti-pattern, with SifWallet.vue, "undefined" as culprit
  // but is alternative to define in vue with empty string?
  if (!sifTransaction.denom) throw "No denom.";
  // https://github.com/tendermint/vue/blob/develop/src/store/cosmos.js#L91
  const msg = {
    type: "cosmos-sdk/MsgSend",
    value: {
      amount: [
        {
          amount: sifTransaction.amount,
          denom: sifTransaction.denom,
        },
      ],
      from_address: sifWalletClient.senderAddress,
      to_address: sifTransaction.to_address,
    },
  };

  const fee = {
    amount: coins(500, sifTransaction.denom),
    gas: "200000",
  };

  return await sifWalletClient.signAndBroadcast([msg], fee, "cool");
}
