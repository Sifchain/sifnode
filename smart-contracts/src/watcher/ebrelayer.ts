import {filter, map} from 'rxjs/operators';
import * as fs from 'fs';
import * as readline from 'readline'
import {Observable, ReplaySubject} from "rxjs";
import {jsonParseSimple, readableStreamToObservable} from "./utilities";

interface EbRelayerEvmEvent {
    kind: "EbRelayerEvmEvent"
    data: object
}

interface EbRelayerEthBridgeClaimArray {
    kind: "EbRelayerEthBridgeClaimArray"
    data: {
        claims: {
            network_descriptor: number
            bridge_contract_address: string
            nonce: number
            symbol: string
            // "token_contract_address": "0x0000000000000000000000000000000000000000",
            // "ethereum_sender": "0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65",
            // "cosmos_receiver": "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace",
            // "validator_address": "sifvaloper1n77760y8rvs4f5y77ssk3ak0c7p6efcyva2f48",
            amount: string
            claim_type: number
            token_name: string
            decimals: number
            denom_hash: string
        }[]
    }
}

interface EbRelayerEvmStateTransition {
    kind: "EbRelayerEvmStateTransition"
    data: object
}

interface EbRelayerError {
    kind: "EbRelayerError"
    data: object
}

export type EbRelayerEvent =
    EbRelayerEvmEvent
    | EbRelayerEvmStateTransition
    | EbRelayerError
    | EbRelayerEthBridgeClaimArray

export function toEvmRelayerEvent(x: any): EbRelayerEvent | undefined {
    if (x["M"] === "peggytest") {
        switch (x["kind"]) {
            case "EthereumEvent":
                return {kind: "EbRelayerEvmEvent", data: x}
                break
            case "EthBridgeClaimArray":
                return {
                    kind: "EbRelayerEthBridgeClaimArray",
                    data: x
                } as EbRelayerEthBridgeClaimArray
            default:
                return {kind: "EbRelayerEvmStateTransition", data: x}
                break
        }
    } else if (x["L"] === "ERROR") {
        return {kind: "EbRelayerError", data: x}
    } else {
        return undefined
    }
}
