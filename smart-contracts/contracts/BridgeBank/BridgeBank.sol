// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

import "./CosmosBank.sol";
import "./EthereumWhitelist.sol";
import "./CosmosWhiteList.sol";
import "../Oracle.sol";
import "../CosmosBridge.sol";
import "./BankStorage.sol";
import "./Pausable.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

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
    EthereumWhiteList,
    CosmosWhiteList,
    Pausable {

    using SafeERC20 for IERC20;

    bool private _initialized;

    /*
     * @dev: Initializer
     */
    function initialize(
        address _operator,
        address _cosmosBridgeAddress,
        address _owner,
        address _pauser,
        uint256 _networkDescriptor
    ) public {
        require(!_initialized, "Init");

        CosmosWhiteList._cosmosWhitelistInitialize();
        EthereumWhiteList.initialize();
        Pausable._pausableInitialize(_pauser);

        operator = _operator;
        cosmosBridge = _cosmosBridgeAddress;
        owner = _owner;
        networkDescriptor = _networkDescriptor;
        _initialized = true;
        contractName[address(0)] = "Ethereum";
        contractSymbol[address(0)] = "ETH";
    }

    /*
     * @dev: Modifier to restrict access to operator
     */
    modifier onlyOperator() {
        require(msg.sender == operator, "!operator");
        _;
    }

    /*
     * @dev: Modifier to restrict access to owner
     */
    modifier onlyOwner {
        require(msg.sender == owner, "!owner");
        _;
    }

    /*
     * @dev: Modifier to restrict access to the cosmos bridge
     */
    modifier onlyCosmosBridge {
        require(
            msg.sender == cosmosBridge,
            "!cosmosbridge"
        );
        _;
    }

    /*
     * @dev: Modifier to only allow valid sif addresses
     */
    modifier validSifAddress(bytes calldata sifAddress) {
        require(verifySifAddress(sifAddress) == true, "INV_SIF_ADDR");
        _;
    }

    /*
     * @dev: Set the token address in Eth whitelist
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
        // Do not allow a token with the same address to be whitelisted
        if (_inList) {
            // if we want to add it to the whitelist, make sure it's not there yet
            require(!getTokenInEthWhiteList(_token), "whitelisted");
        } else {
            // if we want to de-whitelist it, make sure that the token is already whitelisted 
            require(getTokenInEthWhiteList(_token), "!whitelisted");
        }
        return setTokenInEthWhiteList(_token, _inList);
    }

    /*
     * @dev: Set the token address in whitelist
     *
     * @param token: ERC 20's address
     * @param inList: set the token in list or not
     * @return: new value of if token in whitelist
     */
    function setTokenInCosmosWhiteList(address token, bool inList)
        internal returns (bool)
    {
        _cosmosTokenWhiteList[token] = inList;
        emit LogWhiteListUpdate(token, inList);
        return inList;
    }

    function changeOwner(address newOwner) public onlyOwner {
        require(newOwner != address(0), "invalid address");
        owner = newOwner;
    }

    function changeOperator(address _newOperator) public onlyOperator {
        require(_newOperator != address(0), "invalid address");
        operator = _newOperator;
    }


    /*
     * @dev: function to validate if a sif address has a correct prefix
     */
    function verifySifPrefix(bytes calldata sifAddress) private pure returns (bool) {
        bytes3 sifInHex = 0x736966;

        for (uint256 i = 0; i < sifInHex.length; i++) {
            if (sifInHex[i] != sifAddress[i]) {
                return false;
            }
        }
        return true;
    }

    function verifySifAddressLength(bytes calldata sifAddress) private pure returns (bool) {
        return sifAddress.length == 42;
    }

    function verifySifAddress(bytes calldata sifAddress) private pure returns (bool) {
        return verifySifAddressLength(sifAddress) && verifySifPrefix(sifAddress);
    }

    // function used only for testing
    function VSA(bytes calldata _sifAddress) external pure returns (bool) {
        return verifySifAddress(_sifAddress);
    }

    /*
     * @dev: Creates a new BridgeToken
     *
     * @param _symbol: The new BridgeToken's symbol
     * @return: The new BridgeToken contract's address
     */
    function createNewBridgeToken(
        string calldata name,
        string calldata symbol,
        uint8 decimals
    ) external onlyCosmosBridge returns (address) {
        address newTokenAddress = deployNewBridgeToken(
            name,
            symbol,
            decimals
        );
        setTokenInCosmosWhiteList(newTokenAddress, true);

        return newTokenAddress;
    }

    /*
     * @dev: Creates a new BridgeToken
     *
     * @param contractAddress: The new BridgeToken's address
     */
    function addExistingBridgeToken(
        address contractAddress    
    ) external onlyOwner returns (bool) {
        return setTokenInCosmosWhiteList(contractAddress, true);
    }

    function handleUnpeg(
        address payable ethereumReceiver,
        address tokenAddress,
        uint256 amount
    ) external onlyCosmosBridge whenNotPaused {
        // if this is a bridge controlled token, then we need to mint
        if (getCosmosTokenInWhiteList(tokenAddress)) {
            mintNewBridgeTokens(
                ethereumReceiver,
                tokenAddress,
                amount
            );
        } else {
            // if this is an external token, we unlock
            unlock(ethereumReceiver, tokenAddress, amount);
        }
    }

    /*
     * @dev: Burns BridgeTokens representing native Cosmos assets.
     *
     * @param recipient: bytes representation of destination address.
     * @param token: token address in origin chain (0x0 if ethereum)
     * @param amount: value of deposit
     */
    function burn(
        bytes calldata recipient,
        address token,
        uint256 amount
    ) external validSifAddress(recipient) onlyCosmosTokenWhiteList(token) whenNotPaused {
        // burn the tokens
        BridgeToken(token).burnFrom(msg.sender, amount);
        
        // decimals defaults to 18 if call to decimals fails
        uint8 decimals = getDecimals(token);

        lockBurnNonce = lockBurnNonce + 1;

        emit LogBurn(
            msg.sender,
            recipient,
            token,
            amount,
            lockBurnNonce,
            decimals,
            networkDescriptor
        );
    }

    function getName(address token) private returns (string memory) {
        string memory name = contractName[token];

        // check to see if we already have this name stored in the smart contract
        if (keccak256(abi.encodePacked(name)) != keccak256(abi.encodePacked(""))) {
            return name;
        }

        try BridgeToken(token).name() returns (string memory _name) {
            name = _name;
            contractName[token] = _name;
        } catch {
            // if we can't access the decimals function of this token,
            // assume that it has 18 decimals
            name = "";
        }

        return name;
    }

    function getSymbol(address token) private returns (string memory) {
        string memory symbol = contractSymbol[token];

        // check to see if we already have this name stored in the smart contract
        if (keccak256(abi.encodePacked(symbol)) != keccak256(abi.encodePacked(""))) {
            return symbol;
        }

        try BridgeToken(token).symbol() returns (string memory _symbol) {
            symbol = _symbol;
            contractSymbol[token] = _symbol;
        } catch {
            // if we can't access the decimals function of this token,
            // assume that it has 18 decimals
            symbol = "";
        }

        return symbol;
    }

    function getDecimals(address token) private returns (uint8) {
        uint8 decimals = contractDecimals[token];
        if (decimals > 0) {
            return decimals;
        }

        try BridgeToken(token).decimals() returns (uint8 _decimals) {
            require(decimals < 100, "invalid decimals");
            decimals = _decimals;
            contractDecimals[token] = _decimals;
        } catch {
            // if we can't access the decimals function of this token,
            // assume that it has 18 decimals
            decimals = 18;
        }

        return decimals;
    }

    /*
     * @dev: Locks received Ethereum/ERC20 funds.
     *
     * @param recipient: bytes representation of destination address.
     * @param token: token address in origin chain (0x0 if ethereum)
     * @param amount: value of deposit
     */
    function lock(
        bytes calldata recipient,
        address token,
        uint256 amount
    ) external payable validSifAddress(recipient) whenNotPaused {
        if (token == address(0)) {
            return handleNativeCurrencyLock(recipient, amount);
        }
        require(msg.value == 0, "INV_NATIVE_SEND");

        lockBurnNonce += 1;
        _lockTokens(recipient, token, amount, lockBurnNonce);
    }

    /**
     * Locks multiple tokens in the bridge contract in a single tx.
     * This is used to handle the case where the user is sending tokens
     *
     * @param recipient: bytes representation of destination address.
     * @param token: token address
     * @param amount: value of deposit
     */
    function multiLock(
        bytes[] calldata recipient,
        address[] calldata token,
        uint256[] calldata amount
    ) external whenNotPaused {
        require(recipient.length == token.length, "M_P");
        require(token.length == amount.length, "M_P");

        // use intermediate lock burn nonce to distinguish between different lock calls
        // this alows us to track the lock calls in the logs
        // and to prevent double locking
        // (i.e. if a user calls lock twice with the same amount, we don't want to lock twice)
        // This pattern of using the intermediate value will save us gas
        // by utilizing the stack for all intermediate values
        uint256 intermediateLockBurnNonce = lockBurnNonce;

        for (uint256 i = 0; i < recipient.length; i++) {
            intermediateLockBurnNonce++;

            _lockTokens(
                recipient[i],
                token[i],
                amount[i],
                intermediateLockBurnNonce
            );
        }
        lockBurnNonce = intermediateLockBurnNonce;
    }

    /**
     * Locks multiple tokens in the bridge contract in a single tx.
     * This is used to handle the case where the user is sending tokens
     *
     * @param recipient: bytes representation of destination address.
     * @param token: token address
     * @param amount: value of deposit
     * @param isBurn: value of deposit
     */
    function multiLockBurn(
        bytes[] calldata recipient,
        address[] calldata token,
        uint256[] calldata amount,
        bool[] calldata isBurn
    ) external whenNotPaused {
        // all array inputs must be of the same length
        // else throw malformed params error
        require(recipient.length == token.length, "M_P");
        require(token.length == amount.length, "M_P");
        require(token.length == isBurn.length, "M_P");

        uint256 intermediateLockBurnNonce = lockBurnNonce;

        for (uint256 i = 0; i < recipient.length; i++) {
            intermediateLockBurnNonce++;

            if (isBurn[i]) {
                _burnTokens(
                    recipient[i],
                    token[i],
                    amount[i],
                    intermediateLockBurnNonce
                );
            } else {
                _lockTokens(
                    recipient[i],
                    token[i],
                    amount[i],
                    intermediateLockBurnNonce
                );
            }
        }
        lockBurnNonce = intermediateLockBurnNonce;
    }

    function _lockTokens(
        bytes calldata recipient,
        address tokenAddress,
        uint256 tokenAmount,
        uint256 _lockBurnNonce
    ) private onlyEthTokenWhiteList(tokenAddress) onlyTokenNotInCosmosWhiteList(tokenAddress) validSifAddress(recipient) {
        IERC20 tokenToTransfer = IERC20(tokenAddress);
        // lock tokens
        tokenToTransfer.safeTransferFrom(
            msg.sender,
            address(this),
            tokenAmount
        );

        // decimals defaults to 18 if call to decimals fails
        uint8 decimals = getDecimals(tokenAddress);

        // Get name and symbol
        string memory name = getName(tokenAddress);
        string memory symbol = getSymbol(tokenAddress);

        emit LogLock(
            msg.sender,
            recipient,
            tokenAddress,
            tokenAmount,
            _lockBurnNonce,
            decimals,
            symbol,
            name,
            networkDescriptor
        );
    }

    function _burnTokens(
        bytes calldata recipient,
        address tokenAddress,
        uint256 tokenAmount,
        uint256 _lockBurnNonce
    ) private onlyCosmosTokenWhiteList(tokenAddress) validSifAddress(recipient) {
        BridgeToken tokenToTransfer = BridgeToken(tokenAddress);
        // burn tokens
        tokenToTransfer.burnFrom(
            msg.sender,
            tokenAmount
        );

        // decimals defaults to 18 if call to decimals fails
        uint8 decimals = getDecimals(tokenAddress);

        // Get name and symbol
        string memory name = getName(tokenAddress);
        string memory symbol = getSymbol(tokenAddress);

        emit LogBurn(
            msg.sender,
            recipient,
            tokenAddress,
            tokenAmount,
            _lockBurnNonce,
            decimals,
            networkDescriptor
        );
    }

    /*
     * Locks received Ethereum/ERC20 funds.
     *
     * @param recipient: bytes representation of destination address.
     * @param token: token address in origin chain (0x0 if ethereum)
     * @param amount: value of deposit
     */
    function handleNativeCurrencyLock(
        bytes calldata recipient,
        uint256 amount
    ) internal validSifAddress(recipient) {
        require(msg.value == amount, "amount mismatch");

        address token = address(0);

        // decimals defaults to 18 if call to decimals fails
        uint8 decimals = 18;

        // Get name and symbol
        string memory name = getName(token);
        string memory symbol = getSymbol(token);

        lockBurnNonce = lockBurnNonce + 1;

        emit LogLock(
            msg.sender,
            recipient,
            token,
            amount,
            lockBurnNonce,
            decimals,
            symbol,
            name,
            networkDescriptor
        );
    }

    /**
     *
     * @param recipient: recipient's Ethereum address
     * @param token: token contract address
     * @param amount: wei amount or ERC20 token count
     */
    function unlock(
        address payable recipient,
        address token,
        uint256 amount
    ) internal {
        // Transfer funds to intended recipient
        if (token == address(0)) {
            (bool success,) = recipient.call{value: amount}("");
            require(success, "error sending ether");
        } else {
            IERC20 tokenToTransfer = IERC20(token);
            tokenToTransfer.safeTransfer(recipient, amount);
        }

        emit LogUnlock(recipient, token, amount);
    }
}
