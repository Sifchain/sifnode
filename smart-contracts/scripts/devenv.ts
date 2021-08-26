import { EthereumArguments, HardhatNodeRunner } from "../src/devenv/hardhatNode";
import { GolangBuilder, GolangResults, GolangResultsPromise } from "../src/devenv/golangBuilder";
import { SifnodedResults, SifnodedRunner, ValidatorValues } from "../src/devenv/sifnoded";
import { DeployedContractAddresses } from "../scripts/deploy_contracts";
import { SmartContractDeployer } from "../src/devenv/smartcontractDeployer";
import { EbrelayerRunner } from "../src/devenv/ebrelayer";


async function startHardhat() {
  const node = new HardhatNodeRunner(
    {
      host: "localhost",
      port: 8545,
      nValidators: 1,
      networkId: 1,
      chainId: 1
    })
  const [process, resultsPromise] = node.go()
  const results = await resultsPromise
  console.log(`rsltis: ${JSON.stringify(results, undefined, 2)}`)
  return { process }
}

async function golangBuilder() {
  const node = new GolangBuilder()
  const [process, resultsPromise] = node.go()
  let golangResultsPromise = new GolangResultsPromise(resultsPromise);
  const sifnodeTask = sifnodedBuilder(golangResultsPromise)
  const results = await resultsPromise
  console.log(`golangBuilder: ${JSON.stringify(results, undefined, 2)}`)
  const output = await Promise.all([process, sifnodeTask, results])
  return {
    process: output[0],
    sifnodeTask: output[1],
    results: output[2]
  }
}

async function sifnodedBuilder(golangResults: GolangResultsPromise) {
  console.log('in sifnodedBuilder')
  const node = new SifnodedRunner(
    {
      logfile: "/tmp/sifnoded.log",
      rpcPort: 9000,
      nValidators: 1,
      chainId: "localnet",
      networkConfigFile: "/tmp/sifnodedConfig.yml",
      networkDir: "/tmp/sifnodedNetwork",
      seedIpAddress: "10.10.1.1",
      whitelistFile: "../test/integration/whitelisted-denoms.json"
    },
    golangResults
  )
  const [process, resultsPromise] = node.go()
  const results = await resultsPromise
  console.log(`golangBuilder: ${JSON.stringify(results, undefined, 2)}`)
  return {
    process,
    results
  }
}

async function smartContractDeployer(golangResults: GolangResults, sifnodedResults: SifnodedResults) {
  const node: SmartContractDeployer = new SmartContractDeployer();
  const [process, resultsPromise] = node.go();
  const result = await resultsPromise;
  console.log(`Contracts deployed: ${JSON.stringify(result.contractAddresses, undefined, 2)}`)
  await ebrelayerBuilder(result.contractAddresses, sifnodedResults.validatorValues[0])
  return;
}

async function ebrelayerBuilder(contractAddresses: DeployedContractAddresses, validater: ValidatorValues) {
  const node: EbrelayerRunner = new EbrelayerRunner({
    smartContract: contractAddresses,
    websocketAddress: "ws://localhost:7545/",
    tcpURL: "tcp://0.0.0.0:26657",
    chainNet: "localnet",
    ebrelayerDB: `levelDB.db`,
    relayerdbPath: "",
    validatorValues: validater,
    symbolTranslatorFile: "../test/integration/whitelisted-denom.json"
  });
  const [process, resultsPromise] = node.go();
  const result = await resultsPromise;
  return { process };
}

async function main() {
  try {
    let results = await Promise.all([startHardhat(), golangBuilder()])
    const hardhat = results[0]
    const golang = results[1]
    await smartContractDeployer(golang.results, golang.sifnodeTask.results)
    console.log("Congrats, you did not fail, yay!")
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
