import chai, {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe";
import {SifchainContractFactories} from "../src/tsyringe/contracts";
import {BridgeBank, CosmosBridge} from "../build";
import {BridgeBankMainnetUpgradeAdmin, HardhatRuntimeEnvironmentToken} from "../src/tsyringe/injectionTokens";
import * as hardhat from "hardhat";
import {DeployedBridgeBank, DeployedCosmosBridge} from "../src/contractSupport";
import {impersonateAccount, setupSifchainMainnetDeployment, startImpersonateAccount} from "../src/hardhatFunctions"
import {SifchainAccountsPromise} from "../src/tsyringe/sifchainAccounts";
import web3 from "web3";
import {BigNumber, BigNumberish, ContractTransaction} from "ethers";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";

chai.use(solidity)

describe("BridgeBank and CosmosBridge - updating to latest smart contracts", () => {
    let deployedBridgeBank: BridgeBank

    before('register HardhatRuntimeEnvironmentToken', () => {
        container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    })

    before('use mainnet data', async () => {
        await setupSifchainMainnetDeployment(container, hardhat)
    })

    describe("upgraded BridgeBank", async () => {
        it("should maintain existing stored values", async () => {
            const existingBridgeBank = await container.resolve(DeployedBridgeBank).contract
            const bridgeBankFactory = await container.resolve(SifchainContractFactories).bridgeBank
            const upgradeAdmin = container.resolve(BridgeBankMainnetUpgradeAdmin) as string

            const existingOperator = await existingBridgeBank.operator()
            const existingOracle = await existingBridgeBank.oracle()
            const existingCosmosBridge = await existingBridgeBank.cosmosBridge()
            const existingOwner = await existingBridgeBank.owner()


            const newBridgeBank = await impersonateAccount(hardhat, upgradeAdmin, hardhat.ethers.utils.parseEther("10"), async fakeDeployer => {
                const signedBridgeBankFactory = bridgeBankFactory.connect(fakeDeployer)
                return await hardhat.upgrades.upgradeProxy(existingBridgeBank, signedBridgeBankFactory) as BridgeBank
            })

            expect(existingOperator).to.equal(await newBridgeBank.operator())
            expect(existingOracle).to.equal(await newBridgeBank.oracle())
            expect(existingCosmosBridge).to.equal(await newBridgeBank.cosmosBridge())
            expect(existingOwner).to.equal(await newBridgeBank.owner())
        })

        // TODO this function should track validators added and removed and pay attention to resets.
        // None of those have happened yet on mainnet, so we'll just use validators added.
        async function currentValidators(cosmosBridge: CosmosBridge): Promise<readonly string[]> {
            const validatorsAdded = await cosmosBridge.queryFilter(cosmosBridge.filters.LogValidatorAdded())
            return validatorsAdded.map(t => t.args[0])
        }

        it("should burn ceth via existing validators", async () => {
            const existingCosmosBridge = await container.resolve(DeployedCosmosBridge).contract as CosmosBridge
            const existingValidators = await currentValidators(existingCosmosBridge)

            const upgradeAdmin = container.resolve(BridgeBankMainnetUpgradeAdmin) as string

            await impersonateAccount(hardhat, upgradeAdmin, hardhat.ethers.utils.parseEther("10"), async fakeDeployer => {
                const amount = BigNumber.from(100)
                const accounts = await container.resolve(SifchainAccountsPromise).accounts
                const bridgeBankFactory = await container.resolve(SifchainContractFactories).bridgeBank
                const signedBBFactory = bridgeBankFactory.connect(fakeDeployer)
                const existingBridgeBank = await container.resolve(DeployedBridgeBank).contract
                const operator = await startImpersonateAccount(hardhat, await existingBridgeBank.operator())
                const newBridgeBank = (await hardhat.upgrades.upgradeProxy(existingBridgeBank, signedBBFactory) as BridgeBank).connect(operator)
                const cosmosBridgeFactory = (await container.resolve(SifchainContractFactories).cosmosBridge).connect(fakeDeployer)
                const newCosmosBridge = (await hardhat.upgrades.upgradeProxy(existingCosmosBridge, cosmosBridgeFactory, {unsafeAllowCustomTypes: true}) as CosmosBridge).connect(operator)
                const testTokenFactory = (await container.resolve(SifchainContractFactories).bridgeToken).connect(operator)
                const testToken = await testTokenFactory.deploy("test")
                await testToken.mint(operator.address, amount)
                await testToken.connect(operator).approve(newBridgeBank.address, amount)

                const validators = await currentValidators(existingCosmosBridge)
                expect(validators).to.deep.equal(existingValidators, "validators should not have changed")
                const receiver = accounts.availableAccounts[0]
                const impersonatedValidators = await Promise.all(validators.map(v => startImpersonateAccount(hardhat, v)))

                {
                    const startingBalance = await receiver.getBalance()
                    const prophecyResult = await executeNewProphecyClaimWithTestValues("burn", receiver.address, "eth", amount, newCosmosBridge, impersonatedValidators)
                    console.log("canaryhere")
                    expect(prophecyResult.length).to.equal(validators.length - 1, "we expected one of the validators to fail after the prophecy was completed")
                    expect(await receiver.getBalance()).to.equal(startingBalance.add(amount))
                }
                {
                    // need to add a test token so we can burn it
                    const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")
                    await newBridgeBank.updateEthWhiteList(testToken.address, true)
                    await newBridgeBank.lock(
                        recipient,
                        testToken.address,
                        amount,
                        {
                            value: 0
                        }
                    )

                    const startingBalance = await testToken.balanceOf(receiver.address)
                    const prophecyResult = await executeNewProphecyClaimWithTestValues("burn", receiver.address, "test", amount, newCosmosBridge, impersonatedValidators)
                    expect(prophecyResult.length).to.equal(validators.length - 1, "we expected one of the validators to fail after the prophecy was completed")
                    expect(await testToken.balanceOf(receiver.address)).to.equal(startingBalance.add(amount))
                }
            })
        })

        describe("should impersonate all four relayers", async () => {
            it("should turn ceth to eth in an unlock")
            it("should turn rowan to erowan in a lock")
            it("should turn random ERC20 pegged token back to unlock that token on mainnet")
        })
    })
})

let sequenceNumber = BigNumber.from(0)

async function executeNewProphecyClaimWithTestValues(
    claimType: "burn" | "lock",
    ethereumReceiver: string,
    symbol: string,
    amount: BigNumberish,
    cosmosBridge: CosmosBridge,
    validators: Array<SignerWithAddress>
): Promise<readonly ContractTransaction[]> {
    const cosmosSender = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")
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
                sequenceNumber,
                ethereumReceiver,
                symbol,
                amount
            )
            result.push(tx)
        } catch (e) {
            console.log("errorris: ", e)
            // we expect one of these to fail since the prophecy completes before all validators submit their prophecy claim
            // and we only return the successful transactions
        }
    }
    sequenceNumber = sequenceNumber.add(1)
    return result
}
