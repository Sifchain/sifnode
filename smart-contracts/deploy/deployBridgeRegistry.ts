import {DeployFunction} from 'hardhat-deploy/types';
import {HardhatRuntimeEnvironment} from 'hardhat/types';
import '@openzeppelin/hardhat-upgrades';

import {ContractNames} from "../src/contractNames";

const deployBridgeRegistry: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const cosmosBridge = await hre.deployments.get(ContractNames.CosmosBridge)
    const bridgeBank = await hre.deployments.get(ContractNames.BridgeBank)
    const {operator} = await hre.getNamedAccounts()

    const bridgeRegistry = hre.deployments.deploy(ContractNames.BridgeRegistry, {
        from: operator,
        args: [
            cosmosBridge.address,
            bridgeBank.address
        ]
    })
};

deployBridgeRegistry.tags = [ContractNames.BridgeRegistry]
deployBridgeRegistry.dependencies = [ContractNames.CosmosBridge, ContractNames.BridgeBank]
export default deployBridgeRegistry