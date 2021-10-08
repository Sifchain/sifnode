import chai, {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe";
import {HardhatRuntimeEnvironmentToken} from "../../src/tsyringe/injectionTokens";
import * as hardhat from "hardhat";
import {BigNumber, Wallet} from "ethers";
import {ethereumResultsToSifchainAccounts, readDevEnvObj} from "../../src/tsyringe/devenvUtilities";
import {createSignerWithAddresss} from "../../src/tsyringe/hardhatSupport";
import {SifchainContractFactories} from "../../src/tsyringe/contracts";
import {buildDevEnvContracts} from "../../src/contractSupport";
import web3 from "web3";
import * as ethereumAddress from "../../src/ethereumAddress";

chai.use(solidity)

describe("accounts", async () => {
    it("should return the address matching the private key", async () => {
        let addr = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266";
        const x = createSignerWithAddresss(
            addr,
            "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
        )
        expect(addr.toLowerCase()).to.eq(x.address.toLowerCase())
    })
})

describe("watcher", () => {
    const devEnvObject = readDevEnvObj("environment.json")
    const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace")

    before('register HardhatRuntimeEnvironmentToken', async () => {
        container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
        const sifchainAccounts = await ethereumResultsToSifchainAccounts(devEnvObject.ethResults!)
        const factories = container.resolve(SifchainContractFactories)
        const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)
        console.log("xis: ", contracts.bridgeBank.address)
        console.log("xis: ", await contracts.bridgeBank.owner())
        const sender1 = sifchainAccounts.availableAccounts[0]
        const sender = new Wallet("0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a", hardhat.ethers.provider)
        const smallAmount = BigNumber.from(1017)
        await expect(() => contracts.bridgeBank.connect(sender).lock(
            recipient,
            ethereumAddress.eth.address,
            smallAmount,
            {
                value: smallAmount
            }
        )).to.changeEtherBalance(sender, smallAmount.mul(-1), {includeFee: false})
    })

    it("should get the accounts from devenv")
    it("should send a lock transaction")
    it("should watch evmrelayer logs")
    it("should watch for evm events")
    it("should fail if evmrelayer gets an error")
})
