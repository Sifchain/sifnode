import { Observable } from "rxjs"
import { HardhatRuntimeEnvironment } from "hardhat/types"
import { BridgeBank } from "../../build"
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

export type EthereumMainnetEvent =
  | EthereumMainnetBlock
  | EthereumMainnetLogLock
  | EthereumMainnetLogBurn
  | EthereumMainnetLogUnlock

export function isEthereumMainnetEvent(x: object): x is EthereumMainnetEvent {
  switch ((x as EthereumMainnetEvent).kind) {
    case "EthereumMainnetBlock":
    case "EthereumMainnetLogLock":
    case "EthereumMainnetLogUnlock":
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

    let logUnlockFilter = bridgeBank.filters.LogUnlock()
    // TODO: Naming
    const logUnlockListener = (...args: any[]) => {
      const log: EthereumMainnetLogUnlock = {
        kind: "EthereumMainnetLogUnlock",
        data: {
          kind: "EthereumMainnetLogUnlock",
          to: args[0],
          token: args[1],
          value: args[2],
        },
      }
    }
    bridgeBank.on(logUnlockFilter, logUnlockListener)

    return () => {
      bridgeBank.off(lockLogFilter, logLockListener)
      bridgeBank.off(logBurnFilter, logBurnListener)
      bridgeBank.off(logUnlockFilter, logUnlockListener)
    }
  })
}
