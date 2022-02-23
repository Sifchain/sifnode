import * as chai from "chai"
import {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe"
import {HardhatRuntimeEnvironmentToken} from "../../src/tsyringe/injectionTokens"
import * as hardhat from "hardhat"
import {BigNumber} from "ethers"
import {ethereumResultsToSifchainAccounts, readDevEnvObj} from "../../src/tsyringe/devenvUtilities"
import {SifchainContractFactories, MINTER_ROLE} from "../../src/tsyringe/contracts"
import {buildDevEnvContracts} from "../../src/contractSupport"
import web3 from "web3"
import {EbRelayerAccount, crossChainFeeBase, crossChainBurnFee} from "../../src/devenv/sifnoded"

import * as dotenv from "dotenv"
import "@nomiclabs/hardhat-ethers"
import {SifnodedAdapter, IBC_TOKEN_DENOM} from "./sifnodedAdapter"
import {getDenomHash, ethDenomHash} from "./context"
import {checkSifnodeLockState} from "./sifnode_lock"
import {executeLock, checkEvmLockState} from "./evm_lock"
import {SifchainAccountsPromise} from "../../src/tsyringe/sifchainAccounts"

chai.use(solidity)
const rowan = "rowan"

describe("lock rowan token tests", () => {
  dotenv.config()
  // This test only works when devenv is running, and that requires a connection to localhost
  expect(hardhat.network.name, "please use devenv").to.eq("localhost")

  const devEnvObject = readDevEnvObj("environment.json")
  const networkDescriptor = devEnvObject?.ethResults?.chainId ?? 31337

  const sifnodedAdapter: SifnodedAdapter = new SifnodedAdapter(
    devEnvObject!.sifResults!.adminAddress!.homeDir,
    devEnvObject!.sifResults!.adminAddress!.account,
    process.env["GOBIN"]
  )

  before("register HardhatRuntimeEnvironmentToken", async () => {
    container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
  })

  it.only("should allow lock rowan token to Ethereum", async () => {
    // TODO: Could these be moved out of the test fx? and instantiated via beforeEach?
    const factories = container.resolve(SifchainContractFactories)
    const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)

    const sifchainAccountsPromise = container.resolve(SifchainAccountsPromise)
    const ownerAccount = (await sifchainAccountsPromise.accounts).ownerAccount

    // add rowan contract into whitelist, then bridge bank can mint the token
    await contracts.bridgeBank
      .connect(ownerAccount)
      .addExistingBridgeToken(contracts.rowanContract.address)

    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    )
    const destinationEthereumAddress = ethereumAccounts.availableAccounts[0]

    let testSifAccount: EbRelayerAccount = await sifnodedAdapter.createTestSifAccount(true, true, true)

    const initSenderBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      rowan
    )
    const initReceiverBalance = await contracts.rowanContract.balanceOf(destinationEthereumAddress.address)

    let lockAmount = BigNumber.from("1234")
    let crossChainCethFee = crossChainFeeBase * crossChainBurnFee

    await checkSifnodeLockState(
      sifnodedAdapter,
      contracts,
      testSifAccount,
      destinationEthereumAddress,
      lockAmount,
      rowan,
      String(crossChainCethFee),
      networkDescriptor
    )

    // Here we verify the user balance is correct
    // get the balance after burn
    const finalSenderBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      rowan
    )
    const finalReceiverBalance = await contracts.rowanContract.balanceOf(destinationEthereumAddress.address)

    console.log("Before burn the sender's balance is ", initSenderBalance)
    console.log("Before burn the receiver's balance is ", initReceiverBalance)

    console.log("After burn the sender's balance is ", finalSenderBalance)
    console.log("After burn the receiver's balance is ", finalReceiverBalance)

    // expect(initialErc20SenderBalance.sub(burnAmount), "should be equal ").eq(
    //   finalErc20SenderBalance
    // )
    // expect(initialErc20ReceiverBalance.add(burnAmount), "should be equal ").eq(
    //   finalErc20ReceiverBalance
    // )
  })
})
