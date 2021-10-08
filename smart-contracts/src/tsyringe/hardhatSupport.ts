import {HardhatRuntimeEnvironment} from "hardhat/types";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import {Wallet} from "ethers";

export function isHardhatRuntimeEnvironment(x: any): x is HardhatRuntimeEnvironment {
    return 'hardhatArguments' in x && 'tasks' in x
}

export function createSignerWithAddresss(address: string, privateKey: string): SignerWithAddress {
    return (new Wallet(privateKey) as unknown) as SignerWithAddress
}
