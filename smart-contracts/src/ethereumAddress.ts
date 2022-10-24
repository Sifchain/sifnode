import {ethers} from "ethers";
import {option} from "fp-ts"
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";

interface IHasAddress {
    address: string
}

export class NativeCurrencyAddress implements IHasAddress {
    constructor(readonly address: string) {
    }
}

export class NotNativeCurrencyAddress implements IHasAddress {
    constructor(readonly address: string) {
    }
}

export type EthereumAddress = NativeCurrencyAddress | NotNativeCurrencyAddress

export const eth = new NativeCurrencyAddress("0x0000000000000000000000000000000000000000")
const someEth = option.some(eth)

const nativeAddressRegex = /(0[xX])?0{40}/

function isNativeToken(address: string): boolean {
    return ethers.utils.isAddress(address) && nativeAddressRegex.test(address)
}

function toEthereumAddress(signerWithAddress: SignerWithAddress): option.Option<EthereumAddress>;
function toEthereumAddress(address: string): option.Option<EthereumAddress>;
function toEthereumAddress(address: SignerWithAddress | string): option.Option<EthereumAddress> {
    if (address instanceof SignerWithAddress) {
        return toEthereumAddress(address.address)
    } else {
        if (ethers.utils.isAddress(address)) {
            if (isNativeToken(address)) {
                return someEth
            } else {
                return option.some(new NotNativeCurrencyAddress(address))
            }
        } else {
            return option.none
        }
    }
}

export {toEthereumAddress}
