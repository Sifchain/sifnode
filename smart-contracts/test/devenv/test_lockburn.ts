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
import * as ChildProcess from "child_process"
import { EbRelayerAccount } from "../../src/devenv/sifnoded";
import {v4 as uuidv4} from 'uuid';

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

describe("lock of ethereum", () => {
    // This test only works when devenv is running, and that requires a connection to localhost
    expect(hardhat.network.name, "please use devenv").to.eq("localhost")

    const devEnvObject = readDevEnvObj("environment.json")
    // a generic sif address, nothing special about it
    const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")

    before('register HardhatRuntimeEnvironmentToken', async () => {
        container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    })

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

    // TODO: This ISNT an ebrelayer Account. it is a SIFACCOUNT
    function createTestSifAccount(): EbRelayerAccount {
        // Generate uuid
        let testSifAccountName = uuidv4();
        let cmd: string = `sifnoded keys add ${testSifAccountName} --home ${devEnvObject!.sifResults!.adminAddress!.homeDir} --keyring-backend test --output json 2>&1`
        let responseString: string = ChildProcess.execSync(cmd, { encoding: "utf8"})
        let responseJson = JSON.parse(responseString);

        console.log("CreateTestAccount Response: ", responseJson)
        return {
            name: responseJson.name,
            account: responseJson.address,
            homeDir: "",
        };
    }

    // TODO: Move all these sif TS SDK to it's own class. I think it should go to smart-contract/devenv
    // TODO: Rethink naming. SendToSif?
    function fundSifAccount(adminAccount: string, destination: string, amount: number, symbol: string, homeDir: string): void {
        // sifnoded tx bank send adminAccount testAccountToBeFunded --keyring-backend test --chain-id localnet concat(amount,symbol) --gas-prices=0.5rowan --gas-adjustment=1.5 --home <homeDir> --gas auto -y
        let sifnodedCmd: string = `sifnoded tx bank send ${adminAccount} ${destination} --keyring-backend test --chain-id localnet ${amount}${symbol} --gas-prices=0.5rowan --gas-adjustment=1.5 --home ${homeDir} --gas auto -y`
        let responseString: string = ChildProcess.execSync(sifnodedCmd, { encoding: "utf8"})
        let responseJson = JSON.parse(responseString);

        console.log("FundSifAccount response:", responseJson);

        return;
    }

    // TODO: This is placed here now because devObject is available in this scope
    async function sifTransfer(sender: string, destination: SignerWithAddress, amount: BigNumber,
        symbol: string,
        // TODO: What is correct value for corsschainfee?
        crossChainFee: string, netwrokDescriptor: number) {}

    async function executeSifBurn(sender: string, destination: SignerWithAddress,
        amount: BigNumber,
        symbol: string,
        // TODO: What is correct value for corsschainfee?
        crossChainFee: string, netwrokDescriptor: number) {

            // TODO: Move these out of this test function
            let testSifAccount: EbRelayerAccount = createTestSifAccount();
            fundSifAccount(devEnvObject!.sifResults!.adminAddress!.account, testSifAccount!.account, 10000000000, "ceth", devEnvObject!.sifResults!.adminAddress!.homeDir);
            fundSifAccount(devEnvObject!.sifResults!.adminAddress!.account, testSifAccount!.account, 10000000000, "rowan", devEnvObject!.sifResults!.adminAddress!.homeDir);


            let sifnodedCmd: string = `sifnoded tx ethbridge burn ${testSifAccount.account} ${destination.address} ${amount} ${symbol} ${crossChainFee} --network-descriptor ${netwrokDescriptor} --keyring-backend test --gas-prices=0.5rowan --gas-adjustment=1.5 --chain-id localnet --home ${devEnvObject!.sifResults!.adminAddress!.homeDir} --from ${testSifAccount.name} -y `

            console.log("Executing: ", sifnodedCmd);
            // let responseString = ChildProcess.execSync(sifnodedCmd,
            //     { encoding: "utf8" }
            // )
            // let responseJson = JSON.parse(responseString);

            // console.log("FundSifAccount response:", responseJson);

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

        expect(lv.transactionStep, `did not get CoinsSent, last step was ${JSON.stringify(lv, undefined, 2)}`).to.eq(TransactionStep.CoinsSent)
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



    it.only("should allow ceth to eth tx", async () => {
        const ethereumAccounts = await ethereumResultsToSifchainAccounts(devEnvObject.ethResults!, hardhat.ethers.provider)
        const factories = container.resolve(SifchainContractFactories)
        const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)
        const destinationEthereumAddress = ethereumAccounts.availableAccounts[0]
        const sendAmount = BigNumber.from(3500)

        // Use env to get validator address
        const sifAccount = devEnvObject!.sifResults!.validatorValues[0].address;
        const networkDescriptor = devEnvObject!.ethResults!.chainId;
        console.log("Hardhat network descriptor is: ", networkDescriptor);
        await executeSifBurn(sifAccount, destinationEthereumAddress, sendAmount, "ceth", "1", networkDescriptor)

    })
    it("should watch evmrelayer logs")
    it("should watch for evm events")
    it("should fail if evmrelayer gets an error")
})
