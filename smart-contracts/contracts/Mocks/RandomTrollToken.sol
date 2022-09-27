// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/**
 * @dev This token will report a different symbol, name, decimal ammount, user balance, and total supply for every block.
 */
contract RandomTrollToken is ERC20 {
    string[15] private symbols = [
        "USDT",
        "BNB",
        "USDC",
        "HEX",
        "SHIB",
        "BUSD",
        "MATIC",
        "CRO",
        "WBTC",
        "UST",
        "DAI",
        "LINK",
        "TRX",
        "LEO",
        "OKB"
    ];

    string[15] private names = [
        "Tether USD",
        "BNB",
        "USD Coin",
        "HEX",
        "SHIBA INU",
        "Binance USD",
        "Matic Token",
        "Crypto.com Coin",
        "Wrapped BTC",
        "Wrapped UST Token",
        "Dai Stablecoin",
        "Chainlink Token",
        "Tron",
        "Bitfinex LEO Token",
        "OKB"
    ];

    /**
     * @dev This constructor will prefund the inital accounts
     * @param initialAccounts an array of addresses that should be funded
     * @param quantity an array of initial balances that should match each associated account address
     */
    constructor(address[] memory initialAccounts, uint256[] memory quantity) ERC20("Random Troll Token", "RTT") {
        assert(names.length == symbols.length);
        require(initialAccounts.length == quantity.length, "Accounts and Quantities must be same length");
        for (uint256 i=0; i<initialAccounts.length; ++i) {
            _mint(initialAccounts[i], quantity[i]);
        }
    }

    /**
     * @dev A helper function that will report the same values for all function calls processed on the same block, but
     *      create different values on different block heights 
     * @param nonce A nonce that should be incremented for each subsequent call that should have a different number, 
     *              reuse the same nonce to get the same output on the current block
     * @param value An additional field to get different output associated with the same nonce for different inputs
     *              i.e. if you want the result to be unique for this function but based upon some data such as account
     *              balance you use this field so that a balance matching a prior nonce is still a unique number.
     */
    function _getCurrentBlockNumber(uint256 nonce, uint256 value) private view returns (uint256) {
        return uint256(
            keccak256(
                abi.encodePacked(
                    block.number, 
                    nonce, 
                    value
                )
            )
        );
    }

    /**
     * @dev A helper function that uses the values from _getCurrentBlockNumber to produce the index of the symbols
     *      and names arrays to return to mimic a different token each time.
     */
    function _getIndex() private view returns (uint256) {
        return _getCurrentBlockNumber(0, 0) % symbols.length;
    }

    /**
     * @dev Overriding the symbol function to return a different symbol from a set of valid symbols of other tokens
     */
    function symbol() override public view returns (string memory) {
        return symbols[_getIndex()];
    }

    /**
     * @dev Overriding the name function to return a different name from a set of valid names of other tokens. The
     *      name uses the same nonce as the symbol so the name and symbol will match.
     */
    function name() override public view returns (string memory) {
        return names[_getIndex()];
    }

    /**
     * @dev Overriding the decimal function to produce a different decimal value for each block from 0~255
     */
    function decimals() override public view returns (uint8) {
        return uint8(_getCurrentBlockNumber(1,0));
    }

    /**
     * @dev Overriding the totalSupply function to return a different totalSuppy value for each block
     */
    function totalSupply() override public view returns (uint256) {
        return _getCurrentBlockNumber(2,0);
    }

    /**
     * @dev Overriding the balanceOf function to provide a different balance for each users account. If the user has no 
     *      balance it reports a balance of 0 otherwise it generates the balance from a hash of the blocknumber, nonce of 3, and
     *      the users balance.
     */
    function balanceOf(address account) override public view returns (uint256) {
        uint256 balance = super.balanceOf(account);
        if (balance > 0) {
            return _getCurrentBlockNumber(3, balance);
        } 
        return balance;
    }
}