// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

/**
 * @title Validator set Storage
 * @dev Stores information related to validators
 */
contract ValsetStorage {
  /**
   * @notice Total power of all validators
   */
  uint256 public totalPower;

  /**
   * @notice Current valset version
   */
  uint256 public currentValsetVersion;

  /**
   * @notice validator count
   */
  uint256 public validatorCount;

  /**
   * @notice Keep track of active validator
   */
  mapping(address => mapping(uint256 => bool)) public validators;

  /**
   * @notice operator address that can:
   *         Set BridgeBank's address (if it's not already set)
   *         Add new Validators, remove Validators, and update Validators' powers
   *         Call the function `recoverGas(uint256,address)`
   *         Change the operator
   */
  address public operator;

  /**
   * @notice validator address + uint then hashed equals key mapped to powers
   */
  mapping(address => mapping(uint256 => uint256)) public powers;

  /**
   * @dev gap of storage for future upgrades
   */
  uint256[100] private ____gap;
}
