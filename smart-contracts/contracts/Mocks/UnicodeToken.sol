// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract UnicodeToken is ERC20 {
    // Create a smart contract using unicode strings known to cause computers problems
  constructor() ERC20(unicode"لُلُصّبُلُلصّبُررً ॣ ॣh ॣ ॣ 冗", unicode"ܝܘܚܢܢ ܒܝܬ ܐܦܪܝܡ") {}

  function mint(address account, uint256 amount) public {
    _mint(account, amount);
  }
}
