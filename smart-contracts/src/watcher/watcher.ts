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


// const evmRelayerLines = readableStreamToObservable(fs.createReadStream("/tmp/sifnode/evmrelayer.log"))
const evmRelayerLines = tailFileAsObservable("/tmp/sifnode/evmrelayer.log")

async function main() {
    const evmRelayerEvents = evmRelayerLines.pipe(
        map(jsonParseSimple),
        map(toEvmRelayerEvent),
        filter<EvmStateTransition | undefined, EvmStateTransition>(isNotNullOrUndefined)
    )

    evmRelayerEvents.subscribe({
        next: x => {
            console.log(x)
        },
        error: e => console.log("goterror: ", e),
        complete: () => console.log("alldone")
    })

    const lv = await lastValueFrom(evmRelayerEvents)
    console.log("exitingwatcher")
}

main()
    .catch((error) => {
        console.error(error);
    })
    .finally(() => process.exit(0))
