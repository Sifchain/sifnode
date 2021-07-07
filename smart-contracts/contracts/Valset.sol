pragma solidity 0.8.0;

import "./ValsetStorage.sol";

contract Valset is ValsetStorage {

    bool private _initialized;

    /*
     * @dev: Event declarations
     */
    event LogValidatorAdded(
        address _validator,
        uint256 _power,
        uint256 _currentValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    event LogValidatorPowerUpdated(
        address _validator,
        uint256 _power,
        uint256 _currentValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    event LogValidatorRemoved(
        address _validator,
        uint256 _power,
        uint256 _currentValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    event LogValsetReset(
        uint256 _newValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    event LogValsetUpdated(
        uint256 _newValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    /*
     * @dev: Modifier which restricts access to the operator.
     */
    modifier onlyOperator() {
        require(msg.sender == operator, "Must be the operator.");
        _;
    }

    /*
     * @dev: Constructor
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

        require(
            _initValidators.length == _initPowers.length,
            "Every validator must have a corresponding power"
        );

        resetValset();

        for (uint256 i = 0; i < _initValidators.length; i++) {
            addValidatorInternal(_initValidators[i], _initPowers[i]);
        }

        emit LogValsetUpdated(currentValsetVersion, validatorCount, totalPower);
    }

    /*
     * @dev: addValidator
     */
    function addValidator(address _validatorAddress, uint256 _validatorPower)
        public
        onlyOperator
    {
        addValidatorInternal(_validatorAddress, _validatorPower);
    }

    /*
     * @dev: updateValidatorPower
     */
    function updateValidatorPower(
        address _validatorAddress,
        uint256 _newValidatorPower
    ) public onlyOperator {

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

    /*
     * @dev: removeValidator
     */
    function removeValidator(address _validatorAddress) public onlyOperator {
        require(validators[_validatorAddress][currentValsetVersion], "Can only remove active validators");

        // Update validator count and total power

        // solidity compiler will handle and revert on over or underflows here
        // no need for safemath :)
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

    /*
     * @dev: updateValset
     */
    function updateValset(
        address[] memory _validators,
        uint256[] memory _powers
    ) public onlyOperator {
        require(
            _validators.length == _powers.length,
            "Every validator must have a corresponding power"
        );

        resetValset();

        for (uint256 i = 0; i < _validators.length; i++) {
            addValidatorInternal(_validators[i], _powers[i]);
        }

        emit LogValsetUpdated(currentValsetVersion, validatorCount, totalPower);
    }

    /*
     * @dev: isActiveValidator
     */
    function isActiveValidator(address _validatorAddress)
        public
        view
        returns (bool)
    {
        // Return bool indicating if this address is an active validator
        return validators[_validatorAddress][currentValsetVersion];
    }

    /*
     * @dev: getValidatorPower
     */
    function getValidatorPower(address _validatorAddress)
        public
        view
        returns (uint256)
    {
        return powers[_validatorAddress][currentValsetVersion];
    }

    /*
     * @dev: recoverGas
     */
    function recoverGas(uint256 _valsetVersion, address _validatorAddress)
        external
        onlyOperator
    {
        require(
            _valsetVersion < currentValsetVersion,
            "Gas recovery only allowed for previous validator sets"
        );
        // Delete from mappings and recover gas
        delete (validators[_validatorAddress][currentValsetVersion]);
        delete (powers[_validatorAddress][currentValsetVersion]);
    }

    /*
     * @dev: addValidatorInternal
     */
    function addValidatorInternal(
        address _validatorAddress,
        uint256 _validatorPower
    ) internal {
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

    /*
     * @dev: resetValset
     */
    function resetValset() internal {
        currentValsetVersion = currentValsetVersion + 1;
        validatorCount = 0;
        totalPower = 0;

        emit LogValsetReset(currentValsetVersion, validatorCount, totalPower);
    }
}
