import * as ChildProcess from "child_process"
import { ShellCommand } from "./devEnv"
import { ValidatorValues } from "./sifnoded"
import { DeployedContractAddresses } from "../../scripts/deploy_contracts";

export interface EbrelayerArguments {
  readonly validatorValues: ValidatorValues,
  readonly smartContract: DeployedContractAddresses
}

interface EbrelayerResults {
}

export class EbrelayerRunner extends ShellCommand<EbrelayerResults> {
  constructor(
    readonly args: EbrelayerArguments,
    readonly websocketAddress = "ws://localhost:7545/",
    readonly tcpURL = "tcp://0.0.0.0:26657",
    readonly chainNet = "localnet",
    readonly ebrelayerDB = `levelDB.db`,
    readonly relayerdbPath = "",
    readonly symbolTranslatorFile = "../test/integration/whitelisted-denom.json"
  ) {
    super();
  }

  cmd(): [string, string[]] {
    return ["ebrelayer", [
      "init",
      this.tcpURL,
      this.websocketAddress,
      this.args.smartContract.bridgeRegistry,
      this.args.validatorValues.moniker,
      this.args.validatorValues.mnemonic,
      `--chain-id ${this.chainNet}`,
      `--node ${this.tcpURL}`,
      "--keyring-backend test",
      `--from ${this.args.validatorValues.moniker}`,
      `--symbol-translator-file ${this.symbolTranslatorFile}`,
      `--relayerdb-path ${this.relayerdbPath}`
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
