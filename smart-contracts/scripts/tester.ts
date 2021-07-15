import {container as c} from "../src/tsyringe/sampleContainer"
import * as ethAddrs from "../src/ethereumAddress";
import {NativeCurrencyAddress} from "../src/ethereumAddress";

async function main() {
    const x = ethAddrs.toEthereumAddress("0xFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE")
    console.log(`x is ${x}`)
    const v = ethAddrs.toEthereumAddress("0x0000000000000000000000000000000000000000")
    console.log(`v is ${v}`)
    const z = ethAddrs.toEthereumAddress("0000000000000000000000000000000000000000")
    console.log(`z is ${typeof z}`)
}

main().then()