pragma solidity 0.6.6;

contract CosmosWhiteListStorage {

    /**
    * @notice mapping to keep track of whitelisted tokens
    */
    mapping(address => bool) internal _cosmosTokenWhiteList;

    /**
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}