pragma solidity ^0.5.0;

import "../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./Valset.sol";
import "./BridgeBank/BridgeBank.sol";
import "./CosmosBridgeStorage.sol";


contract CosmosBridge is CosmosBridgeStorage {
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
        bytes _cosmosSender,
        address payable _ethereumReceiver,
        address _validatorAddress,
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
     * @dev: The bridge is not active until oracle and bridge bank are set
     */
    modifier isActive() {
        require(
            hasOracle == true && hasBridgeBank == true,
            "The Operator must set the oracle and bridge bank for bridge activation"
        );
        _;
    }

    /*
     * @dev: Constructor
     */
    function initialize(address _operator, address payable _valset) public {
        require(!_initialized, "Initialized");

        COSMOS_NATIVE_ASSET_PREFIX = "e";
        prophecyClaimCount = 0;
        operator = _operator;
        valset = _valset;
        hasOracle = false;
        hasBridgeBank = false;
        _initialized = true;
    }

    /*
     * @dev: setOracle
     */
    function setOracle(address payable _oracle) public onlyOperator {
        require(
            !hasOracle,
            "The Oracle cannot be updated once it has been set"
        );

        hasOracle = true;
        oracle = _oracle;

        emit LogOracleSet(oracle);
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
        bytes memory _cosmosSender,
        address payable _ethereumReceiver,
        string memory _symbol,
        uint256 _amount
    ) public isActive {
        require(
            Valset(valset).isActiveValidator(msg.sender),
            "Must be an active validator"
        );

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
            _cosmosSender,
            _ethereumReceiver,
            msg.sender,
            tokenAddress,
            symbol,
            _amount,
            Status.Pending
        );

        // Increment count and add the new ProphecyClaim to the mapping
        prophecyClaimCount = prophecyClaimCount.add(1);
        prophecyClaims[prophecyClaimCount] = prophecyClaim;

        emit LogNewProphecyClaim(
            prophecyClaimCount,
            _claimType,
            _cosmosSender,
            _ethereumReceiver,
            msg.sender,
            tokenAddress,
            symbol,
            _amount
        );
    }

    /*
     * @dev: completeProphecyClaim
     *       Allows for the completion of ProphecyClaims once processed by the Oracle.
     *       Burn claims unlock tokens stored by BridgeBank.
     *       Lock claims mint BridgeTokens on BridgeBank's token whitelist.
     */
    function completeProphecyClaim(uint256 _prophecyID)
        public
        isPending(_prophecyID)
    {
        require(
            msg.sender == oracle,
            "Only the Oracle may complete prophecies"
        );

        prophecyClaims[_prophecyID].status = Status.Success;

        ClaimType claimType = prophecyClaims[_prophecyID].claimType;
        if (claimType == ClaimType.Burn) {
            unlockTokens(_prophecyID);
        } else {
            issueBridgeTokens(_prophecyID);
        }

        emit LogProphecyCompleted(_prophecyID, claimType);
    }

    /*
     * @dev: issueBridgeTokens
     *       Issues a request for the BridgeBank to mint new BridgeTokens
     */
    function issueBridgeTokens(uint256 _prophecyID) internal {
        ProphecyClaim memory prophecyClaim = prophecyClaims[_prophecyID];

        BridgeBank(bridgeBank).mintBridgeTokens(
            prophecyClaim.cosmosSender,
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
     * @dev: isProphecyValidatorActive
     *       Returns boolean indicating if the validator that originally
     *       submitted the ProphecyClaim is still an active validator
     */
    function isProphecyClaimValidatorActive(uint256 _prophecyID)
        public
        view
        returns (bool)
    {
        return
            Valset(valset).isActiveValidator(
                prophecyClaims[_prophecyID].originalValidator
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
