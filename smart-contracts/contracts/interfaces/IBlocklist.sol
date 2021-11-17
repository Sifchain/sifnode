pragma solidity 0.8.0;

interface IBlocklist {
	function isBlocklisted(address account) external view returns (bool);
}
