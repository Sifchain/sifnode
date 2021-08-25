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
      "ws://localhost:7545/",
      "tcp://0.0.0.0:26657",
      "localnet",
      "levelDB.db",
      {
        address: "",
        chainID: "",
        ipv4Address: "",
        isSeed: false,
        mnemonic: "",
        moniker: "",
        nodeID: "",
        password: "",
        pubKey: "",
        validatorAddress: "",
        validatorConsensusAddress: ""
      },
      "../test/integration/whitelisted-denoms.json",
      "",
      {
        bridgeBank: "",
        bridgeRegistry: "",
        rowanContract: ""
      }
    )
  }
])
export class EbrelayerArguments {
  constructor(
    readonly websocketAddress: string,
    readonly tcpURL: string,
    readonly chainNet: string,
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

  override async run(): Promise<void> {
    await this.waitForSifAccount()
    const args = this.cmd().slice(1) as string[]
    const commandResult = ChildProcess.spawn(this.cmd()[0], args, { stdio: "inherit" })
    return
  }

  override async results(): Promise<EbrelayerResults> {
    return Promise.resolve({})
  }
}
