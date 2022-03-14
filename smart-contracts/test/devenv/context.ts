import * as chai from "chai"
import {solidity} from "ethereum-waffle"
import {SifEvent} from "../../src/watcher/watcher"
import {EthereumMainnetEvent} from "../../src/watcher/ethereumMainnet"
import {BigNumber} from "ethers"
import {Observable, Subscription} from "rxjs"
import {EbRelayerEvmEvent} from "../../src/watcher/ebrelayer"
import {sha256} from "ethers/lib/utils"

// The hash value for ethereum on mainnet
export const ethDenomHash = "sifBridge99990x0000000000000000000000000000000000000000"

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

  EthereumMainnetLogUnlock = "EthereumMainnetLogUnlock",
}

export function isTerminalState(s: State) {
  switch (s.value.kind) {
    case "success":
    case "failure":
      return true
    default:
      return (
        s.transactionStep === TransactionStep.CoinsSent ||
        s.transactionStep === TransactionStep.EthereumMainnetLogUnlock
      )
  }
}

function isNotTerminalState(s: State) {
  return !isTerminalState(s)
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
