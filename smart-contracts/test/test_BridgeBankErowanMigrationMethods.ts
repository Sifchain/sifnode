import chai, {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe";
import {SifchainContractFactories} from "../src/tsyringe/contracts";
import {BridgeBank, BridgeToken} from "../build";
import {
    BridgeBankMainnetUpgradeAdmin,
    HardhatRuntimeEnvironmentToken
} from "../src/tsyringe/injectionTokens";
import * as hardhat from "hardhat";
import {DeployedBridgeBank, DeployedBridgeToken} from "../src/contractSupport";
import {impersonateAccount, setupSifchainMainnetDeployment} from "../src/hardhatFunctions"
import {SifchainAccountsPromise} from "../src/tsyringe/sifchainAccounts";

chai.use(solidity)

describe("BridgeBank with eRowan migration functionality", () => {
    let deployedBridgeBank: BridgeBank

    before('register HardhatRuntimeEnvironmentToken', () => {
        container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    })


    it("should change balances correctly", async () => {
        // const existingRowanToken = await container.resolve(DeployedBridgeToken).contract
        // existingRowanToken.approve(acc)
        // const bbevents = await bridgeBank.queryFilter({address: undefined} as EventFilter, 12865480)
        // // const events = await existingRowanToken.queryFilter({address: undefined} as EventFilter, 12865480, 12865480)
        // console.log("printme: ", await existingRowanToken.name(), bbevents)
    })
    it("should fire correct events")
    it("should do nothing before setRowanTokens is called", async () => {
        // const accounts = await container.resolve(SifchainAccountsPromise).accounts
        // const testAccount = accounts.availableAccounts[0]
        // await impersonateAccount(
        //     hardhat,
        //     await newBridgeBank.operator(),
        //     hardhat.ethers.utils.parseEther("10"),
        //     async impersonatedBridgeBankOperator => {
        //         let bb = newBridgeBank.connect(impersonatedBridgeBankOperator);
        //         await newRowanToken.mint(testAccount.address, 100)
        //         await bb.mintBridgeTokens(testAccount.address, "erowan", 100)
        //     }
        // )
    })
    it("should be able to terminate after some circumstances") // Do we want to leave this open forever?
    it("should be able to call setRowanTokens", async () => {
        await impersonateAccount(
            hardhat,
            await newBridgeBank.operator(),
            hardhat.ethers.utils.parseEther("10"),
            async impersonatedBridgeBankOperator => {
                const a = await newBridgeBank.connect(impersonatedBridgeBankOperator).setRowanTokens(existingRowanToken.address, newRowanToken.address)
                expect(await existingRowanToken.name()).to.equal("erowan")
                expect(await newRowanToken.name()).to.equal("rowan")
            }
        )
    })
})
