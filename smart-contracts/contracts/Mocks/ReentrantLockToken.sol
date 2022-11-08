// this contract is only
pragma solidity ^0.8;

import "hardhat/console.sol";
import "@openzeppelin/contracts/token/ERC20/presets/ERC20PresetFixedSupply.sol";

import "../BridgeBank/BridgeBank.sol";

contract ReentrantLockToken is ERC20PresetFixedSupply {
    address bridgeBank;
    bytes recursiveLockSifchainDestination; // BridgeBank.lock needs a destination, so we set it in the constructor
    address originalMsgSender;
    bool doLock;
    bool doBurn;

    constructor(
        string memory name_,
        string memory symbol_,
        uint256 initialSupply_,
        address bridgeBank_,
        bytes memory recursiveLockSifchainDestination_
    ) ERC20PresetFixedSupply(name_, symbol_, initialSupply_, msg.sender) {
        bridgeBank = bridgeBank_;
        recursiveLockSifchainDestination = recursiveLockSifchainDestination_;

        // We want BridgeBank to be able to use tokens from this contract itself
        // to do the recursive call, so approve that here
        _approve(address(this), bridgeBank_, initialSupply_);
    }

    function doRecursiveLock() public {
        doLock = true;
    }

    function doRecursiveBurn() public {
        doBurn = true;
    }

    function _transfer(address from, address to, uint256 amount) internal override {
        super._transfer(from, to, amount);
        if (doLock) {
            doLock = false;
            BridgeBank(bridgeBank).lock(recursiveLockSifchainDestination, address(this), 1);
        }
    }
}
