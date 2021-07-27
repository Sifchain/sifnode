// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

/**
 * @title WhiteList
 * @dev WhiteList contract records the ERC 20 list that can be locked in BridgeBank.
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

    /*
     * @dev: Event declarations
     */
    event LogWhiteListUpdate(address _token, bool _value);

    function initialize() public {
        require(!_initialized, "Initialized");
        _ethereumTokenWhiteList[address(0)] = true;
        _initialized = true;
    }

    /*
     * @dev: Modifier to restrict erc20 can be locked
     */
    modifier onlyEthTokenWhiteList(address _token) {
        require(
            getTokenInEthWhiteList(_token),
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
    function setTokenInEthWhiteList(address _token, bool _inList)
        internal
        returns (bool)
    {
        _ethereumTokenWhiteList[_token] = _inList;
        emit LogWhiteListUpdate(_token, _inList);
        return _inList;
    }

    /*
     * @dev: Get if the token in whitelist
     *
     * @param _token: ERC 20's address
     * @return: if _token in whitelist
     */
    function getTokenInEthWhiteList(address _token) public view returns (bool) {
        return _ethereumTokenWhiteList[_token];
    }
}
