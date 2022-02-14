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

export async function checkSifnodeIbcTransferState(
  sifnodedAdapter: SifnodedAdapter,
  srcPort: string,
  srcChannel: string,
  destination: string,
  amount: BigNumber,
  denom: string,
) {
  const evmRelayerEvents: rxjs.Observable<SifEvent> = sifwatch(
    {
      sifnoded: "/tmp/sifnode/sifnoded.log",
    },
    hardhat,
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
          case "SifnodedError":
            // if we get an actual error, that's always a failure
            return {...acc, value: {kind: "failure", value: v, message: "simple error"}}
          case "SifHeartbeat": {
            // we just store the heartbeat
            return {...acc, currentHeartbeat: v.value} as State
          }
          
          // Sifnoded side log assertions
          case "SifnodedPeggyEvent": {
            const sifnodedEvent: any = v.data
            // TODO define events and transition
            console.log("received event", sifnodedEvent.kind)
            // switch (sifnodedEvent.kind) {
              
            // }
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

  await sifnodedAdapter.executeIbcTransfer(
    srcPort,
    srcChannel,
    destination,
    amount,
    denom,
  )

  const lv = await lastValueFrom(states.pipe(takeWhile((x) => x.value.kind !== "terminate")))
  const expectedEndState: TransactionStep = TransactionStep.EthereumMainnetLogUnlock
  expect(
    lv.transactionStep,
    `did not complete, last step was ${JSON.stringify(lv, undefined, 2)}`
  ).to.eq(expectedEndState)
}
