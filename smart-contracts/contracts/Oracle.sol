pragma solidity 0.6.6;

import "@openzeppelin/contracts/math/SafeMath.sol";
import "./Valset.sol";
import "./OracleStorage.sol";

contract Oracle is OracleStorage, Valset {
    using SafeMath for uint256;

    bool private _initialized;

    /*
     * @dev: Event declarations
     */
    event LogNewOracleClaim(
        uint256 _prophecyID,
        address _validatorAddress
    );

    event LogProphecyProcessed(
        uint256 _prophecyID,
        uint256 _prophecyPowerCurrent,
        uint256 _prophecyPowerThreshold,
        address _submitter
    );


    /*
     * @dev: Initialize Function
     */
    function _initialize(
        address _operator,
        uint256 _consensusThreshold,
        address[] memory _initValidators,
        uint256[] memory _initPowers
    ) internal {
        require(!_initialized, "Initialized");
        require(
            _consensusThreshold > 0,
            "Consensus threshold must be positive."
        );
        require(
            _consensusThreshold <= 100,
            "Invalid consensus threshold."
        );
        operator = _operator;
        consensusThreshold = _consensusThreshold;
        _initialized = true;
        Valset._initialize(_operator, _initValidators, _initPowers);
    }

    /*
     * @dev: newOracleClaim
     *       Allows validators to make new OracleClaims on an existing Prophecy
     */
    function newOracleClaim(
        uint256 _prophecyID,
        address validatorAddress
    ) internal
        returns (bool)
    {
        // Confirm that this address has not already made an oracle claim on this prophecy
        require(
            !hasMadeClaim[_prophecyID][validatorAddress],
            "Cannot make duplicate oracle claims from the same address."
        );

        hasMadeClaim[_prophecyID][validatorAddress] = true;
        oracleClaimValidators[_prophecyID] = oracleClaimValidators[_prophecyID].add(
            getValidatorPower(validatorAddress)
        );

        emit LogNewOracleClaim(
            _prophecyID,
            validatorAddress
        );

        // Process the prophecy
        (bool valid, , ) = getProphecyThreshold(_prophecyID);

        return valid;
    }

    /*
     * @dev: processProphecy
     *       Calculates the status of a prophecy. The claim is considered valid if the
     *       combined active signatory validator powers pass the consensus threshold.
     *       The threshold is x% of Total power, where x is the consensusThreshold param.
     */
    function getProphecyThreshold(uint256 _prophecyID)
        public
        view
        returns (bool, uint256, uint256)
    {
        uint256 signedPower = 0;
        uint256 totalPower = totalPower;

        signedPower = oracleClaimValidators[_prophecyID];

        // Prophecy must reach total signed power % threshold in order to pass consensus
        uint256 prophecyPowerThreshold = totalPower.mul(consensusThreshold);
        // consensusThreshold is a decimal multiplied by 100, so signedPower must also be multiplied by 100
        uint256 prophecyPowerCurrent = signedPower.mul(100);
        bool hasReachedThreshold = prophecyPowerCurrent >=
            prophecyPowerThreshold;

        return (
            hasReachedThreshold,
            prophecyPowerCurrent,
            prophecyPowerThreshold
        );
    }
}
