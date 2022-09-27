// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "./PauserRole.sol";

/**
 * @title Pausable
 * @dev Contract module which allows children to implement an emergency stop
 * mechanism that can be triggered by an authorized account.
 *
 * This module is used through inheritance. It will make available the
 * modifiers `whenNotPaused` and `whenPaused`, which can be applied to
 * the functions of your contract. Note that they will not be pausable by
 * simply including this module, only once the modifiers are put in place.
 */
contract Pausable is PauserRole {
  /**
   * @dev Emitted when the pause is triggered by a pauser (`account`).
   */
  event Paused(address account);

  /**
   * @dev Emitted when the pause is lifted by a pauser (`account`).
   */
  event Unpaused(address account);

  bool private _paused;

  /**
   * @dev Initializes adding a new Pauser
   */
  function _pausableInitialize(address _user) internal {
    _addPauser(_user);
    _paused = false;
  }

  /**
   * @notice Is the contract paused?
   * @dev Returns true if the contract is paused, and false otherwise.
   */
  function paused() public view returns (bool) {
    return _paused;
  }

  /**
   * @dev Modifier to make a function callable only when the contract is not paused.
   */
  modifier whenNotPaused() {
    require(!_paused, "Pausable: paused");
    _;
  }

  /**
   * @dev Modifier to make a function callable only when the contract is paused.
   */
  modifier whenPaused() {
    require(_paused, "Pausable: not paused");
    _;
  }

  /**
   * @dev Called by a owner to toggle pause
   */
  function togglePause() private {
    _paused = !_paused;
  }

  /**
   * @notice Pauses the contract
   * @dev Called by a pauser to pause contract
   */
  function pause() external onlyPauser whenNotPaused {
    togglePause();
    emit Paused(msg.sender);
  }

  /**
   * @notice Unpauses the contract
   * @dev Called by a pauser to unpause contract
   */
  function unpause() external onlyPauser whenPaused {
    togglePause();
    emit Unpaused(msg.sender);
  }
}
