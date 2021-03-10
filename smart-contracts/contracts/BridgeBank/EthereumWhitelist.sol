pragma solidity 0.6.9;

/**
 * @notice deprecated contract that only contains storage variables to keep offsets correct
 **/

contract EthereumWhiteList {
    bool private _initialized;

    /**
    * @notice mapping to keep track of whitelisted tokens
    */
    mapping(address => bool) private _ethereumTokenWhiteList;

    /**
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
    /*
     * @dev: Event declarations
     */
    event LogWhiteListUpdate(address _token, bool _value);
}