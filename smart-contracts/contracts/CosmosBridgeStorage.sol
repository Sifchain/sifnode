pragma solidity 0.6.6;

import "./BridgeBank/CosmosBankStorage.sol";
import "./BridgeBank/EthereumBankStorage.sol";

contract CosmosBridgeStorage {
    /**
    * @notice gap of storage for future upgrades
    */
    string COSMOS_NATIVE_ASSET_PREFIX;
    
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
    address payable public bridgeBank;
    
    /**
    * @notice gap of storage for future upgrades
    */
    bool public hasBridgeBank;

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
        address payable ethereumReceiver;
        string symbol;
        uint256 amount;
    }

    /**
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}