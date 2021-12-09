import * as chai from "chai"
import { expect } from "chai"
import { solidity } from "ethereum-waffle"
import { container } from "tsyringe"
import { HardhatRuntimeEnvironmentToken } from "../../src/tsyringe/injectionTokens"
import * as hardhat from "hardhat"
import { BigNumber } from "ethers"
import {
  ethereumResultsToSifchainAccounts,
  readDevEnvObj,
} from "../../src/tsyringe/devenvUtilities"
import { SifchainContractFactories } from "../../src/tsyringe/contracts"
import { buildDevEnvContracts, DevEnvContracts } from "../../src/contractSupport"
import web3 from "web3"
import * as ethereumAddress from "../../src/ethereumAddress"
import { SifEvent, SifHeartbeat, sifwatch, sifwatchReplayable } from "../../src/watcher/watcher"
import * as rxjs from "rxjs"
import {
  defer,
  distinctUntilChanged,
  lastValueFrom,
  Observable,
  scan,
  Subscription,
  takeWhile,
} from "rxjs"
import { EbRelayerEvmEvent } from "../../src/watcher/ebrelayer"
import { EthereumMainnetEvent } from "../../src/watcher/ethereumMainnet"
import { filter } from "rxjs/operators"
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import * as ChildProcess from "child_process"
import {
  EbRelayerAccount,
  crossChainFeeBase,
  crossChainLockFee,
  crossChainBurnFee,
} from "../../src/devenv/sifnoded"
import { v4 as uuidv4 } from "uuid"
import * as dotenv from "dotenv"
import "@nomiclabs/hardhat-ethers"
import deepEqual = require("deep-equal")
import { ethers } from "hardhat"
import { SifnodedAdapter } from "./sifnodedAdapter"
// import { SifnodedAdapter, SifnodedAdapter } from "./sifnodedAdapter"
// import { createTestSifAccount, fundSifAccount } from "./sifnodedAdapter"

// The hash value for ethereum on mainnet
const ethDenomHash = "sif5ebfaf95495ceb5a3efbd0b0c63150676ec71e023b1043c40bcaaf91c00e15b2"

chai.use(solidity)

const GWEI = Math.pow(10, 9)
const ETH = Math.pow(10, 18)

interface Failure {
  kind: "failure"
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
  fromEthereumAddress: string
  ethereumNonce: BigNumber
  denomHash: string
  ethereumLockBurnSequence: BigNumber
  transactionStep: TransactionStep
  uniqueId: string
}

enum TransactionStep {
  Initial = "Initial",
  SawLogLock = "SawLogLock",
  SawProphecyClaim = "SawProphecyClaim",
  SawEthbridgeClaimArray = "SawEthbridgeClaimArray",
  BroadcastTx = "BroadcastTx",
  EthBridgeClaimArray = "EthBridgeClaimArray",
  CreateEthBridgeClaim = "CreateEthBridgeClaim",
  AddNewTokenMetadata = "AddNewTokenMetadata",
  AddTokenMetadata = "AddTokenMetadata",

  AppendValidatorToProphecy = "AppendValidatorToProphecy",
  ProcessSuccessfulClaim = "ProcessSuccessfulClaim",
  CoinsSent = "CoinsSent",

  Burn = "Burn",
  GetTokenMetadata = "GetTokenMetadata",
  CosmosEvent = "CosmosEvent",
  SignProphecy = "SignProphecy",
  PublishedProphecy = "PublishedProphecy",
  LogBridgeTokenMint = "LogBridgeTokenMint",

  GetCrossChainFeeConfig = "GetCrossChainFeeConfig",
  SendCoinsFromAccountToModule = "SendCoinsFromAccountToModule",
  SetProphecy = "SetProphecy",
  // TODO: Burn coin and burn are confusing. One is receivng a burn msg, the other a cosmos burncoin call
  BurnCoins = "BurnCoins",
  PublishCosmosBurnMessage = "PublishCosmosBurnMessage",
  ReceiveCosmosBurnMessage = "ReceiveCosmosBurnMessage",

  // Witness
  WitnessSignProphecy = "WitnessSignProphecy",
  SetWitnessLockBurnNonce = "SetWitnessLockBurnNonce",

  ProphecyStatus = "ProphecyStatus",
  ProphecyClaimSubmitted = "ProphecyClaimSubmitted",
}

function isTerminalState(s: State) {
  switch (s.value.kind) {
    case "success":
    case "failure":
      return true
    default:
      return (
        s.transactionStep === TransactionStep.CoinsSent ||
        s.transactionStep === TransactionStep.ProphecyClaimSubmitted
      )
  }
}

function isNotTerminalState(s: State) {
  return !isTerminalState(s)
}

type VerbosityLevel = "summary" | "full" | "none"

function verbosityLevel(): VerbosityLevel {
  switch (process.env["VERBOSE"]) {
    case undefined:
      return "none"
    case "summary":
      return "summary"
    default:
      return "full"
  }
}

function attachDebugPrintfs<T>(xs: Observable<T>, verbosity: VerbosityLevel): Subscription {
  return xs.subscribe({
    next: (x) => {
      switch (verbosity) {
        case "full":
          console.log("DebugPrintf", JSON.stringify(x))
          break
        case "summary":
          const p = x as any
          console.log(
            `${p.currentHeartbeat}\t${p.transactionStep}\t${p.value?.kind}\t${p.value?.data?.kind}`
          )
          break
      }
    },
    error: (e) => console.log("goterror: ", e),
    complete: () => console.log("alldone"),
  })
}

function hasDuplicateNonce(a: EbRelayerEvmEvent, b: EbRelayerEvmEvent): boolean {
  return a.data.event.Nonce === b.data.event.Nonce
}

// const gobin = process.env["GOBIN"]

describe("lock and burn tests", () => {
  dotenv.config()
  // const INIT_STATE: State = {
  //   value: { kind: "initialState" },
  //   createdAt: 0,
  //   currentHeartbeat: 0,
  //   transactionStep: TransactionStep.Initial,
  // }
  // This test only works when devenv is running, and that requires a connection to localhost
  expect(hardhat.network.name, "please use devenv").to.eq("localhost")

  const devEnvObject = readDevEnvObj("environment.json")
  // a generic sif address, nothing special about it
  const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")
  const networkDescriptor = devEnvObject?.ethResults?.chainId ?? 31337

  // const factories = container.resolve(SifchainContractFactories)
  // const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)

  const sifnodedAdapter: SifnodedAdapter = new SifnodedAdapter(
    devEnvObject!.sifResults!.adminAddress!.homeDir,
    devEnvObject!.sifResults!.adminAddress!.account,
    process.env["GOBIN"]
  )

  before("register HardhatRuntimeEnvironmentToken", async () => {
    container.register(HardhatRuntimeEnvironmentToken, { useValue: hardhat })
  })

  function ensureCorrectTransition(
    acc: State,
    v: SifEvent,
    predecessor: TransactionStep | TransactionStep[],
    successor: TransactionStep
  ): State {
    var stepIsCorrect: boolean
    if (Array.isArray(predecessor)) {
      stepIsCorrect = (predecessor as string[]).indexOf(acc.transactionStep) >= 0
    } else {
      stepIsCorrect = predecessor === acc.transactionStep
    }
    if (stepIsCorrect) {
      // console.log("Setting transactionStep", successor)
      return {
        ...acc,
        value: v,
        createdAt: acc.currentHeartbeat,
        transactionStep: successor,
      }
    } else {
      // console.log("Step is incorrect", successor)
      return buildFailure(
        acc,
        v,
        `bad transition: expected ${predecessor}, got ${acc.transactionStep} before transition to ${successor}`
      )
    }
  }

  function buildFailure(acc: State, v: SifEvent, message: string): State {
    return {
      ...acc,
      value: {
        kind: "failure",
        value: v,
        message: message,
      },
    }
  }

  // Wrap an async function into an Observable<T>
  function deferAsync<T>(fn: () => Promise<T>): Observable<T> {
    return defer(() => rxjs.from(fn()))
  }

  async function executeLock(
    contracts: DevEnvContracts,
    smallAmount: BigNumber,
    sender1: SignerWithAddress,
    sifchainRecipient: string,
    verbose: boolean,
    identifier: string
  ) {
    const [evmRelayerEvents, replayedEvents] = sifwatchReplayable(
      {
        evmrelayer: "/tmp/sifnode/evmrelayer.log",
        sifnoded: "/tmp/sifnode/sifnoded.log",
      },
      hardhat,
      contracts.bridgeBank
    )

    const tx = await contracts.bridgeBank
      .connect(sender1)
      .lock(sifchainRecipient, ethereumAddress.eth.address, smallAmount, {
        value: smallAmount,
      })

    const states: Observable<State> = evmRelayerEvents
      .pipe(filter((x) => x.kind !== "SifnodedInfoEvent"))
      .pipe(
        scan(
          (acc: State, v: SifEvent) => {
            if (isTerminalState(acc))
              // we've reached a decision
              return { ...acc, value: { kind: "terminate" } as Terminate }
            switch (v.kind) {
              case "EbRelayerError":
              case "SifnodedError":
                // if we get an actual error, that's always a failure
                return { ...acc, value: { kind: "failure", value: v, message: "simple error" } }
              case "SifHeartbeat":
                // we just store the heartbeat
                return { ...acc, currentHeartbeat: v.value } as State
              case "EthereumMainnetLogLock":
                // we should see exactly one lock
                let ethBlock = v.data.block as any
                if (ethBlock.transactionHash === tx.hash && v.data.value.eq(smallAmount)) {
                  const newAcc: State = {
                    ...acc,
                    fromEthereumAddress: v.data.from,
                    ethereumNonce: BigNumber.from(v.data.nonce),
                  }
                  return ensureCorrectTransition(
                    newAcc,
                    v,
                    TransactionStep.Initial,
                    TransactionStep.SawLogLock
                  )
                }
                return {
                  ...acc,
                  value: {
                    kind: "failure",
                    value: v,
                    message: "incorrect EthereumMainnetLogLock",
                  },
                }
              case "EbRelayerEvmStateTransition":
                switch ((v.data as any).kind) {
                  case "EthereumProphecyClaim":
                    const d = v.data as any
                    if (
                      d.prophecyClaim.ethereum_sender == acc.fromEthereumAddress &&
                      BigNumber.from(d.event.Nonce).eq(acc.ethereumNonce)
                    ) {
                      return ensureCorrectTransition(
                        {
                          ...acc,
                          denomHash: d.prophecyClaim.denom_hash,
                        },
                        v,
                        TransactionStep.SawLogLock,
                        TransactionStep.SawProphecyClaim
                      )
                    }
                    break
                  case "EthBridgeClaimArray":
                    let claims = (v.data as any).claims as any[]
                    const matchingClaim = claims.find((claim) => claim.denom_hash === acc.denomHash)
                    if (matchingClaim)
                      return ensureCorrectTransition(
                        acc,
                        v,
                        TransactionStep.SawProphecyClaim,
                        TransactionStep.EthBridgeClaimArray
                      )
                    break
                  case "BroadcastTx":
                    const messages = (v.data as any).messages as any[]
                    const matchingMessage = messages.find(
                      (msg) => msg.eth_bridge_claim.denom_hash === acc.denomHash
                    )
                    if (matchingMessage)
                      return ensureCorrectTransition(
                        acc,
                        v,
                        TransactionStep.EthBridgeClaimArray,
                        TransactionStep.BroadcastTx
                      )
                }
              case "SifnodedPeggyEvent":
                switch ((v.data as any).kind) {
                  case "coinsSent":
                    const coins = ((v.data as any).coins as any)[0]
                    if (coins["denom"] === ethDenomHash && smallAmount.eq(coins["amount"]))
                      return ensureCorrectTransition(
                        acc,
                        v,
                        TransactionStep.ProcessSuccessfulClaim,
                        TransactionStep.CoinsSent
                      )
                    else return buildFailure(acc, v, "incorrect hash or amount")
                  // TODO these steps need validation to make sure they're happing in the right order with the right data
                  case "CreateEthBridgeClaim":
                    let newSequenceNumber = (v.data as any).msg.Interface.eth_bridge_claim
                      .ethereum_lock_burn_sequence
                    if (acc.ethereumNonce?.eq(newSequenceNumber))
                      return ensureCorrectTransition(
                        acc,
                        v,
                        [TransactionStep.BroadcastTx, TransactionStep.AppendValidatorToProphecy],
                        TransactionStep.CreateEthBridgeClaim
                      )
                    break
                  case "AppendValidatorToProphecy":
                    return ensureCorrectTransition(
                      acc,
                      v,
                      TransactionStep.CreateEthBridgeClaim,
                      TransactionStep.AppendValidatorToProphecy
                    )
                  case "ProcessSuccessfulClaim":
                    return ensureCorrectTransition(
                      acc,
                      v,
                      TransactionStep.AppendValidatorToProphecy,
                      TransactionStep.ProcessSuccessfulClaim
                    )
                  case "AddTokenMetadata":
                    return ensureCorrectTransition(
                      acc,
                      v,
                      TransactionStep.ProcessSuccessfulClaim,
                      TransactionStep.AddTokenMetadata
                    )
                }
                return { ...acc, value: v, createdAt: acc.currentHeartbeat }
              default:
                // we have a new value (of any kind) and it should use the current heartbeat as its creation time
                return { ...acc, value: v, createdAt: acc.currentHeartbeat }
            }
          },
          {
            value: { kind: "initialState" },
            createdAt: 0,
            currentHeartbeat: 0,
            transactionStep: TransactionStep.Initial,
            uniqueId: identifier,
          } as State
        )
      )

    // it's useful to skip debug prints of states where only the heartbeat changed
    const withoutHeartbeat = states.pipe(
      distinctUntilChanged<State>((a, b) => {
        return deepEqual({ ...a, currentHeartbeat: 0 }, { ...b, currentHeartbeat: 0 })
      })
    )

    const verboseSubscription = attachDebugPrintfs(withoutHeartbeat, verbosityLevel())

    const lv = await lastValueFrom(states.pipe(takeWhile((x) => x.value.kind !== "terminate")))

    expect(
      lv.transactionStep,
      `did not get CoinsSent, last step was ${JSON.stringify(lv, undefined, 2)}`
    ).to.eq(TransactionStep.CoinsSent)

    verboseSubscription.unsubscribe()
    replayedEvents.unsubscribe()
  }

  it.only("should allow ceth to eth tx", async () => {
    // TODO: Could these be moved out of the test fx? and instantiated via beforeEach?
    const factories = container.resolve(SifchainContractFactories)
    const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)

    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    )
    const destinationEthereumAddress = ethereumAccounts.availableAccounts[0]

    // These two can happen together
    const initialBalance = (
      await ethers.provider.getBalance(destinationEthereumAddress.address)
    ).toString()

    const contractInitialBalance = (
      await ethers.provider.getBalance(contracts.bridgeBank.address)
    ).toString()

    const sendAmount = BigNumber.from(3500 * GWEI) // 3500 gwei

    let testSifAccount: EbRelayerAccount = sifnodedAdapter.createTestSifAccount()
    // sifnodedAdapter.fundSifAccount(testSifAccount!.account, 10000000000, "rowan")

    // TODO: This is temporary. I think the right thing is for them to accept a verbose level
    let originalVerboseLevel: string | undefined = process.env["VERBOSE"]
    process.env["VERBOSE"] = "summary"
    // Need to have a burn of eth happen at least once or there's no data about eth in the token metadata
    await executeLock(
      contracts,
      sendAmount,
      ethereumAccounts.availableAccounts[0],
      web3.utils.utf8ToHex(testSifAccount.account),
      false,
      "ceth to eth"
    )

    const intermediateBalance = (
      await ethers.provider.getBalance(destinationEthereumAddress.address)
    ).toString()
    let contractIntermediateBalance = (
      await ethers.provider.getBalance(contracts.bridgeBank.address)
    ).toString()

    // These are temporarily added to make the logging lvl lower
    process.env["VERBOSE"] = originalVerboseLevel

    console.log("Lock complete")

    const evmRelayerEvents: rxjs.Observable<SifEvent> = sifwatch(
      {
        evmrelayer: "/tmp/sifnode/evmrelayer.log",
        sifnoded: "/tmp/sifnode/sifnoded.log",
        witness: "/tmp/sifnode/witness.log",
      },
      hardhat,
      contracts.bridgeBank
    ).pipe(filter((x) => x.kind !== "SifnodedInfoEvent"))

    evmRelayerEvents.subscribe((event) => console.log("Subscription", event))

    let receivedCosmosBurnmsg: boolean = false
    let witnessSignedProphecy: boolean = false
    const states: Observable<State> = evmRelayerEvents.pipe(
      scan(
        (acc: State, v: SifEvent) => {
          console.log("State assertion machine", v)
          if (isTerminalState(acc)) {
            // we've reached a decision
            console.log("Reached terminate state", acc)
            return { ...acc, value: { kind: "terminate" } as Terminate }
          }
          switch (v.kind) {
            case "EbRelayerError":
            case "SifnodedError":
              // if we get an actual error, that's always a failure
              return { ...acc, value: { kind: "failure", value: v, message: "simple error" } }
            case "SifHeartbeat": {
              // we just store the heartbeat
              return { ...acc, currentHeartbeat: v.value } as State
            }

            // Ebrelayer side log assertions
            case "EbRelayerEvmStateTransition": {
              let ebrelayerEvent: any = v.data
              switch (ebrelayerEvent.kind) {
                case "ReceiveCosmosBurnMessage": {
                  // console.log("Seeing ReceiveCosmosBurnMessage")
                  if (!receivedCosmosBurnmsg) {
                    // console.log("Receiving ReceiveCosmosBurnMessage for the first time")
                    // Ignore subsequence occurrences, witness will reprocess until keeper updates nonce
                    receivedCosmosBurnmsg = true
                    return ensureCorrectTransition(
                      acc,
                      v,
                      TransactionStep.PublishCosmosBurnMessage,
                      TransactionStep.ReceiveCosmosBurnMessage
                    )
                  } else {
                    return { ...acc, value: v, createdAt: acc.currentHeartbeat }
                  }
                }
                case "WitnessSignProphecy": {
                  // console.log("Seeing WitnessSignProphecy")
                  if (!witnessSignedProphecy) {
                    // console.log("Receiving WitnessSignProphecy for the first time")
                    witnessSignedProphecy = true
                    return ensureCorrectTransition(
                      acc,
                      v,
                      TransactionStep.ReceiveCosmosBurnMessage,
                      TransactionStep.WitnessSignProphecy
                    )
                  } else {
                    return { ...acc, value: v, createdAt: acc.currentHeartbeat }
                  }
                }

                case "ProphecyClaimSubmitted": {
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.ProphecyStatus,
                    TransactionStep.ProphecyClaimSubmitted
                  )
                }
              }
            }
            // Sifnoded side log assertions
            case "SifnodedPeggyEvent": {
              const sifnodedEvent: any = v.data
              switch (sifnodedEvent.kind) {
                case "Burn": {
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.Initial,
                    TransactionStep.Burn
                  )
                }

                // case "GetTokenMetadata":
                //   return ensureCorrectTransition(
                //     acc,
                //     v,
                //     TransactionStep.Burn,
                //     TransactionStep.GetTokenMetadata
                //   )

                case "GetCrossChainFeeConfig": {
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.Burn,
                    TransactionStep.GetCrossChainFeeConfig
                  )
                }

                case "SendCoinsFromAccountToModule": {
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.GetCrossChainFeeConfig,
                    TransactionStep.SendCoinsFromAccountToModule
                  )
                }

                case "BurnCoins": {
                  // TODO: Add assertion on expected amount, and expected denom
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.SendCoinsFromAccountToModule,
                    TransactionStep.BurnCoins
                  )
                }

                /**
                 * We comment this out because SetProphecy is the crUd operation, gets invoked multiple times throughout
                 * the call,
                 * But we still want to assert it has created a prophecy between BurnCoin and PublishCosmosBurnMessage
                 * TODO: Option 1. Refine the instrumentation statement in SetProphecy
                 *       Option 2. ???
                 */
                // case "SetProphecy":
                //   return ensureCorrectTransition(
                //     acc,
                //     v,
                //     TransactionStep.BurnCoins,
                //     TransactionStep.SetProphecy
                //   )

                case "PublishCosmosBurnMessage": {
                  // console.log("Received PublishCosmosBurnMessage")
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.BurnCoins,
                    TransactionStep.PublishCosmosBurnMessage
                  )
                }

                case "SetWitnessLockBurnNonce": {
                  // console.log("Receiving SetWitnessLockBurnNonce. Acc,", acc)
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.WitnessSignProphecy,
                    TransactionStep.SetWitnessLockBurnNonce
                  )
                }

                case "ProphecyStatus": {
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.SetWitnessLockBurnNonce,
                    TransactionStep.ProphecyStatus
                  )
                }
              }
            }

            default: {
              // we have a new value (of any kind) and it should use the current heartbeat as its creation time
              return { ...acc, value: v, createdAt: acc.currentHeartbeat }
            }
          }
        },
        {
          value: { kind: "initialState" },
          createdAt: 0,
          currentHeartbeat: 0,
          transactionStep: TransactionStep.Initial,
        } as State
      )
    )

    // it's useful to skip debug prints of states where only the heartbeat changed
    const withoutHeartbeat = states.pipe(
      distinctUntilChanged<State>((a, b) => {
        return deepEqual({ ...a, currentHeartbeat: 0 }, { ...b, currentHeartbeat: 0 })
      })
    )

    const verboseSubscription = attachDebugPrintfs(withoutHeartbeat, verbosityLevel())

    let crossChainCethFee = crossChainFeeBase * crossChainBurnFee

    let newSendAmount = BigNumber.from(2300 * Math.pow(10, 9)) // 2300 gwei
    await sifnodedAdapter.executeSifBurn(
      testSifAccount,
      destinationEthereumAddress,
      newSendAmount.sub(crossChainCethFee),
      ethDenomHash,
      String(crossChainCethFee),
      networkDescriptor
    )

    const lv = await lastValueFrom(states.pipe(takeWhile((x) => x.value.kind !== "terminate")))
    const expectedEndState: TransactionStep = TransactionStep.ProphecyClaimSubmitted
    expect(
      lv.transactionStep,
      `did not complete, last step was ${JSON.stringify(lv, undefined, 2)}`
    ).to.eq(expectedEndState)

    // Here we verify the user balance is correct
    const finalBalance = (
      await ethers.provider.getBalance(destinationEthereumAddress.address)
    ).toString()
    let contractFinalBalance = (
      await ethers.provider.getBalance(contracts.bridgeBank.address)
    ).toString()

    console.log("Initial Balance     ", initialBalance)
    console.log("intermediate Balance", intermediateBalance)
    console.log("final Balance       ", finalBalance)

    console.log("Contract Initial Balance     ", contractInitialBalance)
    console.log("Contract intermediate Balance", contractIntermediateBalance)
    console.log("Contract Final Balance       ", contractFinalBalance)

    verboseSubscription.unsubscribe()
  })

  it("should send two locks of ethereum", async () => {
    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    )
    const factories = container.resolve(SifchainContractFactories)
    const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)
    const sender1 = ethereumAccounts.availableAccounts[0]
    const smallAmount = BigNumber.from(1017)

    // Do two locks of ethereum
    await executeLock(contracts, smallAmount, sender1, recipient, true, "lock of eth")
    await executeLock(contracts, smallAmount, sender1, recipient, true, "second lock of eth")
  })
})
