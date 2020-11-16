pragma solidity ^0.5.0;

import "./CosmosBankStorage.sol";
import "./EthereumBankStorage.sol";
import "./WhiteListStorage.sol";

contract BankStorage is 
    CosmosBankStorage,
    EthereumBankStorage,
    WhiteListStorage {

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

    /**
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}