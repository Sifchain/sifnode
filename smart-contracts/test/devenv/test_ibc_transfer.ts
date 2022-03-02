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
import {getDenomHash, nullContractAddress} from "./context"
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
    const bridgeToken = await factories.bridgeToken

    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    )
    const destinationEthereumAddress = ethereumAccounts.availableAccounts[0]

    let testSifAccount: EbRelayerAccount = await sifnodedAdapter.createTestSifAccount(true, true, true)
    // the source address copy from registry json file
    const ibcTokenSourceAddress = "0x1111111111111111111111111111111111111111"
    const destinationAddress = await contracts.cosmosBridge.sourceAddressToDestinationAddress(ibcTokenSourceAddress)
    
    console.log("mapped destinationAddress is ", destinationAddress)
    const initSenderBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      IBC_TOKEN_DENOM
    )

    let initReceiverBalance = BigNumber.from(0)
    let createNewBridgeToken = destinationAddress == nullContractAddress
    if (!createNewBridgeToken) {
      initReceiverBalance = await bridgeToken.attach(destinationAddress).balanceOf(destinationEthereumAddress.address)
    }

    let lockAmount = BigNumber.from("1234")
    let crossChainCethFee = crossChainFeeBase * crossChainBurnFee

    const [newContractAddress, mintContractAddress] = await checkSifnodeLockState(
      sifnodedAdapter,
      contracts,
      testSifAccount,
      destinationEthereumAddress,
      lockAmount,
      IBC_TOKEN_DENOM,
      String(crossChainCethFee),
      networkDescriptor,
      createNewBridgeToken
    )

    console.log("New contract address is ", newContractAddress)
    console.log("Mint contract address is ", mintContractAddress)
    if (createNewBridgeToken) {
      expect(initReceiverBalance, "should be equal ").eq(
        BigNumber.from(0))
    } else {
      expect(newContractAddress, "should be equal ").eq(
        nullContractAddress)
      expect(mintContractAddress, "should be equal ").eq(
        destinationAddress)
    }

    // Here we verify the user balance is correct
    // get the balance after burn
    const finalSenderBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      IBC_TOKEN_DENOM
    )
    const finalReceiverBalance = await bridgeToken.attach(mintContractAddress).balanceOf(destinationEthereumAddress.address)

    console.log("Before burn the sender's balance is ", initSenderBalance)
    console.log("Before burn the receiver's balance is ", initReceiverBalance)

    console.log("After burn the sender's balance is ", finalSenderBalance)
    console.log("After burn the receiver's balance is ", finalReceiverBalance)


    expect(initSenderBalance.sub(lockAmount), "should be equal ").eq(
      finalSenderBalance
    )
    expect(initReceiverBalance.add(lockAmount), "should be equal ").eq(
      finalReceiverBalance
    )
  })
})
