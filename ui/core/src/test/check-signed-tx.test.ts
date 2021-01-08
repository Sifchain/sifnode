import {
  makeCosmoshubPath,
  SigningCosmosClient,
  Secp256k1HdWallet,
  isBroadcastTxFailure,
} from "@cosmjs/launchpad";

import { resolve } from "path";
import { exec } from "child_process";
import { promisify } from "util";
import axios from "axios";

const execProm = promisify(exec);

const akasha = {
  mnemonic:
    "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard",
  address: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
  name: "akasha",
};

async function runCmd(cmd: string) {
  const smartContractsPath = resolve(__dirname, "../../../../smart-contracts");

  const result = await execProm(cmd, { cwd: smartContractsPath });
  console.log(result.stdout);
  console.log(result.stderr);
}

// test for https://gist.github.com/ryardley/d81e475cad2a08d2ddfddc6bbc9496ec
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

  async function getNextSequence(address: string) {
    const currentSequence = (
      await axios.get(`http://127.0.0.1:1317/auth/accounts/${address}`)
    ).data.result.value.sequence;

    return `${parseInt(currentSequence) + 1}`;
  }

  await runCmd(
    `yarn peggy:lock ${akasha.address} 0x0000000000000000000000000000000000000000 2000000000000000000`
  );
  await runCmd(`sleep 5`);
  await runCmd(`yarn advance 200`);
  await runCmd(`sleep 5`);

  // This works!
  // await runCmd(
  //   `sifnodecli tx ethbridge burn ${akasha.address} 0x627306090abaB3A6e1400e9345bC60c78a8BEf57 2000000000000000000 ceth --ethereum-chain-id=5777 --from=akasha --yes`
  // );

  // The following is the JS/REST equivalent of the above
  const burnPayload = {
    ethereum_receiver: "0x627306090abaB3A6e1400e9345bC60c78a8BEf57",
    base_req: {
      chain_id: "sifchain",
      from: akasha.address,
      sequence: await getNextSequence(akasha.address),
    },
    amount: "2000000000000000000",
    symbol: "ceth",
    cosmos_sender: akasha.address,
    ethereum_chain_id: "5777",
    token_contract_address: "0x0000000000000000000000000000000000000000",
  };

  const burnResult = await axios.post(
    "http://127.0.0.1:1317/ethbridge/burn",
    burnPayload
  );

  const msg: any[] = burnResult.data.value.msg;

  const fee = {
    amount: [],
    gas: "200000",
  };

  // This makes the call to http://localhost:1317/tx
  const txHash = await signingClient.signAndBroadcast(msg, fee, "");

  if (isBroadcastTxFailure(txHash)) {
    // This fails!
    //   "unauthorized: signature verification failed; verify correct account sequence and chain-id"
    throw new Error(txHash.rawLog);
  }

  // end JS/REST equiv

  expect(!!"it should pass without error!").toBe(true);
});
