import { filter, map } from "rxjs/operators"
import { connectable, interval, merge, Observable, ReplaySubject, Subscription } from "rxjs"
import { isNotNullOrUndefined, jsonParseSimple, tailFileAsObservable } from "./utilities"
import { EbRelayerEvent, toEvmRelayerEvent } from "./ebrelayer"
import { SifnodedEvent, toSifnodedEvent } from "./sifnoded"
import { HardhatRuntimeEnvironment } from "hardhat/types"
import { BridgeBank, CosmosBridge } from "../../build"
import {
  EthereumMainnetEvent,
  subscribeToEthereumCosmosBridgeEvents,
  subscribeToEthereumEvents,
} from "./ethereumMainnet"
import * as path from "path"

export interface SifwatchLogs {
  evmrelayer: string
  sifnoded: string
  witness?: string
}

export interface SifHeartbeat {
  kind: "SifHeartbeat"
  value: number
}

export type SifEvent = EbRelayerEvent | SifnodedEvent | EthereumMainnetEvent | SifHeartbeat

export function sifwatch(
  logs: SifwatchLogs,
  hre: HardhatRuntimeEnvironment,
  bridgeBank: BridgeBank,
  cosmosBridge?: CosmosBridge
): Observable<SifEvent> {
  // TODO: Const?
  const observables: Observable<SifEvent>[] = new Array()

  const evmRelayerLines = tailFileAsObservable(logs.evmrelayer)
  const evmRelayerEvents: Observable<EbRelayerEvent> = evmRelayerLines.pipe(
    map(jsonParseSimple),
    map(toEvmRelayerEvent),
    filter<EbRelayerEvent | undefined, EbRelayerEvent>(isNotNullOrUndefined)
  )

  observables.push(evmRelayerEvents)

  const sifnodedLines = tailFileAsObservable(logs.sifnoded)
  const sifnodedEvents: Observable<SifnodedEvent> = sifnodedLines.pipe(
    map(jsonParseSimple),
    map(toSifnodedEvent),
    filter<SifnodedEvent | undefined, SifnodedEvent>(isNotNullOrUndefined)
  )

  observables.push(sifnodedEvents)

  const heartbeat = interval(1000).pipe(
    map((i) => {
      return { kind: "SifHeartbeat", value: i } as SifHeartbeat
    })
  )
  observables.push(heartbeat)

  const ethereumEvents = subscribeToEthereumEvents(bridgeBank)
  observables.push(ethereumEvents)

  if (logs.witness != undefined) {
    const witnessLines = tailFileAsObservable(logs.witness)
    const witnessEvents: Observable<EbRelayerEvent> = witnessLines.pipe(
      map(jsonParseSimple),
      map(toEvmRelayerEvent),
      filter<EbRelayerEvent | undefined, EbRelayerEvent>(isNotNullOrUndefined)
    )
    observables.push(witnessEvents)
  }

  if (cosmosBridge != undefined) {
    console.log("Cosmosbridge subscription")
    observables.push(subscribeToEthereumCosmosBridgeEvents(cosmosBridge))
  }

  return merge(...observables)
}

export function sifwatchReplayable(
  logs: SifwatchLogs,
  hre: HardhatRuntimeEnvironment,
  bridgeBank: BridgeBank
): [Observable<SifEvent>, Subscription] {
  const eventStream = connectable(sifwatch(logs, hre, bridgeBank), {
    connector: () => new ReplaySubject(),
    resetOnDisconnect: false,
  })
  const subscription = eventStream.connect()
  return [eventStream, subscription]
}

/**
 * Given a base directory, return a new SifwatchLogs
 * containing basedir + relayer.log for evmrelayer, etc.
 * @param basedir - the base directory containing logs
 */
export function defaultSifwatchLogs(basedir: string = "/tmp/sifnode"): SifwatchLogs {
  return new (class implements SifwatchLogs {
    evmrelayer = path.join(basedir, "relayer.log")
    sifnoded = path.join(basedir, "sifnoded.log")
    witness = path.join(basedir, "witness0.log")
  })()
}
