pragma solidity 0.5.16;

contract OracleStorage {
    /*
     * @dev: Public variable declarations
     */
    address public cosmosBridge;

    /**
    * @notice Tracks the number of OracleClaims made on an individual BridgeClaim
    */
    address public valset;

    /**
    * @notice Tracks the number of OracleClaims made on an individual BridgeClaim
    */
    address public operator;

    /**
    * @notice Tracks the number of OracleClaims made on an individual BridgeClaim
    */
    uint256 public consensusThreshold; // e.g. 75 = 75%

    /**
    * @notice Tracks the number of OracleClaims made on an individual BridgeClaim
    */
    mapping(uint256 => address[]) public oracleClaimValidators;

    /**
    * @notice mapping of prophecyid to validator address to boolean
    */
    mapping(uint256 => mapping(address => bool)) public hasMadeClaim;

    /**
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}