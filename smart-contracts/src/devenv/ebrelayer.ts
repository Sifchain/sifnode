import * as ChildProcess from "child_process"
import { ShellCommand, EthereumAddressAndKey } from "./devEnv"
import { ValidatorValues } from "./sifnoded"
import { DeployedContractAddresses } from "../../scripts/deploy_contracts";
import notifier from 'node-notifier';
import * as path from "path"
import { GolangResults } from "./golangBuilder";

export interface EbrelayerArguments {
  readonly validatorValues: ValidatorValues,
  readonly account: EthereumAddressAndKey,
  readonly smartContract: DeployedContractAddresses,
  readonly golangResults: GolangResults
}

interface EbrelayerResults {
  process: ChildProcess.ChildProcess;
}

export class EbrelayerRunner extends ShellCommand<EbrelayerResults> {
  private output: Promise<EbrelayerResults>;
  private outputResolve: any;
  constructor(
    readonly args: EbrelayerArguments,
    readonly websocketAddress = "ws://localhost:7545/",
    readonly tcpURL = "tcp://0.0.0.0:26657",
    readonly chainNet = "localnet",
    readonly ebrelayerDB = `levelDB.db`,
    readonly relayerdbPath = "./relayerdb",
    readonly symbolTranslatorFile = "../test/integration/config/symbol_translator.json"
  ) {
    super();
    this.output = new Promise<EbrelayerResults>((res, rej) => {
      this.outputResolve = res;
    })
  }

  cmd(): [string, string[]] {
    return ["ebrelayer", [
      "init",
      this.tcpURL,
      this.websocketAddress,
      this.args.smartContract.bridgeRegistry,
      this.args.validatorValues.moniker,
      `'${this.args.validatorValues.mnemonic}'`,
      "--chain-id",
      String(this.chainNet),
      "--node",
      String(this.tcpURL),
      "--keyring-backend",
      "test",
      "--from",
      this.args.validatorValues.moniker,
      "--symbol-translator-file",
      this.symbolTranslatorFile,
      // "--relayerdb-path",
      // this.relayerdbPath
    ]]
  }

  async waitForSifAccount() {
    const scriptArgs = [
      "FirstOptionIsIgnored",
      this.args.validatorValues.address
    ]
    const child = ChildProcess.execFileSync(
      "./src/devenv/wait_for_sif_account.py",
      scriptArgs
    )
  }

  override async run(): Promise<void> {
    await this.waitForSifAccount()
    // const args: string[] = this.cmd()[1]// as string[]
    const spawncmd = path.join(this.args.golangResults.goBin, this.cmd()[0] + " " + this.cmd()[1].join(" "));
    process.env["ETHEREUM_PRIVATE_KEY"] = this.args.account.privateKey;
    process.env["ETHEREUM_ADDRESS"] = this.args.account.address;
    const commandResult = ChildProcess.spawn(
      spawncmd,
      {
        shell: true,
        stdio: "inherit",
      }
    )
    commandResult.on('exit', (code) => {
      notifier.notify({
        title: "Ebrelayer Notice",
        message: `Ebrelayer has just exited with exit code: ${code}`
      })
    })
    this.outputResolve(
      {
        process: commandResult
      }
    )
  }

  override async results(): Promise<EbrelayerResults> {
    return this.output;
  }
}
