pragma solidity ^0.5.0;


contract BridgeRegistry {
    address public cosmosBridge;
    address public bridgeBank;
    address public oracle;
    address public valset;

    bool private _initialized;

    uint256[100] private ____gap;

    event LogContractsRegistered(
        address _cosmosBridge,
        address _bridgeBank,
        address _oracle,
        address _valset
    );

    function initialize(
        address _cosmosBridge,
        address _bridgeBank,
        address _oracle,
        address _valset
    ) public {
        require(!_initialized, "Initialized");

        cosmosBridge = _cosmosBridge;
        bridgeBank = _bridgeBank;
        oracle = _oracle;
        valset = _valset;
        _initialized = true;

        emit LogContractsRegistered(cosmosBridge, bridgeBank, oracle, valset);
    }
}
