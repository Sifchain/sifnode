// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

import "./CosmosWhiteListStorage.sol";

/**
 * @title WhiteList
 * @dev WhiteList contract records the ERC 20 list that can be locked in BridgeBank.
 **/

contract CosmosWhiteList is CosmosWhiteListStorage {
    bool private _initialized;

    /*
     * @dev: Event declarations
     * @notice: This event is in EthereumWhitelist.sol
     */
    function _cosmosWhitelistInitialize() internal {
        require(!_initialized, "Initialized");
        _initialized = true;
    }

    /*
     * @dev: Modifier to restrict erc20 can be locked
     */
    modifier onlyCosmosTokenWhiteList(address _token) {
        require(
            getCosmosTokenInWhiteList(_token),
            "Only token in cosmos whitelist can be burned"
        );
        _;
    }

    /*
     * @dev: Modifier to restrict erc20 can be locked
     */
    modifier onlyTokenNotInCosmosWhiteList(address _token) {
        require(
            !getCosmosTokenInWhiteList(_token),
            "Only token not in whitelist can be locked"
        );
        _;
    }

    /*
     * @dev: Get if the token in whitelist
     *
     * @param _token: ERC 20's address
     * @return: if _token in whitelist
     */
    function getCosmosTokenInWhiteList(address _token) public view returns (bool) {
        return _cosmosTokenWhiteList[_token];
    }
}
