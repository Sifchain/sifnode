import chai, {expect} from "chai"
import {BigNumber} from "ethers"
import {solidity} from "ethereum-waffle"
import web3 from "web3"
import * as ethereumAddress from "../src/ethereumAddress"
import {container} from "tsyringe";
import {BridgeBankProxy, SifchainContractFactories} from "../src/tsyringe/contracts";
import {BridgeBank, BridgeToken} from "../build";
import {SifchainAccounts, SifchainAccountsPromise} from "../src/tsyringe/sifchainAccounts";
import {HardhatRuntimeEnvironmentToken} from "../src/tsyringe/injectionTokens";
import * as hardhat from "hardhat";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";

chai.use(solidity)

describe("BridgeBank", () => {
    let bridgeBank: BridgeBank

    before('register HardhatRuntimeEnvironmentToken', () => {
        container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    })

    before('get BridgeBank', async () => {
        bridgeBank = await container.resolve(BridgeBankProxy).contract
        expect(bridgeBank).to.exist
    })

    it("should deploy the BridgeBank, correctly setting the owner", async function () {
        const bridgeBankOwner = await bridgeBank.owner()
        const accounts = await container.resolve(SifchainAccountsPromise).accounts
        expect(bridgeBankOwner).to.equal(accounts.ownerAccount.address);
    })

    it("should correctly set initial values", async function () {
        expect(await bridgeBank.lockBurnNonce()).to.equal(0);
        expect(await bridgeBank.bridgeTokenCount()).to.equal(0)
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
        let operator: SignerWithAddress;
        let accounts: SifchainAccounts;
        let amount: BigNumber
        let smallAmount: BigNumber
        let testToken: BridgeToken
        const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")
        const invalidRecipient = web3.utils.utf8ToHex("esif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")

        before('create test token', async () => {
            accounts = await container.resolve(SifchainAccountsPromise).accounts
            sender = accounts.availableAccounts[0]
            operator = accounts.operatorAccount
            amount = hardhat.ethers.utils.parseEther("100") as BigNumber
            smallAmount = amount.div(100)
            const testTokenFactory = (await container.resolve(SifchainContractFactories).bridgeToken).connect(sender)
            testToken = await testTokenFactory.deploy("TEST token", "test", 18)
            await testToken.mint(sender.address, amount)
            await testToken.approve(bridgeBank.address, hardhat.ethers.constants.MaxUint256)
        })

        it("should lock a test token", async () => {
            // Add the token into white list
            await bridgeBank.connect(operator).updateEthWhiteList(testToken.address, true)

            const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")
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
            )).to.be.revertedWith("INV_SIF_ADDR")
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
})
