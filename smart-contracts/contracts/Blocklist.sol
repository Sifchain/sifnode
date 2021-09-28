pragma solidity 0.5.16;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract Blocklist is Ownable {
  mapping(address => bool) private _isBlocklisted;

  event addedToBlocklist(address indexed account, address by);
  event removedFromBlocklist(address indexed account, address by);

  modifier onlyInBlocklist(address account) {
    require(_isBlocklisted[account], "Not in blocklist");
    _;
  }

  modifier onlyNotInBlocklist(address account) {
    require(!_isBlocklisted[account], "Already in blocklist");
    _;
  }

  function _addToBlocklist(address account) private onlyNotInBlocklist(account) returns(bool) {
    _isBlocklisted[account] = true;

    emit addedToBlocklist(account, msg.sender);

    return true;
  }

  function batchAddToBlocklist(address[] memory accounts) public onlyOwner {
    for (uint256 i = 0; i < accounts.length; i++) {
      require(_addToBlocklist(accounts[i]));
    }
  }

  function addToBlocklist(address account) public onlyOwner returns(bool) {
    return _addToBlocklist(account);
  }


  function _removeFromBlocklist(address account) private onlyInBlocklist(account) returns(bool) {
    _isBlocklisted[account] = false;

    emit removedFromBlocklist(account, msg.sender);
    
    return true;
  }

  function batchRemoveFromBlocklist(address[] memory accounts) public onlyOwner {
    for (uint256 i = 0; i < accounts.length; i++) {
      require(_removeFromBlocklist(accounts[i]));
    }
  }

  function removeFromBlocklist(address account) public onlyOwner returns(bool) {
    return _removeFromBlocklist(account);
  }

  function isBlocklisted(address account) external view returns(bool) {
    return _isBlocklisted[account];
  }
}