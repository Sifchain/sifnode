import {HardhatRuntimeEnvironment} from "hardhat/types";

export function isHardhatRuntimeEnvironment(x: any): x is HardhatRuntimeEnvironment {
    return 'hardhatArguments' in x && 'tasks' in x
}
