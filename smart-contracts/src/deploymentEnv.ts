export interface DeploymentEnv {
    consensusThreshold: number,
    initialPowers: Array<number>,
}

export function validate(obj: DeploymentEnv): DeploymentEnv {
    if (obj.initialPowers.length < 1)
        throw new Error('initialPowers is empty')
    if (obj.consensusThreshold <= 0)
        throw new Error('consensusThreshold must be positive')
    return obj
}

export function loadDeploymentEnvWithDotenv(): DeploymentEnv {
    const environment = require("dotenv").config();
    let result = {
        consensusThreshold: environment.CONSENSUS_THRESHOLD || 100,
        initialPowers: (environment.INITIAL_VALIDATOR_POWERS || "100").split(",")
    };
    return validate(result)
}