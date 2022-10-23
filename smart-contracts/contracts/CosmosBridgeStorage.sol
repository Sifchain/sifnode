// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

/**
 * @title Cosmos Bridge Storage
 * @dev Stores the operator's address,
        BridgeBank's address,
        networkDescriptor,
        cosmosDenomToDestinationAddress of a pegged token
 **/
contract CosmosBridgeStorage {
  /**
   * @dev {DEPRECATED}
   */
  string private COSMOS_NATIVE_ASSET_PREFIX;

  /**
   * @dev {DEPRECATED}
   */
  address private operator;

  /**
   * @dev {DEPRECATED}
   */
  address payable private valset;

  /**
   * @dev {DEPRECATED}
   */
  address payable private oracle;

  /**
   * @notice Address of the BridgeBank contract
   */
  address payable public bridgeBank;

  /**
   * @notice Has the BridgeBank contract been registered yet?
   */
  bool public hasBridgeBank;

  /**
   * @dev {DEPRECATED}
   */
  mapping(uint256 => ProphecyClaim) private prophecyClaims;

  /**
   * @dev {DEPRECATED}
   */
  enum Status {
    Null,
    Pending,
    Success,
    Failed
  }

  /**
   * @dev {DEPRECATED}
   */
  enum ClaimType {
    Unsupported,
    Burn,
    Lock
  }

  /**
   * @notice {DEPRECATED}
   */
  struct ProphecyClaim {
    address payable ethereumReceiver;
    string symbol;
    uint256 amount;
  }

  /**
   * @dev gap of storage for future upgrades
   */
  uint256[100] private ____gap;
}
