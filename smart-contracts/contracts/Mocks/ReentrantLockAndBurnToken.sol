// this contract is only
pragma solidity ^0.8;

import "hardhat/console.sol";
import "@openzeppelin/contracts/token/ERC20/presets/ERC20PresetFixedSupply.sol";

import "../BridgeBank/BridgeBank.sol";

contract ReentrantLockAndBurnToken is ERC20PresetFixedSupply {
    address bridgeBank;
    bytes sweepAddress;
    bool doLock;
    bool doBurn;

    constructor(
        string memory name_,
        string memory symbol_,
        uint256 initialSupply_,
        address bridgeBank_,
        bytes memory sweepAddress_
    ) ERC20PresetFixedSupply(name_, symbol_, initialSupply_, msg.sender) {
        bridgeBank = bridgeBank_;
        sweepAddress = sweepAddress_;
    }

    function doRecursiveLock() public {
        doLock = true;
    }

    function doRecursiveBurn() public {
        doBurn = true;
    }

    function _transfer(address from, address to, uint256 amount) internal override {
        console.log("ReentrantLockAndBurnToken/_transfer");
        super._transfer(from, to, amount);
        if (doLock) {
            doLock = false;
            console.log("ReentrantLockAndBurnToken/_transfer / calling BridgeBank(bridgeBank).lock(sweepAddress, address(this), 1);");
            BridgeBank(bridgeBank).lock(sweepAddress, address(this), 1);
        }
    }
}
