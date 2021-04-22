pragma solidity 0.8.0;

import "./Oracle.sol";
import "./BridgeBank/BridgeBank.sol";
import "./CosmosBridgeStorage.sol";

contract CosmosBridge is CosmosBridgeStorage, Oracle {
    bool private _initialized;
    uint256[100] private ___gap;

    /*
     * @dev: Event declarations
     */
    event LogBridgeBankSet(address bridgeBank);

    event LogNewProphecyClaim(
        uint256 indexed prophecyID,
        ClaimType claimType,
        address indexed ethereumReceiver,
        uint256 indexed amount
    );

    event LogNewBridgeTokenCreated(
        uint8 decimals,
        uint256 indexed sourceChainDescriptor,
        string name,
        string symbol,
        address indexed sourceContractAddress,
        address indexed bridgeTokenAddress
    );

    event LogProphecyCompleted(uint256 prophecyID, ClaimType claimType);

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
        bytes calldata _cosmosSender,
        uint256 _cosmosSenderSequence,
        address payable _ethereumReceiver,
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

    /**
     * function: newProphecyClaim
     *       Creates a new burn or lock prophecy claim, adding it to the prophecyClaims mapping.
     *       Burn claims require that there are enough locked Ethereum assets to complete the prophecy.
     *       Lock claims have a new token contract deployed or use an existing contract based on symbol.
     *
     * @param _claimType type of claim, either lock or burn
     * @param _cosmosSender sifchain sender's address
     * @param _cosmosSenderSequence nonce of the cosmos sender
     * @param _ethereumReceiver ethereum address of the receiver
     * @param _tokenAddress address of the token to send
     * @param _amount amount of token to send
     * @param _doublePeg whether or not this is a double peg transaction
     *
     */
    function newProphecyClaim(
        ClaimType _claimType,
        bytes calldata _cosmosSender,
        uint256 _cosmosSenderSequence,
        address payable _ethereumReceiver,
        address _tokenAddress,
        uint256 _amount,
        bool _doublePeg
    ) external onlyValidator validClaimType(_claimType) {

        uint256 _prophecyID = getProphecyID(
            _claimType, 
            _cosmosSender,
            _cosmosSenderSequence,
            _ethereumReceiver,
            _tokenAddress,
            _amount
        );

        require(!prophecyRedeemed[_prophecyID], "prophecy already redeemed");

        if (oracleClaimValidators[_prophecyID] == 0) {
            emit LogNewProphecyClaim(
                _prophecyID,
                _claimType,
                _ethereumReceiver,
                _amount
            );
        }
    
        bool claimComplete = newOracleClaim(_prophecyID, msg.sender);

        if (claimComplete) {
            // you cannot redeem this prophecy again
            prophecyRedeemed[_prophecyID] = true;

            _tokenAddress = _doublePeg ? sourceAddressToDestinationAddress[_tokenAddress] : _tokenAddress;

            completeProphecyClaim(
                _claimType,
                _prophecyID,
                _ethereumReceiver,
                _tokenAddress,
                _amount
            );
        }
    }
    
    /**
     *
     * @param _cosmosSender address of the sifchain address
     * @param _symbol symbol of the ERC20 token on the source chain
     * @param _name name of the ERC20 token on the source chain
     * @param _sourceChainTokenAddress address of the ERC20 token on the source chain
     * @param _decimals of the ERC20 token on the source chain
     * @param chainDescriptor descriptor of the source chain
     */
    function createNewBridgeToken(
        bytes calldata _cosmosSender,
        string calldata _symbol,
        string calldata _name,
        address _sourceChainTokenAddress,
        uint8 _decimals,
        uint256 chainDescriptor
    ) external onlyValidator {
        // need to make a business decision on what this symbol should be
        // First lock of this asset, deploy new contract and get new symbol/token address
        address tokenAddress = BridgeBank(bridgeBank)
            .createNewBridgeToken(
                _name,
                _symbol,
                _decimals
            );

        sourceAddressToDestinationAddress[_sourceChainTokenAddress] = tokenAddress;

        emit LogNewBridgeTokenCreated(
            _decimals,
            chainDescriptor,
            _name,
            _symbol,
            _sourceChainTokenAddress,
            tokenAddress
        );
    }

    // struct prophecyBundle {
    //     ClaimType _claimType;
    //     bytes _cosmosSender;
    //     string _symbol;
    //     uint256 _cosmosSenderSequence;
    //     address payable _ethereumReceiver;
    //     address _tokenAddress;
    //     uint256 _amount;
    // }

    // function batchSubmitProphecies(
    //     prophecyBundle[] calldata _prophecies
    // ) external onlyValidator {
    //     for (uint256 i = 0; i < _prophecies.length; i++) {
            
    //     }
    // }

    /*
     * @dev: completeProphecyClaim
     *       Allows for the completion of ProphecyClaims once processed by the Oracle.
     *       Burn claims unlock tokens stored by BridgeBank.
     *       Lock claims mint BridgeTokens on BridgeBank's token whitelist.
     */
    function completeProphecyClaim(
        ClaimType _claimType,
        uint256 _prophecyID,
        address payable ethereumReceiver,
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
