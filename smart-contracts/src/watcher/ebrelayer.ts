import {filter, map} from 'rxjs/operators';
import * as fs from 'fs';
import * as readline from 'readline'
import {Observable, ReplaySubject} from "rxjs";
import {jsonParseSimple, readableStreamToObservable} from "./utilities";

export class EvmStateTransition {
    constructor(readonly data: object) {
    }
}

export class EvmEvent extends EvmStateTransition {
}

export class EvmError extends EvmStateTransition {
}

export function toEvmRelayerEvent(x: any): EvmStateTransition | undefined {
    if (x["M"] === "devenv") {
        switch (x["kind"]) {
            case "EthereumEvent":
                return new EvmEvent(x)
                break
            default:
                return new EvmStateTransition(x)
                break
        }
    } else if (x["L"] === "ERROR") {
        return new EvmError(x)
    } else {
        return undefined
    }
}
