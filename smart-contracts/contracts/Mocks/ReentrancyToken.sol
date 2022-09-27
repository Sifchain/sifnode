// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

import "../CosmosBridge.sol";

contract ReentrancyToken is ERC20 {
  CosmosBridge cosmosBridge;

  constructor(
    string memory _name,
    string memory _symbol,
    address _cosmosBridgeAddress,
    address attackerUser,
    uint256 mintAmount
  ) ERC20(_name, _symbol) {
    cosmosBridge = CosmosBridge(_cosmosBridgeAddress);
    _mint(attackerUser, mintAmount);
  }

  // transfer will try a reentrancy attack
  function transfer(address recipient, uint256 amount) public override returns (bool) {
    bytes32 hashDigest = 0x8a68aee7fbbbed476def7430bacd3579c38ade9eddc9a0597c98f3530f21e918;

    CosmosBridge.ClaimData memory claimData = CosmosBridge.ClaimData(
      "0x736966316e78363530733871397732386632673374397a74787967343875676c64707475777a70616365", // cosmosSender
      1, // cosmosSenderSequence
      payable(0x70997970C51812dc3A010C7d01b50e0d17dc79C8), // ethereumReceiver
      0xa48a285BAb4061e9104EeA29f968b1B801423E32, // tokenAddress
      100, // amount
      "Reentrancy Token", // tokenName
      "RTK", // tokenSymbol
      18, // tokenDecimals
      1, // networkDescriptor
      false, // doublePeg
      1, // nonce
      "" // cosmosDenom
    );

    CosmosBridge.SignatureData[] memory sigData = new CosmosBridge.SignatureData[](3);

    CosmosBridge.SignatureData memory sig1;
    sig1.signer = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;
    sig1._v = 27;
    sig1._r = 0xef89b2121cc5579e7909ac78160d1488a24e1898237ba0dec57056c53ed602ca;
    sig1._s = 0x06eb1a1375a81e26f45987597f97cac20952ca4ab18ac9928c4e269619ee818a;
    sigData[0] = sig1;

    CosmosBridge.SignatureData memory sig2;
    sig2.signer = 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC;
    sig2._v = 27;
    sig2._r = 0x19627f1bbbd3d5ca11b112da9da0e7bf5dd33cceef5581555ec2a3f5b332fe97;
    sig2._s = 0x0fc1c5564a942e24d078ad24a5fefc514be43b0e4077d448f036712cc6a0e039;
    sigData[1] = sig2;

    CosmosBridge.SignatureData memory sig3;
    sig3.signer = 0x90F79bf6EB2c4f870365E785982E1f101E93b906;
    sig3._v = 27;
    sig3._r = 0x48321cc08333eb832c4797a276d317f74b636fda5db8f7ba92604931fbe0f2a8;
    sig3._s = 0x76ca07a2ec6ca237ede8b9ce4574b22bc7e4dec941978bc2fd62853ba28a8d63;
    sigData[2] = sig3;

    // doesn't revert, but user doesn't get the funds either
    cosmosBridge.submitProphecyClaimAggregatedSigs(hashDigest, claimData, sigData);
  }

  function mint(address account, uint256 amount) public {
    _mint(account, amount);
  }
}
