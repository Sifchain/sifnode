import * as ChildProcess from "child_process"
import { ShellCommand } from "./devEnv"
import { GolangResults } from "./golangBuilder";
import * as path from "path"
import * as fs from "fs";
import YAML from 'yaml'
import notifier from 'node-notifier';

export interface ValidatorValues {
  chain_id: string,
  node_id: string,
  ipv4_address: string,
  moniker: string,
  password: string,
  address: string,
  pub_key: string,
  mnemonic: string,
  validator_address: string,
  validator_consensus_address: string,
  is_seed: boolean,
}
export interface SifnodedResults {
  validatorValues: ValidatorValues[];
  adminAddress: string;
  process: ChildProcess.ChildProcess;
  tcpurl: string;
}

export class SifnodedRunner extends ShellCommand<SifnodedResults> {
  output: Promise<SifnodedResults>;
  private outputResolve: any;
  private sifnodedCommand: string;
  private sifgenCommand: string;

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
    this.output = new Promise<SifnodedResults>((res, _) => {
      this.outputResolve = res;
    });
    this.sifgenCommand = path.join(this.golangResults.goBin, "sifgen")
    this.sifnodedCommand = path.join(this.golangResults.goBin, "sifnoded")
  }

  cmd(): [string, string[]] {
    return ["sifgen", [
      "node"
    ]]
  }

  async sifgenNetworkCreate(): Promise<SifnodedResults> {
    // Missing mint amount. Although it has default value
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

    // sifgen network create
    const sifgenOutput = ChildProcess.execFileSync(
      this.sifgenCommand,
      sifgenArgs,
      { encoding: "utf8" }
    )

    // Debug log
    // TODO: Add formal loglevel aware logging
    console.log("SifgenOutput", sifgenOutput)

    const file = fs.readFileSync(this.networkConfigFile, 'utf8')
    const networkConfig: ValidatorValues[] = YAML.parse(file)
    let homeDir: string = "";

    for (const validator of networkConfig) {
      const moniker = validator["moniker"]
      const mnemonic = validator["mnemonic"]
      const password = validator["password"]
      let chainDir: string = path.join(
        this.networkDir,
        "validators",
        this.chainId,
        moniker
      )

      homeDir = path.join(chainDir, ".sifnoded")
      await this.addValidatorKeyToTestKeyring(
        moniker,
        mnemonic,
      )

      const valOperKey = await this.readValoperKey(moniker, chainDir)

      const stdout = await this.addGenesisValidator(chainDir, valOperKey)
      console.log("Add genesis validator output", stdout)
      console.log("Added validator", valOperKey)
    }

    // Sifchain start daemon
    // TODO: Unnecessary formatting
    let sifnodeadmincmd = `${this.sifnodedCommand} keys add sifnodeadmin --keyring-backend test --output json`;
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

    // TODO: Homedir would contain value of last assignment. Might need to be fixed when we support more than 1 acc
    ChildProcess.execSync(
      `${this.sifnodedCommand} add-genesis-account ${sifnodedAdminAddress} 100000000000000000000rowan --home ${homeDir}`,
      { encoding: "utf8" }
    ).trim()
    ChildProcess.execSync(
      `${this.sifnodedCommand} set-genesis-oracle-admin ${sifnodedAdminAddress} --home ${homeDir}`,
      { encoding: "utf8" }
    ).trim()
    ChildProcess.execSync(
      `${this.sifnodedCommand} set-genesis-whitelister-admin ${sifnodedAdminAddress} --home ${homeDir}`,
      { encoding: "utf8" }
    ).trim()
    ChildProcess.execSync(
      `${this.sifnodedCommand} set-gen-denom-whitelist ${this.whitelistFile} --home ${homeDir}`,
      { encoding: "utf8" }
    ).trim()
    let sifnodedDaemonCmd = `${this.sifnodedCommand} start --log_format json --minimum-gas-prices 0.5rowan --rpc.laddr tcp://0.0.0.0:26657 --home ${homeDir}`;

    const sifnoded = ChildProcess.spawn(
      sifnodedDaemonCmd,
      { shell: true, stdio: "inherit" }
    )

    sifnoded.on('exit', (code) => {
      notifier.notify({
        title: "Sifnoded Notice",
        message: `Sifnoded has just exited with exit code: ${code}`
      })
    });

    return {
      validatorValues: networkConfig,
      adminAddress: sifnodedAdminAddress,
      process: sifnoded,
      tcpurl: "tcp://0.0.0.0:26657"
    }
    //    return lastValueFrom(eventEmitterToObservable(sifnoded, "sifnoded"))
  }

  async addValidatorKeyToTestKeyring(moniker: string, mnemonic: string) {
    const sifnodedArgs = [
      "keys",
      "add",
      moniker,
      "--keyring-backend",
      "test",
      "--recover",
    ]

    console.log("Add Validator with mnemonics: ", mnemonic);

    let child = ChildProcess.execFileSync(
      this.sifnodedCommand,
      sifnodedArgs,
      {
        encoding: "utf8",
        shell: false,
        input: `${mnemonic}\n`
      }
    );
    console.log("Add Validator key to test ring output:", child)
  }

  // TODO: args Position
  async readValoperKey(moniker: string, homeDirectory: string): Promise<string> {
    const sifgenArgs = [
      "keys",
      "show",
      "-a",
      "--bech",
      "val",
      moniker,
      "--keyring-backend", "test",
      "--home", path.join(homeDirectory, ".sifnoded")
    ]
    return ChildProcess.execFileSync(
      this.sifnodedCommand,
      sifgenArgs,
      { encoding: "utf8" }
    ).trim()
  }

  // sifnoded add-genesis-validators $valoper --home $CHAINDIR/.sifnoded
  async addGenesisValidator(chainDir: string, valoper: string): Promise<string> {
    const sifgenArgs = [
      "add-genesis-validators",
      "1",
      valoper,
      "100",
      "--home", path.join(chainDir, ".sifnoded"),
    ]

    console.log("Add genesis validator")
    return ChildProcess.execFileSync(
      this.sifnodedCommand,
      sifgenArgs,
      { encoding: "utf8" }
    )
  }

  // TODO: This function is incomplete. it is extracted from sifchain_start_daemon.sh
  // Currently fails in CLI
  async whitelistValidators(moniker: string): Promise<string> {
    const sifnodedArgs = [
      "keys",
      "show",
      "--keyring-backend",
      "file",
      "-a",
      "--bech", "val",
      moniker,
    ]
    return "";
  }

  override async run(): Promise<void> {
    const output = await this.sifgenNetworkCreate();
    this.outputResolve(output)
  }

  override async results(): Promise<SifnodedResults> {
    return this.output;
  }
}
