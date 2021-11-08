import {filter, map} from 'rxjs/operators';
import * as fs from 'fs';
import * as readline from 'readline'
import {from, Observable, ReplaySubject} from "rxjs";
import {jsonParseSimple, readableStreamToObservable} from "./utilities";
import {HardhatRuntimeEnvironment} from "hardhat/types";
import {BridgeBank} from "../../build";
import {BigNumber} from "ethers";

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
    }
}

export type EthereumMainnetEvent = EthereumMainnetBlock | EthereumMainnetLogLock

export function isEthereumMainnetEvent(x: object): x is EthereumMainnetEvent {
    switch ((x as EthereumMainnetEvent).kind) {
        case "EthereumMainnetBlock":
        case "EthereumMainnetLogLock":
            return true
        default:
            return false
    }
}

export function isNotEthereumMainnetEvent(x: object): x is EthereumMainnetEvent {
    return !isEthereumMainnetEvent(x)
}

export function subscribeToEthereumEvents(hre: HardhatRuntimeEnvironment, bridgeBank: BridgeBank): Observable<EthereumMainnetEvent> {
    const events = new ReplaySubject<EthereumMainnetEvent>()

    // subscribe to new blocks, just as the evm liveness indicator
    bridgeBank.on(bridgeBank.filters.LogLock(), (...args) => {
        events.next({
            kind: "EthereumMainnetLogLock",
            data: {
                kind: "EthereumMainnetLogLock",
                from: args[0],
                to: args[1],
                token: args[2],
                value: args[3],
                nonce: args[4],
                decimals: args[5],
                symbol: args[6],
                name: args[7],
                networkDescriptor: args[8],
                block: args[9]
            }
        } as EthereumMainnetLogLock)
    })
    return events
}
