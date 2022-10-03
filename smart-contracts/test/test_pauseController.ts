import { ethers } from "hardhat";
import { use, expect } from "chai";
import { solidity } from "ethereum-waffle";
import { BridgeBank } from "../build";

use(solidity);

describe("Test Pause Controller", function () {
    let BridgeBank: BridgeBank;
    beforeEach(async () => {
        const bridgeBankFactory = ethers.
    })
})