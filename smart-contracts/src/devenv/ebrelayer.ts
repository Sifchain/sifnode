import * as ChildProcess from "child_process"
import { ShellCommand } from "./devEnv"
import { ValidatorValues } from "./sifnoded"
import { DeployedContractAddresses } from "../../scripts/deploy_contracts";

export interface EbrelayerArguments {
  readonly websocketAddress: string,
  readonly tcpURL: string,
  readonly chainNet: string,
  readonly ebrelayerDB: string,
  readonly validatorValues: ValidatorValues,
  readonly symbolTranslatorFile: string,
  readonly relayerdbPath: string,
  readonly smartContract: DeployedContractAddresses
}

interface EbrelayerResults {
}

export class EbrelayerRunner extends ShellCommand<EbrelayerResults> {
  constructor(
    readonly args: EbrelayerArguments,
  ) {
    super();
  }

  cmd(): [string, string[]] {
    return ["ebrelayer", [
      "init",
      this.args.tcpURL,
      this.args.websocketAddress,
      this.args.smartContract.bridgeRegistry,
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
