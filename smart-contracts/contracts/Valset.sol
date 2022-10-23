// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "./ValsetStorage.sol";

/**
 * @title Validator set Storage
 * @dev Manages validators
 */
contract Valset is ValsetStorage {
  /**
   * @dev has the contract been initialized?
   */
  bool private _initialized;

  /**
   * @dev Event emitted when a new validator is added to the list
   */
  event LogValidatorAdded(
    address _validator,
    uint256 _power,
    uint256 _currentValsetVersion,
    uint256 _validatorCount,
    uint256 _totalPower
  );

  /**
   * @dev Event emitted when the power of a validator has been updated
   */
  event LogValidatorPowerUpdated(
    address _validator,
    uint256 _power,
    uint256 _currentValsetVersion,
    uint256 _validatorCount,
    uint256 _totalPower
  );

  /**
   * @dev Event emitted when a validator is removed from the list
   */
  event LogValidatorRemoved(
    address _validator,
    uint256 _power,
    uint256 _currentValsetVersion,
    uint256 _validatorCount,
    uint256 _totalPower
  );

  /**
   * @dev Event emitted when values have been reset
   */
  event LogValsetReset(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower);

  /**
   * @dev Event emitted when values have been updated
   */
  event LogValsetUpdated(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower);

  /**
   * @dev Modifier which restricts access to the operator.
   */
  modifier onlyOperator() {
    require(msg.sender == operator, "Must be the operator.");
    _;
  }

  /**
   * @dev Initializer
   */
  function _initialize(
    address _operator,
    address[] memory _initValidators,
    uint256[] memory _initPowers
  ) internal {
    require(!_initialized, "Initialized");

    operator = _operator;
    currentValsetVersion = 0;
    _initialized = true;
    
    uint256 initValLength = _initValidators.length;

    require(
      initValLength == _initPowers.length,
      "Every validator must have a corresponding power"
    );

    resetValset();

    for (uint256 i; i < initValLength;) {
      _addValidatorInternal(_initValidators[i], _initPowers[i]);
      unchecked { ++i; }
    }

    emit LogValsetUpdated(currentValsetVersion, validatorCount, totalPower);
  }

  /**
   * @notice Adds `_validatorAddress` to the list with `_validatorPower` power
   * @dev Can only be called by the operator
   * @param _validatorAddress Address of the new validator
   * @param _validatorPower The power this validator has
   */
  function addValidator(address _validatorAddress, uint256 _validatorPower) external onlyOperator {
    _addValidatorInternal(_validatorAddress, _validatorPower);
  }

  /**
   * @notice Updates the power of validator `_validatorAddress` to `_newValidatorPower`
   * @dev Can only be called by the operator
   * @param _validatorAddress Address of the validator
   * @param _newValidatorPower The power this validator has
   */
  function updateValidatorPower(address _validatorAddress, uint256 _newValidatorPower)
    external
    onlyOperator
  {
    require(
      validators[_validatorAddress][currentValsetVersion],
      "Can only update the power of active valdiators"
    );

    // Adjust total power by new validator power
    uint256 priorPower = powers[_validatorAddress][currentValsetVersion];
    // solidity compiler will handle and revert on over or underflows here
    // no need for safemath :)
    totalPower = totalPower - priorPower;
    totalPower = totalPower + _newValidatorPower;

    // Set validator's new power
    powers[_validatorAddress][currentValsetVersion] = _newValidatorPower;

    emit LogValidatorPowerUpdated(
      _validatorAddress,
      _newValidatorPower,
      currentValsetVersion,
      validatorCount,
      totalPower
    );
  }

  /**
   * @notice Removes validator `_validatorAddress` from the list
   * @dev Can only be called by the operator
   * @param _validatorAddress Address of the validator
   */
  function removeValidator(address _validatorAddress) external onlyOperator {
    require(
      validators[_validatorAddress][currentValsetVersion],
      "Can only remove active validators"
    );

    // Update validator count and total power
    validatorCount = validatorCount - 1;
    totalPower = totalPower - powers[_validatorAddress][currentValsetVersion];

    // Delete validator and power
    delete validators[_validatorAddress][currentValsetVersion];
    delete powers[_validatorAddress][currentValsetVersion];

    emit LogValidatorRemoved(
      _validatorAddress,
      0,
      currentValsetVersion,
      validatorCount,
      totalPower
    );
  }

  /**
   * @notice Replaces the list of validators with `_validators`, each with `_powers` power
   * @dev Can only be called by the operator; lists must have the same length
   * @param _validators List of validator addresses
   * @param _powers List of validator powers
   */
  function updateValset(address[] memory _validators, uint256[] memory _powers)
    external
    onlyOperator
  {
    uint256 valLength = _validators.length;
    require(
      valLength == _powers.length,
      "Every validator must have a corresponding power"
    );

    resetValset();

    for (uint256 i; i < valLength;) {
      _addValidatorInternal(_validators[i], _powers[i]);
      unchecked{ ++i; }
    }

    emit LogValsetUpdated(currentValsetVersion, validatorCount, totalPower);
  }

  /**
   * @notice Consults whether `_validatorAddress` is an active validator or not
   * @param _validatorAddress Address of the validator
   * @return Boolean: is it an active validator?
   */
  function isActiveValidator(address _validatorAddress) public view returns (bool) {
    // Return bool indicating if this address is an active validator
    return validators[_validatorAddress][currentValsetVersion];
  }

  /**
   * @notice Consults how much validation power `_validatorAddress` has
   * @param _validatorAddress Address of the validator
   * @return The validator's power
   */
  function getValidatorPower(address _validatorAddress) public view returns (uint256) {
    return powers[_validatorAddress][currentValsetVersion];
  }

  /**
   * @notice Deletes an old validator, recovering some gas in the process
   * @dev Can only be part of an execution flow started by the operator
   * @param _valsetVersion Address of the validator
   * @param _validatorAddress Address of the validator
   */
  function recoverGas(uint256 _valsetVersion, address _validatorAddress) external onlyOperator {
    require(
      _valsetVersion < currentValsetVersion,
      "Gas recovery only allowed for previous validator sets"
    );
    // Delete from mappings and recover gas
    delete (validators[_validatorAddress][currentValsetVersion]);
    delete (powers[_validatorAddress][currentValsetVersion]);
  }

  /**
   * @dev Adds a new validator to the list
   * @param _validatorAddress Address of the validator
   * @param _validatorPower The power this validator has
   */
  function _addValidatorInternal(address _validatorAddress, uint256 _validatorPower) internal {
    require(validators[_validatorAddress][currentValsetVersion] == false, "Already a validator");
    
    validatorCount = validatorCount + 1;
    totalPower = totalPower + _validatorPower;

    // Set validator as active and set their power
    validators[_validatorAddress][currentValsetVersion] = true;
    powers[_validatorAddress][currentValsetVersion] = _validatorPower;

    emit LogValidatorAdded(
      _validatorAddress,
      _validatorPower,
      currentValsetVersion,
      validatorCount,
      totalPower
    );
  }

  /**
   * @dev Resets variables and bumps currentValsetVersion
   */
  function resetValset() internal {
    currentValsetVersion = currentValsetVersion + 1;
    validatorCount = 0;
    totalPower = 0;

    emit LogValsetReset(currentValsetVersion, validatorCount, totalPower);
  }
}
