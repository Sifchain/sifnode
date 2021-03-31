pragma solidity 0.5.16;

contract EthereumBankStorage {

    /**
    * @notice current lock and or burn nonce
    */
    uint256 public lockBurnNonce;

    /**
    * @notice This mapping has been deprecated to save gas when unpegging
    * @notice In the past, this was used to track how much funds we had of a certain token
    */
    mapping(address => uint256) public lockedFunds;

    /**
    * @notice map the token symbol to the token address
    */
    mapping(string => address) public lockedTokenList;

    /**
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}