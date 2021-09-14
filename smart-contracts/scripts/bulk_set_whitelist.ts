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
 
   const addressList = [];
   for (const addr of whitelistData.array) {
     if(await bridgeBankAsOperator.getTokenInEthWhiteList(addr.address)) {
       // this token is already in the whitelist;
       // the contract will not blow up on us, so we just skip this one.
       console.log(`\x1b[31mToken ${addr.address} NOT added to the whitelist: already there, no transaction sent\x1b[0m`);
       continue;
     }
 
     addressList.push(addr.address);
   }
 
   if(addressList.length > 0) {
     // Force ABI:
     const factory = await hardhat.ethers.getContractFactory("BridgeBank");
     const encodedData = factory.interface.encodeFunctionData('bulkWhitelistUpdateLimits', [addressList]);
     
     const receipt = await (
       await operatorSigner.sendTransaction({ data: encodedData, to: bridgeBank.address })
     ).wait();
 
     logResult(addressList, receipt);
   } else {
     // logs in red
     console.log(`\x1b[31mFailed to whitelist tokens: the final token list is empty. Were all tokens already whitelisted?\x1b[0m`);
   }
 
   console.log('~~~ DONE ~~~');
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