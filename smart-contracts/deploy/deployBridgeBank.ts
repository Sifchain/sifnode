import {DeployFunction} from 'hardhat-deploy/types';
import {HardhatRuntimeEnvironment} from 'hardhat/types';
import '@openzeppelin/hardhat-upgrades';

import {ContractNames} from "../src/contractNames";
import {loadDeploymentEnvWithDotenv} from "../src/deploymentEnv";
import deployBridgeRegistry from "./deployBridgeRegistry";
import {DeployOptions} from "hardhat-deploy/dist/types";

const deployBridgeBank: DeployFunction = async function (hre: HardhatRuntimeEnvironment): Promise<void> {
    console.log("starting deployBridgeBank")
    // const BridgeBank = await hre.ethers.getContractFactory(ContractNames.BridgeBank)
    console.log("starting deployBridgeBank1")
    const {owner} = await hre.getNamedAccounts()
    const {operator} = await hre.getNamedAccounts()
    const {pauser} = await hre.getNamedAccounts()
    const {validator1} = await hre.getNamedAccounts()
    const deploymentEnv = loadDeploymentEnvWithDotenv()
    const CosmosBridge = await hre.deployments.get(ContractNames.CosmosBridge)
    let bridgeBankArgs = [
        CosmosBridge.address,
        owner,
        pauser
    ];
    await hre.deployments.deploy(ContractNames.BridgeBank, {
        from: operator,
        args: bridgeBankArgs,
        proxy: true
    })
};

deployBridgeBank.tags = [ContractNames.BridgeBank]
deployBridgeBank.dependencies = [ContractNames.CosmosBridge]
export default deployBridgeBank;