// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

/**
 * @title Cosmos Whitelist Storage
 * @dev Records the Cosmos whitelisted tokens
 **/
contract CosmosWhiteListStorage {
  /**
   * @dev mapping to keep track of whitelisted tokens
   */
  mapping(address => bool) internal _cosmosTokenWhiteList;

  /**
   * @dev gap of storage for future upgrades
   */
  uint256[100] private ____gap;
}
