// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract TrollToken is ERC20 {
    constructor(string memory _name, string memory _symbol)
        ERC20(_name, _symbol) {}

    // transfer will never succeed. Need to ensure that this doesn't break unpegging this token
    function transfer(address recipient, uint256 amount) public override returns (bool) {
        // trolololololololololol
        for (uint256 i = 0; i < 1e30; i++) {}
    }

    function mint(address account, uint256 amount) public {
        _mint(account, amount);
    }
}
