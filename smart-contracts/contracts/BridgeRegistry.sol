// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

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

  /**
   * @dev {DEPRECATED}
   */
  address private oracle;

  /**
   * @dev {DEPRECATED}
   */
  address private valset;

  /**
   * @dev has this contract been initialized?
   */
  bool private _initialized;

  /**
   * @dev gap of storage for future upgrades
   */
  uint256[100] private ____gap;

  /**
   * @dev Event emitted when this contract is initialized
   */
  event LogContractsRegistered(address _cosmosBridge, address _bridgeBank);

  /**
   * @notice Initializer
   * @param _cosmosBridge Address of the CosmosBridge contract
   * @param _bridgeBank Address of the BridgeBank contract
   */
  function initialize(address _cosmosBridge, address _bridgeBank) public {
    require(!_initialized, "Initialized");

    cosmosBridge = _cosmosBridge;
    bridgeBank = _bridgeBank;
    _initialized = true;

    emit LogContractsRegistered(cosmosBridge, bridgeBank);
  }
}
