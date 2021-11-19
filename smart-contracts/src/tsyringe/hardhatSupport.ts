import { HardhatRuntimeEnvironment } from "hardhat/types"
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import { Wallet } from "ethers"
import * as ethers from "ethers"

export function isHardhatRuntimeEnvironment(x: any): x is HardhatRuntimeEnvironment {
  return "hardhatArguments" in x && "tasks" in x
}

export function createSignerWithAddresss(
  privateKey: string,
  provider: ethers.providers.JsonRpcProvider
): SignerWithAddress {
  return new Wallet(privateKey, provider) as unknown as SignerWithAddress
}
