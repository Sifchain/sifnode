pragma solidity 0.8.17;

interface Pausable {
    function paused() external view returns (bool);
    
    function pause() external;

    function unpause() external;
}
