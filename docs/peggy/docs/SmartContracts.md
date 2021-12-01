# Smart Contracts

A full set of smart contracts is deployed on every supported EVM chain.
Each smart contract has its own EVM address.

## CosmosBridge
- sumbitProphecyClaimAggregatedSigs()

#### verifySignature()

- This function uses `ecrecover` to find the public key of the message signer, like so:
```
function verifySignature(
    address signer,
    bytes32 hashDigest,
    uint8 _v,
    bytes32 _r,
    bytes32 _s
  ) private pure returns (bool) {
    bytes32 messageDigest = keccak256(
      abi.encodePacked("\x19Ethereum Signed Message:\n32", hashDigest)
    );
    return signer == ecrecover(messageDigest, _v, _r, _s);
  }
```

- Then, in order to process a prophecy claim, we check if the resulting public key is a known relayer public key. If it is, we add its validation power to the prophecy.

## BridgeBank
- lock()

## TokenRegistry
- @TODO@ Decribe usage, essential fields and methods

## BridgeToken
- @TODO@ Decribe usage, essential fields and methods
