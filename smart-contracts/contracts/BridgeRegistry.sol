pragma solidity 0.6.6;


contract BridgeRegistry {
    address public cosmosBridge;
    address public bridgeBank;

    // these variables are now deprecated and are made private
    // so that the getter helper method is not available.
    address private oracle;
    address private valset;

    bool private _initialized;

    uint256[100] private ____gap;

    event LogContractsRegistered(
        address _cosmosBridge,
        address _bridgeBank
    );

    function initialize(
        address _cosmosBridge,
        address _bridgeBank
    ) public {
        require(!_initialized, "Initialized");

        cosmosBridge = _cosmosBridge;
        bridgeBank = _bridgeBank;
        _initialized = true;

        emit LogContractsRegistered(cosmosBridge, bridgeBank);
    }
}
