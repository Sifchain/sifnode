pragma solidity 0.5.16;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./BridgeToken.sol";
import "./CosmosBankStorage.sol";
import "./ToLower.sol";

/**
 * @title CosmosBank
 * @dev Manages the deployment and minting of ERC20 compatible BridgeTokens
 *      which represent assets based on the Cosmos blockchain.
 **/

contract CosmosBank is CosmosBankStorage, ToLower {
    using SafeMath for uint256;

    /*
     * @dev: Event declarations
     */
    event LogNewBridgeToken(address _token, string _symbol);

    event LogBridgeTokenMint(
        address _token,
        string _symbol,
        uint256 _amount,
        address _beneficiary
    );

    /*
     * @dev: Get a token symbol's corresponding bridge token address.
     *
     * @param _symbol: The token's symbol/denom without 'e' prefix.
     * @return: Address associated with the given symbol. Returns address(0) if none is found.
     */
    function getBridgeToken(string memory _symbol)
        public
        view
        returns (address)
    {
        return (controlledBridgeTokens[_symbol]);
    }

    function safeLowerToUpperTokens(string memory _symbol)
        public
        view
        returns (string memory)
    {
        string memory retrievedSymbol = lowerToUpperTokens[_symbol];
        return keccak256(abi.encodePacked(retrievedSymbol)) == keccak256("") ? _symbol : retrievedSymbol;
    }

    /*
     * @dev: Deploys a new BridgeToken contract
     *
     * @param _symbol: The BridgeToken's symbol
     */
    function deployNewBridgeToken(string memory _symbol)
        internal
        returns (address)
    {
        bridgeTokenCount = bridgeTokenCount.add(1);

        // Deploy new bridge token contract
        BridgeToken newBridgeToken = (new BridgeToken)(_symbol);

        // Set address in tokens mapping
        address newBridgeTokenAddress = address(newBridgeToken);
        controlledBridgeTokens[_symbol] = newBridgeTokenAddress;
        lowerToUpperTokens[toLower(_symbol)] = _symbol;

        emit LogNewBridgeToken(newBridgeTokenAddress, _symbol);
        return newBridgeTokenAddress;
    }

    /*
     * @dev: Deploys a new BridgeToken contract
     *
     * @param _symbol: The BridgeToken's symbol
     *
     * @note the Rowan token symbol needs to be "Rowan" so that it integrates correctly with the cosmos bridge
     */
    function useExistingBridgeToken(address _contractAddress)
        internal
        returns (address)
    {
        bridgeTokenCount = bridgeTokenCount.add(1);

        string memory _symbol = BridgeToken(_contractAddress).symbol();
        // Set address in tokens mapping
        address newBridgeTokenAddress = _contractAddress;
        controlledBridgeTokens[_symbol] = newBridgeTokenAddress;
        lowerToUpperTokens[toLower(_symbol)] = _symbol;

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
        address payable _intendedRecipient,
        address _bridgeTokenAddress,
        string memory _symbol,
        uint256 _amount
    ) internal {
        require(
            controlledBridgeTokens[_symbol] == _bridgeTokenAddress,
            "Token must be a controlled bridge token"
        );

        // Mint bridge tokens
        require(
            BridgeToken(_bridgeTokenAddress).mint(_intendedRecipient, _amount),
            "Attempted mint of bridge tokens failed"
        );

        emit LogBridgeTokenMint(
            _bridgeTokenAddress,
            _symbol,
            _amount,
            _intendedRecipient
        );
    }
}
