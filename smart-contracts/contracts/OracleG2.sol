pragma solidity 0.8.0;

import "./Valset.sol";
import "./OracleStorage.sol";

contract Oracle is OracleStorage, Valset {
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
     * @dev: processProphecy
     *       Calculates the status of a prophecy. The claim is considered valid if the
     *       combined active signatory validator powers pass the consensus threshold.
     *       The threshold is x% of Total power, where x is the consensusThreshold param.
     */
    function getProphecyStatus(uint256 signedPower)
        public
        view
        returns (bool)
    {
        // Prophecy must reach total signed power % threshold in order to pass consensus
        uint256 prophecyPowerThreshold = totalPower * consensusThreshold;
        // consensusThreshold is a decimal multiplied by 100, so signedPower must also be multiplied by 100
        uint256 prophecyPowerCurrent = signedPower * 100;
        bool hasReachedThreshold = prophecyPowerCurrent >= prophecyPowerThreshold;

        return hasReachedThreshold;
    }
}
