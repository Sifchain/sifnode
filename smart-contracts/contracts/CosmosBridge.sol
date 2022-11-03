// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "./Oracle.sol";
import "./BridgeBank/BridgeBank.sol";
import "./CosmosBridgeStorage.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

error OutOfOrderSigner(uint256 index);
error DuplicateSigner(uint256 index, address signer);

/**
 * @title Cosmos Bridge
 * @dev Processes Prophecy Claims and communicates with the
 *      BridgeBank contract to deploy, mint or unlock BridgeTokens.
 */
contract CosmosBridge is CosmosBridgeStorage, Oracle {
  /**
   * @dev has the contract been initialized?
   */
  bool private _initialized;

  /**
   * @dev gap of storage for future upgrades
   */
  uint256[100] private ___gap;

  /**
   * @notice Maps the cosmos denom to its bridge token address
   */
  mapping(string => address) public cosmosDenomToDestinationAddress;

  /**
   * @notice network descriptor
   */
  int32 public networkDescriptor;

  /**
   * @notice mapping of validator address to last nonce submitted
   */
  uint256 public lastNonceSubmitted;

  /**
   * @dev Event emitted when BridgeBank's address has been set
   */
  event LogBridgeBankSet(address bridgeBank);

  /**
   * @dev Event emitted when a ProphecyClaim has been accepted
   */
  event LogNewProphecyClaim(
    uint256 indexed prophecyID,
    address indexed ethereumReceiver,
    uint256 indexed amount
  );

  /**
   * @dev Event emitted when a new BridgeToken has been created
   */
  event LogNewBridgeTokenCreated(
    uint8 decimals,
    int32 indexed _networkDescriptor,
    string name,
    string symbol,
    address indexed sourceContractAddress,
    address indexed bridgeTokenAddress,
    string cosmosDenom
  );

  /**
   * @dev Event emitted when a ProphecyClaim has been completed
   */
  event LogProphecyCompleted(uint256 prophecyID, bool success);

 /**
  * @dev Event emitted when operator account is changed
  */
  event OperatorAccountChange(address newOperator);

  /**
   * @notice Initializer
   * @param _operator Address of the operator
   * @param _consensusThreshold Minimum required power for a valid prophecy
   * @param _initValidators List of initial validators
   * @param _initPowers List of numbers representing the power of each validator in the above list
   * @param _networkDescriptor Unique identifier of the network that this contract cares about
   */
  function initialize(
    address _operator,
    uint256 _consensusThreshold,
    address[] calldata _initValidators,
    uint256[] calldata _initPowers,
    int32 _networkDescriptor
  ) external {
    require(!_initialized, "Initialized");

    operator = _operator;
    networkDescriptor = _networkDescriptor;
    hasBridgeBank = false;
    _initialized = true;
    Oracle._initialize(_operator, _consensusThreshold, _initValidators, _initPowers);
  }

  /**
   * @notice Transfers the operator role to `_newOperator`
   * @dev Cannot transfer role to the zero address
   * @param _newOperator: the new operator's address
   */
  function changeOperator(address _newOperator) external onlyOperator {
    require(_newOperator != address(0), "invalid address");
    operator = _newOperator;
    emit OperatorAccountChange(_newOperator);
  }

  /**
   * @notice Sets the brigeBank address to `_bridgeBank`
   * @param _bridgeBank The address of BridgeBank
   */
  function setBridgeBank(address payable _bridgeBank) external onlyOperator {
    require(!hasBridgeBank, "The Bridge Bank cannot be updated once it has been set");

    hasBridgeBank = true;
    bridgeBank = _bridgeBank;

    emit LogBridgeBankSet(bridgeBank);
  }

  /**
   * @notice Calculates the ID of a Prophecy based on its properties
   * @param cosmosSender Address of the Cosmos account sending this prophecy
   * @param cosmosSenderSequence Nonce of the Cosmos account sending this prophecy
   * @param ethereumReceiver Destination address
   * @param tokenAddress Original address
   * @param amount token transferred in this prophecy
   * @param tokenName token name in bridge token contract
   * @param tokenSymbol token symbol in bridge token contract
   * @param tokenDecimals token decimal in bridge token contract
   * @param _networkDescriptor Unique identifier of the network
   * @param bridgeToken if the token created by cosmos bridge
   * @param nonce lock burn sequence recorded in sifnode side
   * @param denom token identity in sifnode bank system
   * @return A hash that uniquely identifies this Prophecy
   */

  function getProphecyID(
    bytes memory cosmosSender,
    uint256 cosmosSenderSequence,
    address payable ethereumReceiver,
    address tokenAddress,
    uint256 amount,
    string memory tokenName,
    string memory tokenSymbol,
    uint8 tokenDecimals,
    int32 _networkDescriptor,
    bool bridgeToken,
    uint128 nonce,
    string memory denom
  ) public pure returns (uint256) {
    return
      uint256(
        keccak256(
          abi.encode(
            cosmosSender,
            cosmosSenderSequence,
            ethereumReceiver,
            tokenAddress,
            amount,
            tokenName,
            tokenSymbol,
            tokenDecimals,
            _networkDescriptor,
            bridgeToken,
            nonce,
            denom
          )
        )
      );
  }

  /**
   * @dev Guarantees that the signature is correct
   * @param signer Address that supposedly signed the message
   * @param hashDigest Hash of the message
   * @param _v The signature's recovery identifier
   * @param _r The signature's random value
   * @param _s The signature's proof
   * @return Boolean: has the message been signed by `signer`?
   */
  function verifySignature(
    address signer,
    bytes32 hashDigest,
    uint8 _v,
    bytes32 _r,
    bytes32 _s
  ) private pure returns (bool) {
    bytes32 messageDigest = keccak256(
      abi.encodePacked("\x19Ethereum Signed Message:\n32", hashDigest)
    );
    return signer == ECDSA.recover(messageDigest, _v, _r, _s);
  }

  /**
   * @dev Runs verifications on a ProphecyClaim
   * @dev Prevents duplicates signers, makes sure validators are active,
   * @dev Validates signatures and calculates the total validation power
   * @param _validators List of validators for this ProphecyClaim
   * @param hashDigest The message in this prophecy
   * @return pow : aggregated signing power of all validators
   */
  function getSignedPowerAndFindDup(SignatureData[] calldata _validators, bytes32 hashDigest)
    private
    view
    returns (uint256 pow)
  {
    uint256 validatorLength = _validators.length;
    for (uint256 i; i < validatorLength;) {
      SignatureData calldata validator = _validators[i];

      require(isActiveValidator(validator.signer), "INV_SIGNER");

      require(
        verifySignature(validator.signer, hashDigest, validator._v, validator._r, validator._s),
        "INV_SIG"
      );

      unchecked {
        pow += getValidatorPower(validator.signer);
      }

      // Signatures must be sorted on the relayer side, so we just
      // need to make sure that the next witness in the array
      // (if we're not at the end) isn't a duplicate and is
      // sorted correctly
      if (i + 1 <= validatorLength - 1) {
        if (validator.signer == _validators[i + 1].signer) {
          revert DuplicateSigner(i + 1, validator.signer);
        }
        if (validator.signer > _validators[i + 1].signer) {
          revert OutOfOrderSigner(i);
        }
      }

      unchecked { ++i; }
    }
  }

  /**
   * @dev Information on a signature: address, r, s, and v
   */
  struct SignatureData {
    address signer;
    uint8 _v;
    bytes32 _r;
    bytes32 _s;
  }

  /**
   * @dev Data structure of a claim
   */
  struct ClaimData {
    bytes cosmosSender;
    uint256 cosmosSenderSequence;
    address payable ethereumReceiver;
    address tokenAddress;
    uint256 amount;
    string tokenName;
    string tokenSymbol;
    uint8 tokenDecimals;
    int32 networkDescriptor;
    bool bridgeToken;
    uint128 nonce;
    string cosmosDenom;
  }

  /**
   * @notice Submits a list of ProphecyClaims in a batch
   * @dev All arrays must have the same length
   * @param sigs List of hashed messages
   * @param claims List of claims
   * @param signatureData List of signature data
   */
  function batchSubmitProphecyClaimAggregatedSigs(
    bytes32[] calldata sigs,
    ClaimData[] calldata claims,
    SignatureData[][] calldata signatureData
  ) external {
    uint256 sigsLength = sigs.length;
    uint256 claimLength = claims.length;
    require(sigsLength == claimLength, "INV_CLM_LEN");
    require(sigsLength == signatureData.length, "INV_SIG_LEN");

    uint256 intermediateNonce = lastNonceSubmitted;
    lastNonceSubmitted += claimLength;

    for (uint256 i; i < sigsLength;) {
      require(intermediateNonce + 1 + i == claims[i].nonce, "INV_ORD");
      _submitProphecyClaimAggregatedSigs(sigs[i], claims[i], signatureData[i]);
      unchecked { ++i; }
    }
  }

  /**
   * @notice Submits a ProphecyClaim
   * @param hashDigest The hashed message
   * @param claimData The claim itself
   * @param signatureData The signature data
   */
  function submitProphecyClaimAggregatedSigs(
    bytes32 hashDigest,
    ClaimData calldata claimData,
    SignatureData[] calldata signatureData
  ) external {
    uint256 previousNonce = lastNonceSubmitted;
    require(previousNonce + 1 == claimData.nonce, "INV_ORD");

    // update the nonce
    lastNonceSubmitted = claimData.nonce;

    _submitProphecyClaimAggregatedSigs(hashDigest, claimData, signatureData);
  }

  /**
   * @dev Submits a ProphecyClaim
   * @param hashDigest The hashed message
   * @param claimData The claim itself
   * @param signatureData The signature data
   */
  function _submitProphecyClaimAggregatedSigs(
    bytes32 hashDigest,
    ClaimData calldata claimData,
    SignatureData[] calldata signatureData
  ) private {
    uint256 prophecyID = getProphecyID(
      claimData.cosmosSender,
      claimData.cosmosSenderSequence,
      claimData.ethereumReceiver,
      claimData.tokenAddress,
      claimData.amount,
      claimData.tokenName,
      claimData.tokenSymbol,
      claimData.tokenDecimals,
      claimData.networkDescriptor,
      claimData.bridgeToken,
      claimData.nonce,
      claimData.cosmosDenom
    );

    require(uint256(hashDigest) == prophecyID, "INV_DATA");

    // ensure signature lengths are correct
    require(signatureData.length > 0 && signatureData.length <= validatorCount, "INV_SIG_LEN");

    // ensure the networkDescriptor matches
    if (!claimData.bridgeToken) {
      require(_verifyNetworkDescriptor(claimData.networkDescriptor), "INV_NET_DESC");
    }

    uint256 pow = getSignedPowerAndFindDup(signatureData, hashDigest);
    require(getProphecyStatus(pow), "INV_POW");

    address tokenAddress;

    // bridgeToken means the token from other EVM chain or ibc token
    // we should deploy bridge token for them automatically
    if (claimData.bridgeToken) {
      if (_isBridgeTokenCreated(claimData.cosmosDenom)) {
        tokenAddress = cosmosDenomToDestinationAddress[claimData.cosmosDenom];
      } else {
          tokenAddress = _createNewBridgeToken(
            claimData.tokenSymbol,
            claimData.tokenName,
            claimData.tokenAddress,
            claimData.tokenDecimals,
            claimData.networkDescriptor,
            claimData.cosmosDenom
          );
      }
    }
    else {
        tokenAddress = claimData.tokenAddress;
    }

    completeProphecyClaim(prophecyID, claimData.ethereumReceiver, tokenAddress, claimData.amount);

    emit LogNewProphecyClaim(prophecyID, claimData.ethereumReceiver, claimData.amount);
  }

  /**
   * @dev Verifies if `cosmosDenom` is a bridge token for the cosmos denom created
   * @param cosmosDenom The cosmos denom of the token
   * @return Boolean: is `cosmosDenom` is a bridge token for the cosmos denom created?
   */
  function _isBridgeTokenCreated(string calldata cosmosDenom) private view returns (bool) {
    return cosmosDenomToDestinationAddress[cosmosDenom] != address(0);
  }

  /**
   * @dev Verifies if `_networkDescriptor` matches this contract's networkDescriptor
   * @param _networkDescriptor Unique identifier of the network
   * @return Boolean: is `_networkDescriptor` what we expected?
   */
  function _verifyNetworkDescriptor(int32 _networkDescriptor) private view returns (bool) {
    return _networkDescriptor == networkDescriptor;
  }

  /**
   * @dev Deploys a new BridgeToken, delegating this responsibility to BridgeBank
   * @param symbol The symbol of the token
   * @param name The name of the token
   * @param sourceChainTokenAddress Address of the token on its original chain
   * @param decimals The number of decimals this token has
   * @param _networkDescriptor Unique identifier of the network
   * @param cosmosDenom The token's Cosmos denom
   * @return tokenAddress : The address of the newly deployed BridgeToken
   */
  function _createNewBridgeToken(
    string memory symbol,
    string memory name,
    address sourceChainTokenAddress,
    uint8 decimals,
    int32 _networkDescriptor,
    string calldata cosmosDenom
  ) internal returns (address tokenAddress) {
    require(
      cosmosDenomToDestinationAddress[cosmosDenom] == address(0),
      "INV_SRC_ADDR"
    );
    // need to make a business decision on what this symbol should be
    // First lock of this asset, deploy new contract and get new symbol/token address
    tokenAddress = BridgeBank(bridgeBank).createNewBridgeToken(
      name,
      symbol,
      decimals,
      cosmosDenom
    );

    cosmosDenomToDestinationAddress[cosmosDenom] = tokenAddress;

    emit LogNewBridgeTokenCreated(
      decimals,
      _networkDescriptor,
      name,
      symbol,
      sourceChainTokenAddress,
      tokenAddress,
      cosmosDenom
    );
  }

  /**
   * @dev completeProphecyClaim
   *      Allows for the completion of ProphecyClaims once processed by the Oracle.
   *      Burn claims unlock tokens stored by BridgeBank.
   *      Lock claims mint BridgeTokens on BridgeBank's token whitelist.
   * @param prophecyID The calculated prophecyID
   * @param ethereumReceiver The Recipient's address
   * @param tokenAddress The tokens address
   * @param amount How much should be transacted
   */
  function completeProphecyClaim(
    uint256 prophecyID,
    address payable ethereumReceiver,
    address tokenAddress,
    uint256 amount
  ) internal {
    (bool success, ) = bridgeBank.call{ gas: 120000 }(
      abi.encodeWithSignature(
        "handleUnpeg(address,address,uint256)",
        ethereumReceiver,
        tokenAddress,
        amount
      )
    );

    // prophecy completed and whether or not the call to bridgebank was successful
    emit LogProphecyCompleted(prophecyID, success);
  }
}
