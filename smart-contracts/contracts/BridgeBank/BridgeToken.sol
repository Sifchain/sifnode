// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

/**
 * @title BridgeToken
 * @dev Mintable, ERC20Burnable, ERC20 compatible BankToken for use by BridgeBank
 **/
contract BridgeToken is ERC20Burnable, AccessControl {
  bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

  /**
   * @dev Number of decimals this token uses
   */
  uint8 private _decimals;

  /**
   * @dev The Cosmos denom of this token
   */
  string public cosmosDenom;

  constructor(
    string memory _name,
    string memory _symbol,
    uint8 _tokenDecimals,
    string memory _cosmosDenom
  ) ERC20(_name, _symbol) {
    _decimals = _tokenDecimals;
    cosmosDenom = _cosmosDenom;
    _setupRole(DEFAULT_ADMIN_ROLE, msg.sender);
    _setupRole(MINTER_ROLE, msg.sender);
  }

  /**
   * @notice If sender is a Minter, mints `amount` to `user`
   * @param user Address of the recipient
   * @param amount How much should be minted
   * @return true if the operation succeeds
   */
  function mint(address user, uint256 amount) external onlyRole(MINTER_ROLE) returns (bool) {
    _mint(user, amount);
    return true;
  }

  /**
   * @notice Number of decimals this token has
   */
  function decimals() public view override returns (uint8) {
    return _decimals;
  }

  /**
   * @notice Sets the cosmosDenom
   * @param denom The new cosmos denom
   * @return true if the operation succeeds
   */
  function setDenom(string calldata denom) external onlyRole(DEFAULT_ADMIN_ROLE) returns (bool) {
    cosmosDenom = denom;
    return true;
  }
}
