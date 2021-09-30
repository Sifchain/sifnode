pragma solidity 0.5.16;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract Blocklist is Ownable {
  struct UserStruct {
    uint256 index;
  }

  mapping(address => UserStruct) private _userStructs;
  address[] private _userIndex;

  event addedToBlocklist(address indexed account, address by);
  event removedFromBlocklist(address indexed account, address by);

  modifier onlyInBlocklist(address account) {
    require(isBlocklisted(account), "Not in blocklist");
    _;
  }

  modifier onlyNotInBlocklist(address account) {
    require(!isBlocklisted(account), "Already in blocklist");
    _;
  }

  function _addToBlocklist(address account) private onlyNotInBlocklist(account) returns(bool) {
    //_isBlocklisted[account] = true;
    _userStructs[account].index = _userIndex.length;
    _userIndex.push(account);

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
    //_isBlocklisted[account] = false;
    uint rowToDelete = _userStructs[account].index;
    address keyToMove = _userIndex[_userIndex.length-1];
    _userIndex[rowToDelete] = keyToMove;
    _userStructs[keyToMove].index = rowToDelete; 
    _userIndex.length--;

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

  function isBlocklisted(address account) public view returns(bool) {
    //return _isBlocklisted[account];
    if(_userIndex.length == 0) return false;

    return _userIndex[_userStructs[account].index] == account;
  }
}