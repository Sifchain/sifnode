// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title BridgeToken
 * @dev Mintable, ERC20Burnable, ERC20 compatible BankToken for use by BridgeBank
 **/

contract BridgeToken is ERC20Burnable, Ownable {

    uint8 private _decimals;

    constructor(string memory _name, string memory _symbol, uint8 _tokenDecimals)
        public
        ERC20(_name, _symbol)
        Ownable()
    {
       _decimals = _tokenDecimals;
    }

    function mint(address user, uint256 amount) public onlyOwner returns (bool) {
        _mint(user, amount);
        return true;
    }

    function decimals() public override view returns (uint8) {
        return _decimals;
    }
}
