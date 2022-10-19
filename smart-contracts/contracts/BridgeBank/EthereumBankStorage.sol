// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

/**
 * @title Ethereum Bank Storage
 * @dev Stores nonces, locked tokens, token data (name, symbol, decimals, and denom)
 **/
contract EthereumBankStorage {
  /**
   * @notice current lock and or burn nonce
   */
  uint256 public lockBurnNonce;

  /**
   * @dev {DEPRECATED}
   */
  mapping(address => uint256) private lockedFunds;

  /**
   * @dev {DEPRECATED}
   */
  mapping(string => address) private lockedTokenList;

  /**
   * @dev gap of storage for future upgrades
   */
  uint256[100] private ____gap;

  /**
   * @dev Event emitted when tokens are burned
   */
  event LogBurn(
    address _from,
    bytes _to,
    address _token,
    uint256 _value,
    uint256 indexed _nonce,
    uint8 _decimals,
    int32 _networkDescriptor,
    string _denom
  );

  /**
   * @dev Event emitted when tokens are locked
   */
  event LogLock(
    address _from,
    bytes _to,
    address _token,
    uint256 _value,
    uint256 indexed _nonce,
    uint8 _decimals,
    string _symbol,
    string _name,
    int32 _networkDescriptor
  );

  /**
   * @dev Event emitted when tokens are unlocked
   */
  event LogUnlock(address _to, address _token, uint256 _value);
}
