// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "./BridgeToken.sol";

/**
 * @title Rowan
 * @dev Mintable, ERC20Burnable, ERC20 compatible, Migration-enabled BankToken for use by BridgeBank
 **/
contract Rowan is BridgeToken {
  /**
   * @dev Instance of the old erowan contract
   */
  BridgeToken erowan;

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
    string memory _cosmosDenom,
    address _erowanAddress
  ) BridgeToken(_name, _symbol, _tokenDecimals, _cosmosDenom) {
    erowan = BridgeToken(_erowanAddress);
  }

  /**
   * @notice Migrates the user's balance from the old erowan to this contract
   * @notice Assumes 100% allowance has been given to this contract
   */
  function migrate() external {
    uint256 balance = erowan.balanceOf(msg.sender);
    erowan.burnFrom(msg.sender, balance);
    _mint(msg.sender, balance);

    emit MigrationComplete(msg.sender, balance);
  }
}
