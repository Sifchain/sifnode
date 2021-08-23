// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

/**
 * @title Bridge Registry
 * @dev Stores the addresses of BridgeBank and CosmosBridge
 */
contract BridgeRegistry {
    /**
    * @notice Address of the CosmosBridge contract
    */
    address public cosmosBridge;

    /**
    * @notice Address of the BridgeBank contract
    */
    address public bridgeBank;

    // these variables are now deprecated and are made private
    // so that the getter helper method is not available.
    // [DEPRECATED]
    address private oracle;
    // [DEPRECATED]
    address private valset;

    bool private _initialized;

    uint256[100] private ____gap;

    event LogContractsRegistered(
        address _cosmosBridge,
        address _bridgeBank
    );

    /**
     * @notice Initializer
     * @param _cosmosBridge Address of the CosmosBridge contract
     * @param _bridgeBank Address of the BridgeBank contract
     */
    function initialize(
        address _cosmosBridge,
        address _bridgeBank
    ) public {
        require(!_initialized, "Initialized");

        cosmosBridge = _cosmosBridge;
        bridgeBank = _bridgeBank;
        _initialized = true;

        emit LogContractsRegistered(cosmosBridge, bridgeBank);
    }
}
