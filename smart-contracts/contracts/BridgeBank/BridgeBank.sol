pragma solidity 0.5.16;

import "./CosmosBank.sol";
import "./EthereumBank.sol";
import "./EthereumWhitelist.sol";
import "./CosmosWhiteList.sol";
import "../Oracle.sol";
import "../CosmosBridge.sol";
import "./BankStorage.sol";
import "./Pausable.sol";

/*
 * @title BridgeBank
 * @dev Bank contract which coordinates asset-related functionality.
 *      CosmosBank manages the minting and burning of tokens which
 *      represent Cosmos based assets, while EthereumBank manages
 *      the locking and unlocking of Ethereum and ERC20 token assets
 *      based on Ethereum. WhiteList records the ERC20 token address 
 *      list that can be locked.
 **/

contract BridgeBank is BankStorage,
    CosmosBank,
    EthereumBank,
    EthereumWhiteList,
    CosmosWhiteList,
    Pausable {

    bool private _initialized;

    using SafeMath for uint256;

    /*
     * @dev: Initializer, sets operator
     */
    function initialize(
        address _operatorAddress,
        address _cosmosBridgeAddress,
        address _owner,
        address _pauser
    ) public {
        require(!_initialized, "Init");

        EthereumWhiteList.initialize();
        CosmosWhiteList.initialize();
        Pausable.initialize(_pauser);

        operator = _operatorAddress;
        cosmosBridge = _cosmosBridgeAddress;
        owner = _owner;
        _initialized = true;

        // hardcode since this is the first token
        lowerToUpperTokens["erowan"] = "erowan";
        lowerToUpperTokens["eth"] = "eth";
    }

    /*
     * @dev: Modifier to restrict access to operator
     */
    modifier onlyOperator() {
        require(msg.sender == operator, "!operator");
        _;
    }

    /*
     * @dev: Modifier to restrict access to operator
     */
    modifier onlyOwner() {
        require(msg.sender == owner, "!owner");
        _;
    }

    /*
     * @dev: Modifier to restrict access to the cosmos bridge
     */
    modifier onlyCosmosBridge() {
        require(
            msg.sender == cosmosBridge,
            "!cosmosbridge"
        );
        _;
    }

    /*
     * @dev: Modifier to only allow valid sif addresses
     */
    modifier validSifAddress(bytes memory _sifAddress) {
        require(_sifAddress.length == 42, "Invalid len");
        require(verifySifPrefix(_sifAddress) == true, "Invalid sif address");
        _;
    }

    function changeOwner(address _newOwner) public onlyOwner {
        require(_newOwner != address(0), "invalid address");
        owner = _newOwner;
    }

    function changeOperator(address _newOperator) public onlyOperator {
        require(_newOperator != address(0), "invalid address");
        operator = _newOperator;
    }

    /*
     * @dev: function to validate if a sif address has a correct prefix
     */
    function verifySifPrefix(bytes memory _sifAddress) public pure returns (bool) {
        bytes3 sifInHex = 0x736966;

        for (uint256 i = 0; i < sifInHex.length; i++) {
            if (sifInHex[i] != _sifAddress[i]) {
                return false;
            }
        }
        return true;
    }

    /*
     * @dev: Creates a new BridgeToken
     *
     * @param _symbol: The new BridgeToken's symbol
     * @return: The new BridgeToken contract's address
     */
    function createNewBridgeToken(string memory _symbol)
        public
        onlyCosmosBridge
        returns (address)
    {
        address newTokenAddress = deployNewBridgeToken(_symbol);
        setTokenInCosmosWhiteList(newTokenAddress, true);

        return newTokenAddress;
    }

    /*
     * @dev: Creates a new BridgeToken
     *
     * @param _symbol: The new BridgeToken's symbol
     * @return: The new BridgeToken contract's address
     */
    function addExistingBridgeToken(
        address _contractAddress
    ) public onlyOwner returns (address) {
        setTokenInCosmosWhiteList(_contractAddress, true);

        return useExistingBridgeToken(_contractAddress);
    }

    /*
     * @dev: Set the token address in whitelist
     *
     * @param _token: ERC 20's address
     * @param _inList: set the _token in list or not
     * @return: new value of if _token in whitelist
     */
    function updateEthWhiteList(address _token, bool _inList)
        public
        onlyOperator
        returns (bool)
    {
        string memory symbol = BridgeToken(_token).symbol();
        address listAddress = lockedTokenList[symbol];
        
        // Do not allow a token with the same symbol to be whitelisted
        if (_inList) {
            // if we want to add it to the whitelist, make sure that the address
            // is 0, meaning we have not seen that symbol in the whitelist before
            require(listAddress == address(0), "whitelisted");
        } else {
            // if we want to de-whitelist it, make sure that the symbol is 
            // in fact stored in our locked token list before we set to false
            require(uint256(listAddress) > 0, "!whitelisted");
        }
        lowerToUpperTokens[toLower(symbol)] = symbol;
        return setTokenInEthWhiteList(_token, _inList);
    }

    // Method that is only for doing the setting of the mapping
    // private so that it is not inheritable or able to be called
    // by anyone other than this contract
    function _updateTokenLimits(address _token, uint256 _amount) private {
        string memory symbol = _token == address(0) ? "eth" : BridgeToken(_token).symbol();
        maxTokenAmount[symbol] = _amount;
    }

    function updateTokenLockBurnLimit(address _token, uint256 _amount)
        public
        onlyOperator
        returns (bool)
    {
        _updateTokenLimits(_token, _amount);
        return true;
    }

    function bulkWhitelistUpdateLimits(
        address[] calldata tokenAddresses,
        uint256[] calldata tokenLimit
        ) external
        onlyOperator
        returns (bool)
    {
        require(tokenAddresses.length == tokenLimit.length, "!same length");
        for (uint256 i = 0; i < tokenAddresses.length; i++) {
            _updateTokenLimits(tokenAddresses[i], tokenLimit[i]);
            setTokenInEthWhiteList(tokenAddresses[i], true);
            string memory symbol = BridgeToken(tokenAddresses[i]).symbol();
            lowerToUpperTokens[toLower(symbol)] = symbol;
        }
        return true;
    }

    /*
     * @dev: Burns BridgeTokens representing native Cosmos assets.
     *
     * @param _recipient: bytes representation of destination address.
     * @param _token: token address in origin chain (0x0 if ethereum)
     * @param _amount: value of deposit
     */
    function burn(
        bytes calldata _recipient,
        address _token,
        uint256 _amount
    ) external validSifAddress(_recipient) onlyCosmosTokenWhiteList(_token) whenNotPaused {
        string memory symbol = BridgeToken(_token).symbol();

        if (_amount > maxTokenAmount[symbol]) {
            revert("Amount being transferred is over the limit for this token");
        }

        BridgeToken(_token).burnFrom(msg.sender, _amount);
        burnFunds(msg.sender, _recipient, _token, symbol, _amount);
    }

    /*
     * @dev: Locks received Ethereum/ERC20 funds.
     *
     * @param _recipient: bytes representation of destination address.
     * @param _token: token address in origin chain (0x0 if ethereum)
     * @param _amount: value of deposit
     */
    function lock(
        bytes calldata _recipient,
        address _token,
        uint256 _amount
    ) external payable onlyEthTokenWhiteList(_token) validSifAddress(_recipient) whenNotPaused {
        string memory symbol;

        // Ethereum deposit
        if (msg.value > 0) {
            require(
                _token == address(0),
                "!address(0)"
            );
            require(
                msg.value == _amount,
                "incorrect eth amount"
            );
            symbol = "eth";
            // ERC20 deposit
        } else {
            IERC20 tokenToTransfer = IERC20(_token);
            tokenToTransfer.safeTransferFrom(
                msg.sender,
                address(this),
                _amount
            );
            symbol = BridgeToken(_token).symbol();
        }

        if (_amount > maxTokenAmount[symbol]) {
            revert("Amount being transferred is over the limit");
        }
        lockFunds(msg.sender, _recipient, _token, symbol, _amount);
    }

    /*
     * @dev: Unlocks Ethereum and ERC20 tokens held on the contract.
     *
     * @param _recipient: recipient's Ethereum address
     * @param _token: token contract address
     * @param _symbol: token symbol
     * @param _amount: wei amount or ERC20 token count
     */
    function unlock(
        address payable _recipient,
        address tokenAddress,
        uint256 _amount
    ) external onlyCosmosBridge whenNotPaused {
        // address tokenAddress = getExternalAddressByLowerCaseSymbol(_symbol);

        unlockFunds(_recipient, tokenAddress, _amount);
    }

    /*
     * @dev: Mints new BankTokens
     *
     * @param _cosmosSender: The sender's Cosmos address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _cosmosTokenAddress: The currency type
     * @param _symbol: comsos token symbol
     * @param _amount: number of comsos tokens to be minted
     */
    function mintBridgeTokens(
        address payable _intendedRecipient,
        address _bridgeTokenAddress,
        uint256 _amount
    ) external onlyCosmosBridge whenNotPaused {
        // address _bridgeTokenAddress = getBridgeTokenAddressByLowerCaseSymbol(_symbol);
        return
            mintNewBridgeTokens(
                _intendedRecipient,
                _bridgeTokenAddress,
                _amount
            );
    }


    function getExternalAddressByLowerCaseSymbol(string memory symbol)
        public
        view
        returns (address)
    {
        symbol = safeLowerToUpperTokens(symbol);
        return lockedTokenList[symbol];
    }
    
    function getBridgeTokenAddressByLowerCaseSymbol(string memory symbol)
        public
        view
        returns (address)
    {
        symbol = safeLowerToUpperTokens(symbol);
        return controlledBridgeTokens[symbol];
    }

    /*
    * @dev fallback function for ERC223 tokens so that we can receive these tokens in our contract
    * Don't need to do anything to handle these tokens
    */
    function tokenFallback(address _from, uint _value, bytes memory _data) public {}
}
