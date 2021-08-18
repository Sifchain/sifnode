import {registry, singleton} from "tsyringe";
import * as childProcess from "child_process"
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import * as hre from "hardhat"

export abstract class ShellCommand {
    public run(): Promise<void> {
        let cmd = this.cmd();
        const result = childProcess.execSync(cmd)
        console.log("resultis: ", result.toString("UTF-8"))
        return Promise.resolve()
    }

    abstract cmd(): string

    abstract results(): Promise<EthereumResults>
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
    httpPort: number
    accounts: EthereumAccounts
}

@singleton()
export class HardhatNodeRunner extends ShellCommand {
    constructor(
        readonly args: EthereumArguments
    ) {
        super();
    }

    cmd(): string {
        return `node_modules/.bin/hardhat node --hostname ${this.args.host} --port ${this.args.port}`
    }

    override run(): Promise<void> {
        return new Promise((resolve, reject) => {
            childProcess.exec(this.cmd(), (err, stdout, stderr) => {
                if (err == null)
                    resolve()
                else
                    reject(err)
            })
        })
    }

    override async results(): Promise<EthereumResults> {
        const hardhatSigners = await hre.ethers.getSigners()
        let ethereumAccounts = signerArrayToEthereumAccounts(hardhatSigners, this.args.nValidators);
        return new class implements EthereumResults {
            accounts: EthereumAccounts = ethereumAccounts,
            httpHost: "fnord"
            httpPort: 1
        }
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
