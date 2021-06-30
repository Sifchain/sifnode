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

    function changeOperator(address _newOperator) public onlyOperator {
        require(_newOperator != address(0), "invalid address");
        operator = _newOperator;
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

    function getProphecyID(
        ClaimType _claimType,
        bytes calldata _cosmosSender,
        uint256 _cosmosSenderSequence,
        address payable _ethereumReceiver,
        string calldata _symbol,
        uint256 _amount
    ) external pure returns (uint256) {
        return uint256(keccak256(abi.encodePacked(_claimType, _cosmosSender, _cosmosSenderSequence, _ethereumReceiver, _symbol, _amount)));
    }

    /*
     * @dev: newProphecyClaim
     *       Creates a new burn or lock prophecy claim, adding it to the prophecyClaims mapping.
     *       Burn claims require that there are enough locked Ethereum assets to complete the prophecy.
     *       Lock claims have a new token contract deployed or use an existing contract based on symbol.
     */
    function newProphecyClaim(
        ClaimType _claimType,
        bytes memory _cosmosSender,
        uint256 _cosmosSenderSequence,
        address payable _ethereumReceiver,
        string memory _symbol,
        uint256 _amount
    ) public onlyValidator {
        uint256 _prophecyID = uint256(keccak256(abi.encodePacked(_claimType, _cosmosSender, _cosmosSenderSequence, _ethereumReceiver, _symbol, _amount)));
        (bool prophecyCompleted, , ) = getProphecyThreshold(_prophecyID);
        require(!prophecyCompleted, "prophecyCompleted");

        if (oracleClaimValidators[_prophecyID] == 0) {
            string memory symbol = BridgeBank(bridgeBank).safeLowerToUpperTokens(_symbol);

            if (_claimType == ClaimType.Burn) {
                address tokenAddress = BridgeBank(bridgeBank).getLockedTokenAddress(symbol);
                if (tokenAddress == address(0) && uint256(keccak256(abi.encodePacked(symbol))) != uint256(keccak256("eth"))) {
                    revert("Invalid token address");
                }
            } else if (_claimType == ClaimType.Lock) {
                address bridgeTokenAddress = BridgeBank(bridgeBank).getBridgeToken(symbol);
                if (bridgeTokenAddress == address(0)) {
                    // First lock of this asset, deploy new contract and get new symbol/token address
                    BridgeBank(bridgeBank).createNewBridgeToken(symbol);
                }
            } else {
                revert("Invalid claim type, only burn and lock are supported.");
            }

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
            completeProphecyClaim(
                _prophecyID,
                _claimType,
                _ethereumReceiver,
                _symbol,
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
        uint256 _prophecyID,
        ClaimType claimType,
        address payable ethereumReceiver,
        string memory symbol,
        uint256 amount
    ) internal {

        if (claimType == ClaimType.Burn) {
            unlockTokens(ethereumReceiver, symbol, amount);
        } else {
            issueBridgeTokens(ethereumReceiver, symbol, amount);
        }

        emit LogProphecyCompleted(_prophecyID, claimType);
    }

    /*
     * @dev: issueBridgeTokens
     *       Issues a request for the BridgeBank to mint new BridgeTokens
     */
    function issueBridgeTokens(
        address payable ethereumReceiver,
        string memory symbol,
        uint256 amount
    ) internal {
        BridgeBank(bridgeBank).mintBridgeTokens(
            ethereumReceiver,
            symbol,
            amount
        );
    }

    /*
     * @dev: unlockTokens
     *       Issues a request for the BridgeBank to unlock funds held on contract
     */
    function unlockTokens(
        address payable ethereumReceiver,
        string memory symbol,
        uint256 amount
    ) internal {
        BridgeBank(bridgeBank).unlock(
            ethereumReceiver,
            symbol,
            amount
        );
    }
}
