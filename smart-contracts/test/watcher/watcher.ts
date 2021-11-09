import chai, {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe";
import {HardhatRuntimeEnvironmentToken} from "../../src/tsyringe/injectionTokens";
import * as hardhat from "hardhat";
import {BigNumber} from "ethers";
import {ethereumResultsToSifchainAccounts, readDevEnvObj} from "../../src/tsyringe/devenvUtilities";
import {SifchainContractFactories} from "../../src/tsyringe/contracts";
import {buildDevEnvContracts, DevEnvContracts} from "../../src/contractSupport";
import web3 from "web3";
import * as ethereumAddress from "../../src/ethereumAddress";
import {SifEvent, SifHeartbeat, sifwatch} from "../../src/watcher/watcher";
import {distinctUntilChanged, lastValueFrom, Observable, scan, takeWhile} from "rxjs";
import {EbRelayerEvmEvent} from "../../src/watcher/ebrelayer";
import {EthereumMainnetEvent} from "../../src/watcher/ethereumMainnet";
import {filter} from "rxjs/operators";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import deepEqual = require("deep-equal");

// The hash value for ethereum on mainnet
const ethDenomHash = "sif5ebfaf95495ceb5a3efbd0b0c63150676ec71e023b1043c40bcaaf91c00e15b2"

chai.use(solidity)

interface Failure {
    kind: "failure",
    value: SifEvent | "timeout"
    message: string
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
    value: SifEvent | EthereumMainnetEvent | Success | Failure | InitialState | Terminate
    createdAt: number
    currentHeartbeat: number
    transactionStep: TransactionStep
}

enum TransactionStep {
    Initial = "Initial",
    SawLogLock = "SawLogLock",
    SawProphecyClaim = "SawProphecyClaim",
    SawEthbridgeClaimArray = "SawEthbridgeClaimArray",
    BroadcastTx = "BroadcastTx",
    CreateEthBridgeClaim = "CreateEthBridgeClaim",
    AppendValidatorToProphecy = "AppendValidatorToProphecy",
    ProcessSuccessfulClaim = "ProcessSuccessfulClaim",
    CoinsSent = "CoinsSent",
}

function isTerminalState(s: State) {
    switch (s.value.kind) {
        case "success":
        case "failure":
            return true
        default:
            return s.transactionStep === TransactionStep.CoinsSent
    }
}

function isNotTerminalState(s: State) {
    return !isTerminalState(s)
}

function attachDebugPrintfs<T>(xs: Observable<T>, summary: boolean) {
    xs.subscribe({
        next: x => {
            const p = x as any
            if (summary)
                console.log(`${p.currentHeartbeat}\t${p.transactionStep}\t${p.value?.kind}\t${p.value?.data?.kind}`)
            else
                console.log(JSON.stringify(x))
        },
        error: e => console.log("goterror: ", e),
        complete: () => console.log("alldone")
    })
}

function hasDuplicateNonce(a: EbRelayerEvmEvent, b: EbRelayerEvmEvent): boolean {
    return a.data.event.Nonce === b.data.event.Nonce
}

describe("watcher", () => {
    const devEnvObject = readDevEnvObj("environment.json")
    // a generic sif address, nothing special about it
    const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")

    before('register HardhatRuntimeEnvironmentToken', async () => {
        container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    })

    it("should get the accounts from devenv")

    function ensureCorrectTransition(acc: State, v: SifEvent, predecessor: TransactionStep | TransactionStep[], successor: TransactionStep): State {
        var stepIsCorrect: boolean
        if (Array.isArray(predecessor)) {
            stepIsCorrect = (predecessor as string[]).indexOf(acc.transactionStep) >= 0
        } else {
            stepIsCorrect = predecessor === acc.transactionStep
        }
        if (stepIsCorrect)
            return {
                ...acc,
                value: v,
                createdAt: acc.currentHeartbeat,
                transactionStep: successor
            }
        else
            return buildFailure(acc, v, `bad transition: expected ${predecessor}, got ${acc.transactionStep} before transition to ${successor}`)
    }

    function buildFailure(acc: State, v: SifEvent, message: string): State {
        return {
            ...acc,
            value: {
                kind: "failure",
                value: v,
                message: message
            }
        }
    }

    async function executeLock(contracts: DevEnvContracts, smallAmount: BigNumber, sender1: SignerWithAddress) {
        const evmRelayerEvents = sifwatch({
            evmrelayer: "/tmp/sifnode/evmrelayer.log",
            sifnoded: "/tmp/sifnode/sifnoded.log"
        }, hardhat, contracts.bridgeBank).pipe(filter(x => x.kind !== "SifnodedInfoEvent"))

        const states: Observable<State> = evmRelayerEvents.pipe(scan((acc: State, v: SifEvent) => {
            if (isTerminalState(acc))
                // we've reached a decision
                return {...acc, value: {kind: "terminate"} as Terminate}
            switch (v.kind) {
                case "EbRelayerError":
                case "SifnodedError":
                    // if we get an actual error, that's always a failure
                    return {...acc, value: {kind: "failure", value: v, message: "simple error"}}
                case "SifHeartbeat":
                    // we just store the heartbeat
                    return {...acc, currentHeartbeat: v.value} as State
                case "EthereumMainnetLogLock":
                    // we should see exactly one lock
                    if (v.data.value.eq(smallAmount) && acc.transactionStep == TransactionStep.Initial)
                        return {...acc, value: v, transactionStep: TransactionStep.SawLogLock}
                    else
                        return {
                            ...acc,
                            value: {
                                kind: "failure",
                                value: v,
                                message: "incorrect EthereumMainnetLogLock"
                            }
                        }
                case "EbRelayerEvmStateTransition":
                    switch ((v.data as any).kind) {
                        case "EthereumProphecyClaim":
                            return {
                                ...acc,
                                value: v,
                                transactionStep: TransactionStep.SawProphecyClaim
                            }
                        case "EthBridgeClaimArray":
                            return {
                                ...acc,
                                value: v,
                                transactionStep: TransactionStep.SawEthbridgeClaimArray
                            }
                        case "BroadcastTx":
                            return {...acc, value: v, transactionStep: TransactionStep.BroadcastTx}
                    }
                case "SifnodedPeggyEvent":
                    switch ((v.data as any).kind) {
                        case "coinsSent":
                            const coins = ((v.data as any).coins as any)[0]
                            if (coins["denom"] === ethDenomHash && smallAmount.eq(coins["amount"]))
                                return ensureCorrectTransition(acc, v, TransactionStep.ProcessSuccessfulClaim, TransactionStep.CoinsSent)
                            else
                                return buildFailure(acc, v, "incorrect hash or amount")
                        // TODO these steps need validation to make sure they're happing in the right order with the right data
                        case "CreateEthBridgeClaim":
                            return ensureCorrectTransition(
                                acc,
                                v,
                                [TransactionStep.BroadcastTx, TransactionStep.AppendValidatorToProphecy],
                                TransactionStep.CreateEthBridgeClaim
                            )
                        case "AppendValidatorToProphecy":
                            return ensureCorrectTransition(acc, v, TransactionStep.CreateEthBridgeClaim, TransactionStep.AppendValidatorToProphecy)
                        case "ProcessSuccessfulClaim":
                            return ensureCorrectTransition(acc, v, TransactionStep.AppendValidatorToProphecy, TransactionStep.ProcessSuccessfulClaim)
                    }
                    return {...acc, value: v, createdAt: acc.currentHeartbeat}
                default:
                    // we have a new value (of any kind) and it should use the current heartbeat as its creation time
                    return {...acc, value: v, createdAt: acc.currentHeartbeat}
            }
        }, {
            value: {kind: "initialState"},
            createdAt: 0,
            currentHeartbeat: 0,
            transactionStep: TransactionStep.Initial
        } as State))

        // it's useful to skip debug prints of states where only the heartbeat changed
        const withoutHeartbeat = states.pipe(distinctUntilChanged<State>((a, b) => {
            return deepEqual({...a, currentHeartbeat: 0}, {...b, currentHeartbeat: 0})
        }))

        // attachDebugPrintfs(evmRelayerEvents)
        // attachDebugPrintfs(evmRelayerEvents.pipe(filter(isNotSifnodedEvent)))
        attachDebugPrintfs(withoutHeartbeat, true)

        await contracts.bridgeBank.connect(sender1).lock(
            recipient,
            ethereumAddress.eth.address,
            smallAmount,
            {
                value: smallAmount
            }
        )

        const lv = await lastValueFrom(states.pipe(takeWhile(x => x.value.kind !== "terminate")))

        console.debug("last step was: ", JSON.stringify(lv, undefined, 2))

        expect(lv.transactionStep).to.eq(TransactionStep.CoinsSent)
    }

    it("should send two locks of ethereum", async () => {
        const ethereumAccounts = await ethereumResultsToSifchainAccounts(devEnvObject.ethResults!, hardhat.ethers.provider)
        const factories = container.resolve(SifchainContractFactories)
        const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)
        const sender1 = ethereumAccounts.availableAccounts[0]
        const smallAmount = BigNumber.from(1017)

        // Do two locks of ethereum
        await executeLock(contracts, smallAmount, sender1);
        await executeLock(contracts, smallAmount, sender1);
    })

    it("should watch evmrelayer logs")
    it("should watch for evm events")
    it("should fail if evmrelayer gets an error")
})
