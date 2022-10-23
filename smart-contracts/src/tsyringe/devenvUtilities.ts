import { EthereumResults } from "../devenv/devEnv"
import fs from "fs"
import { SifchainAccounts } from "./sifchainAccounts"
import { createSignerWithAddresss } from "./hardhatSupport"
import { DevEnvObject } from "../devenv/outputWriter"
import * as ethers from "ethers"

export function waitForFileToExist(filename: string): boolean {
  while (!fs.existsSync(filename)) {}
  return true
}

export function readDevEnvObj(devenvJsonPath: string): DevEnvObject {
  waitForFileToExist(devenvJsonPath)
  const contents = fs.readFileSync(devenvJsonPath, "utf8")
  const jsonObj = JSON.parse(contents)
  return jsonObj as DevEnvObject
}

export function devenvEthereumResults(devenvJsonPath: string): EthereumResults {
  const contents = fs.readFileSync(devenvJsonPath, "utf8")
  const jsonObj = JSON.parse(contents)
  return jsonObj["ethResults"] as EthereumResults
}

export async function ethereumResultsToSifchainAccounts(
  e: EthereumResults,
  provider: ethers.providers.JsonRpcProvider
): Promise<SifchainAccounts> {
  const operator = createSignerWithAddresss(e.accounts.operator.privateKey, provider)
  const owner = createSignerWithAddresss(e.accounts.owner.privateKey, provider)
  const pauser = createSignerWithAddresss(e.accounts.pauser.privateKey, provider)
  const validators = e.accounts.validators.map((k) =>
    createSignerWithAddresss(k.privateKey, provider)
  )
  const av = e.accounts.available.map((k) => createSignerWithAddresss(k.privateKey, provider))
  return new SifchainAccounts(operator, owner, pauser, validators, av)
}
