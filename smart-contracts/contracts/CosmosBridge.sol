pragma solidity =0.6.6;

import "@openzeppelin/contracts/math/SafeMath.sol";
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
        address _ethereumReceiver,
        uint256 _amount
    );

    event LogProphecyCompleted(uint256 _prophecyID, ClaimType _claimType);

    /*
     * @dev: Modifier to restrict access to current ValSet validators
     */
    modifier onlyValidator {
        require(
            isActiveValidator(msg.sender),
            "Must be an active validator"
        );
        _;
    }

    modifier validClaimType(ClaimType _claimType) {
        require(
            (_claimType == ClaimType.Lock || _claimType == ClaimType.Burn),
            "Invalid claim type"
        );
        _;
    }

    /*
     * @dev: Constructor
     */
    function initialize(
        address _operator,
        uint256 _consensusThreshold,
        address[] calldata _initValidators,
        uint256[] calldata _initPowers
    ) external {
        require(!_initialized, "Initialized");

        operator = _operator;
        hasBridgeBank = false;
        _initialized = true;
        Oracle._initialize(
            _operator,
            _consensusThreshold,
            _initValidators,
            _initPowers
        );
    }

    function changeOperator(address _newOperator) external onlyOperator {
        require(_newOperator != address(0), "invalid address");
        operator = _newOperator;
    }

    /*
     * @dev: setBridgeBank
     */
    function setBridgeBank(address payable _bridgeBank) external onlyOperator {
        require(
            !hasBridgeBank,
            "The Bridge Bank cannot be updated once it has been set"
        );

        hasBridgeBank = true;
        bridgeBank = _bridgeBank;

        emit LogBridgeBankSet(bridgeBank);
    }

    function getProphecyID(
        ClaimType _claimType,
        bytes memory _cosmosSender,
        uint256 _cosmosSenderSequence,
        address _ethereumReceiver,
        address _tokenAddress,
        uint256 _amount
    ) public pure returns (uint256) {
        return uint256(
            keccak256(
                abi.encodePacked(
                    _claimType,
                    _cosmosSender,
                    _cosmosSenderSequence,
                    _ethereumReceiver,
                    _tokenAddress,
                    _amount
                )
            )
        );
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
        string calldata _symbol,
        uint256 _cosmosSenderSequence,
        address _ethereumReceiver,
        address _tokenAddress,
        uint256 _amount
    ) external onlyValidator validClaimType(_claimType) {

        uint256 _prophecyID = getProphecyID(
            _claimType, 
            _cosmosSender,
            _cosmosSenderSequence,
            _ethereumReceiver,
            _tokenAddress,
            _amount
        );

        (bool prophecyCompleted, , ) = getProphecyThreshold(_prophecyID);
        require(!prophecyCompleted, "prophecyCompleted");

        if (oracleClaimValidators[_prophecyID] == 0) {
            if (_claimType == ClaimType.Lock && _tokenAddress == address(0)) {
                // need to make a business decision on what this symbol should be
                // First lock of this asset, deploy new contract and get new symbol/token address
                _tokenAddress = BridgeBank(bridgeBank).createNewBridgeToken(_symbol);
            }

            emit LogNewProphecyClaim(
                _prophecyID,
                _claimType,
                _ethereumReceiver,
                _amount
            );
        }
    
        bool claimComplete = newOracleClaim(_prophecyID, msg.sender);

        if (claimComplete) {
            completeProphecyClaim(
                _claimType,
                _prophecyID,
                _ethereumReceiver,
                _tokenAddress,
                _amount
            );
        }
    }

    /*
     * @dev: completeProphecyClaim
     *       Allows for the completion of ProphecyClaims once processed by the Oracle.
     *       Burn claims unlock tokens stored by BridgeBank.
     *       Lock claims mint BridgeTokens on BridgeBank's token whitelist.
     */
    function completeProphecyClaim(
        ClaimType _claimType,
        uint256 _prophecyID,
        address ethereumReceiver,
        address _tokenAddress,
        uint256 amount
    ) internal {
        BridgeBank(bridgeBank).handleUnpeg(
            ethereumReceiver,
            _tokenAddress,
            amount
        );

        emit LogProphecyCompleted(_prophecyID, _claimType);
    }
}
