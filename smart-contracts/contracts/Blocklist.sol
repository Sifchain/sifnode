pragma solidity 0.5.16;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract Blocklist is Ownable {
  struct UserStruct {
    uint256 index;
  }

  mapping(address => uint256) private _userIndex;
  address[] private _userList;

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
    _userIndex[account] = _userList.length;
    _userList.push(account);

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
    uint rowToDelete = _userIndex[account];
    address keyToMove = _userList[_userList.length-1];
    _userList[rowToDelete] = keyToMove;
    _userIndex[keyToMove] = rowToDelete; 
    _userList.length--;

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
    if(_userList.length == 0) return false;
    if(_userIndex[account] >= _userList.length) return false;

    return _userList[_userIndex[account]] == account;
  }

  function getFullList() public view returns(address[] memory) {
    return _userList;
  }
}