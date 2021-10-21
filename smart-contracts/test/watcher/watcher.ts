import chai from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe";
import {HardhatRuntimeEnvironmentToken} from "../../src/tsyringe/injectionTokens";
import * as hardhat from "hardhat";
import {BigNumber, Wallet} from "ethers";
import {ethereumResultsToSifchainAccounts, readDevEnvObj} from "../../src/tsyringe/devenvUtilities";
import {SifchainContractFactories} from "../../src/tsyringe/contracts";
import {buildDevEnvContracts} from "../../src/contractSupport";
import web3 from "web3";
import * as ethereumAddress from "../../src/ethereumAddress";
import {sifwatch} from "../../src/watcher/watcher";
import {lastValueFrom} from "rxjs";
import {filter} from "rxjs/operators";
import {EvmError} from "../../src/watcher/ebrelayer";

chai.use(solidity)

describe("watcher", () => {
    const devEnvObject = readDevEnvObj("environment.json")
    // a generic sif address, nothing special about it
    const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")

    before('register HardhatRuntimeEnvironmentToken', async () => {
        container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    })

    it("should get the accounts from devenv")

    it("should send a lock transaction", async () => {
        const ethereumAccounts = await ethereumResultsToSifchainAccounts(devEnvObject.ethResults!, hardhat.ethers.provider)
        const factories = container.resolve(SifchainContractFactories)
        const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)
        const sender1 = ethereumAccounts.availableAccounts[0]
        const smallAmount = BigNumber.from(1017)

        const evmRelayerEvents = sifwatch("/tmp/sifnode/evmrelayer.log")
        evmRelayerEvents.subscribe({
            next: x => {
                console.log(x)
            },
            error: e => console.log("goterror: ", e),
            complete: () => console.log("alldone")
        })
        evmRelayerEvents.pipe(
            filter(x => x instanceof EvmError)
        ).subscribe(t => {
            throw Error(`got error: ${JSON.stringify(t)}`)
        })

        await contracts.bridgeBank.connect(sender1).lock(
            recipient,
            ethereumAddress.eth.address,
            smallAmount,
            {
                value: smallAmount
            }
        )

        console.log("lock sent")

        const lv = await lastValueFrom(evmRelayerEvents)
    })

    it("should watch evmrelayer logs")
    it("should watch for evm events")
    it("should fail if evmrelayer gets an error")
})
