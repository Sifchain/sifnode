pragma solidity 0.5.16;

contract ValsetStorage {

    /**
     * @dev: Total power of all validators
     */
    uint256 public totalPower;

    /**
     * @dev: Current valset version
     */
    uint256 public currentValsetVersion;

    /**
     * @dev: validator count
     */
    uint256 public validatorCount;

    /**
     * @dev: Keep track of active validator
     */
    mapping(bytes32 => bool) public validators;

    /**
     * @dev: operator address
     */
    address public operator;

    /**
     * @dev: validator address + uint then hashed equals key mapped to powers
     */
    mapping(bytes32 => uint256) public powers;

    /**
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}