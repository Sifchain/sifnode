import {filter, map} from 'rxjs/operators';
import {interval, merge, Observable} from "rxjs";
import {isNotNullOrUndefined, jsonParseSimple, tailFileAsObservable} from "./utilities";
import {EbRelayerEvent, toEvmRelayerEvent} from "./ebrelayer";
import {SifnodedEvent, SifnodedInfoEvent, toSifnodedEvent} from "./sifnoded";

export interface SifwatchLogs {
    evmrelayer: string
    sifnoded: string
}

export interface SifHeartbeat {
    kind: "SifHeartbeat"
    value: number
}

export type SifEvent = EbRelayerEvent | SifnodedEvent | SifHeartbeat

export function sifwatch(logs: SifwatchLogs): Observable<SifEvent> {
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
    return merge(evmRelayerEvents, sifnodedEvents, heartbeat)
}
