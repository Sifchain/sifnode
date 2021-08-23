pragma solidity 0.5.16;

interface IBridgeToken {
    function mint(address to, uint256 amount) external;
    function burnFrom(address account, uint256 amount) external;
}
