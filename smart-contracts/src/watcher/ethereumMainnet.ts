import { Observable } from "rxjs"
import { HardhatRuntimeEnvironment } from "hardhat/types"
import { BridgeBank, CosmosBridge } from "../../build"
import { BigNumber } from "ethers"

export interface EthereumMainnetBlock {
  kind: "EthereumMainnetBlock"
  blockNumber: number
}

export interface EthereumMainnetLogLock {
  kind: "EthereumMainnetLogLock"
  data: {
    kind: "EthereumMainnetLogLock"
    from: string
    to: string
    token: string
    value: BigNumber
    nonce: BigNumber
    decimals: number
    symbol: string
    name: string
    networkDescriptor: number
    block: object
  }
}

export interface EthereumMainnetLogBurn {
  kind: "EthereumMainnetLogBurn"
  data: {
    kind: "EthereumMainnetLogBurn"
    from: string
    to: string
    token: string
    value: BigNumber
    nonce: BigNumber
    decimals: number
    symbol: string
    name: string
    networkDescriptor: number
    block: object
  }
}

export interface EthereumMainnetLogUnlock {
  kind: "EthereumMainnetLogUnlock"
  data: {
    kind: "EthereumMainnetLogUnlock"
    to: string
    token: string
    value: string
  }
}

export interface EthereumMainnetLogBridgeTokenMint {
  kind: "EthereumMainnetLogBridgeTokenMint"
  data: {
    kind: "EthereumMainnetLogBridgeTokenMint"
    token: string
    symbol: string
    cosmosDenom: string
  }
}

export interface EthereumMainnetNewProphecyClaim {
  kind: "EthereumMainnetNewProphecyClaim"
  data: {
    kind: "EthereumMainnetNewProphecyClaim"
    // Commented out because they are indexed, these live inside topic
    // prophecyId: string
    // ethereumReceiver: string
    // amount: BigNumber
  }
}

export interface EthereumMainnetLogProphecyCompleted {
  kind: "EthereumMainnetLogProphecyCompleted"
  data: {
    kind: "EthereumMainnetLogProphecyCompleted"
    // prophecyId: string
  }
}

export type EthereumMainnetEvent =
  | EthereumMainnetBlock
  | EthereumMainnetLogLock
  | EthereumMainnetLogBurn
  | EthereumMainnetLogUnlock
  | EthereumMainnetLogBridgeTokenMint
  | EthereumMainnetNewProphecyClaim
  | EthereumMainnetLogProphecyCompleted

export function isEthereumMainnetEvent(x: object): x is EthereumMainnetEvent {
  switch ((x as EthereumMainnetEvent).kind) {
    case "EthereumMainnetBlock":
    case "EthereumMainnetLogLock":
    case "EthereumMainnetLogBurn":
    case "EthereumMainnetLogUnlock":
    case "EthereumMainnetLogBridgeTokenMint":
    case "EthereumMainnetNewProphecyClaim":
    case "EthereumMainnetLogProphecyCompleted":
      return true
    default:
      return false
  }
}

export function isNotEthereumMainnetEvent(x: object): x is EthereumMainnetEvent {
  return !isEthereumMainnetEvent(x)
}

export function subscribeToEthereumEvents(
  bridgeBank: BridgeBank
): Observable<EthereumMainnetEvent> {
  return new Observable<EthereumMainnetEvent>((subscriber) => {
    const logLockListener = (...args: any[]) => {
      const newVar: EthereumMainnetLogLock = {
        kind: "EthereumMainnetLogLock",
        data: {
          kind: "EthereumMainnetLogLock",
          from: args[0],
          to: args[1],
          token: args[2],
          value: BigNumber.from(args[3]),
          nonce: BigNumber.from(args[4]),
          decimals: parseInt(args[5]),
          symbol: args[6],
          name: args[7],
          networkDescriptor: parseInt(args[8]),
          block: args[9],
        },
      }
      subscriber.next(newVar)
    }
    let lockLogFilter = bridgeBank.filters.LogLock()
    bridgeBank.on(lockLogFilter, logLockListener)

    const logBurnListener = (...args: any[]) => {
      let newVar: EthereumMainnetLogBurn = {
        kind: "EthereumMainnetLogBurn",
        data: {
          kind: "EthereumMainnetLogBurn",
          from: args[0],
          to: args[1],
          token: args[2],
          value: BigNumber.from(args[3]),
          nonce: BigNumber.from(args[4]),
          decimals: parseInt(args[5]),
          symbol: args[6],
          name: args[7],
          networkDescriptor: parseInt(args[8]),
          block: args[9],
        },
      }
      subscriber.next(newVar)
    }
    let logBurnFilter = bridgeBank.filters.LogBurn()
    bridgeBank.on(logBurnFilter, logBurnListener)

    const logUnlockListener = (...args: any[]) => {
      const event: EthereumMainnetLogUnlock = {
        kind: "EthereumMainnetLogUnlock",
        data: {
          kind: "EthereumMainnetLogUnlock",
          to: args[0],
          token: args[1],
          value: args[2],
        },
      }
      subscriber.next(event)
    }
    let logUnlockFilter = bridgeBank.filters.LogUnlock()
    bridgeBank.on(logUnlockFilter, logUnlockListener)

    const logBridgeTokenMintListener = (...args: any[]) => {
      console.log("Received token mint")
      const log: EthereumMainnetLogBridgeTokenMint = {
        kind: "EthereumMainnetLogBridgeTokenMint",
        data: {
          kind: "EthereumMainnetLogBridgeTokenMint",
          token: args[0],
          symbol: args[1],
          cosmosDenom: args[2],
        },
      }
      subscriber.next(log)
    }
    let logBridgeTokenMintFilter = bridgeBank.filters.LogBridgeTokenMint()
    bridgeBank.on(logBridgeTokenMintFilter, logBridgeTokenMintListener)

    return () => {
      bridgeBank.off(lockLogFilter, logLockListener)
      bridgeBank.off(logBurnFilter, logBurnListener)
      bridgeBank.off(logUnlockFilter, logUnlockListener)
      bridgeBank.off(logBridgeTokenMintFilter, logBridgeTokenMintListener)
    }
  })
}

// TODO: Consider using function overloading to make this user-friendly
export function subscribeToEthereumCosmosBridgeEvents(
  cosmosBridge: CosmosBridge
): Observable<EthereumMainnetEvent> {
  return new Observable<EthereumMainnetEvent>((subscriber) => {
    const logNewProphecyClaimListener = (...args: any[]) => {
      console.log("Received new prophecy claim event")
      const log: EthereumMainnetNewProphecyClaim = {
        kind: "EthereumMainnetNewProphecyClaim",
        data: {
          kind: "EthereumMainnetNewProphecyClaim",
        },
      }
      subscriber.next(log)
    }
    let logNewProphecyClaimFilter = cosmosBridge.filters.LogNewProphecyClaim()
    cosmosBridge.on(logNewProphecyClaimFilter, logNewProphecyClaimListener)

    const LogProphecyCompleted = (...args: any[]) => {
      console.log("Receive prophecy completed")
      const log: EthereumMainnetLogProphecyCompleted = {
        kind: "EthereumMainnetLogProphecyCompleted",
        data: {
          kind: "EthereumMainnetLogProphecyCompleted",
        },
      }
      subscriber.next(log)
    }
    let logProphecyCompletedFilter = cosmosBridge.filters.LogProphecyCompleted()
    cosmosBridge.on(logProphecyCompletedFilter, LogProphecyCompleted)

    return () => {
      cosmosBridge.off(logNewProphecyClaimFilter, logNewProphecyClaimListener)
      cosmosBridge.off(logProphecyCompletedFilter, LogProphecyCompleted)
    }
  })
}
