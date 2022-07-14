// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.4;

interface IBlocklist {
  function isBlocklisted(address account) external view returns (bool);
}
