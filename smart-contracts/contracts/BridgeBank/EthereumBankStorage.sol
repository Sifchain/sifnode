pragma solidity 0.8.0;

contract EthereumBankStorage {

    /**
    * @notice current lock and or burn nonce
    */
    uint256 public lockBurnNonce;

    /*
    * @notice {DEPRECATED}
    */
    mapping(address => uint256) private lockedFunds;

    /*
    * @notice {DEPRECATED}
    */
    mapping(string => address) private lockedTokenList;

    /**
    * @notice contract decimals based off of contract address
    */
    mapping (address => uint8) public contractDecimals;

    /**
    * @notice contract symbol based off of address
    */
    mapping (address => string) public contractSymbol;

    /**
    * @notice contract name based off of address
    */
    mapping (address => string) public contractName;

    /*
    * @notice gap of storage for future upgrades
    */
    uint256[97] private ____gap;

    /*
     * @dev: Event declarations
     */
    event LogBurn(
        address _from,
        bytes _to,
        address _token,
        uint256 _value,
        uint256 _nonce,
        uint256 _decimals,
        uint256 _chainId
    );

    event LogLock(
        address _from,
        bytes _to,
        address _token,
        uint256 _value,
        uint256 _nonce,
        uint256 _decimals,
        string _symbol,
        string _name,
        uint256 _chainId
    );

    event LogUnlock(
        address _to,
        address _token,
        uint256 _value
    );
}
