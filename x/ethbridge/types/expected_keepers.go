package types

import (
	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(sdk.Context, sdk.AccAddress) authtypes.AccountI
	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
}

// BankKeeper defines the expected supply keeper
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// OracleKeeper defines the expected oracle keeper
type OracleKeeper interface {
	ProcessClaim(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, prophecyID []byte, address string) (oracletypes.StatusText, error)
	GetProphecy(ctx sdk.Context, prophecyID []byte) (oracletypes.Prophecy, bool)
	GetProphecyInfo(ctx sdk.Context, prophecyID []byte) (oracletypes.ProphecyInfo, bool)
	ProcessUpdateWhiteListValidator(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, cosmosSender sdk.AccAddress, validator sdk.ValAddress, power uint32) error
	IsAdminAccount(ctx sdk.Context, cosmosSender sdk.AccAddress) bool
	GetAdminAccount(ctx sdk.Context) sdk.AccAddress
	SetAdminAccount(ctx sdk.Context, cosmosSender sdk.AccAddress)
	ProcessSignProphecy(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, prophecyID []byte, cosmosSender, tokenAddress, ethereumAddress, signature string) error
	GetCrossChainFee(ctx sdk.Context, networkIdentity oracletypes.NetworkIdentity) (string, error)
	GetCrossChainFeeConfig(ctx sdk.Context, networkIdentity oracletypes.NetworkIdentity) (oracletypes.CrossChainFeeConfig, error)
	SetFeeInfo(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, feeCurrency string, gas, burnCost, lockCost, firstBurnDoublePeggyCost sdk.Int) error
	SetProphecyWithInitValue(ctx sdk.Context, prophecyID []byte)
	SetProphecyInfo(ctx sdk.Context, prophecyID []byte, networkDescriptor oracletypes.NetworkDescriptor, cosmosSender string,
		cosmosSenderSequence uint64, ethereumReceiver, tokenDenomHash, tokenContractAddress string,
		tokenAmount sdk.Int, crosschainFee sdk.Int, bridgeToken bool, globalNonce uint64,
		tokenDecimal uint8, tokenName, tokenSymbol string) error
	GetProphecyInfoWithScopeGlobalSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, startGlobalNonce uint64) []*oracletypes.ProphecyInfo
	GetWitnessLockBurnSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) uint64
	ProcessUpdateConsensusNeeded(ctx sdk.Context, cosmosSender sdk.AccAddress, networkDescriptor oracletypes.NetworkDescriptor, ConsensusNeeded uint32) error
}

type AdminKeeper interface {
	IsAdminAccount(ctx sdk.Context, moduleName admintypes.AdminType, adminAccount sdk.AccAddress) bool
}
