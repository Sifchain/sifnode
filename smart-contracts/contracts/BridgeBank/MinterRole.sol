// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

/**
 * @title Minter Role
 * @dev Manages a list of Minters. Minters can mint new tokens.
 *      Only the owner of this contract has permission to manage Minters.
 */
contract MinterRole {
    /**
     * @dev List of accounts that have the Minter role
     */
    mapping (address => bool) public minters;

    /**
     * @dev Event emitted when the list of Minters is updated
     */
    event MinterUpdate(address indexed account, bool isMinter);

    /**
     * @dev Modifier to restrict access to Minters
     */
    modifier onlyMinter() {
        require(minters[msg.sender], "MinterRole: caller does not have the Minter role");
        _;
    }

    /**
     * @dev Adds `account` to the list of Minters
     * @param account The address of the new Minter
     */
    function _addMinter(address account) internal {
        minters[account] = true;

        emit MinterUpdate(account, true);
    }

    /**
     * @dev Removes `account` from the list of Minters
     * @param account The address of the Minter to be removed
     */
    function _removeMinter(address account) internal {
        minters[account] = false;

        emit MinterUpdate(account, false);
    }
}
