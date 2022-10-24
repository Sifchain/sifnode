pragma solidity 0.5.16;

interface IBlocklist {
    function isBlocklisted(address account) external view returns(bool);
}

