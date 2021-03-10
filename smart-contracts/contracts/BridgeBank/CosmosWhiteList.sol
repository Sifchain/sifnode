pragma solidity 0.6.9;

import "./CosmosWhiteListStorage.sol";

/**
 * @title WhiteList
 * @dev WhiteList contract records the ERC 20 list that can be locked in BridgeBank.
 **/

contract CosmosWhiteList is CosmosWhiteListStorage {
    bool private _initialized;

    /*
     * @dev: Event declarations
     */
    event LogWhiteListUpdate(address _token, bool _value);

    function _cosmosWhitelistInitialize() internal {
        require(!_initialized, "Initialized");
        _cosmosTokenWhiteList[address(0)] = true;
        _initialized = true;
    }

    /*
     * @dev: Modifier to restrict erc20 can be locked
     */
    modifier onlyCosmosTokenWhiteList(address _token) {
        require(
            getCosmosTokenInWhiteList(_token),
            "Only token in whitelist can be transferred to cosmos"
        );
        _;
    }

    /*
     * @dev: Set the token address in whitelist
     *
     * @param _token: ERC 20's address
     * @param _inList: set the _token in list or not
     * @return: new value of if _token in whitelist
     */
    function setTokenInCosmosWhiteList(address _token, bool _inList)
        internal virtual
        returns (bool)
    {
        _cosmosTokenWhiteList[_token] = _inList;
        emit LogWhiteListUpdate(_token, _inList);
        return _inList;
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