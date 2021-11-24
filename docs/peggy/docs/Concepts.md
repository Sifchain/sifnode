# Concepts

## Token Denom Hashing for EVM-native assets

To uniquely identify EVM-native assets such as the native currency or ERC20 tokens, we internally use a denomHash for
the relayers and sifnode banking modules. The EVM Native denom is calculated as a sha256 hash of a concatenated string
consisting of the base ten integer (int32) of the Network descriptor (an internal enum for different EVM chains) and the
base 16 token contract address (prefixed with '0x0' and characters converted to lower case) each field separated by a
'/' separator. For the native currency of an EVM chain (Ethereum), the null token address (0x0000000000000000000000000000000000000000)
is used for the token contract address field. After the hash is computed, the string is concatenated with the prefix 'sif' to
comply with cosmos denom requirements. 


### Example 1
We want to calculate denom hash of Ethereum native currency, Ether.
1. We take the network descriptor of '1' (as defined in @TODO@ where) 
1. We take contract address of '0x0000000000000000000000000000000000000000'
1. We calculate the SHA256 of the concatetnated string '1/0x0000000000000000000000000000000000000000'
   ```
   echo -n "1/0x0000000000000000000000000000000000000000" | sha256sum
   ```
1. The result is 'ffd2528d90e15f2ebb0eabe66957aab2a822373b8dda8cd9a47dacb3e6419991':
1. We add a 'sif' prefix and get the denom hash 'sifffd2528d90e15f2ebb0eabe66957aab2a822373b8dda8cd9a47dacb3e6419991'.


### Example 2
We want to calculate denom hash of an ERC20 token 'JimmyToken' on an EVM network.
1. We take network descriptor of '20' (as defined in @TODO@ where)
1. We take the EVM address of deployed token's smart contract, e.g. '0xbF45BFc92ebD305d4C0baf8395c4299bdFCE9EA2' (in lowercase)
1. We calculate the SHA256 of concatenated string '20/0xbf45bfc92ebd305d4c0baf8395c4299bdfce9ea2'
   ```
   echo -n "20/0xbf45bfc92ebd305d4c0baf8395c4299bdfce9ea2" | sha256sum
   ```
1. The result is '490f0dba89dd63c9c766af860241da51f88bd9f1e1e73bd409589b5adb25f695'.
1. We add a 'sif' prefix and get the denom hash 'sif490f0dba89dd63c9c766af860241da51f88bd9f1e1e73bd409589b5adb25f695'.
