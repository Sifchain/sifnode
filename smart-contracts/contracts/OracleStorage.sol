// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

/**
 * @title Oracle Storage
 * @dev Stores prophecy-related information and the CosmosBridge address
 */
contract OracleStorage {
    /**
     * @notice Address of the Cosmos Bridge contract
     */
    address public cosmosBridge;

    /**
    * @dev {DEPRECATED}
    */
    address private _operator;

    /**
    * @notice Tracks the number of OracleClaims made on an individual BridgeClaim
    */
    uint256 public consensusThreshold; // e.g. 75 = 75%

    /**
    * @dev {DEPRECATED}
    */
    mapping(uint256 => uint256) private oracleClaimValidators;

    /**
    * @dev {DEPRECATED}
    */
    mapping(uint256 => mapping(address => bool)) private hasMadeClaim;

    /**
    * @notice mapping of validator address to last nonce submitted
    */
    uint256 public lastNonceSubmitted;

    /**
    * @dev gap of storage for future upgrades
    */
    uint256[99] private ____gap;
}
