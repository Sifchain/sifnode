// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";

/**
 * @title FailHardToken
 * @dev This will always revert after having worked just fine once
 **/
contract FailHardToken is ERC20Burnable {
  uint8 private _decimals;
  bool hasTransferredOnce = false;
  bool hasTransferredFromOnce = false;

  constructor(
    string memory _name,
    string memory _symbol,
    address _user,
    uint256 _amountToMint
  ) ERC20(_name, _symbol) {
    _mint(_user, _amountToMint);
  }

  function name() public view override returns (string memory) {
    revert();
  }

  function symbol() public view override returns (string memory) {
    revert();
  }

  function decimals() public view override returns (uint8) {
    revert();
  }

  function totalSupply() public view override returns (uint256) {
    revert();
  }

  //function balanceOf() public view returns (uint256) {
  //	revert();
  //}

  function cosmosDenom() public view returns (string memory) {
    revert();
  }

  function transfer(address to, uint256 amount) public override returns (bool) {

    revert();
  }

  function transferFrom(
    address from,
    address to,
    uint256 value
  ) public override returns (bool) {
    if (!hasTransferredFromOnce) {

      _transfer(from, to, value);
      hasTransferredFromOnce = true;
      return true;
    }


    revert();
  }

  function mint(address user, uint256 amount) external returns (bool) {
    revert();
  }

  function burn(address user, uint256 amount) external returns (bool) {
    revert();
  }

  function burnFrom(address user, uint256 amount) public override {
    revert();
  }

  function setDenom(string calldata denom) external returns (bool) {
    revert();
  }
}
