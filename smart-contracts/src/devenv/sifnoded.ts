import {registry, singleton} from "tsyringe";
import * as ChildProcess from "child_process"
import {SpawnSyncReturns} from "child_process"
import {ShellCommand} from "./devEnv"
import {GolangResultsPromise} from "./golangBuilder";
import * as path from "path"
import events from "events";
import {lastValueFrom, ReplaySubject} from "rxjs";
import * as fs from "fs";
import YAML from 'yaml'

class ErrorEvent {
    constructor(readonly errorObject: any) {
    }
}

function eventEmitterToObservable(eventEmitter: events.EventEmitter) {
    const subject = new ReplaySubject<"exit" | ErrorEvent>(1)
    eventEmitter.on('error', e => {
        subject.error(new ErrorEvent(e))
    })
    eventEmitter.on('exit', e => {
        console.log("in eventEmitter")
        switch (e) {
            case 0:
                subject.next("exit")
                subject.complete()
                break
            default:
                subject.error(new ErrorEvent(e))
                break
        }
    })
    return subject.asObservable()
}

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
        if (result.error || (result?.stderr ?? "") != "") {
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
            "--keyring-backend",
            "test",
            this.args.chainId,
            this.args.nValidators.toString(),
            this.args.networkDir,
            this.args.seedIpAddress,
            this.args.networkConfigFile
        ]
        this.ensureCorrectExecution(ChildProcess.spawnSync(
                path.join((await this.golangResults.results).goBin, "sifgen"),
                sifgenArgs,
                {encoding: "utf8"}
            )
        )
        const file = fs.readFileSync(this.args.networkConfigFile, 'utf8')
        const networkConfig = YAML.parse(file)
        await this.addValidatorKeyToTestKeyring(
            networkConfig[0]["moniker"],
            this.args.networkDir,
            networkConfig[0]["mnemonic"],
        )
        console.log("finishedaddval")
    }

    // echo "$MNEMONIC" | sifnoded keys add $MONIKER --keyring-backend test --recover
    // valoper=$(sifnoded keys show -a --bech val $MONIKER --home $CHAINDIR/.sifnoded --keyring-backend test)
    // sifnoded add-genesis-validators $valoper --home $CHAINDIR/.sifnoded
    async addValidatorKeyToTestKeyring(moniker: string, chainDir: string, mnemonic: string) {
        const sifgenArgs = [
            "keys",
            "add",
            moniker,
            "--keyring-backend",
            "test",
        ]
        let child = ChildProcess.spawn(
            path.join((await this.golangResults.results).goBin, "sifnoded"),
            sifgenArgs,
            {stdio: "pipe"}
        );
        child.stdin.end(mnemonic)
        await lastValueFrom(eventEmitterToObservable(child))
    }

    async execute() {
        await this.sifgenNetworkCreate()
    }

    override run(): Promise<void> {
        console.log("inrun")
        return this.execute()
    }

    override async results(): Promise<SifnodedResults> {
        return Promise.resolve({})
    }
}
