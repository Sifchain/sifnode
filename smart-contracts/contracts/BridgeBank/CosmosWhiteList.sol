// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "./CosmosWhiteListStorage.sol";

/**
 * @title WhiteList
 * @dev WhiteList contract records the ERC 20 list that can be locked in BridgeBank.
 **/
contract CosmosWhiteList is CosmosWhiteListStorage {
  bool private _initialized;

  /**
   * @dev Initializer
   */
  function _cosmosWhitelistInitialize() internal {
    require(!_initialized, "Initialized");
    _initialized = true;
  }

  /**
   * @dev Modifier to restrict erc20 can be locked
   */
  modifier onlyCosmosTokenWhiteList(address _token) {
    require(getCosmosTokenInWhiteList(_token), "Token is not in Cosmos whitelist");
    _;
  }

  /**
   * @dev Modifier to restrict erc20 can be locked
   */
  modifier onlyTokenNotInCosmosWhiteList(address _token) {
    require(!getCosmosTokenInWhiteList(_token), "Only token not in cosmos whitelist can be locked");
    _;
  }

  /**
   * @notice Is `_token` in Cosmos Whitelist?
   * @dev Get if the token in whitelist
   * @param _token: ERC 20's address
   * @return if _token in whitelist
   */
  function getCosmosTokenInWhiteList(address _token) public view returns (bool) {
    return _cosmosTokenWhiteList[_token];
  }
}
