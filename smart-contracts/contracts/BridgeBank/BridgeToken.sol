pragma solidity 0.8.0;

import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";

/**
 * @title BridgeToken
 * @dev Mintable, ERC20Burnable, ERC20 compatible BankToken for use by BridgeBank
 **/

contract BridgeToken is ERC20Burnable {
    constructor(string memory _symbol)
        public
        ERC20(_symbol, _symbol)
    {
        // Intentionally left blank
    }
}
