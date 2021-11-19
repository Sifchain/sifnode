import {filter, map} from 'rxjs/operators';
import {connectable, interval, merge, Observable, ReplaySubject, Subscription} from "rxjs";
import {isNotNullOrUndefined, jsonParseSimple, tailFileAsObservable} from "./utilities";
import {EbRelayerEvent, toEvmRelayerEvent} from "./ebrelayer";
import {SifnodedEvent, toSifnodedEvent} from "./sifnoded";
import {HardhatRuntimeEnvironment} from "hardhat/types";
import {BridgeBank} from "../../build";
import {EthereumMainnetEvent, subscribeToEthereumEvents} from "./ethereumMainnet";

export interface SifwatchLogs {
    evmrelayer: string
    sifnoded: string
}

export interface SifHeartbeat {
    kind: "SifHeartbeat"
    value: number
}

export type SifEvent = EbRelayerEvent | SifnodedEvent | EthereumMainnetEvent | SifHeartbeat

export function sifwatch(logs: SifwatchLogs, hre: HardhatRuntimeEnvironment, bridgeBank: BridgeBank): Observable<SifEvent> {
    // const evmRelayerLines = readableStreamToObservable(fs.createReadStream("/tmp/sifnode/evmrelayer.log"))
    const evmRelayerLines = tailFileAsObservable(logs.evmrelayer)
    const evmRelayerEvents: Observable<EbRelayerEvent> = evmRelayerLines.pipe(
        map(jsonParseSimple),
        map(toEvmRelayerEvent),
        filter<EbRelayerEvent | undefined, EbRelayerEvent>(isNotNullOrUndefined)
    )
    const sifnodedLines = tailFileAsObservable(logs.sifnoded)
    const sifnodedEvents: Observable<SifnodedEvent> = sifnodedLines.pipe(
        map(jsonParseSimple),
        map(toSifnodedEvent),
        filter<SifnodedEvent | undefined, SifnodedEvent>(isNotNullOrUndefined)
    )
    const heartbeat = interval(1000).pipe(map(i => {
        return {kind: "SifHeartbeat", value: i} as SifHeartbeat
    }))
    const ethereumEvents = subscribeToEthereumEvents(hre, bridgeBank)
    return merge(evmRelayerEvents, sifnodedEvents, ethereumEvents, heartbeat)

    // TODO: Add cosmosEvent
}

export function sifwatchReplayable(logs: SifwatchLogs, hre: HardhatRuntimeEnvironment, bridgeBank: BridgeBank): [Observable<SifEvent>, Subscription] {
    const eventStream = connectable(
        sifwatch(logs, hre, bridgeBank),
        { connector: () => new ReplaySubject(), resetOnDisconnect: false }
    )
    const subscription = eventStream.connect()
    return [eventStream, subscription]
}
