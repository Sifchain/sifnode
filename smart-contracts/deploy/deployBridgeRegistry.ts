import {DeployFunction} from 'hardhat-deploy/types';
import {HardhatRuntimeEnvironment} from 'hardhat/types';
import '@openzeppelin/hardhat-upgrades';

import {ContractNames} from "../src/contractNames";

const deployBridgeRegistry: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const BridgeRegistry = await hre.ethers.getContractFactory(ContractNames.BridgeRegistry)
    const cosmosBridge = await hre.deployments.get(ContractNames.CosmosBridge)
    const bridgeBank = await hre.deployments.get(ContractNames.BridgeBank)

    const bridgeRegistry = await hre.upgrades.deployProxy(BridgeRegistry, [
        cosmosBridge.address,
        bridgeBank.address
    ]);
};

deployBridgeRegistry.tags = [ContractNames.BridgeRegistry]
deployBridgeRegistry.dependencies = [ContractNames.CosmosBridge, ContractNames.BridgeBank]
export default deployBridgeRegistry