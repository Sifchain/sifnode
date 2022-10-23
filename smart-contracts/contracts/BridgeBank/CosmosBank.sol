// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "./BridgeToken.sol";
import "./CosmosBankStorage.sol";

/**
 * @title Cosmos Bank
 * @dev Manages the deployment and minting of ERC20 compatible BridgeTokens
 *      which represent assets based on the Cosmos blockchain.
 */
contract CosmosBank is CosmosBankStorage {
  /**
   * @dev Event emitted when a new BridgeToken is deployed
   */
  event LogNewBridgeToken(
    address indexed _token,
    string indexed _symbol,
    string indexed _cosmosDenom
  );

  /**
   * @dev Event emitted when a mint happens
   */
  event LogBridgeTokenMint(address _token, uint256 _amount, address _beneficiary);

  /**
   * @dev Deploys a new BridgeToken contract
   * @param _name The BridgeToken's name
   * @param _symbol The BridgeToken's symbol
   * @param _decimals The BridgeToken's decimals
   * @param _cosmosDenom The BridgeToken's Cosmos denom
   * @return The address of the newly deployed token
   */
  function deployNewBridgeToken(
    string memory _name,
    string memory _symbol,
    uint8 _decimals,
    string memory _cosmosDenom
  ) internal returns (address) {
    // Deploy new bridge token contract
    BridgeToken newBridgeToken = new BridgeToken(_name, _symbol, _decimals, _cosmosDenom);

    // Set address in tokens mapping
    address newBridgeTokenAddress = address(newBridgeToken);

    emit LogNewBridgeToken(newBridgeTokenAddress, _symbol, _cosmosDenom);
    return newBridgeTokenAddress;
  }

  /**
   * @dev Mints new Cosmos tokens
   * @param _intendedRecipient The intended recipient's Ethereum address.
   * @param _bridgeTokenAddress The currency type
   * @param _amount Number of Cosmos tokens to be minted
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

    emit LogBridgeTokenMint(_bridgeTokenAddress, _amount, _intendedRecipient);
  }
}
