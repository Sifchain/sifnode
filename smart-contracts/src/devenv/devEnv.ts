import {registry, singleton} from "tsyringe";
import * as childProcess from "child_process"
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import * as hre from "hardhat"

export abstract class ShellCommand {
    abstract run(): Promise<void>

    abstract cmd(): [string, string[]]

    abstract results(): Promise<EthereumResults>

    /**
     * A combination of run and results
     */
    go(): [Promise<void>, Promise<EthereumResults>] {
        return [this.run(), this.results()]
    }
}

@registry([{
    token: EthereumArguments,
    useValue: new EthereumArguments("localhost", 8545, 1, 1, 1)
}])
export class EthereumArguments {
    constructor(
        readonly host: string,
        readonly port: number,
        readonly nValidators: number,
        readonly networkId: number,
        readonly chainId: number,
    ) {
    }
}


export interface EthereumAccount {
    address: string
    privateKey: string
}

export interface EthereumAccounts {
    operator: string,
    owner: string,
    pauser: string,
    proxyAdmin: string,
    validators: string[],
    available: string[]
}

export interface EthereumResults {
    httpHost: string
    httpPort: number,
    chainId: number,  // note that hardhat doesn't believe networkId exists...
    accounts: EthereumAccounts
}

@singleton()
export class HardhatNodeRunner extends ShellCommand {
    constructor(
        readonly args: EthereumArguments
    ) {
        super();
    }

    cmd(): [string, string[]] {
        return ["node_modules/.bin/hardhat", [
            "node",
            "--hostname", this.args.host,
            "--port", this.args.port.toString()
        ]]
    }

    override run(): Promise<void> {
        return new Promise((resolve, reject) => {
            const [c, args] = this.cmd()
            const process = childProcess.spawn(c, args, {stdio: "inherit"})
            process.on("error", err => {
                reject(err)
            })
            process.on("exit", e => {
                if (e && e != 0) {
                    reject(e)
                } else
                    resolve()
            })
        })
    }

    override async results(): Promise<EthereumResults> {
        const hardhatSigners = await hre.ethers.getSigners()
        let ethereumAccounts = signerArrayToEthereumAccounts(hardhatSigners, this.args.nValidators)
        if (hre.network.config.chainId) {
            return {
                accounts: ethereumAccounts,
                httpHost: this.args.host,
                httpPort: this.args.port,
                chainId: hre.network.config.chainId
            }
        } else throw "unknown chainId"
    }
}

function signerArrayToEthereumAccounts(accounts: SignerWithAddress[], nValidators: number): EthereumAccounts {
    const [operator, owner, pauser, ...moreAccounts] = accounts.map(x => x.address)
    const validators = moreAccounts.slice(0, nValidators)
    const available = moreAccounts.slice(nValidators)
    return {
        proxyAdmin: operator,
        operator,
        owner,
        pauser,
        validators: validators,
        available,
    }
}
