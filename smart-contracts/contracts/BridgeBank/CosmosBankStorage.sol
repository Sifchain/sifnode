pragma solidity 0.6.9;

contract CosmosBankStorage {

    /**
    * @notice Cosmos deposit struct
    */
    struct CosmosDeposit {
        bytes cosmosSender;
        address payable ethereumRecipient;
        address bridgeTokenAddress;
        uint256 amount;
        bool locked;
    }

    /**
    * @notice number of bridge tokens
    */
    uint256 public bridgeTokenCount;

    /**
    * @notice cosmos deposit nonce
    */
    uint256 public cosmosDepositNonce;
    
    /**
    * @notice mapping of symbols to token addresses
    */
    mapping(string => address) controlledBridgeTokens;
    
    /**
    * @notice mapping of lowercase symbols to properly capitalized tokens
    */
    mapping(string => string) public lowerToUpperTokens;

    /**
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}