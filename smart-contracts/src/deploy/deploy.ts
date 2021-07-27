import * as hardhat from "hardhat"
import {container, instanceCachingFactory} from "tsyringe";
import {
    BridgeRegistryProxy, BridgeTokenSetup,
    CosmosBridgeArguments, CosmosBridgeArgumentsPromise,
    defaultCosmosBridgeArguments
} from "../tsyringe/contracts";
import {HardhatRuntimeEnvironmentToken} from "../tsyringe/injectionTokens";
import {SifchainAccounts, SifchainAccountsPromise} from "../tsyringe/sifchainAccounts";
import {BridgeToken} from "../../build";

async function main() {
    container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    const sifchainAccounts = container.resolve(SifchainAccountsPromise)
    container.register(CosmosBridgeArgumentsPromise, {
        useFactory: instanceCachingFactory<CosmosBridgeArgumentsPromise>(c => {
            const accountsPromise = c.resolve(SifchainAccountsPromise)
            return new CosmosBridgeArgumentsPromise(accountsPromise.accounts.then(accts => {
                return defaultCosmosBridgeArguments(accts)
            }))
        })
    })
    const bridgeRegistry = await container.resolve(BridgeRegistryProxy).contract
    const x = await bridgeRegistry.bridgeBank()
    const bts = container.resolve(BridgeTokenSetup)
    await bts.complete
}

main()
    .catch((error) => {
        console.error(error);
    })
    .finally(() => process.exit(0))
