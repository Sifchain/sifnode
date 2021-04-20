pragma solidity 0.8.0;

contract EthereumBankStorage {

    /**
    * @notice current lock and or burn nonce
    */
    uint256 public lockBurnNonce;

    /**
    * @notice how much funds we have stored of a certain token
    */
    mapping(address => uint256) public lockedFunds;

    /**
    * @notice map the token symbol to the token address
    */
    mapping(string => address) public lockedTokenList;

    /*
    * @notice gap of storage for future upgrades
    */
    uint256[100] private ____gap;
}