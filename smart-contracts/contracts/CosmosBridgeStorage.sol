pragma solidity 0.8.0;

contract CosmosBridgeStorage {
    /*
    * @notice gap of storage for future upgrades
    * @dev {DEPRECATED}
    */
    string private COSMOS_NATIVE_ASSET_PREFIX;

    /*
     * @dev: Public variable declarations
     */
    address private _operator;

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


    mapping (address => address) public sourceAddressToDestinationAddress;

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

    /*
    * @notice gap of storage for future upgrades
    */
    uint256[99] private ____gap;
}