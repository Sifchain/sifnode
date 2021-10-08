import * as ChildProcess from "child_process"
import { ShellCommand } from "./devEnv"
import { GolangResults } from "./golangBuilder";
import * as path from "path"
import * as fs from "fs";
import YAML from 'yaml'
import { boolean } from "yargs";
import notifier from 'node-notifier';
import { EbrelayerArguments } from "./ebrelayer";

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
export interface EbRelayerAccount {
  name: string;
  account: string;
}
export interface SifnodedResults {
  validatorValues: ValidatorValues[];
  relayerAddresses: EbRelayerAccount[];
  witnessAddresses: EbRelayerAccount[];
  adminAddress: EbRelayerAccount;
  process: ChildProcess.ChildProcess;
  tcpurl: string;
}

export class SifnodedRunner extends ShellCommand<SifnodedResults> {
  output: Promise<SifnodedResults>;
  private outputResolve: any;
  private sifnodedCommand: string;

  constructor(
    readonly golangResults: GolangResults,
    readonly logfile = "/tmp/sifnoded.log",
    readonly rpcPort = 9000,
    readonly nValidators = 1,
    readonly nRelayers = 1,
    readonly nWitnesses = 1,
    readonly chainId = "localnet",
    readonly networkConfigFile = "/tmp/sifnodedConfig.yml",
    readonly networkDir = "/tmp/sifnodedNetwork",
    readonly seedIpAddress = "10.10.1.1",
    readonly whitelistFile = "../test/integration/whitelisted-denoms.json"
  ) {
    super();
    this.sifnodedCommand = path.join(this.golangResults.goBin, "sifnoded")
    this.output = new Promise<SifnodedResults>((res, _) => {
      this.outputResolve = res;
    });
  }

  cmd(): [string, string[]] {
    return ["sifgen", [
      "node"
    ]]
  }

  async sifgenNetworkCreate(): Promise<SifnodedResults> {
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
    const networkConfig: ValidatorValues[] = YAML.parse(file)
    let chainDir: string = "";
    let homeDir: string = "";
    for (const validator of networkConfig) {
      const moniker = validator["moniker"]
      const mnemonic = validator["mnemonic"]
      const password = validator["password"]
      chainDir = path.join(
        this.networkDir,
        "validators",
        this.chainId,
        moniker
      )
      homeDir = path.join(chainDir, ".sifnoded")
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
        `${this.sifnodedCommand} keys show -a --bech val ${moniker} --keyring-backend test`,
        { encoding: "utf8", input: password }
      ).trim()
    }

    // Create an ADMIN account on sifnode with name sifnodeadmin
    const sifnodedAdminAddress = this.addRelayerWitnessAccount("sifnodeadmin", homeDir);
    // Create an account for each relayer as requested
    const relayerAddresses = Array.from({ length: this.nRelayers },
      (_, relayer) => this.addRelayerWitnessAccount(`relayer-${relayer}`, homeDir));
    // Create an account for each witness as requested
    const witnessAddresses = Array.from({ length: this.nWitnesses },
      (_, witness) => this.addRelayerWitnessAccount(`witness-${witness}`, homeDir));

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
      relayerAddresses: relayerAddresses,
      witnessAddresses: witnessAddresses,
      process: sifnoded,
      tcpurl: "tcp://0.0.0.0:26657"
    }
    //    return lastValueFrom(eventEmitterToObservable(sifnoded, "sifnoded"))
  }

  addRelayerWitnessAccount(name: string, homeDir: string): EbRelayerAccount {
    let accountAddCmd = `${this.sifnodedCommand} keys add ${name} --keyring-backend test --output json`;
    const accountJSON = ChildProcess.execSync(
      accountAddCmd,
      { encoding: "utf8", input: "yes\nyes" }
    ).trim()
    const accountAddress = JSON.parse(accountJSON)["address"]
    // const q = ChildProcess.execSync(
    //     `${sifnodedCommand} add-genesis-validators ${whitelistedValidator} --home ${homeDir}`,
    //     {encoding: "utf8", input: password}
    // ).trim()
    // sifnoded add-genesis-account $adminuser 100000000000000000000rowan --home $CHAINDIR/.sifnoded
    // sifnoded set-genesis-oracle-admin $adminuser --home $CHAINDIR/.sifnoded
    // sifnoded set-genesis-whitelister-admin $adminuser --home $CHAINDIR/.sifnoded
    // sifnoded set-gen-denom-whitelist $SCRIPT_DIR/whitelisted-denoms.json --home $CHAINDIR/.sifnoded
    ChildProcess.execSync(
      `${this.sifnodedCommand} add-genesis-account ${accountAddress} 100000000000000000000rowan --home ${homeDir}`,
      { encoding: "utf8" }
    ).trim()
    ChildProcess.execSync(
      `${this.sifnodedCommand} set-genesis-oracle-admin ${accountAddress} --home ${homeDir}`,
      { encoding: "utf8" }
    ).trim()
    ChildProcess.execSync(
      `${this.sifnodedCommand} set-genesis-whitelister-admin ${accountAddress} --home ${homeDir}`,
      { encoding: "utf8" }
    ).trim()
    ChildProcess.execSync(
      `${this.sifnodedCommand} set-gen-denom-whitelist ${this.whitelistFile} --home ${homeDir}`,
      { encoding: "utf8" }
    ).trim()
    return {
      account: accountAddress,
      name: name
    };
  }

  async addValidatorKeyToTestKeyring(moniker: string, chainDir: string, mnemonic: string) {
    const sifgenArgs = [
      "keys",
      "add",
      moniker,
      "--keyring-backend",
      "test",
      "--recover",
    ]
    let child = ChildProcess.execFileSync(
      path.join(this.golangResults.goBin, "sifnoded"),
      sifgenArgs,
      { encoding: "utf8", shell: false, input: `${mnemonic}\n` }
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
      "1",
      valoper,
      "100",
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
    const output = await this.sifgenNetworkCreate();
    this.outputResolve(output)
  }

  override async results(): Promise<SifnodedResults> {
    return this.output;
  }
}
