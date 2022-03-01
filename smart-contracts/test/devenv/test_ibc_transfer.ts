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
import {SifchainAccountsPromise} from "../../src/tsyringe/sifchainAccounts"

chai.use(solidity)

describe("burn ibc token tests", () => {
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

  it.only("should allow burn ibc token to Ethereum", async () => {
    // TODO: Could these be moved out of the test fx? and instantiated via beforeEach?
    const factories = container.resolve(SifchainContractFactories)
    const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)

    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    )
    const destinationEthereumAddress = ethereumAccounts.availableAccounts[0]

    let testSifAccount: EbRelayerAccount = await sifnodedAdapter.createTestSifAccount(true, true, true)

    const initSenderBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      IBC_TOKEN_DENOM
    )

    let lockAmount = BigNumber.from("1234")
    let crossChainCethFee = crossChainFeeBase * crossChainBurnFee

    await checkSifnodeLockState(
      sifnodedAdapter,
      contracts,
      testSifAccount,
      destinationEthereumAddress,
      lockAmount,
      IBC_TOKEN_DENOM,
      String(crossChainCethFee),
      networkDescriptor
    )

    // Here we verify the user balance is correct
    // get the balance after burn
    const finalErc20SenderBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      IBC_TOKEN_DENOM
    )
    // const finalErc20ReceiverBalance = await erc20.balanceOf(destinationEthereumAddress.address)

    // console.log("Before burn the sender's balance is ", initialErc20SenderBalance)
    // console.log("Before burn the receiver's balance is ", initialErc20ReceiverBalance)

    console.log("After burn the sender's balance is ", finalErc20SenderBalance)
    // console.log("After burn the receiver's balance is ", finalErc20ReceiverBalance)

    // expect(initialErc20SenderBalance.sub(burnAmount), "should be equal ").eq(
    //   finalErc20SenderBalance
    // )
    // expect(initialErc20ReceiverBalance.add(burnAmount), "should be equal ").eq(
    //   finalErc20ReceiverBalance
    // )
  })
})
