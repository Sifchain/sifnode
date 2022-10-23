import { DevEnvContracts } from "../../src/contractSupport"
import { BridgeToken } from "../../build"
import { BigNumber, ContractTransaction } from "ethers"
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import { SifEvent, sifwatchReplayable } from "../../src/watcher/watcher"
import * as hardhat from "hardhat"
import * as ethereumAddress from "../../src/ethereumAddress"
import deepEqual = require("deep-equal")
import { expect } from "chai"

import { distinctUntilChanged, lastValueFrom, Observable, scan, takeWhile } from "rxjs"
import { filter } from "rxjs/operators"

import {
  State,
  Terminate,
  isTerminalState,
  ensureCorrectTransition,
  TransactionStep,
  buildFailure,
  attachDebugPrintfs,
  verbosityLevel,
} from "./context"
import { exec } from "child_process"

/**
 * Executes a lock on a Ethereum Transfer to send Ether or EVM native currency to Sifchain
 * @param contracts The ethers contract instances to interact with (e.g. Bridgebank)
 * @param amount The amount of ether to send to sifchain
 * @param sender Who is sending the ether to sifcahin
 * @param sifchainRecipient What sifchain address is recieving the Ether
 */
export async function executeLock(
  contracts: DevEnvContracts,
  amount: BigNumber,
  sender: SignerWithAddress,
  sifchainRecipient: string,
): Promise<ContractTransaction>;
/**
 * Executes a lock of ERC20 tokens on EVM chain to sifchain
 * @param contracts The ethers contract instancs to interact with (e.g. Bridgebank)
 * @param amount The ammount of ether to send to sifchain
 * @param sender Who is sending the ether to sifchain
 * @param sifchainRecipient What sifchain address is recieving the ERC20 tokens
 * @param tokenContract The ERC20 contract that is being bridged
 * @returns
 */
export async function executeLock(
  contracts: DevEnvContracts,
  amount: BigNumber,
  sender: SignerWithAddress,
  sifchainRecipient: string,
  tokenContract: BridgeToken,
): Promise<ContractTransaction>;
export async function executeLock(
  contracts: DevEnvContracts,
  amount: BigNumber,
  sender: SignerWithAddress,
  sifchainRecipient: string,
  tokenContract?: BridgeToken,
): Promise<ContractTransaction> {
  let tx: ContractTransaction
  if (tokenContract === undefined) {
    tx = await contracts.bridgeBank
      .connect(sender)
      .lock(sifchainRecipient, ethereumAddress.eth.address, amount, {
        value: amount,
      })
  } else {
    await tokenContract.connect(sender).approve(contracts.bridgeBank.address, amount)
    tx = await contracts.bridgeBank
      .connect(sender)
      .lock(sifchainRecipient, tokenContract.address, amount, {
        value: 0,
      })
  }
  return tx
}

export async function checkEvmLockState(
  contracts: DevEnvContracts,
  tx: ContractTransaction,
  sendAmount: BigNumber,
  denomHash: string
) {
  const [evmRelayerEvents, replayedEvents] = sifwatchReplayable(
    {
      evmrelayer: "/tmp/sifnode/relayer.log",
      sifnoded: "/tmp/sifnode/sifnoded.log",
    },
    hardhat,
    contracts.bridgeBank
  )

  const states: Observable<State> = evmRelayerEvents
    .pipe(filter((x) => x.kind !== "SifnodedInfoEvent"))
    .pipe(
      scan(
        (acc: State, v: SifEvent) => {
          if (isTerminalState(acc))
            // we've reached a decision
            return { ...acc, value: { kind: "terminate" } as Terminate }
          switch (v.kind) {
            case "EbRelayerError":
            case "SifnodedError":
              // if we get an actual error, that's always a failure
              return { ...acc, value: { kind: "failure", value: v, message: "simple error" } }
            case "SifHeartbeat":
              // we just store the heartbeat
              return { ...acc, currentHeartbeat: v.value } as State
            case "EthereumMainnetLogLock":
              // we should see exactly one lock
              let ethBlock = v.data.block as any
              if (ethBlock.transactionHash === tx.hash && v.data.value.eq(sendAmount)) {
                const newAcc: State = {
                  ...acc,
                  fromEthereumAddress: v.data.from,
                  ethereumNonce: BigNumber.from(v.data.nonce),
                }
                return ensureCorrectTransition(
                  newAcc,
                  v,
                  TransactionStep.Initial,
                  TransactionStep.SawLogLock
                )
              }
              return {
                ...acc,
                value: {
                  kind: "failure",
                  value: v,
                  message: "incorrect EthereumMainnetLogLock",
                },
              }
            case "EbRelayerEvmStateTransition":
              switch ((v.data as any).kind) {
                case "EthereumBridgeClaim":
                  const d = v.data as any
                  if (
                    d.prophecyClaim.ethereum_sender == acc.fromEthereumAddress &&
                    BigNumber.from(d.event.Nonce).eq(acc.ethereumNonce)
                  ) {
                    return ensureCorrectTransition(
                      {
                        ...acc,
                        denomHash: d.prophecyClaim.denom_hash,
                      },
                      v,
                      TransactionStep.SawLogLock,
                      TransactionStep.SawProphecyClaim
                    )
                  }
                  break
                case "EthBridgeClaimArray":
                  let claims = (v.data as any).claims as any[]
                  const matchingClaim = claims.find((claim) => claim.denom_hash === acc.denomHash)
                  if (matchingClaim)
                    return ensureCorrectTransition(
                      acc,
                      v,
                      TransactionStep.SawProphecyClaim,
                      TransactionStep.EthBridgeClaimArray
                    )
                  break
                case "BroadcastTx":
                  const messages = (v.data as any).messages as any[]
                  const matchingMessage = messages.find(
                    (msg) => msg.eth_bridge_claim.denom_hash === acc.denomHash
                  )
                  if (matchingMessage)
                    return ensureCorrectTransition(
                      acc,
                      v,
                      TransactionStep.EthBridgeClaimArray,
                      TransactionStep.BroadcastTx
                    )
              }
            case "SifnodedPeggyEvent":
              switch ((v.data as any).kind) {
                case "coinsSent":
                  const coins = ((v.data as any).coins as any)[0]
                  if (coins["denom"] === denomHash && sendAmount.eq(coins["amount"]))
                    return ensureCorrectTransition(
                      acc,
                      v,
                      TransactionStep.ProcessSuccessfulClaim,
                      TransactionStep.CoinsSent
                    )
                  else return buildFailure(acc, v, "incorrect hash or amount")
                // TODO these steps need validation to make sure they're happing in the right order with the right data
                case "CreateEthBridgeClaim":
                  let newSequenceNumber = (v.data as any).msg.Interface.eth_bridge_claim
                    .ethereum_lock_burn_sequence
                  if (acc.ethereumNonce?.eq(newSequenceNumber))
                    return ensureCorrectTransition(
                      acc,
                      v,
                      [TransactionStep.BroadcastTx, TransactionStep.AppendValidatorToProphecy],
                      TransactionStep.CreateEthBridgeClaim
                    )
                  break
                case "AppendValidatorToProphecy":
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.CreateEthBridgeClaim,
                    TransactionStep.AppendValidatorToProphecy
                  )
                case "ProcessSuccessfulClaim":
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.AppendValidatorToProphecy,
                    TransactionStep.ProcessSuccessfulClaim
                  )
                case "AddTokenMetadata":
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.ProcessSuccessfulClaim,
                    TransactionStep.AddTokenMetadata
                  )
              }
              return { ...acc, value: v, createdAt: acc.currentHeartbeat }
            default:
              // we have a new value (of any kind) and it should use the current heartbeat as its creation time
              return { ...acc, value: v, createdAt: acc.currentHeartbeat }
          }
        },
        {
          value: { kind: "initialState" },
          createdAt: 0,
          currentHeartbeat: 0,
          transactionStep: TransactionStep.Initial,
          uniqueId: "eth to ceth",
        } as State
      )
    )

  // it's useful to skip debug prints of states where only the heartbeat changed
  const withoutHeartbeat = states.pipe(
    distinctUntilChanged<State>((a, b) => {
      return deepEqual({ ...a, currentHeartbeat: 0 }, { ...b, currentHeartbeat: 0 })
    })
  )

  const verboseSubscription = attachDebugPrintfs(withoutHeartbeat, verbosityLevel())

  const lv = await lastValueFrom(states.pipe(takeWhile((x) => x.value.kind !== "terminate")))

  expect(
    lv.transactionStep,
    `did not get CoinsSent, last step was ${JSON.stringify(lv, undefined, 2)}`
  ).to.eq(TransactionStep.CoinsSent)

  verboseSubscription.unsubscribe()
  replayedEvents.unsubscribe()
}
