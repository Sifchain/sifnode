pragma solidity 0.8.0;

import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title BridgeToken
 * @dev Mintable, ERC20Burnable, ERC20 compatible BankToken for use by BridgeBank
 **/

contract BridgeToken is ERC20Burnable, Ownable {
    constructor(string memory _symbol)
        public
        ERC20(_symbol, _symbol)
        Ownable()
    {
        // Intentionally left blank
    }

    function mint(address user, uint256 amount) public onlyOwner returns (bool) {
        _mint(user, amount);
        return true;
    }
}
