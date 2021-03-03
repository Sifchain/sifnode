pragma solidity >=0.4.22 <0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20Detailed.sol";

contract UsdCoin is ERC20, ERC20Detailed {
    constructor() public ERC20Detailed("UsdCoin", "usdc", 18) {
        _mint(msg.sender, 10000 * (10**uint256(decimals())));
    }
}
