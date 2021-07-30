import chai, {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe";
import {SifchainContractFactories} from "../src/tsyringe/contracts";
import {BridgeBank, BridgeToken, CosmosBridge} from "../build";
import {
    BridgeBankMainnetUpgradeAdmin,
    HardhatRuntimeEnvironmentToken
} from "../src/tsyringe/injectionTokens";
import * as hardhat from "hardhat";
import {DeployedBridgeBank, DeployedBridgeToken, DeployedCosmosBridge} from "../src/contractSupport";
import {
    impersonateAccount,
    setupSifchainMainnetDeployment,
    startImpersonateAccount
} from "../src/hardhatFunctions"
import {SifchainAccountsPromise} from "../src/tsyringe/sifchainAccounts";
import web3 from "web3";
import {BigNumber, BigNumberish, ContractTransaction} from "ethers";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";

chai.use(solidity)

describe("BridgeBank and CosmosBridge after updating to latest smart contracts", () => {
    let deployedBridgeBank: BridgeBank

    before('register HardhatRuntimeEnvironmentToken', () => {
        container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    })

    before('use mainnet data', async () => {
        await setupSifchainMainnetDeployment(container, hardhat)
    })

    describe("upgrade the BridgeBank contract", async () => {
        let existingRowanToken: BridgeToken
        let existingBridgeBank: BridgeBank
        let newBridgeBank: BridgeBank
        let newCosmosBridge: CosmosBridge
        let newRowanToken: BridgeToken
        let upgradeAdmin: string

        it("should maintain existing values", async () => {
            existingRowanToken = await container.resolve(DeployedBridgeToken).contract
            existingBridgeBank = await container.resolve(DeployedBridgeBank).contract
            const bridgeBankFactory = await container.resolve(SifchainContractFactories).bridgeBank
            upgradeAdmin = container.resolve(BridgeBankMainnetUpgradeAdmin) as string

            const existingOperator = await existingBridgeBank.operator()
            const existingOracle = await existingBridgeBank.oracle()
            const existingCosmosBridge = await existingBridgeBank.cosmosBridge()
            const existingOwner = await existingBridgeBank.owner()

            await impersonateAccount(hardhat, upgradeAdmin, hardhat.ethers.utils.parseEther("10"), async fakeDeployer => {
                const signedBBFactory = bridgeBankFactory.connect(fakeDeployer)
                newBridgeBank = await hardhat.upgrades.upgradeProxy(existingBridgeBank, signedBBFactory) as BridgeBank

                // Deploy CosmosBridge
            })

            expect(existingOperator).to.equal(await newBridgeBank.operator())
            expect(existingOracle).to.equal(await newBridgeBank.oracle())
            expect(existingCosmosBridge).to.equal(await newBridgeBank.cosmosBridge())
            expect(existingOwner).to.equal(await newBridgeBank.owner())
        })

        async function executeNewProphecyClaim(
            claimType: "burn" | "lock",
            ethereumReceiver: string,
            symbol: string,
            amount: BigNumberish,
            validators: Array<SignerWithAddress>
        ): Promise<readonly ContractTransaction[]> {
            const cosmosSender = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")
            const cosmosBridge = await container.resolve(DeployedCosmosBridge).contract as CosmosBridge
            const claimTypeValue = {
                "burn": 1,
                "lock": 2
            }[claimType]
            const result = new Array<ContractTransaction>()
            for (const validator of validators) {
                try {
                    const tx = await cosmosBridge.connect(validator).newProphecyClaim(
                        claimTypeValue,
                        cosmosSender,
                        17,
                        ethereumReceiver,
                        symbol,
                        amount
                    )
                    result.push(tx)
                } catch (e) {
                    // we expect one of these to fail since the prophecy completes before all validators submit their prophecy claim
                }
            }
            return result
        }

        // TODO this function should track validators added and removed and pay attention to resets.
        // None of those have happened yet on mainnet, so we'll just use validators added.
        async function currentValidators(cosmosBridge: CosmosBridge): Promise<readonly string[]> {
            const validatorsAdded = await cosmosBridge.queryFilter(cosmosBridge.filters.LogValidatorAdded())
            return validatorsAdded.map(t => t.args[0])
        }

        it("should burn ceth via existing validators", async () => {
            const accounts = await container.resolve(SifchainAccountsPromise).accounts
            const cosmosBridge = await container.resolve(DeployedCosmosBridge).contract as CosmosBridge

            const upgradeAdmin = container.resolve(BridgeBankMainnetUpgradeAdmin) as string

            const existingOperator = await existingBridgeBank.operator()
            const existingOracle = await existingBridgeBank.oracle()
            const existingCosmosBridge = await existingBridgeBank.cosmosBridge()
            const existingOwner = await existingBridgeBank.owner()

            await impersonateAccount(hardhat, upgradeAdmin, hardhat.ethers.utils.parseEther("10"), async fakeDeployer => {
                const bridgeBankFactory = await container.resolve(SifchainContractFactories).bridgeBank
                const signedBBFactory = bridgeBankFactory.connect(fakeDeployer)
                newBridgeBank = await hardhat.upgrades.upgradeProxy(existingBridgeBank, signedBBFactory) as BridgeBank
                const cosmosBridgeFactory = (await container.resolve(SifchainContractFactories).cosmosBridge).connect(fakeDeployer)
                newCosmosBridge = (await hardhat.upgrades.upgradeProxy(cosmosBridge, cosmosBridgeFactory, {unsafeAllowCustomTypes: true}) as CosmosBridge).connect(fakeDeployer)
                // Deploy CosmosBridge
            })

            const validators = await currentValidators(cosmosBridge)
            const receiver = accounts.availableAccounts[0]
            const amount = BigNumber.from(100)
            const impersonatedValidators = await Promise.all(validators.map(v => startImpersonateAccount(hardhat, v)))
            const startingBalance = await receiver.getBalance()
            const prophecyResult = await executeNewProphecyClaim("burn", receiver.address, "eth", amount, impersonatedValidators)
            expect(prophecyResult.length).to.equal(validators.length - 1, "we expected one of the validators to fail after the prophecy was completed")
            expect(await receiver.getBalance()).to.equal(startingBalance.add(amount))
        })

        describe("should impersonate all four relayers", async () => {
            it("should turn ceth to eth in an unlock")
            it("should turn rowan to erowan in a lock")
            it("should turn random ERC20 pegged token back to unlock that token on mainnet")
        })
    })
})
