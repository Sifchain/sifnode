import chai, {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe";
import {SifchainContractFactories} from "../src/tsyringe/contracts";
import {HardhatRuntimeEnvironmentToken} from "../src/tsyringe/injectionTokens";
import * as hardhat from "hardhat";
import {DeployedBridgeBank} from "../src/contractSupport";
import {impersonateAccount, setupSifchainMainnetDeployment} from "../src/hardhatFunctions"
import {buildIbcTokens, readTokenData} from "../src/ibcMatchingTokens";

chai.use(solidity)

describe("IBC matching tokens", () => {
    before('register HardhatRuntimeEnvironmentToken', () => {
        container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    })

    before('use mainnet data', async () => {
        await setupSifchainMainnetDeployment(container, hardhat, "sifchain")
    })

    describe("creating and reassigning ibc tokens", async () => {
        it("should create and use ibc tokens", async () => {
            const existingBridgeBank = await container.resolve(DeployedBridgeBank).contract
            const bridgeBankOperator = await existingBridgeBank.operator()
            const ibcTokenFactory = await container.resolve(SifchainContractFactories).ibcToken
            await impersonateAccount(hardhat, bridgeBankOperator, hardhat.ethers.utils.parseEther("10"), async fakeDeployer => {
                const tokenData = await readTokenData("../test_data/ibc_token_data.json")
                expect(tokenData).to.have("size").eq(2)
                const x = buildIbcTokens(ibcTokenFactory, tokenData, existingBridgeBank.connect(fakeDeployer), true)
                expect(x).to.eq(3)
            })
        })
    })
})
