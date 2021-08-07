import chai, {expect} from "chai"
import {BigNumber, BigNumberish, Contract, ContractFactory} from "ethers"
import {solidity} from "ethereum-waffle"
import web3 from "web3"
import * as ethereumAddress from "../src/ethereumAddress"
import {container, DependencyContainer} from "tsyringe";
import {
    BridgeBankProxy,
    BridgeTokenSetup, CosmosBridgeProxy,
    RowanContract,
    SifchainContractFactories
} from "../src/tsyringe/contracts";
import {BridgeBank, BridgeBank__factory, BridgeToken, CosmosBridge} from "../build";
import {SifchainAccounts, SifchainAccountsPromise} from "../src/tsyringe/sifchainAccounts";
import {
    DeploymentChainId,
    DeploymentDirectory,
    DeploymentName,
    HardhatRuntimeEnvironmentToken
} from "../src/tsyringe/injectionTokens";
import * as hardhat from "hardhat";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import {DeployedBridgeBank, DeployedBridgeToken} from "../src/contractSupport";
import {HardhatRuntimeEnvironment} from "hardhat/types";
import {impersonateAccount, setupSifchainMainnetDeployment} from "../src/hardhatFunctions";
import {buildTestToken} from "../src/sifchainHelpers";

chai.use(solidity)

describe("BridgeBank", () => {
    const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")

    let bridgeBank: BridgeBank

    before('register HardhatRuntimeEnvironmentToken', () => {
        container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    })

    before('get BridgeBank', async () => {
        await container.resolve(BridgeTokenSetup).complete
        bridgeBank = await container.resolve(BridgeBankProxy).contract
        expect(bridgeBank).to.exist
    })

    it("should deploy the BridgeBank, correctly setting the owner", async function () {
        const accounts = await container.resolve(SifchainAccountsPromise).accounts
        const bridgeBank = await container.resolve(BridgeBankProxy).contract
        const bridgeBankOwner = await bridgeBank.connect(accounts.ownerAccount).owner()
        expect(bridgeBankOwner).to.equal(accounts.ownerAccount.address);
        expect(bridgeBankOwner).to.equal(accounts.ownerAccount.address);
    })

    it("should correctly set initial values", async function () {
        expect(await bridgeBank.lockBurnNonce()).to.equal(0);
        expect(await bridgeBank.bridgeTokenCount()).to.equal(1)
    });

    it("should not allow a user to send ethereum directly to the contract", async function () {
        const accounts = await container.resolve(SifchainAccountsPromise).accounts
        await expect(hardhat.network.provider.send('eth_sendTransaction', [{
            to: bridgeBank.address,
            from: accounts.availableAccounts[0].address,
            value: "0x1"
        }])).to.be.reverted
    })

    describe("locking and burning", function () {
        let sender: SignerWithAddress;
        let accounts: SifchainAccounts;
        let amount: BigNumber
        let smallAmount: BigNumber
        let testToken: BridgeToken
        const invalidRecipient = web3.utils.utf8ToHex("esif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")

        before('create test token', async () => {
            accounts = await container.resolve(SifchainAccountsPromise).accounts
            sender = accounts.availableAccounts[0]
            amount = hardhat.ethers.utils.parseEther("100") as BigNumber
            smallAmount = amount.div(100)
            const testTokenFactory = (await container.resolve(SifchainContractFactories).bridgeToken).connect(sender)
            testToken = await testTokenFactory.deploy("test")
            await testToken.mint(sender.address, amount)
            await testToken.approve(bridgeBank.address, hardhat.ethers.constants.MaxUint256)
        })

        it("should lock a test token", async () => {
            const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")
            await bridgeBank.updateEthWhiteList(testToken.address, true)
            await expect(() => bridgeBank.connect(sender).lock(
                recipient,
                testToken.address,
                smallAmount,
                {
                    value: 0
                }
            )).to.changeTokenBalance(testToken, sender, smallAmount.mul(-1))
        })

        it("should fail to lock a test token to an invalid address", async () => {
            await expect(bridgeBank.connect(sender).lock(
                invalidRecipient,
                testToken.address,
                smallAmount,
                {
                    value: 0
                }
            )).to.be.revertedWith("Invalid len")
        })

        it("should lock eth", async () => {
            await expect(() => bridgeBank.connect(sender).lock(
                recipient,
                ethereumAddress.eth.address,
                smallAmount,
                {
                    value: smallAmount
                }
            )).to.changeEtherBalance(sender, smallAmount.mul(-1), {includeFee: false})
        })
    })

    it("should allow swapping erowan for rowan", async () => {
        const amount = 17

        const accounts = await container.resolve(SifchainAccountsPromise).accounts
        const migrator = accounts.availableAccounts[0]
        const bridgeBank = await container.resolve(BridgeBankProxy).contract
        const cosmosBridge = await container.resolve(CosmosBridgeProxy).contract as CosmosBridge
        await container.resolve(BridgeTokenSetup).complete
        const erowan = await container.resolve(RowanContract).contract
        await erowan.mint(migrator.address, 10000)
        await erowan.connect(migrator).approve(bridgeBank.address, hardhat.ethers.constants.MaxUint256)
        await bridgeBank.connect(migrator).burn(recipient, erowan.address, amount)
        const newRowanToken = await buildTestToken(hardhat, bridgeBank, "Rowan", accounts.ownerAccount, 1)
        await newRowanToken.addMinter(bridgeBank.address)
        await newRowanToken.mint(migrator.address, amount)
        await bridgeBank.connect(accounts.ownerAccount).addExistingBridgeToken(erowan.address)
        await bridgeBank.connect(accounts.ownerAccount).addExistingBridgeToken(newRowanToken.address)
        const tx = await bridgeBank.connect(migrator).migrateERowan(10)
        const wx = await tx.wait()
        await expect(() => bridgeBank.connect(migrator).migrateERowan(amount))
            .to.changeTokenBalance(newRowanToken, migrator, amount)
    })
})

