pragma solidity 0.5.16;

import "./CosmosBank.sol";
import "./EthereumBank.sol";
import "./EthereumWhitelist.sol";
import "./CosmosWhiteList.sol";
import "../Oracle.sol";
import "../CosmosBridge.sol";
import "./BankStorage.sol";

/**
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
    CosmosWhiteList {

    bool private _initialized;

    using SafeMath for uint256;

    /*
     * @dev: Initializer, sets operator
     */
    function initialize(
        address _operatorAddress,
        address _oracleAddress,
        address _cosmosBridgeAddress,
        address _owner
    ) public {
        require(!_initialized, "Initialized");

        EthereumWhiteList.initialize();
        CosmosWhiteList.initialize();

        operator = _operatorAddress;
        oracle = _oracleAddress;
        cosmosBridge = _cosmosBridgeAddress;
        owner = _owner;
        _initialized = true;
    }

    /*
     * @dev: Modifier to restrict access to operator
     */
    modifier onlyOperator() {
        require(msg.sender == operator, "Must be BridgeBank operator.");
        _;
    }

    /*
     * @dev: Modifier to restrict access to the oracle
     */
    modifier onlyOracle() {
        require(
            msg.sender == oracle,
            "Access restricted to the oracle"
        );
        _;
    }

    /*
     * @dev: Modifier to restrict access to operator
     */
    modifier onlyOwner() {
        require(msg.sender == owner, "Must be Owner.");
        _;
    }



    /*
     * @dev: Modifier to restrict access to the cosmos bridge
     */
    modifier onlyCosmosBridge() {
        require(
            msg.sender == cosmosBridge,
            "Access restricted to the cosmos bridge"
        );
        _;
    }

    /*
     * @dev: Modifier to only allow valid sif addresses
     */
    modifier validSifAddress(bytes memory _sifAddress) {
        require(_sifAddress.length == 42, "Invalid sif address length");
        require(verifySifPrefix(_sifAddress) == true, "Invalid sif address prefix");
        _;
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
            require(listAddress == address(0), "Token already whitelisted");
        } else {
            // if we want to de-whitelist it, make sure that the symbol is 
            // in fact stored in our locked token list before we set to false
            require(uint256(listAddress) > 0, "Token not whitelisted");
        }
        return setTokenInEthWhiteList(_token, _inList);
    }

    // Method that is only for doing the setting of the mapping
    // private so that it is not inheritable or able to be called
    // by anyone other than this contract
    function _updateTokenLimits(address _token, uint256 _amount) private {
        string memory symbol = _token == address(0) ? "ETH" : BridgeToken(_token).symbol();
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
        }
        return true;
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
        bytes memory _cosmosSender,
        address payable _intendedRecipient,
        address _bridgeTokenAddress,
        string memory _symbol,
        uint256 _amount
    ) public onlyCosmosBridge {
        return
            mintNewBridgeTokens(
                _cosmosSender,
                _intendedRecipient,
                _bridgeTokenAddress,
                _symbol,
                _amount
            );
    }

    /*
     * @dev: Burns BridgeTokens representing native Cosmos assets.
     *
     * @param _recipient: bytes representation of destination address.
     * @param _token: token address in origin chain (0x0 if ethereum)
     * @param _amount: value of deposit
     */
    function burn(
        bytes memory _recipient,
        address _token,
        uint256 _amount
    ) public validSifAddress(_recipient) onlyCosmosTokenWhiteList(_token) {
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
        bytes memory _recipient,
        address _token,
        uint256 _amount
    ) public payable onlyEthTokenWhiteList(_token) validSifAddress(_recipient) {
        string memory symbol;

        // Ethereum deposit
        if (msg.value > 0) {
            require(
                _token == address(0),
                "Ethereum deposits require the 'token' address to be the null address"
            );
            require(
                msg.value == _amount,
                "The transactions value must be equal the specified amount (in wei)"
            );
            symbol = "ETH";
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
        string memory _symbol,
        uint256 _amount
    ) public onlyCosmosBridge {
        // Confirm that the bank has sufficient locked balances of this token type
        require(
            getLockedFunds(_symbol) >= _amount,
            "The Bank does not hold enough locked tokens to fulfill this request."
        );

        // Confirm that the bank holds sufficient balances to complete the unlock
        address tokenAddress = lockedTokenList[_symbol];
        if (tokenAddress == address(0)) {
            // uint256 contractBalance = ;
            // revert("no error before 299 for eth unlock");
            require(
                ((address(this)).balance) >= _amount,
                "Insufficient ethereum balance for delivery."
            );
        } else {
            require(
                BridgeToken(tokenAddress).balanceOf(address(this)) >= _amount,
                "Insufficient ERC20 token balance for delivery."
            );
        }
        unlockFunds(_recipient, tokenAddress, _symbol, _amount);
    }

    /*
     * @dev: Exposes an item's current status.
     *
     * @param _id: The item in question.
     * @return: Boolean indicating the lock status.
     */
    function getCosmosDepositStatus(bytes32 _id) public view returns (bool) {
        return isLockedCosmosDeposit(_id);
    }

    /*
     * @dev: Allows access to a Cosmos deposit's information via its unique identifier.
     *
     * @param _id: The deposit to be viewed.
     * @return: Original sender's Ethereum address.
     * @return: Intended Cosmos recipient's address in bytes.
     * @return: The lock deposit's currency, denoted by a token address.
     * @return: The amount locked in the deposit.
     * @return: The deposit's unique nonce.
     */
    function viewCosmosDeposit(bytes32 _id)
        public
        view
        returns (
            bytes memory,
            address payable,
            address,
            uint256
        )
    {
        return getCosmosDeposit(_id);
    }
}
