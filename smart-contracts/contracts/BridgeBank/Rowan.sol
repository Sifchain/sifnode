// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

import "./BridgeToken.sol";

/**
 * @title Rowan
 * @dev Mintable, ERC20Burnable, ERC20 compatible, Migration-enabled BankToken for use by BridgeBank
 **/
contract Rowan is BridgeToken {
  /**
   * @notice Address of the old erowan contract
   */
  address public constant erowanAddress = 0x07baC35846e5eD502aA91AdF6A9e7aA210F2DcbE;

  /**
   * @dev Instance of the old erowan contract 
   */
  BridgeToken erowan = BridgeToken(erowanAddress);

  /**
   * @notice Event emitted when a user migrates their balance
   *         from the old erowan contract to this contract
   * @param account Address of the user who migrated their balance
   * @param amount How many tokens have been migrated
   */
  event MigrationComplete(address indexed account, uint256 amount);
  
  constructor(
      string memory _name,
      string memory _symbol,
      uint8 _tokenDecimals,
      string memory _cosmosDenom
  ) BridgeToken(_name, _symbol, _tokenDecimals, _cosmosDenom) {

  }

  /**
   * @notice Migrates the user's balance from the old erowan to this contract
   * @notice Assumes 100% allowance has been given to this contract
   */
  function migrate() external {
    uint balance = erowan.balanceOf(msg.sender);
    erowan.burnFrom(msg.sender, balance);
    _mint(msg.sender, balance);

    emit MigrationComplete(msg.sender, balance);
  }
}
