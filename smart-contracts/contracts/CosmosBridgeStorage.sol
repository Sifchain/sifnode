// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.0;

/**
 * @title Cosmos Bridge Storage
 * @dev Stores the operator's address,
        BridgeBank's address,
        networkDescriptor,
        sourceAddressToDestinationAddress of a pegged token
 **/
contract CosmosBridgeStorage {
    /**
    * @dev {DEPRECATED}
    */
    string private COSMOS_NATIVE_ASSET_PREFIX;

    /**
     * @dev Public variable declarations
     */
    address private _operator;

    /**
    * @dev {DEPRECATED}
    */
    address payable private valset;

    /**
    * @dev {DEPRECATED}
    */
    address payable private oracle;

    /**
    * @notice Address of the BridgeBank contract
    */
    address payable public bridgeBank;
    
    /**
    * @notice Has the BridgeBank contract been registered yet?
    */
    bool public hasBridgeBank;

    /**
    * @dev {DEPRECATED}
    */
    mapping(uint256 => ProphecyClaim) private prophecyClaims;

    /**
    * @notice Maps the original address of a token to its address in another network
    */
    mapping (address => address) public sourceAddressToDestinationAddress;

    /**
    * @dev {DEPRECATED}
    */
    enum Status {Null, Pending, Success, Failed}

    /**
    * @dev {DEPRECATED}
    */
    enum ClaimType {Unsupported, Burn, Lock}

    /**
    * @notice {DEPRECATED}
    */
    struct ProphecyClaim {
        address payable ethereumReceiver;
        string symbol;
        uint256 amount;
    }

    /**
    * @notice network descriptor
    */
    uint256 public networkDescriptor;

    /**
    * @dev gap of storage for future upgrades
    */
    uint256[98] private ____gap;
}