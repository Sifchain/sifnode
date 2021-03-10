pragma solidity 0.6.9;

import "@openzeppelin/contracts/math/SafeMath.sol";
import "./BridgeToken.sol";
import "./CosmosBankStorage.sol";

/**
 * @title CosmosBank
 * @dev Manages the deployment and minting of ERC20 compatible BridgeTokens
 *      which represent assets based on the Cosmos blockchain.
 **/

contract CosmosBank is CosmosBankStorage {
    using SafeMath for uint256;

    /*
     * @dev: Event declarations
     */
    event LogNewBridgeToken(address _token, string _symbol);

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
    function deployNewBridgeToken(string memory _symbol)
        internal
        returns (address)
    {
        // Deploy new bridge token contract
        BridgeToken newBridgeToken = (new BridgeToken)(_symbol);

        // Set address in tokens mapping
        address newBridgeTokenAddress = address(newBridgeToken);

        emit LogNewBridgeToken(newBridgeTokenAddress, _symbol);
        return newBridgeTokenAddress;
    }

    /*
     * @dev: Mints new cosmos tokens
     *
     * @param _cosmosSender: The sender's Cosmos address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _cosmosTokenAddress: The currency type
     * @param _symbol: comsos token symbol
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
