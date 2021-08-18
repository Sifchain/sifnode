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

export interface EthereumAddressAndKey {
    privateKey: string
    address: string
}

export interface EthereumAccounts {
    operator: EthereumAddressAndKey,
    owner: EthereumAddressAndKey,
    pauser: EthereumAddressAndKey,
    proxyAdmin: EthereumAddressAndKey,
    validators: EthereumAddressAndKey[],
    available: EthereumAddressAndKey[]
}

export interface EthereumResults {
    httpHost: string
    httpPort: number,
    chainId: number,  // note that hardhat doesn't believe networkId exists...
    accounts: EthereumAccounts
}
