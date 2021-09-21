/**
 * Adds tokens to the whitelist in a batch
 * Please read LimitUpdating.md for instructions
 * 
 * @dev We're setting gasPrice explicitly, in accordance with the received ask.
 *      If this causes problems, please remove gasPrice from the transaction,
 *      and consult https://github.com/ethers-io/ethers.js/issues/1610 to understand why.
 *      In principle, it should work without it (as soon as EIP-1559 settles everywhere).
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

// Will estimate gas and multiply the result by this value (wiggle room)
const GAS_PRICE_BUFFER = 1.2;

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
  console.log(`\x1b[36mRunning bulk_set_whitelist script. Please wait...\x1b[0m`);

  container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat});

  const deploymentName = requiredEnvVar("DEPLOYMENT_NAME");

  container.register(DeploymentName, {useValue: deploymentName});

  switch (hardhat.network.name) {
    case "ropsten":
        await setupRopstenDeployment(container, hardhat, deploymentName);
        break;
    case "mainnet":
    case "hardhat":
    case "localhost":
        await setupSifchainMainnetDeployment(container, hardhat, deploymentName);
        break;
  }

  const useForking = !!process.env["USE_FORKING"];
  if (useForking)
    await impersonateBridgeBankAccounts(container, hardhat, deploymentName);

  const whitelistData = await readTokenData(process.env["WHITELIST_DATA"] ?? "/tmp/nothing");

  const bridgeBank = (await container.resolve(DeployedBridgeBank).contract);

  const operator = await bridgeBank.operator();
  console.log(`\x1b[36mOperator account is ${operator}\x1b[0m`);

  const operatorSigner = await hardhat.ethers.getSigner(operator);
  const bridgeBankAsOperator = bridgeBank.connect(operatorSigner);

  const addressList = [];
  for (const addr of whitelistData.array) {
    if(await bridgeBankAsOperator.getTokenInEthWhiteList(addr.address)) {
      // this token is already in the whitelist;
      // the contract will not blow up on us, so we just skip this one.
      console.log(`\x1b[31mToken ${addr.address} NOT added to the whitelist: already there, no transaction sent\x1b[0m`);
      continue;
    }

    addressList.push(addr.address);
    console.log(`\x1b[36mToken ${addr.address} will be added to the whitelist\x1b[0m`);
  }

  if(addressList.length > 0) {
    // Force ABI:
    const factory = await hardhat.ethers.getContractFactory("BridgeBank");
    const encodedData = factory.interface.encodeFunctionData('bulkWhitelistUpdateLimits', [addressList]);
    
    // Estimate gasPrice:
    const gasPrice = await estimateGasPrice();
    
    // UX
    console.log(`\x1b[46m\x1b[30mSending transaction. This may take a while, please wait...\x1b[0m`);

    const receipt = await (
      await operatorSigner.sendTransaction({
        data: encodedData,
        to: bridgeBank.address,
        gasPrice
      })
    ).wait();

    logResult(addressList, receipt);
  } else {
    // logs in red
    console.log(`\x1b[31mFailed to whitelist tokens: the final token list is empty. Were all tokens already whitelisted?\x1b[0m`);
  }

  console.log('~~~ DONE ~~~');
}

async function estimateGasPrice() {
  console.log('Estimating ideal Gas price, please wait...');

  const gasPrice = await hardhat.ethers.provider.getGasPrice();
  const finalGasPrice = Math.round(gasPrice.toNumber() * GAS_PRICE_BUFFER);

  console.log(`Using ideal Gas price: ${hardhat.ethers.utils.formatUnits(finalGasPrice, 'gwei')} GWEI`);
  
  return finalGasPrice;
}

function logResult(addressList:Array<String>, receipt:any) {
  if(receipt?.logs?.length > 0) {
    // logs success in green
    console.log(`\x1b[32mTokens added to the whitelist:\x1b[0m`);
    console.log(`\x1b[32m${addressList.join('\n')}\x1b[0m`);
  } else {
    // logs failure in red 
    console.log(`\x1b[31mFAILED: either got no tx receipt, or the receipt had no events.\x1b[0m`);
  }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });