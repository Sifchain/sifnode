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
