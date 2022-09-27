// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "./../CosmosBridge.sol";
import "./ERC20UNSAFE.sol";

/// @notice Add a token to the cosmos bridge to test that upgrades work correctly
contract MockCosmosBridgeUpgrade is CosmosBridge, ERC20UNSAFE {
  function tokenFaucet() public {
    _mint(msg.sender, 100000000000);
  }
}
