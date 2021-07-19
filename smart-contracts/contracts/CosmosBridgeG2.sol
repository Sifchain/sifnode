pragma solidity 0.8.0;

import "./Oracle.sol";
import "./BridgeBank/BridgeBank.sol";
import "./CosmosBridgeStorage.sol";
import "hardhat/console.sol";

contract CosmosBridge is CosmosBridgeStorage, Oracle {
    bool private _initialized;
    uint256[100] private ___gap;

    /*
     * @dev: Event declarations
     */
    event LogBridgeBankSet(address bridgeBank);

    event LogNewProphecyClaim(
        uint256 indexed prophecyID,
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

    event LogProphecyCompleted(uint256 prophecyID, bool success);

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
        bytes calldata cosmosSender,
        uint256 cosmosSenderSequence,
        address payable ethereumReceiver,
        address tokenAddress,
        uint256 amount,
        bool doublePeg,
        uint128 nonce
    ) public pure returns (uint256) {
        return uint256(
            keccak256(
                abi.encode(
                    cosmosSender,
                    cosmosSenderSequence,
                    ethereumReceiver,
                    tokenAddress,
                    amount,
                    doublePeg,
                    nonce
                )
            )
        );
    }

    function verifySignature(
        address signer,
        bytes32 hashDigest,
        uint8 _v,
		bytes32 _r,
		bytes32 _s
    ) private pure returns (bool) {
		bytes32 messageDigest = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", hashDigest));
		return signer == ecrecover(messageDigest, _v, _r, _s);
	}
    
    // this is unfortunately the best we can do to ensure no duplicate validators are calling
    // it is possible to build a hashmap in memory, but I'm unsure of how much that saves and
    // it would require some pretty low level work for this very simple function
    // Alternatively, cast addresses to UINT's and possibly do some bitwise operations
    // to ensure there are no duplicate numbers
    function findDup(SignatureData[] calldata validators) public pure returns (bool) {
        for (uint256 i = 0; i < validators.length; i++) {
            for (uint256 j = i + 1; j < validators.length; j++) {
                if (validators[i].signer == validators[j].signer) {
                    return true;
                }
            }
        }

        return false;
    }

    function getSignedPower(SignatureData[] calldata validators) public view returns(uint256) {
        uint256 totalPower = 0;
        for (uint256 i = 0; i < validators.length; i++) {
            totalPower += getValidatorPower(validators[i].signer);
        }

        return totalPower;
    }

    struct SignatureData {
        address signer;
        uint8 _v;
		bytes32 _r;
		bytes32 _s;
    }

    struct ClaimData {
        bytes cosmosSender;
        uint256 cosmosSenderSequence;
        address payable ethereumReceiver;
        address tokenAddress;
        uint256 amount;
        bool doublePeg;
        uint128 nonce;
    }

    function batchSubmitProphecyClaimAggregatedSigs(
        bytes32[] calldata sigs,
        ClaimData[] calldata claims,
        SignatureData[][] calldata signatureData
    ) external {
        require(sigs.length == claims.length, "INV_CLM_LEN");
        require(sigs.length == signatureData.length, "INV_SIG_LEN");

        for (uint256 i = 0; i < sigs.length; i++) {
            _submitProphecyClaimAggregatedSigs(sigs[i], claims[i], signatureData[i]);
        }
    }

    function submitProphecyClaimAggregatedSigs(
        bytes32 hashDigest,
        ClaimData calldata claimData,
        SignatureData[] calldata signatureData
    ) external {
        _submitProphecyClaimAggregatedSigs(hashDigest, claimData, signatureData);
    }

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
            claimData.doublePeg,
            claimData.nonce
        );

        require(
            uint256(hashDigest) == prophecyID,
            "INV_DATA"
        );

        // ensure signature lengths are correct
        require(
            signatureData.length > 0 && signatureData.length <= validatorCount,
            "INV_SIG_LEN"
        );

        // ensure there are no duplicate signers
        require(
            !findDup(signatureData), "DUP_SIGNER"
        );

        // check that all signers are validators and are unique
        for (uint256 i = 0; i < signatureData.length; i++) {
            require(isActiveValidator(signatureData[i].signer), "INV_SIGNER");
            require(
                verifySignature(
                    signatureData[i].signer,
                    hashDigest,
                    signatureData[i]._v,
                    signatureData[i]._r,
                    signatureData[i]._s
                ) == true,
                "INV_SIG"
            );
        }

        uint256 signedPower = getSignedPower(signatureData);

        require(getProphecyStatus(signedPower), "INV_POW");

        uint256 previousNonce = lastNonceSubmitted;
        require(
            // assert nonce is correct
            previousNonce + 1 == claimData.nonce,
            "INV_ORD"
        );
        lastNonceSubmitted = claimData.nonce;

        emit LogNewProphecyClaim(
            prophecyID,
            claimData.ethereumReceiver,
            claimData.amount
        );

        // if we are double pegging, then we are going to need to get the address on this chain
        address tokenAddress = claimData.doublePeg ? sourceAddressToDestinationAddress[claimData.tokenAddress] : claimData.tokenAddress;

        completeProphecyClaim(
            prophecyID,
            claimData.ethereumReceiver,
            tokenAddress,
            claimData.amount
        );
    }

    /**
     * @param symbol symbol of the ERC20 token on the source chain
     * @param name name of the ERC20 token on the source chain
     * @param sourceChainTokenAddress address of the ERC20 token on the source chain
     * @param decimals of the ERC20 token on the source chain
     * @param chainDescriptor descriptor of the source chain
     */
    function createNewBridgeToken(
        string calldata symbol,
        string calldata name,
        address sourceChainTokenAddress,
        uint8 decimals,
        uint256 chainDescriptor
    ) external onlyValidator {
        require(
            sourceAddressToDestinationAddress[sourceChainTokenAddress] == address(0),
            "INV_SRC_ADDR"
        );
        // need to make a business decision on what this symbol should be
        // First lock of this asset, deploy new contract and get new symbol/token address
        address tokenAddress = BridgeBank(bridgeBank)
            .createNewBridgeToken(
                name,
                symbol,
                decimals
            );

        sourceAddressToDestinationAddress[sourceChainTokenAddress] = tokenAddress;

        emit LogNewBridgeTokenCreated(
            decimals,
            chainDescriptor,
            name,
            symbol,
            sourceChainTokenAddress,
            tokenAddress
        );
    }

    /*
     * @dev: completeProphecyClaim
     *       Allows for the completion of ProphecyClaims once processed by the Oracle.
     *       Burn claims unlock tokens stored by BridgeBank.
     *       Lock claims mint BridgeTokens on BridgeBank's token whitelist.
     */
    function completeProphecyClaim(
        uint256 prophecyID,
        address payable ethereumReceiver,
        address tokenAddress,
        uint256 amount
    ) internal {
        (bool success, ) = bridgeBank.call{gas: 120000}(
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
