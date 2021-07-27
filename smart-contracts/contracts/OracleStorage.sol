// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

contract OracleStorage {
    /*
     * @dev: Public variable declarations
     */
    address public cosmosBridge;

    /*
    * @notice [DEPRECATED] variable
    */
    address private _operator;

    /**
    * @notice Tracks the number of OracleClaims made on an individual BridgeClaim
    */
    uint256 public consensusThreshold; // e.g. 75 = 75%

    /**
    * @notice Tracks the number of OracleClaims made on an individual BridgeClaim
    */
    mapping(uint256 => uint256) public oracleClaimValidators;

    /**
    * @notice mapping of prophecyid to validator address to boolean
    */
    mapping(uint256 => mapping(address => bool)) public hasMadeClaim;

    /**
    * @notice mapping of prophecyid to check if it has been redeemed
    */
    mapping(uint256 => bool) public prophecyRedeemed;

    /**
    * @notice mapping of validator address to last nonce submitted
    */
    uint256 public lastNonceSubmitted;

    /*
    * @notice gap of storage for future upgrades
    */
    uint256[98] private ____gap;
}
