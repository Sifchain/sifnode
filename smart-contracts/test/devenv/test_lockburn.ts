import * as chai from "chai"
import {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe"
import {HardhatRuntimeEnvironmentToken} from "../../src/tsyringe/injectionTokens"
import * as hardhat from "hardhat"
import {BigNumber} from "ethers"
import {ethereumResultsToSifchainAccounts, readDevEnvObj} from "../../src/tsyringe/devenvUtilities"
import {SifchainContractFactories} from "../../src/tsyringe/contracts"
import {buildDevEnvContracts} from "../../src/contractSupport"
import web3 from "web3"
import {EbRelayerAccount, crossChainFeeBase, crossChainBurnFee} from "../../src/devenv/sifnoded"
import * as dotenv from "dotenv"
import "@nomiclabs/hardhat-ethers"
import {ethers} from "hardhat"
import {SifnodedAdapter} from "./sifnodedAdapter"
import {checkSifnodeBurnState} from "./sifnode_lock_burn"
import {ethDenomHash} from "./context"

import {executeLock, checkEvmLockState} from "./evm_lock_burn"

chai.use(solidity)

describe("lock and burn tests", () => {
  dotenv.config()
  // This test only works when devenv is running, and that requires a connection to localhost
  expect(hardhat.network.name, "please use devenv").to.eq("localhost")

  const devEnvObject = readDevEnvObj("environment.json")
  // a generic sif address, nothing special about it
  const networkDescriptor = devEnvObject?.ethResults?.chainId ?? 9999

  const sifnodedAdapter: SifnodedAdapter = new SifnodedAdapter(
    devEnvObject!.sifResults!.adminAddress!.homeDir,
    devEnvObject!.sifResults!.adminAddress!.account,
    process.env["GOBIN"]
  )

  before("register HardhatRuntimeEnvironmentToken", async () => {
    container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
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
    // const sendAmount = BigNumber.from(5 * ETH) // 3500 gwei
    const sendAmount = BigNumber.from("5000000000000000000") // 3500 gwei

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

    let burnAmount = BigNumber.from("2300000000000000000") // 2300 gwei
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
})
