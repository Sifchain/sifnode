pragma solidity 0.6.9;

import "./CosmosBankStorage.sol";
import "./EthereumBankStorage.sol";
import "./CosmosWhiteListStorage.sol";

contract BankStorage is 
    CosmosBankStorage,
    EthereumBankStorage,
    CosmosWhiteListStorage {

    /**
    * @notice operator address that can update the smart contract
    */
    address public operator;

    /**
    * @notice address of the Oracle smart contract
    */
    address public oracle;

    /**
    * @notice address of the Cosmos Bridge smart contract
    */
    address public cosmosBridge;

    /**
    * @notice owner address that can use the admin API
    */
    address public owner;

    mapping (string => uint256) public maxTokenAmount;

    /**
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}