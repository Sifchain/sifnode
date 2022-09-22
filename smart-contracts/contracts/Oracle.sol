// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "./Valset.sol";
import "./OracleStorage.sol";

/**
 * @title Oracle
 * @dev Calculates a prophecy status
 */
contract Oracle is OracleStorage, Valset {
  /**
   * @dev has the contract been initialized?
   */
  bool private _initialized;

  /**
   * @dev {DEPRECATED}
   */
  event LogNewOracleClaim(uint256 _prophecyID, address _validatorAddress);

  /**
   * @dev {DEPRECATED}
   */
  event LogProphecyProcessed(
    uint256 _prophecyID,
    uint256 _prophecyPowerCurrent,
    uint256 _prophecyPowerThreshold,
    address _submitter
  );

  /**
   * @dev Initializer
   * @param _operator Address of the operator
   * @param _consensusThreshold Minimum required power for a valid prophecy
   * @param _initValidators List of initial validators
   * @param _initPowers List of numbers representing the power of each validator in the above list
   */
  function _initialize(
    address _operator,
    uint256 _consensusThreshold,
    address[] memory _initValidators,
    uint256[] memory _initPowers
  ) internal {
    require(!_initialized, "Initialized");
    require(_consensusThreshold > 0, "Consensus threshold must be positive.");
    require(_consensusThreshold <= 100, "Invalid consensus threshold.");
    operator = _operator;
    consensusThreshold = _consensusThreshold;
    _initialized = true;
    Valset._initialize(_operator, _initValidators, _initPowers);
  }

  /**
   * @dev Calculates the status of a prophecy. The claim is considered valid if the
   *      combined active signatory validator powers pass the consensus threshold.
   *      The threshold is x% of Total power, where x is the consensusThreshold param.
   * @param signedPower aggregated power of signers signing the prophecy
   * @return Boolean: has this prophecy reached the threshold?
   */
  function getProphecyStatus(uint256 signedPower) public view returns (bool) {
    // Prophecy must reach total signed power % threshold in order to pass consensus
    uint256 prophecyPowerThreshold = totalPower * consensusThreshold;
    // consensusThreshold is a decimal multiplied by 100, so signedPower must also be multiplied by 100
    uint256 prophecyPowerCurrent = signedPower * 100;
    bool hasReachedThreshold = prophecyPowerCurrent >= prophecyPowerThreshold;

    return hasReachedThreshold;
  }

  function updateConsensusThreshold(uint256 _consensusThreshold) public onlyOperator {
    require(_consensusThreshold > 0, "Consensus threshold must be positive.");
    require(_consensusThreshold <= 100, "Invalid consensus threshold.");
    consensusThreshold = _consensusThreshold;
  }
}
