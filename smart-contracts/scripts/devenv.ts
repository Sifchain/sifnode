import { HardhatNodeRunner } from "../src/devenv/hardhatNode";
import { GolangBuilder, GolangResults } from "../src/devenv/golangBuilder";
import { SifnodedRunner, ValidatorValues } from "../src/devenv/sifnoded";
import { DeployedContractAddresses } from "../scripts/deploy_contracts";
import { SmartContractDeployer } from "../src/devenv/smartcontractDeployer";
import { EbrelayerRunner } from "../src/devenv/ebrelayer";
import { EthereumAddressAndKey } from "../src/devenv/devEnv";

async function startHardhat() {
  const node = new HardhatNodeRunner()
  const resultsPromise = node.go()
  const results = await resultsPromise
  console.log(`rsltis: ${JSON.stringify(results, undefined, 2)}`)
  return { process, results }
}

async function golangBuilder() {
  const node = new GolangBuilder()
  const resultsPromise = node.go()
  const results = await resultsPromise
  console.log(`golangBuilder: ${JSON.stringify(results, undefined, 2)}`)
  const output = await Promise.all([process, results])
  return {
    process: output[0],
    results: output[1]
  }
}

async function sifnodedBuilder(golangResults: GolangResults) {
  console.log('in sifnodedBuilder')
  const node = new SifnodedRunner(golangResults)
  const resultsPromise = node.go()
  const results = await resultsPromise
  console.log(`golangBuilder: ${JSON.stringify(results, undefined, 2)}`)
  return {
    process,
    results
  }
}

async function smartContractDeployer() {
  const node: SmartContractDeployer = new SmartContractDeployer();
  const resultsPromise = node.go();
  const result = await resultsPromise;
  console.log(`Contracts deployed: ${JSON.stringify(result.contractAddresses, undefined, 2)}`)
  return { process, result };
}

async function ebrelayerBuilder(contractAddresses: DeployedContractAddresses, ethereumAccount: EthereumAddressAndKey, validater: ValidatorValues) {
  const node: EbrelayerRunner = new EbrelayerRunner({
    smartContract: contractAddresses,
    account: ethereumAccount,
    validatorValues: validater,
  });
  const resultsPromise = node.go();
  const result = await resultsPromise;
  return { process, result };
}

async function main() {
  try {
    const sigterm = new Promise((res, _) => {
      process.on('SIGINT', res);
      process.on('SIGTERM', res);
    });
    const [hardhat, golang] = (await Promise.all([startHardhat(), golangBuilder()]))
    const sifnode = await sifnodedBuilder(golang.results);
    const smartcontract = await smartContractDeployer()
    const ebrelayer = await ebrelayerBuilder(
      smartcontract.result.contractAddresses,
      hardhat.results.accounts.validators[0],
      sifnode.results.validatorValues[0]
    )
    await sigterm
    console.log("Caught interrupt signal, cleaning up.");
    sifnode.process.kill(sifnode.process.pid);
    hardhat.process.kill(hardhat.process.pid);
    ebrelayer.process.kill(ebrelayer.process.pid);
    console.log("All child process terminated, goodbye.");
  } catch (error) {
    console.log("Deployment failed. Lets log where it broke: ", error);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    if (typeof error == "number")
      process.exit(error)
    else {
      console.error(error);
      process.exit(1)
    }
  });
