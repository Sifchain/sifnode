import {filter, map} from 'rxjs/operators';
import * as fs from 'fs';
import * as readline from 'readline'
import {Observable, ReplaySubject} from "rxjs";
import {
    isNotNullOrUndefined,
    jsonParseSimple,
    readableStreamToObservable, tailFileAsObservable
} from "./utilities";
import {EvmStateTransition, toEvmRelayerEvent} from "./ebrelayer";
import {lastValueFrom} from "rxjs";

export function sifwatch(filename: string): Observable<EvmStateTransition> {
    // const evmRelayerLines = readableStreamToObservable(fs.createReadStream("/tmp/sifnode/evmrelayer.log"))
    const evmRelayerLines = tailFileAsObservable(filename)
    const evmRelayerEvents = evmRelayerLines.pipe(
        map(jsonParseSimple),
        map(toEvmRelayerEvent),
        filter<EvmStateTransition | undefined, EvmStateTransition>(isNotNullOrUndefined)
    )
    return evmRelayerEvents
}
