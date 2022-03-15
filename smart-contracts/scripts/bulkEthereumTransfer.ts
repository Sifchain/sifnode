import * as hardhat from "hardhat";
import { container } from "tsyringe";
import { buildDevEnvContracts } from "../src/contractSupport";
import { HardhatRuntimeEnvironmentToken } from "../src/tsyringe/injectionTokens";
import { SifchainContractFactories } from "../src/tsyringe/contracts";
import * as fs from "fs";
import { executeLock } from "../test/devenv/evm_lock_burn";
import { readDevEnvObj } from "../src/tsyringe/devenvUtilities";
import { SifchainAccountsPromise } from "../src/tsyringe/sifchainAccounts";
import { BigNumber, ContractTransaction } from "ethers";
import web3 from "web3";

async function main() {
  container.register(HardhatRuntimeEnvironmentToken, { useValue: hardhat });
  const devEnvObject = readDevEnvObj("environment.json");
  const factories = container.resolve(SifchainContractFactories);
  const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories);

  const sifchainAccounts = container.resolve(SifchainAccountsPromise);

  const sender = (await sifchainAccounts.accounts).availableAccounts[0];

  // read in all sifchain accounts
  const destinationAccountsString = fs.readFileSync("test_data/test_keychain.json", "utf8");
  const destinationAccounts = JSON.parse(destinationAccountsString);
  // Array.slice is exclusive
  const locks = destinationAccounts.slice(0, 300)
    .map((account: any) => web3.utils.utf8ToHex(account["address"]))
    .map((sifAddr: string) => executeLock(contracts, BigNumber.from(1000), sender, sifAddr))
    .map(async (transactionPromise: Promise<ContractTransaction>) => await (await transactionPromise).wait(50))
  await Promise.all(locks);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
