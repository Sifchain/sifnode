pragma solidity ^0.5.0;

contract WhiteListStorage {

    /**
    * @notice mapping to keep track of whitelisted tokens
    */
    mapping(address => bool) whiteList;

    /**
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}