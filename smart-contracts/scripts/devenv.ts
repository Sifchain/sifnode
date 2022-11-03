import { HardhatNodeRunner } from "../src/devenv/hardhatNode"
import { GolangBuilder, GolangResults } from "../src/devenv/golangBuilder"
import {
  SifnodedResults,
  SifnodedRunner,
  ValidatorValues,
  EbRelayerAccount,
} from "../src/devenv/sifnoded"
import { DeployedContractAddresses } from "./deploy_contracts_dev"
import {
  SmartContractDeployer,
  SmartContractDeployResult,
} from "../src/devenv/smartcontractDeployer"
import { RelayerRunner, WitnessRunner, EbrelayerArguments } from "../src/devenv/ebrelayer"
import { EthereumAddressAndKey, EthereumResults } from "../src/devenv/devEnv"
import path from "path"
import { notify } from "node-notifier"
import { strict, string } from "yargs"
import { ContractFactory } from "ethers"
import { EnvJSONWriter } from "../src/devenv/outputWriter"
import fs from "fs"

async function startHardhat() {
  const node = new HardhatNodeRunner()
  const resultsPromise = node.go()
  const results = await resultsPromise
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
    results: output[1],
  }
}

async function sifnodedBuilder(golangResults: GolangResults) {
  console.log("in sifnodedBuilder")
  const node = new SifnodedRunner(golangResults)
  const resultsPromise = node.go()
  const results = await resultsPromise
  console.log(`golangBuilder: ${JSON.stringify(results, undefined, 2)}`)
  return {
    process,
    results,
  }
}

async function smartContractDeployer() {
  const node: SmartContractDeployer = new SmartContractDeployer()
  const resultsPromise = node.go()
  const result = await resultsPromise
  console.log(`Contracts deployed: ${JSON.stringify(result.contractAddresses, undefined, 2)}`)
  return { process, result }
}

async function relayerBuilder(args: EbrelayerArguments) {
  const node: RelayerRunner = new RelayerRunner(args)
  const resultsPromise = node.go()
  const result = await resultsPromise
  return { process, result }
}

async function witnessBuilder(args: EbrelayerArguments) {
  const node: WitnessRunner = new WitnessRunner(args)
  const resultsPromise = node.go()
  const result = await resultsPromise
  return { process, result }
}

async function ebrelayerWitnessBuilder(
  contractAddresses: DeployedContractAddresses,
  ethereumAccount: EthereumAddressAndKey,
  validater: ValidatorValues,
  relayerAccount: EbRelayerAccount,
  witnessAccount: EbRelayerAccount,
  golangResults: GolangResults,
  chainId: number
) {
  const relayerArgs: EbrelayerArguments = {
    smartContract: contractAddresses,
    account: ethereumAccount,
    validatorValues: validater,
    sifnodeAccount: relayerAccount,
    golangResults,
    chainId,
  }
  const witnessArgs = { ...relayerArgs, sifnodeAccount: witnessAccount }
  const relayerPromise = relayerBuilder(relayerArgs)
  const witnessPromise = witnessBuilder(witnessArgs)
  const [relayer, witness] = await Promise.all([relayerPromise, witnessPromise])
  return {
    relayer,
    witness,
  }
}

async function main() {
  try {
    await fs.promises.mkdir("/tmp/sifnode", { recursive: true })
    const sigterm = new Promise((res, _) => {
      process.on("SIGINT", res)
      process.on("SIGTERM", res)
    })
    const [hardhat, golang] = await Promise.all([startHardhat(), golangBuilder()])
    const sifnode = await sifnodedBuilder(golang.results)
    const smartcontract = await smartContractDeployer()
    const { relayer, witness } = await ebrelayerWitnessBuilder(
      smartcontract.result.contractAddresses,
      hardhat.results.accounts.validators[0],
      sifnode.results.validatorValues[0],
      sifnode.results.relayerAddresses[0],
      sifnode.results.witnessAddresses[0],
      golang.results,
      // we need configure the chain id as hardhat
      // hardhat.results.chainId
      9999
    )
    EnvJSONWriter({
      contractResults: smartcontract.result,
      ethResults: hardhat.results,
      goResults: golang.results,
      sifResults: sifnode.results,
    })
    await sigterm
    console.log("Caught interrupt signal, cleaning up.")
    sifnode.process.kill(sifnode.process.pid)
    hardhat.process.kill(hardhat.process.pid)
    relayer.process.kill(relayer.process.pid)
    witness.process.kill(witness.process.pid)
    console.log("All child process terminated, goodbye.")
    notify({
      title: "Sifchain DevEnvironment Notice",
      message: `Dev Environment has recieved either a SIGINT or SIGTERM signal, all process have exited.`,
    })
  } catch (error) {
    console.log("Deployment failed. Lets log where it broke: ", error)
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    if (typeof error == "number") process.exit(error)
    else {
      console.error(error)
      process.exit(1)
    }
  })
