// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

import "./BridgeToken.sol";
import "./CosmosBankStorage.sol";

/**
 * @title CosmosBank
 * @dev Manages the deployment and minting of ERC20 compatible BridgeTokens
 *      which represent assets based on the Cosmos blockchain.
 **/

contract CosmosBank is CosmosBankStorage {
    /*
     * @dev: Event declarations
     */
    event LogNewBridgeToken(address indexed _token, string indexed _symbol);

    event LogBridgeTokenMint(
        address _token,
        uint256 _amount,
        address _beneficiary
    );

    /*
     * @dev: Deploys a new BridgeToken contract
     *
     * @param _symbol: The BridgeToken's symbol
     */
    function deployNewBridgeToken(
        string memory _name,
        string memory _symbol,
        uint8 _decimals
    )
        internal
        returns (address)
    {
        // Deploy new bridge token contract
        BridgeToken newBridgeToken = new BridgeToken(
            _name,
            _symbol,
            _decimals
        );

        // Set address in tokens mapping
        address newBridgeTokenAddress = address(newBridgeToken);

        emit LogNewBridgeToken(newBridgeTokenAddress, _symbol);
        return newBridgeTokenAddress;
    }

    /*
     * @dev: Mints new cosmos tokens
     *
     * @param _intendedRecipient: The intended recipient's Ethereum address.
     * @param _bridgeTokenAddress: The currency type
     * @param _amount: number of comsos tokens to be minted
     */
    function mintNewBridgeTokens(
        address _intendedRecipient,
        address _bridgeTokenAddress,
        uint256 _amount
    ) internal {
        // Mint bridge tokens
        require(
            BridgeToken(_bridgeTokenAddress).mint(_intendedRecipient, _amount),
            "Attempted mint of bridge tokens failed"
        );

        emit LogBridgeTokenMint(
            _bridgeTokenAddress,
            _amount,
            _intendedRecipient
        );
    }
}
