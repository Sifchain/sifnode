import * as chai from "chai"
import { expect } from "chai"
import { solidity } from "ethereum-waffle"
import { container } from "tsyringe"
import { HardhatRuntimeEnvironmentToken } from "../../src/tsyringe/injectionTokens"
import * as hardhat from "hardhat"
import { BigNumber } from "ethers"
import {
  ethereumResultsToSifchainAccounts,
  readDevEnvObj,
} from "../../src/tsyringe/devenvUtilities"
import { SifchainContractFactories } from "../../src/tsyringe/contracts"
import { buildDevEnvContracts } from "../../src/contractSupport"
import web3 from "web3"
import { EbRelayerAccount, crossChainFeeBase, crossChainBurnFee } from "../../src/devenv/sifnoded"
import * as dotenv from "dotenv"
import "@nomiclabs/hardhat-ethers"
import { ethers } from "hardhat"
import { SifnodedAdapter } from "./sifnodedAdapter"
import { checkSifnodeBurnState } from "./sifnode_lock_burn"
import { ETH, ethDenomHash, isTerminalState, State, Terminate, TransactionStep } from "./context"
import { StateMachineVerifier, StateMachineVerifierBuilder } from "./stateMachineVerifier"

import { executeLock, checkEvmLockState } from "./evm_lock_burn"
import { filter, Observable, scan } from "rxjs"
import { SifEvent, sifwatch } from "../../src/watcher/watcher"
import * as rxjs from "rxjs"
import { lastValueFrom } from "rxjs"
import { stat } from "fs"
import { assert } from "console"

chai.use(solidity)

describe("lock and burn tests", () => {
  dotenv.config()
  // This test only works when devenv is running, and that requires a connection to localhost
  expect(hardhat.network.name, "please use devenv").to.eq("localhost")

  const devEnvObject = readDevEnvObj("environment.json")
  // a generic sif address, nothing special about it
  const networkDescriptor = devEnvObject?.ethResults?.chainId ?? 31337

  const sifnodedAdapter: SifnodedAdapter = new SifnodedAdapter(
    devEnvObject!.sifResults!.adminAddress!.homeDir,
    devEnvObject!.sifResults!.adminAddress!.account,
    process.env["GOBIN"]
  )

  before("register HardhatRuntimeEnvironmentToken", async () => {
    container.register(HardhatRuntimeEnvironmentToken, { useValue: hardhat })
  })

  it("should allow ceth to eth tx", async () => {
    // TODO: Could these be moved out of the test fx? and instantiated via beforeEach?
    const factories = container.resolve(SifchainContractFactories)
    const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)

    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    )
    const destinationEthereumAddress = ethereumAccounts.availableAccounts[0]
    const sendAmount = BigNumber.from(5 * ETH) // 3500 gwei
    // const sendAmount = BigNumber.from("5000000000000000000") // 3500 gwei

    let testSifAccount: EbRelayerAccount = sifnodedAdapter.createTestSifAccount()
    process.env["VERBOSE"] = "summary"
    // Need to have a burn of eth happen at least once or there's no data about eth in the token metadata
    let tx = await executeLock(
      contracts,
      sendAmount,
      ethereumAccounts.availableAccounts[1],
      web3.utils.utf8ToHex(testSifAccount.account)
    )

    await checkEvmLockState(contracts, tx, sendAmount, ethDenomHash)

    // record the init balance before burn
    const initialEthSenderBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      ethDenomHash
    )
    const initialEthReceiverBalance = await ethers.provider.getBalance(
      destinationEthereumAddress.address
    )

    let crossChainCethFee = crossChainFeeBase * crossChainBurnFee

    let burnAmount = BigNumber.from("2300000") // 2300 gwei
    await checkSifnodeBurnState(
      sifnodedAdapter,
      contracts,
      testSifAccount,
      destinationEthereumAddress,
      burnAmount.sub(crossChainCethFee),
      ethDenomHash,
      String(crossChainCethFee),
      networkDescriptor
    )

    // Here we verify the user balance is correct
    const finalEthSenderBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      ethDenomHash
    )
    const finalEthReceiverBalance = await ethers.provider.getBalance(
      destinationEthereumAddress.address
    )

    console.log("Before burn the sender's balance is ", initialEthSenderBalance)
    console.log("Before burn the receiver's balance is ", initialEthReceiverBalance)

    console.log("After burn the sender's balance is ", finalEthSenderBalance)
    console.log("After burn the receiver's balance is ", finalEthReceiverBalance)

    expect(initialEthSenderBalance.sub(burnAmount), "should be equal ").eq(finalEthSenderBalance)
    expect(initialEthReceiverBalance.add(burnAmount.sub(crossChainCethFee)), "should be equal ").eq(
      finalEthReceiverBalance
    )
  })

  it("should send two locks of ethereum", async () => {
    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    )
    const factories = container.resolve(SifchainContractFactories)
    const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)
    const sender1 = ethereumAccounts.availableAccounts[0]
    const smallAmount = BigNumber.from(1017)
    const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")

    // Do two locks of ethereum
    let tx = await executeLock(contracts, smallAmount, sender1, recipient)
    await checkEvmLockState(contracts, tx, smallAmount, ethDenomHash)

    tx = await executeLock(contracts, smallAmount, sender1, recipient)
    await checkEvmLockState(contracts, tx, smallAmount, ethDenomHash)
  })

  it.only("Should send ceth back to eth using builder", async () => {
    const factories = container.resolve(SifchainContractFactories)
    const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)

    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    )
    const destinationEthereumAddress = ethereumAccounts.availableAccounts[0]

    // const sendAmount = BigNumber.from(5 * ETH) // 3500 gwei
    const sendAmount = BigNumber.from("500000000000") // 3500 gwei

    let testSifAccount: EbRelayerAccount = sifnodedAdapter.createTestSifAccount()
    let originalVerboseLevel: string | undefined = process.env["VERBOSE"]
    process.env["VERBOSE"] = "summary"
    // Need to have a burn of eth happen at least once or there's no data about eth in the token metadata
    await executeLock(
      contracts,
      sendAmount,
      ethereumAccounts.availableAccounts[1],
      web3.utils.utf8ToHex(testSifAccount.account)
    )

    const evmRelayerEvents: rxjs.Observable<SifEvent> = sifwatch(
      {
        evmrelayer: "/tmp/sifnode/evmrelayer.log",
        sifnoded: "/tmp/sifnode/sifnoded.log",
        witness: "/tmp/sifnode/witness.log",
      },
      hardhat,
      contracts.bridgeBank,
      contracts.cosmosBridge
    ) //.pipe(filter((x) => x.kind !== "SifnodedInfoEvent"))

    const smVerifierBuilder: StateMachineVerifierBuilder = new StateMachineVerifierBuilder()
    const smVerifier: StateMachineVerifier = smVerifierBuilder
      // .initial(TransactionStep.Initial)
      .then(TransactionStep.SawLogLock)
      .then(TransactionStep.SawProphecyClaim)
      .build()

    const states: Observable<State> = evmRelayerEvents
      .pipe(filter((x) => x.kind !== "SifnodedInfoEvent"))
      .pipe(
        scan(
          (acc: State, v: SifEvent) => {
            if (isTerminalState(acc)) {
              return {
                ...acc,
                value: { kind: "terminate" } as Terminate,
              }
            }
            switch (v.kind) {
              case "EbRelayerError":
              case "SifnodedError":
                return { ...acc, value: { kind: "failure", value: v, message: "simple error" } }
              case "SifHeartbeat":
                return { ...acc, currentHeartbeat: v.value } as State
              case "EthereumMainnetLogLock":
              case "EbRelayerEvmStateTransition":
              case "SifnodedPeggyEvent": {
                console.log("Verifying: ", v)

                let newAcc: State = smVerifier.verify(v)
                console.log("Acc: ", acc)
                console.log("NewAcc: ", newAcc)
                return newAcc
                // return {
                //   ...newAcc,
                //   // value: v,
                //   // createdAt: newAcc.currentHeartbeat,
                //   // transactionStep: newAcc.transactionStep,
                // }
              }
              default:
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

    console.log(states)
    const lastValue = await lastValueFrom(
      states.pipe(rxjs.takeWhile((x) => x.value.kind != "terminate"))
    )

    // expect(lastValue.transactionStep).to.eq(TransactionStep.SawLogLock)
    expect(lastValue.value.kind).to.eq("success")
    console.log(lastValue)
  })
})
