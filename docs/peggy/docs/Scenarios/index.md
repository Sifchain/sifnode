# EVM Import/Export flows

Lock and burn operations on tokens and currency native to an EVM chain operate differently than on Cosmos and double
pegged assets. When looking at the design flow of the native EVM tokens and assets, all tokens and currencies on the EVM
side will receive a lock when imported in through the bridgebank or unlock when exported out of the bridgebank. On the
sifnode side, the created assets will be minted when the bridgebank locks or burned when the bridgebank unlocks.

The following flows are supported by Peggy 2.0 when interacting with EVM chains and Sifchain:
 - [Import EVM assets into Sifchain](Scenarios/evmImport)
 - [Export EVM assets out of Sifchain](Scenarios/evmExport)
 - [Export EVM assets into non-native EVM chain from sifchain](Scenarios/doublePegging)
 - [Export Cosmos assets into EVM chain](Scenarios/cosmosExport)
 - [Import Cosmos asstes from EVM chain](Scenarios/cosmosImport)