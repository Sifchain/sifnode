pragma solidity ^0.5.0;

import "./BridgeBank/CosmosBankStorage.sol";
import "./BridgeBank/EthereumBankStorage.sol";

contract CosmosBridgeStorage {
    /**
    * @notice gap of storage for future upgrades
    */
    string COSMOS_NATIVE_ASSET_PREFIX;

    /*
     * @dev: Public variable declarations
     */
    address public operator;
    
    /**
    * @notice gap of storage for future upgrades
    */
    address payable public valset;
    
    /**
    * @notice gap of storage for future upgrades
    */
    address payable public oracle;
    
    /**
    * @notice gap of storage for future upgrades
    */
    bool public hasOracle;
    
    /**
    * @notice gap of storage for future upgrades
    */
    address payable public bridgeBank;
    
    /**
    * @notice gap of storage for future upgrades
    */
    bool public hasBridgeBank;

    /**
    * @notice gap of storage for future upgrades
    */
    uint256 public prophecyClaimCount;

    /**
    * @notice Tracks the number of OracleClaims made on an individual BridgeClaim
    */
    mapping(uint256 => bool) public usedNonce;

    /**
    * @notice gap of storage for future upgrades
    */
    mapping(uint256 => ProphecyClaim) public prophecyClaims;

    /**
    * @notice prophecy status enum
    */
    enum Status {Null, Pending, Success, Failed}

    /**
    * @notice claim type enum
    */
    enum ClaimType {Unsupported, Burn, Lock}

    /**
    * @notice Prophecy claim struct
    */
    struct ProphecyClaim {
        ClaimType claimType;
        bytes cosmosSender;
        address payable ethereumReceiver;
        address originalValidator;
        address tokenAddress;
        string symbol;
        uint256 amount;
        Status status;
    }

    /**
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}