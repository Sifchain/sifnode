import {EthereumResults} from "../devenv/devEnv";
import fs from "fs";
import {SifchainAccounts} from "./sifchainAccounts";
import {createSignerWithAddresss} from "./hardhatSupport";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import {DevEnvObject} from "../devenv/outputWriter";

export function readDevEnvObj(devenvJsonPath: string): DevEnvObject {
    const contents = fs.readFileSync(devenvJsonPath, 'utf8')
    const jsonObj = JSON.parse(contents)
    return jsonObj as DevEnvObject
}

export function devenvEthereumResults(devenvJsonPath: string): EthereumResults {
    const contents = fs.readFileSync(devenvJsonPath, 'utf8')
    const jsonObj = JSON.parse(contents)
    return jsonObj["ethResults"] as EthereumResults
}

export async function ethereumResultsToSifchainAccounts(e: EthereumResults): Promise<SifchainAccounts> {
    const operator = createSignerWithAddresss(e.accounts.operator.address, e.accounts.operator.privateKey)
    const owner = createSignerWithAddresss(e.accounts.owner.address, e.accounts.owner.privateKey)
    const pauser = createSignerWithAddresss(e.accounts.pauser.address, e.accounts.pauser.privateKey)
    const validators = e.accounts.validators.map(k => createSignerWithAddresss(k.address, k.privateKey))
    const av = e.accounts.available.map(k => createSignerWithAddresss(k.address, k.privateKey))
    return new SifchainAccounts(operator, owner, pauser, validators, av)
}
