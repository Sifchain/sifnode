pragma solidity ^0.5.0;

import "../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./Valset.sol";
import "./CosmosBridge.sol";
import "./OracleStorage.sol";


contract Oracle is OracleStorage {
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
     * @dev: Modifier to restrict access to the operator.
     */
    modifier onlyOperator() {
        require(msg.sender == operator, "Must be the operator.");
        _;
    }

    /*
     * @dev: Modifier to restrict access to current ValSet validators
     */
    modifier onlyValidator(address _user) {
        require(
            Valset(valset).isActiveValidator(_user),
            "Must be an active validator"
        );
        _;
    }

    /*
     * @dev: Modifier to restrict access to current ValSet validators
     */
    modifier onlyCosmosBridge() {
        require(
            msg.sender == cosmosBridge,
            "Must be Cosmos Bridge"
        );
        _;
    }

    /*
     * @dev: Modifier to restrict access to current ValSet validators
     */
    modifier isPending(uint256 _prophecyID) {
        require(
            CosmosBridge(cosmosBridge).isProphecyClaimActive(_prophecyID) == true,
            "The prophecy must be pending for this operation"
        );
        _;
    }

    /*
     * @dev: Initialize Function
     */
    function initialize(
        address _operator,
        address _valset,
        address _cosmosBridge,
        uint256 _consensusThreshold
    ) public {
        require(!_initialized, "Initialized");
        require(
            _consensusThreshold > 0,
            "Consensus threshold must be positive."
        );
        operator = _operator;
        cosmosBridge = _cosmosBridge;
        valset = _valset;
        consensusThreshold = _consensusThreshold;
        _initialized = true;
    }

    /*
     * @dev: newOracleClaim
     *       Allows validators to make new OracleClaims on an existing Prophecy
     */
    function newOracleClaim(
        uint256 _prophecyID,
        address validatorAddress
    ) public
        onlyCosmosBridge
        onlyValidator(validatorAddress)
        isPending(_prophecyID)
        returns (bool)
    {
        // Confirm that this address has not already made an oracle claim on this prophecy
        require(
            !hasMadeClaim[_prophecyID][validatorAddress],
            "Cannot make duplicate oracle claims from the same address."
        );

        hasMadeClaim[_prophecyID][validatorAddress] = true;
        oracleClaimValidators[_prophecyID].push(validatorAddress);

        emit LogNewOracleClaim(
            _prophecyID,
            validatorAddress
        );

        // Process the prophecy
        (bool valid, , ) = getProphecyThreshold(_prophecyID);

        return valid;
    }

    /*
     * @dev: checkBridgeProphecy
     *       Operator accessor method which checks if a prophecy has passed
     *       the validity threshold, without actually completing the prophecy.
     */
    function checkBridgeProphecy(uint256 _prophecyID)
        public
        view
        onlyOperator
        isPending(_prophecyID)
        returns (bool, uint256, uint256)
    {
        require(
            CosmosBridge(cosmosBridge).isProphecyClaimActive(_prophecyID) == true,
            "Can only check active prophecies"
        );
        return getProphecyThreshold(_prophecyID);
    }

    /*
     * @dev: processProphecy
     *       Calculates the status of a prophecy. The claim is considered valid if the
     *       combined active signatory validator powers pass the consensus threshold.
     *       The threshold is x% of Total power, where x is the consensusThreshold param.
     */
    function getProphecyThreshold(uint256 _prophecyID)
        internal
        view
        returns (bool, uint256, uint256)
    {
        uint256 signedPower = 0;
        uint256 totalPower = Valset(valset).totalPower();

        // Iterate over the signatory addresses
        for (
            uint256 i = 0;
            i < oracleClaimValidators[_prophecyID].length;
            i = i.add(1)
        ) {
            address signer = oracleClaimValidators[_prophecyID][i];

            // Only add the power of active validators
            if (Valset(valset).isActiveValidator(signer)) {
                signedPower = signedPower.add(Valset(valset).getValidatorPower(signer));
            }
        }

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
