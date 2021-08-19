import {registry, singleton} from "tsyringe";
import * as childProcess from "child_process"
import * as hre from "hardhat"
import {
    EthereumAccounts,
    EthereumAddressAndKey,
    EthereumResults,
    ShellCommand
} from "./devEnv"
import {GolangResults, GolangResultsPromise} from "./golangBuilder";
import * as path from "path"
import {SpawnSyncReturns} from "child_process";

@registry([
    {
        token: SifnodedArguments, useValue: new SifnodedArguments(
            "/tmp/sifnoded.log",
            9000,
            1,
            "localnet",
            "/tmp/sifnodedConfig.yml",
            "/tmp/sifnodedNetwork",
            "10.10.1.1"
        )
    }
])
export class SifnodedArguments {
    // "input": {
    //     "basedir": "/sifnode",
    //     "logfile": "/logs/ganache.log",
    //     "configoutputfile": "/configs/sifnoded.json",
    //     "rpc_port": 26657,
    //     "n_validators": 1,
    //     "chain_id": "localnet",
    //     "network_config_file": "/tmp/netconfig.yml",
    //     "seed_ip_address": "10.10.1.1",
    //     "bin_prefix": "/gobin",
    //     "go_build_config_path": "/configs/golang.json",
    //     "sifnode_host": "sifnoded"
    //   },
    constructor(
        readonly logfile: string,
        readonly rpcPort: number,
        readonly nValidators: number,
        readonly chainId: string,
        readonly networkConfigFile: string,
        readonly networkDir: string,
        readonly seedIpAddress: string,
    ) {
    }
}

interface SifnodedResults {

}

@singleton()
export class SifnodedRunner extends ShellCommand<SifnodedResults> {
    constructor(
        readonly args: SifnodedArguments,
        readonly golangResults: GolangResultsPromise
    ) {
        super();
    }

    cmd(): [string, string[]] {
        return ["sifgen", [
            "node"
        ]]
    }

    ensureCorrectExecution(result: SpawnSyncReturns<string>): SpawnSyncReturns<string> {
        if (result.error || result?.stderr != undefined) {
            console.log("error stdout: ", result.stdout)
            console.log("error stderr: ", result.stderr)
            throw result.error
        }
        return result
    }

    async sifgenNetworkCreate() {
        const sifgenArgs = [
            "network",
            "create",
            this.args.chainId,
            this.args.nValidators.toString(),
            this.args.networkDir,
            this.args.seedIpAddress,
            this.args.networkConfigFile
        ]
        this.ensureCorrectExecution(childProcess.spawnSync(
                path.join((await this.golangResults.results).goBin, "sifgen"),
                sifgenArgs,
                {encoding: "utf8"}
            )
        )
    }

    async execute() {
        await this.sifgenNetworkCreate()
    }

    override run(): Promise<void> {
        return this.execute().then(_ => this.run())
    }

    override async results(): Promise<SifnodedResults> {
        throw "not implemented"
    }
}
