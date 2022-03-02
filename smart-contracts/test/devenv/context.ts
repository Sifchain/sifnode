import * as chai from "chai"
import {solidity} from "ethereum-waffle"
import {SifEvent} from "../../src/watcher/watcher"
import {EthereumMainnetEvent} from "../../src/watcher/ethereumMainnet"
import {BigNumber} from "ethers"
import {Observable, Subscription} from "rxjs"
import {EbRelayerEvmEvent} from "../../src/watcher/ebrelayer"
import {sha256} from "ethers/lib/utils"

// zero contract address
export const nullContractAddress = "0x0000000000000000000000000000000000000000"

// The hash value for ethereum on mainnet
export const ethDenomHash = "sif5ebfaf95495ceb5a3efbd0b0c63150676ec71e023b1043c40bcaaf91c00e15b2"

chai.use(solidity)

const GWEI = Math.pow(10, 9)
const ETH = Math.pow(10, 18)

export interface Failure {
  kind: "failure"
  value: SifEvent | "timeout"
  message: string
}

export interface Success {
  kind: "success"
}

export interface InitialState {
  kind: "initialState"
}

export interface Terminate {
  kind: "terminate"
}

export interface State {
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

export enum TransactionStep {
  Initial = "Initial",
  SawLogLock = "SawLogLock",
  SawLogBurn = "SawLogBurn",
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
  Lock = "Lock",
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
  LockCoins = "LockCoins",
  PublishCosmosLockMessage = "PublishCosmosLockMessage",
  PublishCosmosBurnMessage = "PublishCosmosBurnMessage",
  ReceiveCosmosLockMessage = "ReceiveCosmosLockMessage",
  ReceiveCosmosBurnMessage = "ReceiveCosmosBurnMessage",

  // Witness
  WitnessSignProphecy = "WitnessSignProphecy",
  SetWitnessLockBurnNonce = "SetWitnessLockBurnNonce",

  ProphecyStatus = "ProphecyStatus",
  ProphecyClaimSubmitted = "ProphecyClaimSubmitted",

  EthereumMainnetLogUnlock = "EthereumMainnetLogUnlock",
  EthereumMainnetLogBridgeTokenMint = "EthereumMainnetLogBridgeTokenMint",
  EthereumMainnetNewProphecyClaim = "EthereumMainnetNewProphecyClaim",
  EthereumMainnetLogProphecyCompleted = "EthereumMainnetLogProphecyCompleted",
  EthereumMainnetLogNewBridgeTokenCreated = "EthereumMainnetLogNewBridgeTokenCreated",
}

// the last step is different for token import/export
export enum Direction {
  SifnodeToEthereum = "SifnodeToEthereum",
  EthereumToSifchain = "EthereumToSifchain",
}

export function isTerminalState(s: State, direction: Direction) {
  switch (s.value.kind) {
    case "success":
    case "failure":
      return true
    default:
      switch (direction) {
        case "SifnodeToEthereum":
          return (
            s.transactionStep === TransactionStep.CoinsSent ||
            s.transactionStep === TransactionStep.EthereumMainnetLogUnlock ||
            s.transactionStep === TransactionStep.EthereumMainnetLogProphecyCompleted
          )
          case "EthereumToSifchain":
            return (
              s.transactionStep === TransactionStep.CoinsSent
              // s.transactionStep === TransactionStep.ProcessSuccessfulClaim
            )
      }
  }
}

function isNotTerminalState(s: State, direction: Direction) {
  return !isTerminalState(s, direction)
}

type VerbosityLevel = "summary" | "full" | "none"

export function verbosityLevel(): VerbosityLevel {
  switch (process.env["VERBOSE"]) {
    case undefined:
      return "none"
    case "summary":
      return "summary"
    default:
      return "full"
  }
}

export function attachDebugPrintfs<T>(xs: Observable<T>, verbosity: VerbosityLevel): Subscription {
  return xs.subscribe({
    next: (x) => {
      switch (verbosity) {
        case "full": {
          console.log("DebugPrintf", JSON.stringify(x))
          break
        }
        case "summary": {
          const p = x as any
          console.log(
            `${p.currentHeartbeat}\t${p.transactionStep}\t${p.value?.kind}\t${p.value?.data?.kind}`
          )
          break
        }
      }
    },
    error: (e) => console.log("goterror: ", e),
    complete: () => console.log("alldone"),
  })
}

function hasDuplicateNonce(a: EbRelayerEvmEvent, b: EbRelayerEvmEvent): boolean {
  return a.data.event.Nonce === b.data.event.Nonce
}

export function ensureCorrectTransition(
  acc: State,
  v: SifEvent,
  predecessor: TransactionStep | TransactionStep[],
  successor: TransactionStep,
  skipPredecessor: boolean = false,
): State {
  
  var stepIsCorrect: boolean
  if (Array.isArray(predecessor)) {
    stepIsCorrect = (predecessor as string[]).indexOf(acc.transactionStep) >= 0
  } else {
    stepIsCorrect = predecessor === acc.transactionStep
  }
  stepIsCorrect = true
  if (stepIsCorrect || skipPredecessor) {
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

export function buildFailure(acc: State, v: SifEvent, message: string): State {
  return {
    ...acc,
    value: {
      kind: "failure",
      value: v,
      message: message,
    },
  }
}

export function getDenomHash(networkId: number, contract: string) {
  const data = String(networkId) + contract.toLowerCase()

  const enc = new TextEncoder()

  const denom = "sif" + sha256(enc.encode(data)).substring(2)

  return denom
}
