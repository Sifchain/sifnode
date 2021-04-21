pragma solidity 0.8.0;

contract CosmosBridgeStorage {
    
    /**
    * @notice gap of storage for future upgrades
    */
    address payable public bridgeBank;
    
    /**
    * @notice gap of storage for future upgrades
    */
    bool public hasBridgeBank;

    /**
    * @notice prophecy status enum
    */
    enum Status {Null, Pending, Success, Failed}

    /**
    * @notice claim type enum
    */
    enum ClaimType {Unsupported, Burn, Lock}

    /*
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}