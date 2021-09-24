// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

/**
 * @title Cosmos Bank Storage
 * @dev Stores Cosmos deposits, nonces, networkDescriptor
 */
contract CosmosBankStorage {

    /**
    * @dev {DEPRECATED}
    */
    struct CosmosDeposit {
        bytes cosmosSender;
        address payable ethereumRecipient;
        address bridgeTokenAddress;
        uint256 amount;
        bool locked;
    }

    /**
    * @notice number of bridge tokens
    */
    uint256 public bridgeTokenCount;

    /**
    * @dev {DEPRECATED}
    */
    uint256 public cosmosDepositNonce;

    /**
    * @dev {DEPRECATED}
    */
    mapping(string => address) private controlledBridgeTokens;

    /**
    * @dev {DEPRECATED}
    */
    mapping(string => string) private lowerToUpperTokens;

    /**
    * @notice network descriptor
    */
    int32 public networkDescriptor;

    /**
    * @dev gap of storage for future upgrades
    */
    uint256[99] private ____gap;
}
