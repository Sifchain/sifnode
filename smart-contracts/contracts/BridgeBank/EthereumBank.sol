pragma solidity 0.5.16;

import "./BridgeToken.sol";
import "./EthereumBankStorage.sol";
import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
/*
 *  @title: EthereumBank
 *  @dev: Ethereum bank which locks Ethereum/ERC20 token deposits, and unlocks
 *        Ethereum/ERC20 tokens once the prophecy has been successfully processed.
 */
contract EthereumBank is EthereumBankStorage {
    using SafeMath for uint256;
    using SafeERC20 for IERC20;

    /*
     * @dev: Event declarations
     */
    event LogBurn(
        address _from,
        bytes _to,
        address _token,
        string _symbol,
        uint256 _value,
        uint256 _nonce
    );

    event LogLock(
        address _from,
        bytes _to,
        address _token,
        string _symbol,
        uint256 _value,
        uint256 _nonce
    );

    event LogUnlock(
        address _to,
        address _token,
        string _symbol,
        uint256 _value
    );

    /*
     * @dev: Gets the contract address of locked tokens by symbol.
     *
     * @param _symbol: The asset's symbol.
     */
    function getLockedTokenAddress(string memory _symbol)
        public
        view
        returns (address)
    {
        return lockedTokenList[_symbol];
    }

    /*
     * @dev: Gets the amount of locked tokens by symbol.
     *
     * @param _symbol: The asset's symbol.
     */
    function getLockedFunds(string memory _symbol)
        public
        view
        returns (uint256)
    {
        return lockedFunds[lockedTokenList[_symbol]];
    }

    /*
     * @dev: Creates a new Ethereum deposit with a unique id.
     *
     * @param _sender: The sender's ethereum address.
     * @param _recipient: The intended recipient's cosmos address.
     * @param _token: The currency type, either erc20 or ethereum.
     * @param _amount: The amount of erc20 tokens/ ethereum (in wei) to be itemized.
     */
    function burnFunds(
        address payable _sender,
        bytes memory _recipient,
        address _token,
        string memory _symbol,
        uint256 _amount
    ) internal {
        lockBurnNonce = lockBurnNonce.add(1);
        emit LogBurn(_sender, _recipient, _token, _symbol, _amount, lockBurnNonce);
    }

    /*
     * @dev: Creates a new Ethereum deposit with a unique id.
     *
     * @param _sender: The sender's ethereum address.
     * @param _recipient: The intended recipient's cosmos address.
     * @param _token: The currency type, either erc20 or ethereum.
     * @param _amount: The amount of erc20 tokens/ ethereum (in wei) to be itemized.
     */
    function lockFunds(
        address payable _sender,
        bytes memory _recipient,
        address _token,
        string memory _symbol,
        uint256 _amount
    ) internal {
        lockBurnNonce = lockBurnNonce.add(1);

        // Increment locked funds by the amount of tokens to be locked
        lockedTokenList[_symbol] = _token;

        emit LogLock(_sender, _recipient, _token, _symbol, _amount, lockBurnNonce);
    }

    /*
     * @dev: Unlocks funds held on contract and sends them to the
     *       intended recipient
     *
     * @param _recipient: recipient's Ethereum address
     * @param _token: token contract address
     * @param _symbol: token symbol
     * @param _amount: wei amount or ERC20 token count
     */
    function unlockFunds(
        address payable _recipient,
        address _token,
        string memory _symbol,
        uint256 _amount
    ) internal {
        // Transfer funds to intended recipient
        if (_token == address(0)) {
            (bool success,) = _recipient.call.value(_amount).gas(60000)("");
            require(success, "error sending ether");
        } else {
            IERC20 tokenToTransfer = IERC20(_token);
            tokenToTransfer.safeTransfer(_recipient, _amount);
        }

        emit LogUnlock(_recipient, _token, _symbol, _amount);
    }
}
