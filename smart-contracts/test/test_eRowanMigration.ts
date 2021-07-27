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

describe("BridgeBank", () => {
    let deployedBridgeBank: BridgeBank

    before('register HardhatRuntimeEnvironmentToken', () => {
        container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    })

    describe("migration of existing erowan to new rowan", async () => {
        let existingRowanToken: BridgeToken
        let existingBridgeBank: BridgeBank
        let newBridgeBank: BridgeBank
        let newRowanToken: BridgeToken
        let upgradeAdmin: string

        before("upgraded BridgeBank", async () => {
            existingRowanToken = await container.resolve(DeployedBridgeToken).contract
            existingBridgeBank = await container.resolve(DeployedBridgeBank).contract
            const bridgeBankFactory = await container.resolve(SifchainContractFactories).bridgeBank
            upgradeAdmin = container.resolve(BridgeBankMainnetUpgradeAdmin) as string
            await impersonateAccount(hardhat, upgradeAdmin, hardhat.ethers.utils.parseEther("10"), async fakeDeployer => {
                const signedBBFactory = bridgeBankFactory.connect(fakeDeployer)
                newBridgeBank = await hardhat.upgrades.upgradeProxy(existingBridgeBank, signedBBFactory) as BridgeBank
            })
        })

        before("set up new rowan token", async () => {
            const bridgeTokenFactory = await container.resolve(SifchainContractFactories).bridgeToken
            newRowanToken = await bridgeTokenFactory.deploy("rowan")
        })

        it("should have the same owner", async () => {
            const oldOwner = await existingBridgeBank.owner()
            const newOwner = await newBridgeBank.owner()
            expect(oldOwner).to.equal(newOwner)
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
})
