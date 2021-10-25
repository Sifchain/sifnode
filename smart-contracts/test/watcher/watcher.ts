import chai, {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe";
import {HardhatRuntimeEnvironmentToken} from "../../src/tsyringe/injectionTokens";
import * as hardhat from "hardhat";
import {BigNumber} from "ethers";
import {ethereumResultsToSifchainAccounts, readDevEnvObj} from "../../src/tsyringe/devenvUtilities";
import {SifchainContractFactories} from "../../src/tsyringe/contracts";
import {buildDevEnvContracts} from "../../src/contractSupport";
import web3 from "web3";
import * as ethereumAddress from "../../src/ethereumAddress";
import {SifEvent, SifHeartbeat, sifwatch} from "../../src/watcher/watcher";
import {lastValueFrom, Observable, scan, takeUntil, takeWhile} from "rxjs";

chai.use(solidity)

interface Failure {
    kind: "failure",
    value: SifEvent | "timeout"
}

interface Success {
    kind: "success"
}

interface InitialState {
    kind: "initialState"
}

interface Terminate {
    kind: "terminate"
}

interface State {
    value: SifEvent | Success | Failure | InitialState | Terminate
    createdAt: number
    currentHeartbeat: number
}

function isTerminalState(s: State) {
    switch (s.value.kind) {
        case "success":
        case "failure":
            return true
        default:
            return false
    }
}

function isNotTerminalState(s: State) {
    return !isTerminalState(s)
}

function attachDebugPrintfs<T>(xs: Observable<T>) {
    xs.subscribe({
        next: x => {
            console.log(JSON.stringify(x, undefined, 2))
        },
        error: e => console.log("goterror: ", e),
        complete: () => console.log("alldone")
    })
}

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

        const evmRelayerEvents = sifwatch({evmrelayer: "/tmp/sifnode/evmrelayer.log"})

        const states: Observable<State> = evmRelayerEvents.pipe(scan((acc: State, v: SifEvent) => {
            if (isTerminalState(acc))
                return {...acc, value: {kind: "terminate"} as Terminate}
            else if (v.kind == "EbRelayerError") {
                return {...acc, value: {kind: "failure", value: v}} as State
            } else if (v.kind === "SifHeartbeat") {
                if (acc.value.kind === "EbRelayerEvmEvent" && acc.createdAt + 5 < v.value) {
                    return {...acc, value: {kind: "failure", value: "timeout"}} as State
                }
                return {...acc, currentHeartbeat: v.value} as State
            } else if (v.kind == "EbRelayerEvmStateTransition" || v.kind === "EbRelayerEvmEvent") {
                return {...acc, value: v as SifEvent, createdAt: acc.currentHeartbeat} as State
            } else if (v.kind == "EbRelayerEthBridgeClaimArray") {
                if (v.data.claims[0].amount === smallAmount.toString()) {
                    return {...acc, value: {kind: "success"}} as State
                } else {
                    return {...acc, value: v, createdAt: acc.currentHeartbeat} as State
                }
            } else {
                return acc as State
            }
        }, {value: {kind: "initialState"}, createdAt: 0, currentHeartbeat: 0} as State))

        attachDebugPrintfs(evmRelayerEvents)
        attachDebugPrintfs(states)

        await contracts.bridgeBank.connect(sender1).lock(
            recipient,
            ethereumAddress.eth.address,
            smallAmount,
            {
                value: smallAmount
            }
        )

        console.log("lock sent")

        const lv = await lastValueFrom(states.pipe(takeWhile(x => x.value.kind !== "terminate")))

        console.debug("lastValueIs: ", JSON.stringify(lv, undefined, 2))

        expect((lv as State).value.kind).to.eq("success")
    })

    it("should watch evmrelayer logs")
    it("should watch for evm events")
    it("should fail if evmrelayer gets an error")
})
