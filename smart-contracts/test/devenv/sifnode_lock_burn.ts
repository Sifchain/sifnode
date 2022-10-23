import {DevEnvContracts} from "../../src/contractSupport"
import {BigNumber} from "ethers"
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers"
import {SifEvent, sifwatch} from "../../src/watcher/watcher"
import * as hardhat from "hardhat"
import deepEqual = require("deep-equal")
import {expect} from "chai"
import * as rxjs from "rxjs"
import {EbRelayerAccount} from "../../src/devenv/sifnoded"
import {distinctUntilChanged, lastValueFrom, Observable, scan, takeWhile} from "rxjs"
import {filter} from "rxjs/operators"
import {
  State,
  Terminate,
  isTerminalState,
  ensureCorrectTransition,
  TransactionStep,
} from "./context"
import {SifnodedAdapter} from "./sifnodedAdapter"

export async function checkSifnodeBurnState(
  sifnodedAdapter: SifnodedAdapter,
  contracts: DevEnvContracts,
  sender: EbRelayerAccount,
  destination: SignerWithAddress,
  amount: BigNumber,
  symbol: string,
  // TODO: What is correct value for corsschainfee?
  crossChainFee: string,
  networkDescriptor: number
) {
  const evmRelayerEvents: rxjs.Observable<SifEvent> = sifwatch(
    {
      evmrelayer: "/tmp/sifnode/relayer.log",
      sifnoded: "/tmp/sifnode/sifnoded.log",
      witness: "/tmp/sifnode/witness.log",
    },
    hardhat,
    contracts.bridgeBank,
    contracts.cosmosBridge
  ).pipe(filter((x) => x.kind !== "SifnodedInfoEvent"))

  let receivedCosmosBurnmsg: boolean = false
  let witnessSignedProphecy: boolean = false

  let hasSeenEthereumLogUnlcok: boolean = false
  let hasSeenProphecyClaimSubmitted: boolean = false

  const states: Observable<State> = evmRelayerEvents.pipe(
    scan(
      (acc: State, v: SifEvent) => {
        console.log("Event: ", v)
        // if (v.kind == "")
        if (isTerminalState(acc) || (hasSeenEthereumLogUnlcok && hasSeenProphecyClaimSubmitted)) {
          // we've reached a decision
          console.log("Reached terminate state", acc)
          return {...acc, value: {kind: "terminate"} as Terminate}
        }
        switch (v.kind) {
          case "EbRelayerError":
          case "SifnodedError":
            // if we get an actual error, that's always a failure
            return {...acc, value: {kind: "failure", value: v, message: "simple error"}}
          case "SifHeartbeat": {
            // we just store the heartbeat
            return {...acc, currentHeartbeat: v.value} as State
          }
          case "EthereumMainnetLogUnlock": {
            hasSeenEthereumLogUnlcok = true
            return ensureCorrectTransition(
              acc,
              v,
              TransactionStep.ProphecyStatus,
              TransactionStep.EthereumMainnetLogUnlock
            )
          }
          // Ebrelayer side log assertions
          case "EbRelayerEvmStateTransition": {
            let ebrelayerEvent: any = v.data
            switch (ebrelayerEvent.kind) {
              case "ReceiveCosmosBurnMessage": {
                // console.log("Seeing ReceiveCosmosBurnMessage")
                if (!receivedCosmosBurnmsg) {
                  // console.log("Receiving ReceiveCosmosBurnMessage for the first time")
                  // Ignore subsequence occurrences, witness will reprocess until keeper updates nonce
                  receivedCosmosBurnmsg = true
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.PublishCosmosBurnMessage,
                    TransactionStep.ReceiveCosmosBurnMessage
                  )
                } else {
                  return {...acc, value: v, createdAt: acc.currentHeartbeat}
                }
              }
              case "WitnessSignProphecy": {
                // console.log("Seeing WitnessSignProphecy")
                if (!witnessSignedProphecy) {
                  // console.log("Receiving WitnessSignProphecy for the first time")
                  witnessSignedProphecy = true
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.ReceiveCosmosBurnMessage,
                    TransactionStep.WitnessSignProphecy
                  )
                } else {
                  return {...acc, value: v, createdAt: acc.currentHeartbeat}
                }
              }

              case "ProphecyClaimSubmitted": {
                hasSeenProphecyClaimSubmitted = true
                // return ensureCorrectTransition(
                //   acc,
                //   v,
                //   TransactionStep.EthereumMainnetLogUnlock,
                //   TransactionStep.ProphecyClaimSubmitted
                // )
              }
            }
          }
          // Sifnoded side log assertions
          case "SifnodedPeggyEvent": {
            const sifnodedEvent: any = v.data
            switch (sifnodedEvent.kind) {
              case "Burn": {
                return ensureCorrectTransition(
                  acc,
                  v,
                  TransactionStep.Initial,
                  TransactionStep.Burn
                )
              }

              case "GetCrossChainFeeConfig": {
                return ensureCorrectTransition(
                  acc,
                  v,
                  TransactionStep.Burn,
                  TransactionStep.GetCrossChainFeeConfig
                )
              }

              case "SendCoinsFromAccountToModule": {
                return ensureCorrectTransition(
                  acc,
                  v,
                  TransactionStep.GetCrossChainFeeConfig,
                  TransactionStep.SendCoinsFromAccountToModule
                )
              }

              case "BurnCoins": {
                // TODO: Add assertion on expected amount, and expected denom
                return ensureCorrectTransition(
                  acc,
                  v,
                  TransactionStep.SendCoinsFromAccountToModule,
                  TransactionStep.BurnCoins
                )
              }

              /**
               * We comment this out because SetProphecy is the crUd operation, gets invoked multiple times throughout
               * the call,
               * But we still want to assert it has created a prophecy between BurnCoin and PublishCosmosBurnMessage
               * TODO: Option 1. Refine the instrumentation statement in SetProphecy
               *       Option 2. ???
               */
              // case "SetProphecy":
              //   return ensureCorrectTransition(
              //     acc,
              //     v,
              //     TransactionStep.BurnCoins,
              //     TransactionStep.SetProphecy
              //   )

              case "PublishCosmosBurnMessage": {
                // console.log("Received PublishCosmosBurnMessage")
                return ensureCorrectTransition(
                  acc,
                  v,
                  TransactionStep.BurnCoins,
                  TransactionStep.PublishCosmosBurnMessage
                )
              }

              case "SetWitnessLockBurnNonce": {
                // console.log("Receiving SetWitnessLockBurnNonce. Acc,", acc)
                return ensureCorrectTransition(
                  acc,
                  v,
                  TransactionStep.WitnessSignProphecy,
                  TransactionStep.SetWitnessLockBurnNonce
                )
              }

              case "ProphecyStatus": {
                return ensureCorrectTransition(
                  acc,
                  v,
                  TransactionStep.SetWitnessLockBurnNonce,
                  TransactionStep.ProphecyStatus
                )
              }
            }
          }

          default: {
            // we have a new value (of any kind) and it should use the current heartbeat as its creation time
            return {...acc, value: v, createdAt: acc.currentHeartbeat}
          }
        }
      },
      {
        value: {kind: "initialState"},
        createdAt: 0,
        currentHeartbeat: 0,
        transactionStep: TransactionStep.Initial,
      } as State
    )
  )

  // it's useful to skip debug prints of states where only the heartbeat changed
  const withoutHeartbeat = states.pipe(
    distinctUntilChanged<State>((a, b) => {
      return deepEqual({...a, currentHeartbeat: 0}, {...b, currentHeartbeat: 0})
    })
  )

  await sifnodedAdapter.executeSifBurn(
    sender,
    destination,
    amount,
    symbol,
    crossChainFee,
    networkDescriptor
  )

  const lv = await lastValueFrom(states.pipe(takeWhile((x) => x.value.kind !== "terminate")))
  const expectedEndState: TransactionStep = TransactionStep.EthereumMainnetLogUnlock
  expect(
    lv.transactionStep,
    `did not complete, last step was ${JSON.stringify(lv, undefined, 2)}`
  ).to.eq(expectedEndState)
}
