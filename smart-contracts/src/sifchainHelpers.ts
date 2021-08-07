import {container} from "tsyringe";
import {SifchainContractFactories} from "./tsyringe/contracts";
import * as hardhat from "hardhat";
import {HardhatRuntimeEnvironment} from "hardhat/types";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import {BigNumberish} from "ethers";
import {DeployedBridgeBank} from "./contractSupport";
import {BridgeBank} from "../build";

export async function buildTestToken(
    hre: HardhatRuntimeEnvironment,
    bridgeBank: BridgeBank,
    symbol: string,
    owner: SignerWithAddress,
    amount: BigNumberish
) {
    const testTokenFactory = (await container.resolve(SifchainContractFactories).bridgeToken).connect(owner)
    const testToken = await testTokenFactory.deploy(symbol)
    await testToken.mint(owner.address, amount)
    await testToken.approve(bridgeBank.address, hardhat.ethers.constants.MaxUint256)
    return testToken
}
