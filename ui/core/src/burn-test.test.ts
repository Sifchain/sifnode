import {
  makeCosmoshubPath,
  SigningCosmosClient,
  CosmosClient,
  Secp256k1HdWallet,
  isBroadcastTxFailure,
} from "@cosmjs/launchpad";
import { akasha } from "./test/utils/accounts";
import { resolve } from "path";
import { exec } from "child_process";
import { promisify } from "util";
import axios from "axios";

const execProm = promisify(exec);

async function runCmd(cmd: string) {
  const smartContractsPath = resolve(__dirname, "../../../smart-contracts");

  const result = await execProm(cmd, { cwd: smartContractsPath });
  console.log(result.stdout);
  console.log(result.stderr);
}

test("ethbridge::burn", async () => {
  const signer = await Secp256k1HdWallet.fromMnemonic(
    akasha.mnemonic,
    makeCosmoshubPath(0),
    "sif"
  );
  const signingClient = new SigningCosmosClient(
    "http://localhost:1317",
    akasha.address,
    signer
  );

  const unsignedClient = new CosmosClient("http://localhost:1317");

  async function getBalance(address: string, symbol: string) {
    const account = await unsignedClient.getAccount(address);
    return account?.balance.find((coin) => coin.denom === symbol)?.amount;
  }

  await runCmd(
    `yarn peggy:lock ${akasha.address} 0x0000000000000000000000000000000000000000 2000000000000000000`
  );
  await runCmd(`sleep 5`);
  await runCmd(`yarn advance 200`);
  await runCmd(`sleep 5`);

  expect(await getBalance(akasha.address, "ceth")).toEqual(
    "2000000001000000000"
  );

  // This works!
  // await runCmd(
  //   `sifnodecli tx ethbridge burn ${akasha.address} 0x627306090abaB3A6e1400e9345bC60c78a8BEf57 2000000000000000000 ceth --ethereum-chain-id=5777 --from=akasha --yes`
  // );

  // The following is the JS/REST equivalent of the above
  const result = (
    await axios.post("http://127.0.0.1:1317/ethbridge/burn", {
      ethereum_receiver: "0x627306090abaB3A6e1400e9345bC60c78a8BEf57",
      base_req: {
        chain_id: "sifchain",
        from: akasha.address,
      },
      amount: "2000000000000000000",
      symbol: "ceth",
      cosmos_sender: akasha.address,
      ethereum_chain_id: "5777",
      token_contract_address: "0x0000000000000000000000000000000000000000",
    })
  ).data;
  const msg: any[] = result.value.msg;
  const fee = {
    amount: [],
    gas: "200000",
  };
  const txHash = await signingClient.signAndBroadcast(msg, fee, "");

  if (isBroadcastTxFailure(txHash)) {
    // This fails!
    //   "unauthorized: signature verification failed; verify correct account sequence and chain-id"
    throw new Error(txHash.rawLog);
  }
  // end JS/REST equiv

  await runCmd(`sleep 5`);

  expect(await getBalance(akasha.address, "ceth")).toEqual("1000000000");
});
