import * as chai from "chai"
import {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe"
import {HardhatRuntimeEnvironmentToken} from "../../src/tsyringe/injectionTokens"
import * as hardhat from "hardhat"
import {BigNumber} from "ethers"
import {ethereumResultsToSifchainAccounts, readDevEnvObj} from "../../src/tsyringe/devenvUtilities"
import {SifchainContractFactories} from "../../src/tsyringe/contracts"
import {buildDevEnvContracts, DevEnvContracts} from "../../src/contractSupport"
import web3 from "web3"
import * as ethereumAddress from "../../src/ethereumAddress"
import {SifEvent, SifHeartbeat, sifwatch, sifwatchReplayable} from "../../src/watcher/watcher"
import * as rxjs from "rxjs"
import {
  defer,
  distinctUntilChanged,
  lastValueFrom,
  Observable,
  scan,
  Subscription,
  takeWhile,
} from "rxjs"
import {EbRelayerEvmEvent} from "../../src/watcher/ebrelayer"
import {EthereumMainnetEvent} from "../../src/watcher/ethereumMainnet"
import {filter} from "rxjs/operators"
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers"
import * as ChildProcess from "child_process"
import {
  EbRelayerAccount,
  crossChainFeeBase,
  crossChainLockFee,
  crossChainBurnFee,
} from "../../src/devenv/sifnoded"
import * as dotenv from "dotenv"
import "@nomiclabs/hardhat-ethers"
import deepEqual = require("deep-equal")
import {ethers} from "hardhat"
import {SifnodedAdapter} from "./sifnodedAdapter"
import {checkSifnodeBurnState} from "./sifnode_lock_burn"

import {executeLock} from "./evm_lock_burn"
// The hash value for ethereum on mainnet
const ethDenomHash = "sif5ebfaf95495ceb5a3efbd0b0c63150676ec71e023b1043c40bcaaf91c00e15b2"

chai.use(solidity)

describe("lock and burn tests", () => {
  dotenv.config()
  // This test only works when devenv is running, and that requires a connection to localhost
  expect(hardhat.network.name, "please use devenv").to.eq("localhost")

  const devEnvObject = readDevEnvObj("environment.json")
  // a generic sif address, nothing special about it
  const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")
  const networkDescriptor = devEnvObject?.ethResults?.chainId ?? 31337

  const sifnodedAdapter: SifnodedAdapter = new SifnodedAdapter(
    devEnvObject!.sifResults!.adminAddress!.homeDir,
    devEnvObject!.sifResults!.adminAddress!.account,
    process.env["GOBIN"]
  )

  before("register HardhatRuntimeEnvironmentToken", async () => {
    container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
  })

  it.only("should allow ceth to eth tx", async () => {
    // TODO: Could these be moved out of the test fx? and instantiated via beforeEach?
    const factories = container.resolve(SifchainContractFactories)
    const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)

    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    )
    const destinationEthereumAddress = ethereumAccounts.availableAccounts[0]

    // These two can happen together
    const initialBalance = (
      await ethers.provider.getBalance(destinationEthereumAddress.address)
    ).toString()

    const contractInitialBalance = (
      await ethers.provider.getBalance(contracts.bridgeBank.address)
    ).toString()

    // const sendAmount = BigNumber.from(5 * ETH) // 3500 gwei
    const sendAmount = BigNumber.from("5000000000000000000") // 3500 gwei

    let testSifAccount: EbRelayerAccount = sifnodedAdapter.createTestSifAccount()
    let originalVerboseLevel: string | undefined = process.env["VERBOSE"]
    process.env["VERBOSE"] = "summary"
    // Need to have a burn of eth happen at least once or there's no data about eth in the token metadata
    await executeLock(
      contracts,
      undefined,
      sendAmount,
      ethereumAccounts.availableAccounts[1],
      web3.utils.utf8ToHex(testSifAccount.account)
    )

    const intermediateBalance = (
      await ethers.provider.getBalance(destinationEthereumAddress.address)
    ).toString()
    let contractIntermediateBalance = (
      await ethers.provider.getBalance(contracts.bridgeBank.address)
    ).toString()

    // These are temporarily added to make the logging lvl lower
    process.env["VERBOSE"] = originalVerboseLevel

    console.log("Lock complete")

    let crossChainCethFee = crossChainFeeBase * crossChainBurnFee

    let newSendAmount = BigNumber.from("2300000000000000000") // 2300 gwei
    await checkSifnodeBurnState(
      sifnodedAdapter,
      contracts,
      testSifAccount,
      destinationEthereumAddress,
      newSendAmount.sub(crossChainCethFee),
      ethDenomHash,
      String(crossChainCethFee),
      networkDescriptor
    )

    // Here we verify the user balance is correct
    const finalBalance = (
      await ethers.provider.getBalance(destinationEthereumAddress.address)
    ).toString()
    let contractFinalBalance = (
      await ethers.provider.getBalance(contracts.bridgeBank.address)
    ).toString()

    console.log("Initial Balance     ", initialBalance)
    console.log("intermediate Balance", intermediateBalance)
    console.log("final Balance       ", finalBalance)

    console.log("Contract Initial Balance     ", contractInitialBalance)
    console.log("Contract intermediate Balance", contractIntermediateBalance)
    console.log("Contract Final Balance       ", contractFinalBalance)

    // verboseSubscription.unsubscribe()
  })

  // it("should send two locks of ethereum", async () => {
  //   const ethereumAccounts = await ethereumResultsToSifchainAccounts(
  //     devEnvObject.ethResults!,
  //     hardhat.ethers.provider
  //   )
  //   const factories = container.resolve(SifchainContractFactories)
  //   const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)
  //   const sender1 = ethereumAccounts.availableAccounts[0]
  //   const smallAmount = BigNumber.from(1017)

  //   // Do two locks of ethereum
  //   await executeLock(contracts, smallAmount, sender1, recipient, true,)
  //   await executeLock(contracts, smallAmount, sender1, recipient, true, "second lock of eth")
  // })
})
