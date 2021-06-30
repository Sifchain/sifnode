import {DeployFunction} from 'hardhat-deploy/types';
import {HardhatRuntimeEnvironment} from 'hardhat/types';
import '@openzeppelin/hardhat-upgrades';

import {ContractNames} from "../src/contractNames";
import {loadDeploymentEnvWithDotenv} from "../src/deploymentEnv";
import deployBridgeRegistry from "./deployBridgeRegistry";

const deployBridgeBank: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const BridgeBank = await hre.ethers.getContractFactory(ContractNames.BridgeBank)
    const {owner} = await hre.getNamedAccounts()
    const {pauser} = await hre.getNamedAccounts()
    const {validator1} = await hre.getNamedAccounts()
    const deploymentEnv = loadDeploymentEnvWithDotenv()
    const CosmosBridge = await hre.deployments.get(ContractNames.CosmosBridge)
    let bridgeBankArgs = [
        CosmosBridge.address,
        owner,
        pauser
    ];
    const bridgeBank = await hre.upgrades.deployProxy(BridgeBank, bridgeBankArgs);
    await bridgeBank.deployed()
    console.log("deployed BridgeBank to: ", bridgeBank.address, bridgeBankArgs);
};

deployBridgeBank.tags = [ContractNames.BridgeBank]
deployBridgeBank.dependencies = [ContractNames.CosmosBridge]
export default deployBridgeBank;