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

    describe("with updated BridgeBank contract", async () => {
        let existingRowanToken: BridgeToken
        let existingBridgeBank: BridgeBank
        let newBridgeBank: BridgeBank
        let newRowanToken: BridgeToken
        let upgradeAdmin: string

        it("should upgrade BridgeBank, maintaining stored values", async () => {
            existingRowanToken = await container.resolve(DeployedBridgeToken).contract
            existingBridgeBank = await container.resolve(DeployedBridgeBank).contract
            const bridgeBankFactory = await container.resolve(SifchainContractFactories).bridgeBank
            upgradeAdmin = container.resolve(BridgeBankMainnetUpgradeAdmin) as string

            // What values are stored?
            // whitelist
            // operator
            // oracle
            // cosmos
            const existingOperator = await existingBridgeBank.operator()
            const existingOracle = await existingBridgeBank.oracle()
            const existingCosmosBridge = await existingBridgeBank.cosmosBridge()
            const existingOwner = await existingBridgeBank.owner()

            await impersonateAccount(hardhat, upgradeAdmin, hardhat.ethers.utils.parseEther("10"), async fakeDeployer => {
                const signedBBFactory = bridgeBankFactory.connect(fakeDeployer)
                newBridgeBank = await hardhat.upgrades.upgradeProxy(existingBridgeBank, signedBBFactory) as BridgeBank
            })

            expect(existingOperator).to.equal(await newBridgeBank.operator())
            expect(existingOracle).to.equal(await newBridgeBank.oracle())
            expect(existingCosmosBridge).to.equal(await newBridgeBank.cosmosBridge())
            expect(existingOwner).to.equal(await newBridgeBank.owner())
        })
    })
})
