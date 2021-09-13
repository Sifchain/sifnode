/**
 * Adds tokens to the whitelist in a batch
 * Please read LimitUpdating.md for instructions
 */

import * as hardhat from "hardhat";
import {container} from "tsyringe";
import {DeployedBridgeBank, requiredEnvVar} from "../src/contractSupport";
import {DeploymentName, HardhatRuntimeEnvironmentToken} from "../src/tsyringe/injectionTokens";
import {
    impersonateBridgeBankAccounts,
    setupRopstenDeployment,
    setupSifchainMainnetDeployment
} from "../src/hardhatFunctions";
import * as fs from "fs";

interface WhitelistTokenData {
  address: string
}

interface WhitelistData {
  array: Array<WhitelistTokenData>
}

export async function readTokenData(filename: string): Promise<WhitelistData> {
    const result = fs.readFileSync(filename, {encoding: "utf8"});
    return JSON.parse(result) as WhitelistData;
}

async function main() {
  container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})

  const deploymentName = requiredEnvVar("DEPLOYMENT_NAME")

  container.register(DeploymentName, {useValue: deploymentName})

  switch (hardhat.network.name) {
    case "ropsten":
        await setupRopstenDeployment(container, hardhat, deploymentName)
        break
    case "mainnet":
    case "hardhat":
    case "localhost":
        await setupSifchainMainnetDeployment(container, hardhat, deploymentName)
        break
  }

  const useForking = !!process.env["USE_FORKING"];
  if (useForking)
    await impersonateBridgeBankAccounts(container, hardhat, deploymentName)

  const whitelistData = await readTokenData(process.env["WHITELIST_DATA"] ?? "/tmp/nothing")

  const bridgeBank = (await container.resolve(DeployedBridgeBank).contract)

  const operator = await bridgeBank.operator()

  const operatorSigner = await hardhat.ethers.getSigner(operator)
  const bridgeBankAsOperator = bridgeBank.connect(operatorSigner);

  for (const addr of whitelistData.array) {
    if(await bridgeBankAsOperator.getTokenInEthWhiteList(addr.address)) {
      // this token is already in the whitelist;
      // the contract will not blow up on us, so we just skip this one.
      console.log(`\x1b[31mToken ${addr.address} NOT added to the whitelist: already there, no transaction sent\x1b[0m`);
      continue;
    }

    const tx = await bridgeBankAsOperator.updateEthWhiteList(addr.address, true);
    const receipt = await tx.wait();
    logResult(addr.address, receipt);
  }

  console.log('~~~ DONE ~~~');
}

function logResult(address:string, receipt:any) {
  if(receipt.events?.[0].args?.[1]) {
    // logs success in green
    console.log(`\x1b[32mToken ${address} added to the whitelist\x1b[0m`);
  } else {
    // logs failure in red
    console.log(`\x1b[31mToken ${address} NOT added to the whitelist\x1b[0m`);
  }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
