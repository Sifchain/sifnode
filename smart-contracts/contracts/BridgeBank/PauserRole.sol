// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.4;

/**
 * @title Pauser Role
 * @dev Manages pausers
 */
contract PauserRole {
  /**
   * @notice List of pausers
   */
  mapping(address => bool) public pausers;
  mapping(address => bool) public unpausers;

  /**
   * @dev Modifier to restrict functions that can only be called by pausers
   */
  modifier onlyPauser() {
    require(pausers[msg.sender], "PauserRole: caller does not have the pauser role");
    _;
  }

  /**
   * @notice Adds `account` to the list of pausers
   * @param account The address of the new pauser
   */
  function addPauser(address account) public onlyPauser {
    _addPauser(account);
  }

  /**
   * @notice Removes `msg.sender` from the list of pausers
   */
  function renouncePauser() public {
    _removePauser(msg.sender);
  }

  /**
   * @dev Adds `account` to the list of pausers
   * @param account The address of the new pauser
   */
  function _addPauser(address account) internal {
    pausers[account] = true;
  }

  /**
   * @dev Removes `account` from the list of pausers
   * @param account The address of the pauser to be removed
   */
  function _removePauser(address account) internal {
    pausers[account] = false;
  }

  /**
   * @dev Modifier to restrict functions that can only be called by unpausers
   */
  modifier onlyUnpauser() {
      require(unpausers[msg.sender], "PauserRole: caller does not have the unpauser role");
      _;
  }

  /**
   * @notice Adds `account` to the list of unpausers
   * @param account The address of the new unpauser
   */
  function addUnpauser(address account) public onlyUnpauser {
      _addUnpauser(account);
  }

  /**
   * @notice Removes `msg.sender` from the list of unpausers
   */
  function renounceUnpauser() public {
      _removeUnpauser(msg.sender);
  }

  /**
   * @dev Adds `account` to the list of unpausers
   * @param account The address of the new unpauser
   */
  function _addUnpauser(address account) internal {
      unpausers[account] = true;
  }

  /**
   * @dev Removes `account` form the list of unpausers
   * @param account The address of the unpauser to be removed
   */
  function _removeUnpauser(address account) internal {
      unpausers[account] = false;
  }
}
