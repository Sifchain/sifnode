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
        address _tokenAddress,
        string _symbol,
        uint256 _amount
    );

    event LogProphecyCompleted(uint256 _prophecyID, ClaimType _claimType);

    /*
     * @dev: Modifier which only allows access to currently pending prophecies
     */
    modifier isPending(uint256 _prophecyID) {
        require(
            isProphecyClaimActive(_prophecyID),
            "Prophecy claim is not active"
        );
        _;
    }

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
        bool claimComplete;
        uint256 _prophecyID = uint256(keccak256(abi.encodePacked(_claimType, _cosmosSender, _cosmosSenderSequence, _ethereumReceiver, _symbol, _amount)));
        if (usedNonce[_prophecyID]) {
            claimComplete = newOracleClaim(_prophecyID, msg.sender);
        } else {
            address tokenAddress;
            string memory symbol;
            if (_claimType == ClaimType.Burn) {
                require(
                    BridgeBank(bridgeBank).getLockedFunds(_symbol) >= _amount,
                    "Not enough locked assets to complete the proposed prophecy"
                );
                symbol = _symbol;
                tokenAddress = BridgeBank(bridgeBank).getLockedTokenAddress(_symbol);
            } else if (_claimType == ClaimType.Lock) {
                symbol = concat(COSMOS_NATIVE_ASSET_PREFIX, _symbol); // Add 'e' symbol prefix
                address bridgeTokenAddress = BridgeBank(bridgeBank).getBridgeToken(symbol);
                if (bridgeTokenAddress == address(0)) {
                    // First lock of this asset, deploy new contract and get new symbol/token address
                    tokenAddress = BridgeBank(bridgeBank).createNewBridgeToken(symbol);
                } else {
                    // Not the first lock of this asset, get existing symbol/token address
                    tokenAddress = bridgeTokenAddress;
                }
            } else {
                revert("Invalid claim type, only burn and lock are supported.");
            }

            // Create the new ProphecyClaim
            ProphecyClaim memory prophecyClaim = ProphecyClaim(
                _claimType,
                _ethereumReceiver,
                tokenAddress,
                symbol,
                _amount,
                Status.Pending
            );

            // Increment count and add the new ProphecyClaim to the mapping
            prophecyClaimCount = prophecyClaimCount.add(1);
            prophecyClaims[_prophecyID] = prophecyClaim;

            emit LogNewProphecyClaim(
                _prophecyID,
                _claimType,
                _ethereumReceiver,
                tokenAddress,
                symbol,
                _amount
            );

            usedNonce[_prophecyID] = true;
            claimComplete = newOracleClaim(_prophecyID, msg.sender);
        }

        if (claimComplete) {
            completeProphecyClaim(_prophecyID, _cosmosSender);
        }
    }

    /*
     * @dev: completeProphecyClaim
     *       Allows for the completion of ProphecyClaims once processed by the Oracle.
     *       Burn claims unlock tokens stored by BridgeBank.
     *       Lock claims mint BridgeTokens on BridgeBank's token whitelist.
     */
    function completeProphecyClaim(uint256 _prophecyID, bytes memory cosmosSender)
        internal
        isPending(_prophecyID)
    {
        prophecyClaims[_prophecyID].status = Status.Success;

        ClaimType claimType = prophecyClaims[_prophecyID].claimType;
        if (claimType == ClaimType.Burn) {
            unlockTokens(_prophecyID);
        } else {
            issueBridgeTokens(_prophecyID, cosmosSender);
        }

        emit LogProphecyCompleted(_prophecyID, claimType);
    }

    /*
     * @dev: issueBridgeTokens
     *       Issues a request for the BridgeBank to mint new BridgeTokens
     */
    function issueBridgeTokens(uint256 _prophecyID, bytes memory cosmosSender) internal {
        ProphecyClaim memory prophecyClaim = prophecyClaims[_prophecyID];

        BridgeBank(bridgeBank).mintBridgeTokens(
            cosmosSender,
            prophecyClaim.ethereumReceiver,
            prophecyClaim.tokenAddress,
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
     * @dev: isProphecyClaimActive
     *       Returns boolean indicating if the ProphecyClaim is active
     */
    function isProphecyClaimActive(uint256 _prophecyID)
        public
        view
        returns (bool)
    {
        return prophecyClaims[_prophecyID].status == Status.Pending;
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
