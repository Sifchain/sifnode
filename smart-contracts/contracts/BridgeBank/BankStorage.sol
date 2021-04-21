pragma solidity 0.8.0;

import "./EthereumBankStorage.sol";
import "./CosmosWhiteListStorage.sol";

contract BankStorage is 
    EthereumBankStorage,
    CosmosWhiteListStorage {

    /**
    * @notice address of the Cosmos Bridge smart contract
    */
    address public cosmosBridge;

    /**
    * @notice owner address that can use the admin API
    */
    address public owner;

    /*
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}