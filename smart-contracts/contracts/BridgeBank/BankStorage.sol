// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "./CosmosBankStorage.sol";
import "./EthereumBankStorage.sol";
import "./CosmosWhiteListStorage.sol";

/**
 * @title Bank Storage
 * @dev Stores addresses for owner, operator, and CosmosBridge
 **/
contract BankStorage is CosmosBankStorage, EthereumBankStorage, CosmosWhiteListStorage {
  /**
   * @notice Operator address that can:
   *         Reinitialize BridgeBank
   *         Update Eth whitelist
   *         Change the operator
   */
  address public operator;

  /**
   * @dev {DEPRECATED}
   */
  address private oracle;

  /**
   * @notice Address of the Cosmos Bridge smart contract
   */
  address public cosmosBridge;

  /**
   * @notice Owner address that can use the admin API
   */
  address public owner;

  /**
   * @dev {DEPRECATED}
   */
  mapping(string => uint256) private maxTokenAmount;

  /**
   * @dev gap of storage for future upgrades
   */
  uint256[100] private ____gap;
}
