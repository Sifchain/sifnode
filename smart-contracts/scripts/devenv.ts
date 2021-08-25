import { container } from "tsyringe";
import { HardhatNodeRunner } from "../src/devenv/hardhatNode";
import { GolangBuilder, GolangResultsPromise } from "../src/devenv/golangBuilder";
import { SifnodedRunner } from "../src/devenv/sifnoded";
import { SmartContractDeployer } from "../src/devenv/smartcontractDeployer";
import { cons } from "fp-ts/lib/ReadonlyNonEmptyArray";
import { sampleCode } from "../src/devenv/synchronousCommand";
import { EbrelayerRunner } from "../src/devenv/ebrelayer";


async function startHardhat() {
  const node = container.resolve(HardhatNodeRunner)
  const [process, resultsPromise] = node.go()
  const results = await resultsPromise
  console.log(`rsltis: ${JSON.stringify(results, undefined, 2)}`)
  return process
}

async function golangBuilder() {
  const node = container.resolve(GolangBuilder)
  const [process, resultsPromise] = node.go()
  let golangResultsPromise = new GolangResultsPromise(resultsPromise);
  container.register(GolangResultsPromise, { useValue: golangResultsPromise })
  const sifnodeTask = sifnodedBuilder(golangResultsPromise)
  const results = await resultsPromise
  console.log(`golangBuilder: ${JSON.stringify(results, undefined, 2)}`)
  return Promise.all([process, sifnodeTask])
}

async function sifnodedBuilder(golangResults: GolangResultsPromise) {
  console.log('in sifnodedBuilder')
  const node = container.resolve(SifnodedRunner)
  const [process, resultsPromise] = node.go()
  const results = await resultsPromise
  console.log(`golangBuilder: ${JSON.stringify(results, undefined, 2)}`)
  return process
}

async function smartContractDeployer() {
  const node: SmartContractDeployer = container.resolve(SmartContractDeployer);
  const [process, resultsPromise] = node.go();
  const result = await resultsPromise;
  console.log(`Contracts deployed: ${JSON.stringify(result.contractAddresses, undefined, 2)}`)
  await ebrelayerBuilder()
  return;
}

async function ebrelayerBuilder() {
  const node: EbrelayerRunner = container.resolve(EbrelayerRunner);
  const [process, resultsPromise] = node.go();
  const result = await resultsPromise;
  return process;
}

async function main() {
  await Promise.all([startHardhat(), golangBuilder()])
    .then(smartContractDeployer)
    .then(() => {
      console.log("Congrats, you did not fail, yay!")
    })
    .catch((e) => { console.log("Deployment failed. Lets log where it broke: ", e) });
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
