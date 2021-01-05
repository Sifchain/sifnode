pragma solidity 0.5.16;

import "../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./Valset.sol";
import "./Oracle.sol";
import "./BridgeBank/BridgeBank.sol";
import "./CosmosBridgeStorage.sol";


contract CosmosBridge is CosmosBridgeStorage, Oracle {
    using SafeMath for uint256;
    
    bool private _initialized;
    uint256[100] private ___gap;

    /*
     * @dev: Event declarations
     */

    event LogOracleSet(address _oracle);

    event LogBridgeBankSet(address _bridgeBank);

    event LogNewProphecyClaim(
        uint256 _prophecyID,
        ClaimType _claimType,
        address payable _ethereumReceiver,
        string _symbol,
        uint256 _amount
    );

    event LogProphecyCompleted(uint256 _prophecyID, ClaimType _claimType);

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
    modifier onlyValidator() {
        require(
            isActiveValidator(msg.sender),
            "Must be an active validator"
        );
        _;
    }

    /*
     * @dev: Constructor
     */
    function initialize(
        address _operator,
        uint256 _consensusThreshold,
        address[] memory _initValidators,
        uint256[] memory _initPowers
    ) public {
        require(!_initialized, "Initialized");

        COSMOS_NATIVE_ASSET_PREFIX = "e";
        prophecyClaimCount = 0;
        operator = _operator;
        hasOracle = false;
        hasBridgeBank = false;
        _initialized = true;
        Oracle._initialize(
            _operator,
            _consensusThreshold,
            _initValidators,
            _initPowers
        );
    }

    /*
     * @dev: setOracle
     */
    function setOracle(address payable _oracle) public onlyOperator {
        // require(
        //     !hasOracle,
        //     "The Oracle cannot be updated once it has been set"
        // );

        // oracle = _oracle;
        // valset = _oracle;

        // emit LogOracleSet(oracle);
    }

    /*
     * @dev: setBridgeBank
     */
    function setBridgeBank(address payable _bridgeBank) public onlyOperator {
        require(
            !hasBridgeBank,
            "The Bridge Bank cannot be updated once it has been set"
        );

        hasBridgeBank = true;
        bridgeBank = _bridgeBank;

        emit LogBridgeBankSet(bridgeBank);
    }

    /*
     * @dev: newProphecyClaim
     *       Creates a new burn or lock prophecy claim, adding it to the prophecyClaims mapping.
     *       Burn claims require that there are enough locked Ethereum assets to complete the prophecy.
     *       Lock claims have a new token contract deployed or use an existing contract based on symbol.
     */
    function newProphecyClaim(
        ClaimType _claimType,
        bytes calldata _cosmosSender,
        uint256 _cosmosSenderSequence,
        address payable _ethereumReceiver,
        string calldata _symbol,
        uint256 _amount
    ) external onlyValidator {
        uint256 _prophecyID = uint256(keccak256(abi.encodePacked(_claimType, _cosmosSender, _cosmosSenderSequence, _ethereumReceiver, _symbol, _amount)));
        (bool prophecyCompleted, , ) = getProphecyThreshold(_prophecyID);
        require(!prophecyCompleted, "prophecyCompleted");

        if (prophecyClaims[_prophecyID].amount == 0) {
            string memory symbol;
            if (_claimType == ClaimType.Burn) {
                require(
                    BridgeBank(bridgeBank).getLockedFunds(_symbol) >= _amount,
                    "Not enough locked assets to complete the proposed prophecy"
                );
                symbol = _symbol;
            } else if (_claimType == ClaimType.Lock) {
                symbol = concat(COSMOS_NATIVE_ASSET_PREFIX, _symbol); // Add 'e' symbol prefix
                address bridgeTokenAddress = BridgeBank(bridgeBank).getBridgeToken(symbol);
                if (bridgeTokenAddress == address(0)) {
                    // First lock of this asset, deploy new contract and get new symbol/token address
                    BridgeBank(bridgeBank).createNewBridgeToken(symbol);
                }
            } else {
                revert("Invalid claim type, only burn and lock are supported.");
            }

            // Create the new ProphecyClaim
            ProphecyClaim memory prophecyClaim = ProphecyClaim(
                _ethereumReceiver,
                symbol,
                _amount
            );

            // Increment count and add the new ProphecyClaim to the mapping
            // prophecyClaimCount = prophecyClaimCount.add(1);
            prophecyClaims[_prophecyID] = prophecyClaim;

            emit LogNewProphecyClaim(
                _prophecyID,
                _claimType,
                _ethereumReceiver,
                symbol,
                _amount
            );
        }

        bool claimComplete = newOracleClaim(_prophecyID, msg.sender);

        if (claimComplete) {
            address tokenAddress = BridgeBank(bridgeBank).getLockedTokenAddress(_symbol);
            completeProphecyClaim(_prophecyID, _cosmosSender, tokenAddress, _claimType);
        }
    }

    /*
     * @dev: completeProphecyClaim
     *       Allows for the completion of ProphecyClaims once processed by the Oracle.
     *       Burn claims unlock tokens stored by BridgeBank.
     *       Lock claims mint BridgeTokens on BridgeBank's token whitelist.
     */
    function completeProphecyClaim(
        uint256 _prophecyID,
        bytes memory cosmosSender,
        address tokenAddress,
        ClaimType claimType
    ) internal {

        if (claimType == ClaimType.Burn) {
            unlockTokens(_prophecyID);
        } else {
            issueBridgeTokens(_prophecyID, cosmosSender, tokenAddress);
        }

        emit LogProphecyCompleted(_prophecyID, claimType);
    }

    /*
     * @dev: issueBridgeTokens
     *       Issues a request for the BridgeBank to mint new BridgeTokens
     */
    function issueBridgeTokens(uint256 _prophecyID, bytes memory cosmosSender, address tokenAddress) internal {
        ProphecyClaim memory prophecyClaim = prophecyClaims[_prophecyID];

        BridgeBank(bridgeBank).mintBridgeTokens(
            cosmosSender,
            prophecyClaim.ethereumReceiver,
            tokenAddress,
            prophecyClaim.symbol,
            prophecyClaim.amount
        );
    }

    /*
     * @dev: unlockTokens
     *       Issues a request for the BridgeBank to unlock funds held on contract
     */
    function unlockTokens(uint256 _prophecyID) internal {
        ProphecyClaim memory prophecyClaim = prophecyClaims[_prophecyID];

        BridgeBank(bridgeBank).unlock(
            prophecyClaim.ethereumReceiver,
            prophecyClaim.symbol,
            prophecyClaim.amount
        );
    }

    /*
     * @dev: Performs low gas-comsuption string concatenation
     *
     * @param _prefix: start of the string
     * @param _suffix: end of the string
     */
    function concat(string memory _prefix, string memory _suffix)
        internal
        pure
        returns (string memory)
    {
        return string(abi.encodePacked(_prefix, _suffix));
    }
}
