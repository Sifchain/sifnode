// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "./CosmosBank.sol";
import "./EthereumWhitelist.sol";
import "./CosmosWhiteList.sol";
import "../Oracle.sol";
import "../CosmosBridge.sol";
import "./BankStorage.sol";
import "./Pausable.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/**
 * @title Bridge Bank
 * @dev Bank contract which coordinates asset-related functionality.
 *      CosmosBank manages the minting and burning of tokens which
 *      represent Cosmos based assets, while EthereumBank manages
 *      the locking and unlocking of Ethereum and ERC20 token assets
 *      based on Ethereum. WhiteList records the ERC20 token address
 *      list that can be locked.
 */
contract BridgeBank is BankStorage, CosmosBank, EthereumWhiteList, CosmosWhiteList, Pausable {
  using SafeERC20 for IERC20;

  /**
   * @dev Has the contract been initialized?
   */
  bool private _initialized;

  /**
   * @dev the blocklist contract
   */
  IBlocklist public blocklist;

  /**
   * @notice is the blocklist active?
   */
  bool public hasBlocklist;

  /**
   * @notice network descriptor
   */
  int32 public networkDescriptor;

  /**
   * @notice contract decimals based off of contract address
   */
  mapping(address => uint8) public contractDecimals;

  /**
   * @notice contract symbol based off of address
   */
  mapping(address => string) public contractSymbol;

  /**
   * @notice contract name based off of address
   */
  mapping(address => string) public contractName;

  /**
   * @notice contract denom based off of address
   */
  mapping(address => string) public contractDenom;

  /**
   * @dev Has the contract been reinitialized?
   */
  bool private _reinitialized;

  /**
   * @dev The address of the Rowan Token
   */
   address public rowanTokenAddress;

  /**
   * @notice Initializer
   * @param _operator Manages the contract
   * @param _cosmosBridgeAddress The CosmosBridge contract's address
   * @param _owner Manages whitelists
   * @param _pauser Can pause the system
   * @param _networkDescriptor Indentifies the connected network
   * @param _rowanTokenAddress The address of the Rowan ERC20 contract on this network
   */
  function initialize(
    address _operator,
    address _cosmosBridgeAddress,
    address _owner,
    address _pauser,
    int32 _networkDescriptor,
    address _rowanTokenAddress
  ) public {
    require(!_initialized, "Init");

    CosmosWhiteList._cosmosWhitelistInitialize();
    EthereumWhiteList.initialize();

    contractName[address(0)] = "EVMNATIVE";
    contractSymbol[address(0)] = "EVMNATIVE";

    _initialized = true;

    _initialize(_operator, _cosmosBridgeAddress, _owner, _pauser, _networkDescriptor, _rowanTokenAddress);
  }

  /**
   * @notice Allows the current operator to reinitialize the contract
   * @param _operator Manages the contract
   * @param _cosmosBridgeAddress The CosmosBridge contract's address
   * @param _owner Manages whitelists
   * @param _pauser Can pause the system
   * @param _networkDescriptor Indentifies the connected network
   * @param _rowanTokenAddress The address of the Rowan ERC20 contract on this network
   */
  function reinitialize(
    address _operator,
    address _cosmosBridgeAddress,
    address _owner,
    address _pauser,
    int32 _networkDescriptor,
    address _rowanTokenAddress
  ) public onlyOperator {
    require(!_reinitialized, "Already reinitialized");

    _reinitialized = true;

    _initialize(_operator, _cosmosBridgeAddress, _owner, _pauser, _networkDescriptor, _rowanTokenAddress);
  }

  /**
   * @dev Internal function called by initialize() and reinitialize()
   * @param _operator Manages the contract
   * @param _cosmosBridgeAddress The CosmosBridge contract's address
   * @param _owner Manages whitelists
   * @param _pauser Can pause the system
   * @param _networkDescriptor Indentifies the connected network
   * @param _rowanTokenAddress The address of the Rowan ERC20 contract on this network
   */
  function _initialize(
    address _operator,
    address _cosmosBridgeAddress,
    address _owner,
    address _pauser,
    int32 _networkDescriptor,
    address _rowanTokenAddress
  ) private {
    Pausable._pausableInitialize(_pauser);

    operator = _operator;
    cosmosBridge = _cosmosBridgeAddress;
    owner = _owner;
    networkDescriptor = _networkDescriptor;
    rowanTokenAddress = _rowanTokenAddress;
  }

  /**
   * @dev Set or update the rowanTokenAddress Only the operator can call this function
   * @param _rowanTokenAddress The address of the Rowan ERC20 contract on this network
   * @notice Can be set to null address if Rowan on this network is a standard BridgeToken
   */
   function setRowanTokenAddress(address _rowanTokenAddress) public onlyOperator {
    rowanTokenAddress = _rowanTokenAddress;
   }

  /**
   * @dev Modifier to restrict access to operator
   */
  modifier onlyOperator() {
    require(msg.sender == operator, "!operator");
    _;
  }

  /**
   * @dev Modifier to restrict access to owner
   */
  modifier onlyOwner() {
    require(msg.sender == owner, "!owner");
    _;
  }

  /**
   * @dev Modifier to restrict access to the cosmos bridge
   */
  modifier onlyCosmosBridge() {
    require(msg.sender == cosmosBridge, "!cosmosbridge");
    _;
  }

  /**
   * @dev Modifier to restrict EVM addresses
   */
  modifier onlyNotBlocklisted(address account) {
    if (hasBlocklist) {
      require(!blocklist.isBlocklisted(account), "Address is blocklisted");
    }
    _;
  }

  /**
   * @dev Modifier to only allow valid sif addresses
   */
  modifier validSifAddress(bytes calldata sifAddress) {
    require(verifySifAddress(sifAddress) == true, "INV_SIF_ADDR");
    _;
  }

  /**
   * @dev Set the token address in Cosmos whitelist
   * @param token ERC 20's address
   * @param inList Set the token in list or not
   * @return New value of if token is in whitelist
   */
  function setTokenInCosmosWhiteList(address token, bool inList) internal returns (bool) {
    _cosmosTokenWhiteList[token] = inList;
    emit LogWhiteListUpdate(token, inList);
    return inList;
  }

  /**
   * @notice Transfers ownership of this contract to `newOwner`
   * @dev Cannot transfer ownership to the zero address
   * @param newOwner The new owner's address
   */
  function changeOwner(address newOwner) public onlyOwner {
    require(newOwner != address(0), "invalid address");
    owner = newOwner;
  }

  /**
   * @notice Transfers the operator role to `_newOperator`
   * @dev Cannot transfer role to the zero address
   * @param _newOperator: the new operator's address
   */
  function changeOperator(address _newOperator) public onlyOperator {
    require(_newOperator != address(0), "invalid address");
    operator = _newOperator;
  }

  /**
   * @dev Validates if a sif address has a correct prefix
   * @param sifAddress The Sif address to check
   * @return Boolean: does it have the correct prefix?
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

  /**
   * @dev Validates if a sif address has the correct length
   * @param sifAddress The Sif address to check
   * @return Boolean: does it have the correct length?
   */
  function verifySifAddressLength(bytes calldata sifAddress) private pure returns (bool) {
    return sifAddress.length == 42;
  }

  /**
   * @dev Validates if a sif address has a correct prefix and the correct length
   * @param sifAddress The Sif address to be validated
   * @return Boolean: is it a valid Sif address?
   */
  function verifySifAddress(bytes calldata sifAddress) private pure returns (bool) {
    return verifySifAddressLength(sifAddress) && verifySifPrefix(sifAddress);
  }

  /**
   * @notice Validates whether `_sifAddress` is a valid Sif address
   * @dev Function used only for testing
   * @param _sifAddress Bytes representation of a Sif address
   * @return Boolean: is it a valid Sif address?
   */
  function VSA(bytes calldata _sifAddress) external pure returns (bool) {
    return verifySifAddress(_sifAddress);
  }

  /**
   * @notice CosmosBridge calls this function to create a new BridgeToken
   * @dev Only CosmosBridge can create a new BridgeToken
   * @param name The new BridgeToken's name
   * @param symbol The new BridgeToken's symbol
   * @param decimals The new BridgeToken's decimals
   * @param cosmosDenom The new BridgeToken's denom
   * @return The new BridgeToken contract's address
   */
  function createNewBridgeToken(
    string calldata name,
    string calldata symbol,
    uint8 decimals,
    string calldata cosmosDenom
  ) external onlyCosmosBridge returns (address) {
    address newTokenAddress = deployNewBridgeToken(name, symbol, decimals, cosmosDenom);
    setTokenInCosmosWhiteList(newTokenAddress, true);
    contractDenom[newTokenAddress] = cosmosDenom;

    return newTokenAddress;
  }

  /**
   * @notice Allows the owner to add `contractAddress` as an existing BridgeToken
   * @dev Adds the token to Cosmos Whitelist
   * @param contractAddress The token's address
   * @return New value of if token is in whitelist (true)
   */
  function addExistingBridgeToken(address contractAddress) external onlyOwner returns (bool) {
    return setTokenInCosmosWhiteList(contractAddress, true);
  }

  /**
   * @notice Allows the owner to add many contracts as existing BridgeTokens
   * @dev Adds tokens to Cosmos Whitelist in a batch
   * @param contractsAddresses The list of tokens addresses
   * @return true if the operation succeeded
   */
  function batchAddExistingBridgeTokens(address[] calldata contractsAddresses)
    external
    onlyOwner
    returns (bool)
  {
    uint256 contractLength = contractsAddresses.length;
    for (uint256 i = 0; i < contractLength;) {
      setTokenInCosmosWhiteList(contractsAddresses[i], true);
      unchecked{ ++i; }
    }

    return true;
  }

  /**
   * @notice CosmosBridge calls this function to mint or unlock tokens
   * @dev Controlled tokens will be minted, others will be unlocked
   * @param ethereumReceiver Tokens will be sent to this address
   * @param tokenAddress The BridgeToken's address
   * @param amount How much should be minted or unlocked
   */
  function handleUnpeg(
    address payable ethereumReceiver,
    address tokenAddress,
    uint256 amount
  ) external onlyCosmosBridge whenNotPaused onlyNotBlocklisted(ethereumReceiver) {
    // if this is a bridge controlled token, then we need to mint
    if (getCosmosTokenInWhiteList(tokenAddress)) {
      mintNewBridgeTokens(ethereumReceiver, tokenAddress, amount);
    } else {
      // if this is an external token, we unlock
      unlock(ethereumReceiver, tokenAddress, amount);
    }
  }

  /**
   * @notice Burns `amount` `token` tokens for `recipient`
   * @dev Burns BridgeTokens representing native Cosmos assets
   * @param recipient Bytes representation of destination address
   * @param token Token address in origin chain (0x0 if ethereum)
   * @param amount How much will be burned
   */
  function burn(
    bytes calldata recipient,
    address token,
    uint256 amount
  )
    external
    validSifAddress(recipient)
    onlyCosmosTokenWhiteList(token)
    onlyNotBlocklisted(msg.sender)
    whenNotPaused
  {
    uint256 currentLockBurnNonce = lockBurnNonce + 1;
    lockBurnNonce = currentLockBurnNonce;

    _burnTokens(recipient, token, amount, currentLockBurnNonce);
  }

  /**
   * @dev Fetches the name of a token by address
   * @param token The BridgeTokens's address
   * @return The bridgeTokens's name or ''
   */
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
      // if we can't access the name function of this token, return an empty string
      name = "";
    }

    return name;
  }

  /**
   * @dev Fetches the symbol of a token by address
   * @param token The bridgeTokens's address
   * @return The bridgeTokens's symbol or ''
   */
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
      // if we can't access the symbol function of this token, return an empty string
      symbol = "";
    }

    return symbol;
  }

  /**
   * @dev Fetches the decimals of a token by address
   * @param token The bridgeTokens's address
   * @return The bridgeTokens's decimals or 0
   */
  function getDecimals(address token) private returns (uint8) {
    uint8 decimals = contractDecimals[token];
    if (decimals > 0) {
      return decimals;
    }

    try BridgeToken(token).decimals() returns (uint8 _decimals) {
      decimals = _decimals;
      contractDecimals[token] = _decimals;
    } catch {
      // if we can't access the decimals function of this token,
      // assume that it has 0 decimals
      decimals = 0;
    }

    return decimals;
  }

  /**
   * @dev Fetches the current token balance the bridgebank holds of a given token address
   * @param token The bridgeToken's address
   * @return The balance of the bridgebanks account with the bridge token
   */
  function getBalance(address token) private view returns (uint256) {
    uint256 balance;
    try BridgeToken(token).balanceOf(address(this)) returns (uint256 _balance) {
      balance = _balance;
    } catch {
      balance = 0;
    }
    return balance;
  }

  /**
   * @dev Function which transfers the requested amount of tokens from the calling users account to the bridgebanks account
   *      This function checks the balances before and after the transfer and reports the amount that was transfered in total
   *      such that tokens which charge fees on transfer are accurately represented.
   * @param token The bridgeToken's address
   * @param amount The amount of bridgeToken's to transfer to the bridgebank
   * @return The balance that was transfered as reported by the getBalance command
   */
  function transferBalance(address token, uint256 amount) private returns (uint256) {
    //The interface of the ERC20 token to interact with
    IERC20 tokenToTransfer = IERC20(token);

    //The balance before any transfers take place
    uint256 oldBalance = getBalance(token);

    // locking the tokens by transfering them from the user to the bridgebank
    tokenToTransfer.safeTransferFrom(msg.sender, address(this), amount);

    //Fetch the updated balance reported after the transfer
    uint256 newBalance = getBalance(token);

    //Calculate the total amount transfered from the newbalance vs the old balance
    //Since this contract uses solidity 0.8+ overflows from bad acting tokens should
    //revert.
    uint256 transferedAmount = newBalance - oldBalance;

    return transferedAmount;
  }

  /**
   * @dev Fetches the denom of a token by address
   * @param token The bridgeTokens's address
   * @return The bridgeTokens's denom or ''
   */
  function getDenom(address token) private returns (string memory) {
    if (token == rowanTokenAddress) {
      // If it's the old erowan token, set the denom to 'rowan' and move forward
      return "rowan";
    }

    string memory denom = contractDenom[token];

    // check to see if we already have this denom stored in the smart contract
    if (keccak256(abi.encodePacked(denom)) != keccak256(abi.encodePacked(""))) {
      return denom;
    }

    try BridgeToken(token).cosmosDenom() returns (string memory _denom) {
      denom = _denom;
      contractDenom[token] = _denom;
    } catch {
      denom = "";
    }

    return denom;
  }

  /**
   * @notice Locks `amount` `token` tokens for `recipient`
   * @dev Locks received Ethereum/ERC20 funds
   * @param recipient Bytes representation of destination address
   * @param token Token address in origin chain (0x0 if ethereum)
   * @param amount Value of deposit
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
   * @notice Locks or burns multiple tokens in the bridge contract in a single tx
   * @param recipient Bytes representation of destination address
   * @param token Token address
   * @param amount Value of deposit
   * @param isBurn Should the tokens be burned?
   */
  function multiLockBurn(
    bytes[] calldata recipient,
    address[] calldata token,
    uint256[] calldata amount,
    bool[] calldata isBurn
  ) external whenNotPaused {
    uint256 recipientLength = recipient.length;
    uint256 tokenLength = token.length;
    // all array inputs must be of the same length
    // else throw malformed params error
    require(recipientLength == tokenLength, "M_P");
    require(tokenLength == amount.length, "M_P");
    require(tokenLength == isBurn.length, "M_P");


    // lockBurnNonce contains the previous nonce that was
    // sent in the LogLock/LogBurn, so the first one we send
    // should be lockBurnNonce + 1
    uint256 startingLockBurnNonce = lockBurnNonce + 1;

    // This is equivalent of lockBurnNonce = lockBurnNonce + recipientLength,
    // but it avoids a read of storage
    lockBurnNonce = startingLockBurnNonce - 1 + recipientLength;

    for (uint256 i = 0; i < recipientLength;) {
      if (isBurn[i]) {
        _burnTokens(recipient[i], token[i], amount[i], startingLockBurnNonce + i);
      } else {
        _lockTokens(recipient[i], token[i], amount[i], startingLockBurnNonce + i);
      }
      unchecked { ++i; }
    }

    // If we get any reentrant calls from the _{burn,lock}Tokens functions,
    // make sure that lockBurnNonce is what we expect it to be.
    require(lockBurnNonce == startingLockBurnNonce - 1 + recipientLength, "M_P");
  }

  /**
   * @dev Locks a token in the bridge contract
   * @param recipient Bytes representation of destination address
   * @param tokenAddress Token address
   * @param tokenAmount Value of deposit
   * @param _lockBurnNonce A global nonce for locking an burning tokens
   */
  function _lockTokens(
    bytes calldata recipient,
    address tokenAddress,
    uint256 tokenAmount,
    uint256 _lockBurnNonce
  )
    private
    onlyTokenNotInCosmosWhiteList(tokenAddress)
    validSifAddress(recipient)
    onlyNotBlocklisted(msg.sender)
  {
    uint256 transferedAmount = transferBalance(tokenAddress, tokenAmount);
    require(transferedAmount > 0, "No Balance Transferred");

    // decimals defaults to 18 if call to decimals fails
    uint8 decimals = getDecimals(tokenAddress);

    // Get name and symbol
    string memory name = getName(tokenAddress);
    string memory symbol = getSymbol(tokenAddress);

    emit LogLock(
      msg.sender,
      recipient,
      tokenAddress,
      transferedAmount,
      _lockBurnNonce,
      decimals,
      symbol,
      name,
      networkDescriptor
    );
  }

  /**
   * @dev Burns a token
   * @param recipient Bytes representation of destination address
   * @param tokenAddress Token address
   * @param tokenAmount How much should be burned
   * @param _lockBurnNonce A global nonce for locking an burning tokens
   */
  function _burnTokens(
    bytes calldata recipient,
    address tokenAddress,
    uint256 tokenAmount,
    uint256 _lockBurnNonce
  )
    private
    onlyCosmosTokenWhiteList(tokenAddress)
    validSifAddress(recipient)
    onlyNotBlocklisted(msg.sender)
  {
    BridgeToken tokenToTransfer = BridgeToken(tokenAddress);

    // burn tokens
    tokenToTransfer.burnFrom(msg.sender, tokenAmount);

    string memory denom = getDenom(tokenAddress);

    // Explicitly check that the denom is not the empty string
    require(keccak256(abi.encodePacked(denom)) != keccak256(abi.encodePacked("")), "INV_DENOM");

    // decimals defaults to 18 if call to decimals fails
    uint8 decimals = getDecimals(tokenAddress);

    emit LogBurn(
      msg.sender,
      recipient,
      tokenAddress,
      tokenAmount,
      _lockBurnNonce,
      decimals,
      networkDescriptor,
      denom
    );
  }

  /**
   * @dev Locks received EVM native tokens
   * @param recipient Bytes representation of destination address
   * @param amount Value of deposit
   */
  function handleNativeCurrencyLock(bytes calldata recipient, uint256 amount)
    internal
    validSifAddress(recipient)
    onlyNotBlocklisted(msg.sender)
  {
    require(msg.value == amount, "amount mismatch");

    address token = address(0);

    lockBurnNonce = lockBurnNonce + 1;

    emit LogLock(
      msg.sender,
      recipient,
      token,
      amount,
      lockBurnNonce,
      18, // decimals
      "EVMNATIVE", // symbol
      "EVMNATIVE", // name
      networkDescriptor
    );
  }

  /**
   * @dev Unlocks native or ERC20 tokens
   * @param recipient Recipient's Ethereum address
   * @param token Token contract address
   * @param amount Wei amount or ERC20 token count
   */
  function unlock(
    address payable recipient,
    address token,
    uint256 amount
  ) internal {
    // Transfer funds to intended recipient
    if (token == address(0)) {
      bool success = recipient.send(amount);
      require(success, "error sending ether");
    } else {
      IERC20 tokenToTransfer = IERC20(token);

      tokenToTransfer.safeTransfer(recipient, amount);
    }

    emit LogUnlock(recipient, token, amount);
  }

  /**
   * @notice Changes the denom of `_token` to `_denom`
   * @dev Will set the denom both in this contract AND in the token itself
   * @param _token Address of the BridgeToken
   * @param _denom The Cosmos denom to be applied
   * @return true if the operation succeeded
   */
  function setBridgeTokenDenom(address _token, string calldata _denom)
    external
    onlyOwner
    returns (bool)
  {
    return _setBridgeTokenDenom(_token, _denom);
  }

  /**
   * @notice Changes the denom of many tokens in a batch
   * @dev Will set the denom both in this contract AND in each token
   * @param _tokens List of address of BridgeTokens
   * @param _denoms List of Cosmos denoms to be applied
   * @return true if the operation succeeded
   */
  function batchSetBridgeTokenDenom(address[] calldata _tokens, string[] calldata _denoms)
    external
    onlyOwner
    returns (bool)
  {
    uint256 tokenLength = _tokens.length;
    require(tokenLength == _denoms.length, "INV_LEN");

    for (uint256 i ; i < tokenLength;) {
      _setBridgeTokenDenom(_tokens[i], _denoms[i]);
      unchecked { ++i; }
    }

    return true;
  }

  /**
   * @dev Changes the denom of `_token` to `_denom`
   * @dev Will set the denom both in this contract AND in the token itself
   * @param _token Address of the BridgeToken
   * @param _denom The Cosmos denom to be applied
   * @return true if the operation succeeded
   */
  function _setBridgeTokenDenom(address _token, string calldata _denom) private returns (bool) {
    contractDenom[_token] = _denom;
    return BridgeToken(_token).setDenom(_denom);
  }

  /**
   * @notice Sets in this contract the denom of `_token`
   * @dev Will fetch the denom from `_token` and register it in this contract
   * @dev Anyone may call this function
   * @param _token Address of the BridgeToken
   * @return true if the operation succeeded
   */
  function forceSetBridgeTokenDenom(address _token) external returns (bool) {
    return _forceSetBridgeTokenDenom(_token);
  }

  /**
   * @notice Sets in this contract the denom of a list of BridgeTokens     * @dev Will fetch the denom from each token and register it in this contract
   * @dev Will fetch the denom from each token and register it in this contract
   * @dev Anyone may call this function
   * @param _tokens List of address of BridgeTokens
   * @return true if the operation succeeded
   */
  function batchForceSetBridgeTokenDenom(address[] calldata _tokens) external returns (bool) {
    uint256 tokenLength = _tokens.length;
    for (uint256 i = 0; i < tokenLength;) {
      _forceSetBridgeTokenDenom(_tokens[i]);
      unchecked { ++i; }
    }
    return true;
  }

  /**
   * @dev Sets in this contract the denom of `_token`
   * @dev Will fetch the denom from `_token` and register it in this contract
   * @param _token Address of the BridgeToken
   * @return true if the operation succeeded
   */
  function _forceSetBridgeTokenDenom(address _token)
    private
    onlyCosmosTokenWhiteList(_token)
    returns (bool)
  {
    contractDenom[_token] = BridgeToken(_token).cosmosDenom();

    return true;
  }

  /**
   * @notice Lets the operator set the blocklist address
   * @param blocklistAddress The address of the blocklist contract
   */
  function setBlocklist(address blocklistAddress) public onlyOperator {
    blocklist = IBlocklist(blocklistAddress);
    hasBlocklist = true;
  }
}
