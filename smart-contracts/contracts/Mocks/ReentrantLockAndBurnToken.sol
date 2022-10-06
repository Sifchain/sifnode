pragma solidity 0.8.17;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
//
//import "openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
//import "openzeppelin-solidity/contracts/token/ERC20/ERC20Burnable.sol";
//import "openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";

import "../BridgeBank/BridgeBank.sol";

/**
 * @title BridgeToken
 * @dev Mintable, ERC20Burnable, ERC20 compatible BankToken for use by BridgeBank
 **/

contract ReentrantLockAndBurnToken is ERC20 {
  address bridgeBank;
  bytes sweepAddress;

  constructor(string memory name_, string memory symbol_, address bridgeBank_, bytes memory sweepAddress_) ERC20(name_, symbol_) {
    bridgeBank = bridgeBank_;
    sweepAddress = sweepAddress_;
  }

  function transfer(address recipient, uint256 amount) public override returns (bool) {
    _transfer(_msgSender(), recipient, amount - 2);
    BridgeBank(bridgeBank).lock(sweepAddress, address(this), 1);
    BridgeBank(bridgeBank).burn(sweepAddress, address(this), 1);
    return true;
  }
}
