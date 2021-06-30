import {DeployFunction} from 'hardhat-deploy/types';
import {HardhatRuntimeEnvironment} from 'hardhat/types';
import '@openzeppelin/hardhat-upgrades';

import {ContractNames} from "../src/contractNames";
import {loadDeploymentEnvWithDotenv} from "../src/deploymentEnv";
import {CosmosBridge} from "../build";
import {getUnnamedAccounts, ethers} from 'hardhat';

const deployCosmosBridge: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const CosmosBridgeFactory = await hre.ethers.getContractFactory("CosmosBridge")
    const {operator} = await hre.getNamedAccounts()
    const {validator1} = await hre.getNamedAccounts()
    const deploymentEnv = loadDeploymentEnvWithDotenv()
    let cosmosBridgeArgs = [
        operator,
        deploymentEnv.consensusThreshold,
        [validator1],
        deploymentEnv.initialPowers
    ];
    await hre.deployments.deploy(ContractNames.CosmosBridge, {
        from: operator,
        proxy: true
    })
};

deployCosmosBridge.tags = [ContractNames.CosmosBridge]
deployCosmosBridge.id = `deploy_${ContractNames.CosmosBridge}`

export default deployCosmosBridge;