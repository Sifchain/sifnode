import { filter, map } from "rxjs/operators"
import * as fs from "fs"
import * as readline from "readline"
import { Observable, ReplaySubject } from "rxjs"
import { jsonParseSimple, readableStreamToObservable } from "./utilities"

export interface EbRelayerEvmEvent {
  kind: "EbRelayerEvmEvent"
  data: {
    event: {
      To: string
      Symbol: string
      Name: string
      Decimals: number
      NetworkDescriptor: number
      Value: number
      Nonce: number
      ClaimType: number
      BridgeContractAddress: string
      From: string
      Token: string
    }
  }
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

export interface EbRelayerError {
  kind: "EbRelayerError"
  data: object
}

export type EbRelayerEvent =
  | EbRelayerEvmEvent
  | EbRelayerEvmStateTransition
  | EbRelayerError
  | EbRelayerEthBridgeClaimArray

export function toEvmRelayerEvent(x: any): EbRelayerEvent | undefined {
  if (x["M"] === "peggytest") {
    switch (x["kind"]) {
      default:
        return { kind: "EbRelayerEvmStateTransition", data: x }
        break
    }
  } else if (x["L"] === "ERROR") {
    return { kind: "EbRelayerError", data: x }
  } else {
    return undefined
  }
}
