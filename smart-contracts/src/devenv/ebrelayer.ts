import { registry, injectable } from "tsyringe";
import * as ChildProcess from "child_process"
import { SpawnSyncReturns } from "child_process"
import { ShellCommand } from "./devEnv"
import { GolangResultsPromise } from "./golangBuilder";
import * as path from "path"
import events from "events";
import { lastValueFrom, ReplaySubject } from "rxjs";
import { ValidatorValues } from "./sifnoded"
import { DeployedContractAddresses } from "./smartcontractDeployer";
import * as fs from "fs";
import YAML from 'yaml'
import { eventEmitterToObservable } from "./devEnvUtilities"

@registry([
  {
    token: EbrelayerArguments, useValue: new EbrelayerArguments(
      "/tmp/sifnoded.log",
      9000,
      1,
      "localnet",
      "/tmp/sifnodedConfig.yml",
      "/tmp/sifnodedNetwork",
      "10.10.1.1",
      "../test/integration/whitelisted-denoms.json"
    )
  }
])
export class EbrelayerArguments {
  constructor(
    readonly websocketAddress: string,
    readonly bridgeRegistryAddress: string,
    // Interface in hardhatNode readonly bridgeRegistryAddress: string,
    readonly tcpURL: string,
    readonly chainNet: number,
    readonly ebrelayerDB: string,
    readonly validatorValues: ValidatorValues,
    readonly symbolTranslatorFile: string,
    readonly relayerdbPath: string,
    readonly contractAddress: DeployedContractAddresses
  ) {
  }
}

interface EbrelayerResults {
}

@injectable()
export class EbrelayerRunner extends ShellCommand<EbrelayerResults> {
  constructor(
    readonly args: EbrelayerArguments,
    readonly golangResults: GolangResultsPromise
  ) {
    super();
  }

  cmd(): [string, string[]] {
    return ["ebrelayer", [
      "init",
      this.args.tcpURL,
      this.args.websocketAddress,
      this.args.contractAddress.bridgeRegistry,
      this.args.validatorValues.moniker,
      this.args.validatorValues.mnemonic,
      `--chain-id ${this.args.chainNet}`,
      `--node ${this.args.tcpURL}`,
      "--keyring-backend test",
      `--from ${this.args.validatorValues.moniker}`,
      `--symbol-translator-file ${this.args.symbolTranslatorFile}`,
      `--relayerdb-path ${this.args.relayerdbPath}`
    ]]
  }

  async waitForSifAccount() {
    const scriptArgs = [
      "FirstOptionIsIgnored",
      this.args.validatorValues.address
    ]
    const child = ChildProcess.execFileSync(
      "./wait_for_sif_account.py",
      scriptArgs
    )

  }

  async addValidatorKeyToTestKeyring(moniker: string, chainDir: string, mnemonic: string) {
    const sifgenArgs = [
      "keys",
      "add",
      moniker,
      "--keyring-backend",
      "test",
    ]
    let child = ChildProcess.execFileSync(
      path.join((await this.golangResults.results).goBin, "sifnoded"),
      sifgenArgs,
      { input: mnemonic, encoding: "utf8" }
    );
    child
  }

  async readValoperKey(moniker: string, chainDir: string, mnemonic: string): Promise<string> {
    const sifgenArgs = [
      "keys",
      "show",
      "-a",
      "--bech",
      "val",
      moniker,
      // "--home",
      // path.join(chainDir, ".sifnoded"),
      "--keyring-backend",
      "test",
    ]
    return ChildProcess.execFileSync(
      path.join((await this.golangResults.results).goBin, "sifnoded"),
      sifgenArgs,
      { encoding: "utf8" }
    ).trim()
  }

  // sifnoded add-genesis-validators $valoper --home $CHAINDIR/.sifnoded
  async addGenesisValidator(chainDir: string, valoper: string): Promise<string> {
    const sifgenArgs = [
      "add-genesis-validators",
      valoper,
      "--home",
      path.join(chainDir, ".sifnoded"),
    ]
    return ChildProcess.execFileSync(
      path.join((await this.golangResults.results).goBin, "sifnoded"),
      sifgenArgs,
      { encoding: "utf8" }
    )
  }

  async execute() {
    await this.waitForSifAccount()
    const args = this.cmd().slice(1) as string[]
    const commandResult = ChildProcess.spawnSync(this.cmd()[0], args)
  }

  override run(): Promise<void> {
    return this.execute()
  }

  override async results(): Promise<EbrelayerResults> {
    return Promise.resolve({})
  }
}
