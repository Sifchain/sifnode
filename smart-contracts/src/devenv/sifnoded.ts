import * as ChildProcess from "child_process"
import { ShellCommand } from "./devEnv"
import { GolangResults } from "./golangBuilder";
import * as path from "path"
import * as fs from "fs";
import YAML from 'yaml'

export interface ValidatorValues {
  chainID: string;
  nodeID: string;
  ipv4Address: string;
  moniker: string;
  password: string;
  address: string;
  pubKey: string;
  mnemonic: string;
  validatorAddress: string;
  validatorConsensusAddress: string;
  isSeed: boolean;
}
export interface SifnodedResults {
  validatorValues: ValidatorValues[];
  tcpurl: string;
}

export class SifnodedRunner extends ShellCommand<SifnodedResults> {
  constructor(
    readonly golangResults: GolangResults,
    readonly logfile = "/tmp/sifnoded.log",
    readonly rpcPort = 9000,
    readonly nValidators = 1,
    readonly chainId = "localnet",
    readonly networkConfigFile = "/tmp/sifnodedConfig.yml",
    readonly networkDir = "/tmp/sifnodedNetwork",
    readonly seedIpAddress = "10.10.1.1",
    readonly whitelistFile = "../test/integration/whitelisted-denoms.json"
  ) {
    super();
  }

  cmd(): [string, string[]] {
    return ["sifgen", [
      "node"
    ]]
  }

  async sifgenNetworkCreate() {
    const sifnodedCommand = path.join(this.golangResults.goBin, "sifnoded")
    const sifgenArgs = [
      "network",
      "create",
      this.chainId,
      this.nValidators.toString(),
      this.networkDir,
      this.seedIpAddress,
      this.networkConfigFile,
      "--keyring-backend",
      "test",
    ]

    await fs.promises.mkdir(this.networkDir, { recursive: true });

    const sifgenOutput = ChildProcess.execFileSync(
      path.join(this.golangResults.goBin, "sifgen"),
      sifgenArgs,
      { encoding: "utf8" }
    )
    const file = fs.readFileSync(this.networkConfigFile, 'utf8')
    const networkConfig = YAML.parse(file)
    const moniker = networkConfig[0]["moniker"]
    let mnemonic = networkConfig[0]["mnemonic"]
    let password = networkConfig[0]["password"]
    const chainDir = path.join(
      this.networkDir,
      "validators",
      this.chainId,
      moniker
    )
    const homeDir = path.join(chainDir, ".sifnoded")
    await this.addValidatorKeyToTestKeyring(
      moniker,
      this.networkDir,
      mnemonic,
    )
    const valOperKey = await this.readValoperKey(
      moniker,
      this.networkDir,
      mnemonic,
    )
    const stdout = await this.addGenesisValidator(chainDir, valOperKey)
    const whitelistedValidator = ChildProcess.execSync(
      `${sifnodedCommand} keys show -a --bech val ${moniker} --keyring-backend test`,
      { encoding: "utf8", input: password }
    ).trim()
    let sifnodeadmincmd = `${sifnodedCommand} keys add sifnodeadmin --keyring-backend test --output json`;
    const sifnodedadminJson = ChildProcess.execSync(
      sifnodeadmincmd,
      { encoding: "utf8", input: "yes\nyes" }
    ).trim()
    const sifnodedAdminAddress = JSON.parse(sifnodedadminJson)["address"]
    // const q = ChildProcess.execSync(
    //     `${sifnodedCommand} add-genesis-validators ${whitelistedValidator} --home ${homeDir}`,
    //     {encoding: "utf8", input: password}
    // ).trim()
    // sifnoded add-genesis-account $adminuser 100000000000000000000rowan --home $CHAINDIR/.sifnoded
    // sifnoded set-genesis-oracle-admin $adminuser --home $CHAINDIR/.sifnoded
    // sifnoded set-genesis-whitelister-admin $adminuser --home $CHAINDIR/.sifnoded
    // sifnoded set-gen-denom-whitelist $SCRIPT_DIR/whitelisted-denoms.json --home $CHAINDIR/.sifnoded
    ChildProcess.execSync(
      `${sifnodedCommand} add-genesis-account ${sifnodedAdminAddress} 100000000000000000000rowan --home ${homeDir}`,
      { encoding: "utf8" }
    ).trim()
    ChildProcess.execSync(
      `${sifnodedCommand} set-genesis-oracle-admin ${sifnodedAdminAddress} --home ${homeDir}`,
      { encoding: "utf8" }
    ).trim()
    ChildProcess.execSync(
      `${sifnodedCommand} set-genesis-whitelister-admin ${sifnodedAdminAddress} --home ${homeDir}`,
      { encoding: "utf8" }
    ).trim()
    ChildProcess.execSync(
      `${sifnodedCommand} set-gen-denom-whitelist ${this.whitelistFile} --home ${homeDir}`,
      { encoding: "utf8" }
    ).trim()
    let sifnodedDaemonCmd = `${sifnodedCommand} start --minimum-gas-prices 0.5rowan --rpc.laddr tcp://0.0.0.0:26657 --home ${homeDir}`;
    const sifnoded = ChildProcess.spawn(
      sifnodedDaemonCmd,
      { shell: true, stdio: "inherit" }
    )
    return
    //    return lastValueFrom(eventEmitterToObservable(sifnoded, "sifnoded"))
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
      path.join(this.golangResults.goBin, "sifnoded"),
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
      path.join(this.golangResults.goBin, "sifnoded"),
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
      path.join(this.golangResults.goBin, "sifnoded"),
      sifgenArgs,
      { encoding: "utf8" }
    )
  }

  // async execute() {
  //   await this.sifgenNetworkCreate()
  // }

  override async run(): Promise<void> {
    await this.sifgenNetworkCreate();
  }

  override async results(): Promise<SifnodedResults> {
    return Promise.resolve({
      validatorValues: [],
      tcpurl: ""
    })
  }
}
