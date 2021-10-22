// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

import "../interfaces/IBlocklist.sol";

/**
 * @title Ethereum WhiteList
 * @dev WhiteList contract records the ERC 20 list that can be locked in BridgeBank.
 */
contract EthereumWhiteList {
    /**
     * @dev has the contract been initialized?
     */
    bool private _initialized;

    /**
     * @dev [DEPRECATED] mapping to keep track of whitelisted tokens
     */
    mapping(address => bool) private _ethereumTokenWhiteList;

    /**
     * @dev the blocklist contract
     */
    IBlocklist public blocklist;

    /**
     * @dev is the blocklist active?
     */
    bool public hasBlocklist;

    /**
     * @dev gap of storage for future upgrades
     */
    uint256[98] private ____gap;

    /**
     * @notice Event emitted when the whitelist is updated
     */
    event LogWhiteListUpdate(address _token, bool _value);

    /**
     * @notice Initializer
     */
    function initialize() public {
        require(!_initialized, "Initialized");
        _ethereumTokenWhiteList[address(0)] = true;
        _initialized = true;
    }

    /**
     * @dev Modifier to restrict EVM addresses
     */
    modifier onlyNotBlocklisted(address account) {
        if (hasBlocklist) {
            require(
                !blocklist.isBlocklisted(account),
                "Address is blocklisted"
            );
        }
        _;
    }

    /**
     * @notice Lets the operator set the blocklist address
     * @param blocklistAddress The address of the blocklist contract
     */
    function _setBlocklist(address blocklistAddress) internal {
        blocklist = IBlocklist(blocklistAddress);
        hasBlocklist = true;
    }
}
