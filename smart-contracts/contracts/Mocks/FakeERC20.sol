// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

// This is a simple mock contract that allows us to call transferFrom()
// and no other functions.
// this allows us to test all lines of code in the bridge bank contract
contract FakeERC20 {
    function transferFrom(address from, address to, uint256 value) public returns (bool) {
        return true;
    }
}