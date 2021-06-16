pragma solidity 0.8.0;

/**
 * @notice [DEPRECATED] all variables.
 * Contract that only contains storage variables to keep offsets correct
 **/

contract EthereumWhiteList {
    bool private _initialized;

    /*
    * @notice mapping to keep track of whitelisted tokens
    */
    mapping(address => bool) private _ethereumTokenWhiteList;

    /*
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}
