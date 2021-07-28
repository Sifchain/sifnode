import chai, {expect} from "chai"
import {BigNumber, BigNumberish, ContractFactory} from "ethers"
import {solidity} from "ethereum-waffle"
import web3 from "web3"
import * as ethereumAddress from "../src/ethereumAddress"
import {container, DependencyContainer} from "tsyringe";
import {
    BridgeBankProxy,
    BridgeTokenSetup,
    RowanContract,
    SifchainContractFactories
} from "../src/tsyringe/contracts";
import {BridgeBank, BridgeBank__factory, BridgeToken} from "../build";
import {SifchainAccounts, SifchainAccountsPromise} from "../src/tsyringe/sifchainAccounts";
import {
    DeploymentChainId,
    DeploymentDirectory,
    DeploymentName,
    HardhatRuntimeEnvironmentToken
} from "../src/tsyringe/injectionTokens";
import * as hardhat from "hardhat";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import {DeployedBridgeBank} from "../src/contractSupport";
import {HardhatRuntimeEnvironment} from "hardhat/types";
import {impersonateAccount, setupSifchainMainnetDeployment} from "../src/hardhatFunctions";

chai.use(solidity)

describe("General contact functions", () => {

    before('register HardhatRuntimeEnvironmentToken', () => {
        container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    })

    it("should get the proxy owner", async () => {
        // expect.fail()
    })

    it ("should display useful information for the contract", async () => {
        const bb = container.resolve(BridgeBankProxy).contract
        console.log("gotbb: ", bb)
    })
})

